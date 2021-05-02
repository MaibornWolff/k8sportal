package k8sclient

import (
	"fmt"

	"github.com/rs/zerolog/log"

	networkingv1 "k8s.io/api/networking/v1"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

//IngressInform reacts to changed services
func IngressInform(factory informers.SharedInformerFactory) {

	informer := factory.Networking().V1().Ingresses().Informer()

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
