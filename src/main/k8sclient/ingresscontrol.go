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
	if _, ok := newIngress.Labels["clusterPortalShow"]; ok {
		log.Info().Msgf("Received ingress to add: %v", newIngress.Name)
	}
}

func onIngUpdate(old interface{}, new interface{}) {
	log.Info().Msgf("Received ingress to update")
	onIngDelete(old)
	onIngAdd(new)
}

func onIngDelete(obj interface{}) {
	deletedIngess := obj.(*networkingv1.Ingress)
	if _, ok := deletedIngess.Labels["clusterPortalShow"]; ok {
		log.Info().Msgf("Received ingress to delete: %v", deletedIngess.Name)
	}
}

func addIngressRules(service *model.Service) {
    storeList := storeIngresses.List()

    for _, obj := range storeList {
         if !labelmatch(obj, annotation) {
             continue
         }
         k8sIngress :=  obj.(*networkingv1.Ingress)
         
         for _, ingressRule := range k8sIngress.Spec.Rules {
             ingressServiceName := ingressRule.HTTP.Paths[0].Backend.Service.Name
             if service.ServiceName !=  ingressServiceName {
                 continue
             }
             service.IngressRules = append(service.IngressRules,
                 model.IngressRule{
                    IngressHost: ingressRule.Host,
                    IngressPath: ingressRule.HTTP.Paths[0].Path})
            service.IngressExists = true
         }
    }
}
