package k8sclient

import (
	"context"
	"fmt"
	"k8sportal/model"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var mongodbdatabase = "k8sportal"
var mongodbcollection = "portal-services"

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

	if len((*svcList).Items) == 0 {
		log.Info().Msgf("Found no services to show on portal")
	} else {

		err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).Drop(ctx) //TODO Parameterize
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to clean up k8sportal collection in mongodb")
		}

		for _, svcInfo := range (*svcList).Items {

			svc := model.Service{svcInfo.Name, "", "", true}

			_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).InsertOne(ctx, svc) //TODO Parameterize
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to insert service into mongodb")
			}
			log.Info().Msgf("added the service %v to the database\n", svcInfo.Name)
		}
	}

}

//serviceInform reacts to changed services  TODO Add mongodb client, so changes can be made
func ServiceInform(ctx context.Context, kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

	factory := informers.NewSharedInformerFactory(kubeClient, 0)
	informer := factory.Core().V1().Services().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onAdd(ctx, obj, mongoClient)
		},
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onAdd(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	service := obj.(*corev1.Service)

	svc := model.Service{service.Name, "", "", true}

	_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).InsertOne(ctx, svc) //TODO Parameterize
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to insert new added service into mongodb")
	}
	log.Info().Msgf("added the service %v to the database\n", service.Name)

	log.Print("Service Added")

}

func onUpdate(old interface{}, new interface{}) {
	// Cast the obj as Service
	//service := obj.(*corev1.Service)
	log.Print("Service Changed")

}

func onDelete(obj interface{}) {
	// Cast the obj as Service
	//service := obj.(*corev1.Service)
	log.Print("Service Deleted")

}

//TODO onAdd

//TODO onDelte
