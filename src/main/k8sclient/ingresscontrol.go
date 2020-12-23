package k8sclient

import (
	"context"
	"fmt"
	"k8sportal/model"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

//InitIngress Returns all services with the label showOnCLusterPortal: true
func InitIngress(kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

	options := metav1.ListOptions{
		LabelSelector: "showOnClusterPortal=true",
	}

	ctx := context.Background()

	ingList, err := kubeClient.NetworkingV1().Ingresses("").List(ctx, options)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get running Ingresses from kubernetes cluster")
	}

	if len((*ingList).Items) == 0 {
		log.Info().Msgf("Found no Ingresses to show on portal")
	} else {

		for _, ingInfo := range (*ingList).Items {

			svc := model.Service{
				ServiceName:   ingInfo.Name,
				Category:      "",
				ServiceOnline: true,
				IngressHost:   "",
				IngressPath:   "",
				IngressOnline: false,
			}

			_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).InsertOne(ctx, svc) //TODO Parameterize
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to insert service into mongodb")
			}
			log.Info().Msgf("added the service %v to the database\n", ingInfo.Name)
		}
	}

}

//IngressInform reacts to changed services  TODO Add mongodb client, so changes can be made
func IngressInform(ctx context.Context, kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

	factory := informers.NewSharedInformerFactory(kubeClient, 0)
	informer := factory.Networking().V1().Ingresses().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		/*	AddFunc: func(obj interface{}) {
				onIngAdd(ctx, obj, mongoClient)
			},
			UpdateFunc: onIngUpdate,
			DeleteFunc: func(obj interface{}) {
				onIngDelete(ctx, obj, mongoClient)
			},*/
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onIngAdd(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	service := obj.(*networkingv1.Ingress)

	svc := model.Service{
		ServiceName:   service.Name,
		Category:      "",
		ServiceOnline: true,
		IngressHost:   "",
		IngressPath:   "",
		IngressOnline: false,
	}

	_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).InsertOne(ctx, svc) //TODO Parameterize
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to insert new added service into mongodb")
	}
	log.Info().Msgf("added the service %v to the database\n", service.Name)

	log.Print("Ingress Added")

}

func onIngUpdate(old interface{}, new interface{}) {
	// Cast the obj as Service
	//service := obj.(*corev1.Service)
	log.Print("Ingress Changed")

}

func onIngDelete(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	//service := obj.(*corev1.Service)

	//name := service.Name

	//_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).DeleteOne
	log.Print("SIngress Deleted")

}

//TODO onAdd

//TODO onDelte
