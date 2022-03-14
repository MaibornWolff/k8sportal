package k8sclient

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"k8sportal/model"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const annotation string = "clusterPortalShow"

var storeServices cache.Store

var serviceCache serviceCustomCache = serviceCustomCache{
    servicesMap: make(map[string]*model.Service),
}

//ServiceInform reacts to changed services
func ServiceInform(factory informers.SharedInformerFactory) {

	informer := factory.Core().V1().Services().Informer()
	storeServices = informer.GetStore()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onSvcAdd(obj)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			onSvcUpdate(old, new)
		},
		DeleteFunc: func(obj interface{}) {
			onSvcDelete(obj)
		},
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onSvcAdd(obj interface{}) {
	newK8sService := obj.(*corev1.Service)
	if _, ok := newK8sService.Labels["clusterPortalShow"]; !ok {
	       return
	}
    log.Info().Msgf("Received service to add: %v", newK8sService.Name)
    newService := mapToService(newK8sService)
    serviceCache.AddService(newService)

}

func onSvcUpdate(old interface{}, new interface{}) {
	log.Info().Msgf("Received service to update")
	onSvcDelete(old)
	onSvcAdd(new)
}

func onSvcDelete(obj interface{}) {
	deletedK8sService := obj.(*corev1.Service)
	if _, ok := deletedK8sService.Labels["clusterPortalShow"]; !ok {
	       return
	}
    log.Info().Msgf("Received service to delete: %v", deletedK8sService.Name)

    serviceCache.DeleteService(deletedK8sService.Name)

}

func GetAllServices() []*model.Service {
	return serviceCache.ToList()
}

func mapToService(k8sService *corev1.Service) (ret *model.Service) {
    ret = &model.Service{
        ServiceName: k8sService.Name,
        Category: k8sService.Labels["clusterPortalCategory"],
        ServiceExists: true,
        IngressRules: []model.IngressRule{},

    }
   return
}
