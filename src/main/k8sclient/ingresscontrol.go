package k8sclient

import (
	"fmt"

	"github.com/rs/zerolog/log"
    "k8sportal/model"
	networkingv1 "k8s.io/api/networking/v1"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var storeIngresses cache.Store

//IngressInform reacts to changed services
func IngressInform(factory informers.SharedInformerFactory) {

	informer := factory.Networking().V1().Ingresses().Informer()
    storeIngresses = informer.GetStore()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onIngAdd(obj)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			onIngUpdate(old, new)
		},
		DeleteFunc: func(obj interface{}) {
			onIngDelete(obj)
		},
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onIngAdd(obj interface{}) {
	newIngress := obj.(*networkingv1.Ingress)
	if _, ok := newIngress.Labels["clusterPortalShow"]; !ok {
		return
	}
    log.Info().Msgf("Received ingress to add: %v", newIngress.Name)

    for _, ingressRule := range newIngress.Spec.Rules {
        for _, ingressPath := range ingressRule.HTTP.Paths {
            ingressServiceName := ingressPath.Backend.Service.Name
            var service *model.Service
            if tmp, ok := serviceCache.GetService(ingressServiceName); ok {
                service = tmp
            } else {
                log.Info().Msgf("Ingress received without corresponding service")
                service = &model.Service{
                    ServiceName: ingressServiceName,
                    IngressRules: make([]model.IngressRule, 0, 5),
                }
            }
            service.IngressRules = append(service.IngressRules, model.IngressRule{
                IngressHost: ingressRule.Host,
                IngressPath: ingressPath.Path,
            })
            service.IngressExists = true
            serviceCache.AddService(service)
        }
    }
}

func onIngUpdate(old interface{}, new interface{}) {
	log.Info().Msgf("Received ingress to update")
	onIngDelete(old)
	onIngAdd(new)
}

func onIngDelete(obj interface{}) {
	deletedIngress := obj.(*networkingv1.Ingress)
	if _, ok := deletedIngress.Labels["clusterPortalShow"]; !ok {
		return
	}
    log.Info().Msgf("Received ingress to delete: %v", deletedIngress.Name)
    for _, ingressRule := range deletedIngress.Spec.Rules {

        for _, ingressPath := range ingressRule.HTTP.Paths {

            ingressServiceName := ingressPath.Backend.Service.Name
            if service, ok := serviceCache.GetService(ingressServiceName); ok {
                deletedIngressRule := model.IngressRule{
                    IngressHost: ingressRule.Host,
                    IngressPath: ingressPath.Path,
                }

                rules := service.IngressRules
                rules = deleteIngressRule(rules, deletedIngressRule)
                service.IngressRules = rules
                service.IngressExists = len(rules) > 0
                serviceCache.AddService(service)
            }
        }
    }
}

func deleteIngressRule(rules []model.IngressRule, toDeleteRule model.IngressRule) []model.IngressRule {
    i := 0
    for _, rule := range rules {
        if toDeleteRule != rule {
            rules[i] = rule
            i++
        }
    }
    rules = rules[:i]
    return rules
}
