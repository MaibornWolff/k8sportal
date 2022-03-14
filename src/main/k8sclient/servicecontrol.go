package k8sclient

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"k8sportal/model"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const annotation string = "clusterPortalShow"

var storeServices cache.Store

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
	newService := obj.(*corev1.Service)
	if _, ok := newService.Labels["clusterPortalShow"]; ok {
		log.Info().Msgf("Received service to add: %v", newService.Name)
	}
}

func onSvcUpdate(old interface{}, new interface{}) {
	log.Info().Msgf("Received service to update")
	onSvcDelete(old)
	onSvcAdd(new)
}

func onSvcDelete(obj interface{}) {
	deletedService := obj.(*corev1.Service)
	if _, ok := deletedService.Labels["clusterPortalShow"]; ok {
		log.Info().Msgf("Received service to delete: %v", deletedService.Name)
	}
}

func GetAllServices() (serviceList []*model.Service) {
	storeList := storeServices.List()
    for _, obj := range storeList {
        if !labelmatch(obj, annotation) {
            continue
        }
        service := mapToService(obj)
        addIngressRules(service)
        serviceList = append(serviceList, service)
    }
	return
}

func filter(list []interface{}, match func(v1.Object, string) bool) (ret []interface{}) {
	for _, obj := range list {
		// provides an object interface that allows us to filter for metadata easily
		v1obj := obj.(v1.Object)
		// TOOD should be constant
		if match(v1obj, annotation) {
			ret = append(ret, obj)
		}
	}
	return
}

func labelmatch(obj interface{}, input string) (result bool) {
    v1Obj := obj.(v1.Object)
	result = false
	mapresult := v1Obj.GetLabels()
	log.Info().Msgf("mapresult from getLabels: %v", mapresult)
	if mapresult[input] == "true" {
		log.Info().Msgf("mapresult matches %s: %v", input, mapresult)
		result = true
	}
	return
}

func mapToService(obj interface{}) (ret *model.Service) {
    k8sService := obj.(*corev1.Service)
    ret = &model.Service{
        ServiceName: k8sService.Name,
        Category: k8sService.Labels["clusterPortalCategory"],
        ServiceExists: true,

    }
   return
}
