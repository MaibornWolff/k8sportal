package k8sclient

import (
	"context"
	"fmt"
	"k8sportal/model"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var mongodbdatabase = "k8sportal"
var mongodbcollection = "portal-services"

//InitServices Returns all services with the label showOnCLusterPortal: true
/*func InitServices(kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

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

		for _, svcInfo := range (*svcList).Items {

			svc := model.Service{
				ServiceName:   svcInfo.Name,
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
			log.Info().Msgf("added the service %v to the database\n", svcInfo.Name)
		}
	}

}*/

//ServiceInform reacts to changed services  TODO Add mongodb client, so changes can be made
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
		UpdateFunc: func(old interface{}, new interface{}) {
			onUpdate(ctx, old, new, mongoClient)
		},
		DeleteFunc: func(obj interface{}) {
			onDelete(ctx, obj, mongoClient)
		},
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onAdd(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	newService := obj.(*corev1.Service)

	newServiceAnnotations := newService.GetLabels()
	log.Info().Msgf("Received service %v", newService.Name)
	log.Info().Msgf("Services tags  %v", newServiceAnnotations)

	if val, ok := newServiceAnnotations["showOnClusterPortal"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection)

		filter := bson.M{"serviceName": newService.Name}
		update := bson.M{"$set": bson.M{"serviceOnline": true}}

		after := options.After
		upsert := false
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
			Upsert:         &upsert,
		}

		result := serviceCollection.FindOneAndUpdate(ctx, filter, update, &opt)
		if result.Err() != nil {
			if result.Err().Error() == "ErrNoDocuments" {

				svc := model.Service{
					ServiceName:   newService.Name,
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

			} else {
				log.Fatal().Err(result.Err()).Msg("Failed to insert new added service into mongodb")
			}
		}

		log.Info().Msgf("Added the service %v to the database\n", newService.Name)
	} else {
		log.Info().Msgf("Did not add the service %v to the database, no annotation or set on false\n", newService.Name)
	}

}

func onUpdate(ctx context.Context, old interface{}, new interface{}, mongoClient *mongo.Client) {
	//Cast the obj as Service
	//service := obj.(*corev1.Service)
	log.Print("Service Changed\n")

}

func onDelete(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	//service := obj.(*corev1.Service)

	//name := service.Name

	//_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).DeleteOne
	log.Print("Service Deleted\n")

}
