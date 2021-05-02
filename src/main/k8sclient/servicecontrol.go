package k8sclient

import (
	"fmt"

	"github.com/rs/zerolog/log"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var store cache.Store

//ServiceInform reacts to changed services
func ServiceInform(factory informers.SharedInformerFactory) {

	informer := factory.Core().V1().Services().Informer()
	store = informer.GetStore()
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

func GetAllServices() (list []interface{}) {
	storelist := store.List()
	list = filter(storelist, labelmatch)
	return
}

func filter(list []interface{}, match func(v1.Object, string) bool) (ret []interface{}) {
	for _, obj := range list {
		// provides an object interface that allows us to filter for metadata easily
		v1obj := obj.(v1.Object)
		// TOOD should be constant
		input := "clusterPortalShow"
		if match(v1obj, input) {
			ret = append(ret, obj)
		}
	}
	return
}

func labelmatch(obj v1.Object, input string) (result bool) {
	result = false
	mapresult := obj.GetLabels()
	log.Info().Msgf("mapresult from getLabels: %v", mapresult)
	if mapresult[input] == "true" {
		log.Info().Msgf("mapresult matches %s: %v", input, mapresult)
		result = true
	}
	return
}
