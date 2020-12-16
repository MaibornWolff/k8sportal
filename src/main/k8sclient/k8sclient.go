package k8sclient

import (
	"context"
	"fmt"
	"k8sportal/model"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/apimachinery/pkg/util/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

//GetServices Returns all services with the label showOnCLusterPortal: true
func GetServices(kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

	options := metav1.ListOptions{
		LabelSelector: "showOnClusterPortal=true",
	}

	ctx := context.Background()

	svcList, err := kubeClient.CoreV1().Services("").List(ctx, options)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get running services from kubernetes cluster")
	}

	log.Info().Msgf("Found %v services to show on portal", len((*svcList).Items))

	err = mongoClient.Database("k8sportal").Collection("portal-services").Drop(ctx) //TODO Parameterize
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to clean up k8sportal collection in mongodb")
	}

	for _, svcInfo := range (*svcList).Items {

		svc := model.Service{svcInfo.Name, "", "", true}

		_, err = mongoClient.Database("k8sportal").Collection("portal-services").InsertOne(ctx, svc) //TODO Parameterize
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to insert service into mongodb")
		}
		log.Info().Msgf("added the service %v to the database\n", svcInfo.Name)
	}

}

//Inform reacts to changed services  TODO Add mongodb client, so changes can be made
func Inform(kubeClient kubernetes.Interface) {

	factory := informers.NewSharedInformerFactory(kubeClient, 0)
	informer := factory.Core().V1().Services().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: onUpdate,
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onUpdate(old interface{}, new interface{}) {
	// Cast the obj as Service
	//service := obj.(*corev1.Service)
	log.Print("Service Changed")

}

//TODO onAdd

//TODO onDelte
