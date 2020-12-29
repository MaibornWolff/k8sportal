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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var mongodbdatabase = "k8sportal"
var mongodbcollection = "portal-services"

//InitServices Returns all services with the label showOnCLusterPortal: true
func InitServices(kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

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

}

//ServiceInform reacts to changed services  TODO Add mongodb client, so changes can be made
func ServiceInform(ctx context.Context, kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

	factory := informers.NewSharedInformerFactory(kubeClient, 0)

	informer := factory.Core().V1().Services().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onSvcAdd(ctx, obj, mongoClient)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			onSvcUpdate(ctx, old, new, mongoClient)
		},
		DeleteFunc: func(obj interface{}) {
			onSvcDelete(ctx, obj, mongoClient)
		},
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onSvcAdd(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	newService := obj.(*corev1.Service)

	log.Info().Msgf("onSvcAdd: Received service to add: %v", newService.Name)

	newServiceLabels := newService.GetLabels()

	if val, ok := newServiceLabels["showOnClusterPortal"]; ok && val == "true" {

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
					log.Fatal().Err(err).Msg("Failed to insert new added service into database")
				}

			} else {
				log.Fatal().Err(result.Err()).Msg("Failed to insert new added service into database")
			}
		}

		log.Info().Msgf("onSvcAdd: Added service %v to database", newService.Name)
	} else {
		log.Info().Msgf("onSvcAdd: Did not add service %v to database, no label or set on false\n", newService.Name)
	}

}

func onSvcUpdate(ctx context.Context, old interface{}, new interface{}, mongoClient *mongo.Client) {

	onSvcDelete(ctx, old, mongoClient)
	onSvcAdd(ctx, new, mongoClient)

}

func onSvcDelete(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	deletedService := obj.(*corev1.Service)

	log.Info().Msgf("onSvcDelete: Received service to delete: %v", deletedService.Name)

	deletedServiceLabels := deletedService.GetLabels()

	if val, ok := deletedServiceLabels["showOnClusterPortal"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection)

		filter := bson.M{"serviceName": deletedService.Name}

		svcFromDatabase := serviceCollection.FindOne(ctx, filter)

		if svcFromDatabase.Err() != nil {
			if svcFromDatabase.Err().Error() == "ErrNoDocuments" {
				log.Info().Msgf("Could not delete service %v from database. Does not exist ", deletedService.Name)
			} else {
				log.Fatal().Err(svcFromDatabase.Err()).Msg("Failed to get service that should be deleted from database")
			}
		} else {

			var svc model.Service
			err := svcFromDatabase.Decode(&svc)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed marshalling service that should be deleted from database")
			}

			if svc.IngressOnline {

				update := bson.M{"$set": bson.M{"serviceOnline": false}}

				_ = serviceCollection.FindOneAndUpdate(ctx, filter, update)

			} else {

				serviceCollection.FindOneAndDelete(ctx, filter)

			}
		}

	} else {
		log.Info().Msgf("onSvcDelete: Did not delete service %v from database, no label or set on false", deletedService.Name)
	}

}
