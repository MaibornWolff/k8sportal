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
	"k8s.io/client-go/tools/cache"
)

//ServiceInform reacts to changed services
func ServiceInform(ctx context.Context, factory informers.SharedInformerFactory, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	informer := factory.Core().V1().Services().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onSvcAdd(ctx, obj, mongoClient, mongodbDatabase, mongodbCollection)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			onSvcUpdate(ctx, old, new, mongoClient, mongodbDatabase, mongodbCollection)
		},
		DeleteFunc: func(obj interface{}) {
			onSvcDelete(ctx, obj, mongoClient, mongodbDatabase, mongodbCollection)
		},
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onSvcAdd(ctx context.Context, obj interface{}, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	newService := obj.(*corev1.Service)

	log.Info().Msgf("onSvcAdd: Received service to add: %v", newService.Name)

	newServiceLabels := newService.GetLabels()

	if val, ok := newServiceLabels["clusterPortalShow"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbDatabase).Collection(mongodbCollection)

		filter := bson.M{"serviceName": newService.Name}
		update := bson.M{
			"$set": bson.M{
				"serviceExists": true,
				"category":      newServiceLabels["clusterPortalCategory"],
			}}
		after := options.After
		upsert := false
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
			Upsert:         &upsert,
		}

		result := serviceCollection.FindOneAndUpdate(ctx, filter, update, &opt)
		if result.Err() != nil {
			if result.Err().Error() == "mongo: no documents in result" {

				svc := model.Service{
					ServiceName:   newService.Name,
					Category:      newServiceLabels["clusterPortalCategory"],
					ServiceExists: true,
				}

				_, err := serviceCollection.InsertOne(ctx, svc) //TODO Parameterize
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

func onSvcUpdate(ctx context.Context, old interface{}, new interface{}, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	log.Info().Msgf("onIngUpdate: Received Service to update")
	log.Info().Msgf("onIngUpdate: Delete outdated Service")
	onSvcDelete(ctx, old, mongoClient, mongodbDatabase, mongodbCollection)
	log.Info().Msgf("onIngUpdate: Add updated Service")
	onSvcAdd(ctx, new, mongoClient, mongodbDatabase, mongodbCollection)

}

func onSvcDelete(ctx context.Context, obj interface{}, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	deletedService := obj.(*corev1.Service)

	log.Info().Msgf("onSvcDelete: Received service to delete: %v", deletedService.Name)

	deletedServiceLabels := deletedService.GetLabels()

	if val, ok := deletedServiceLabels["clusterPortalShow"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbDatabase).Collection(mongodbCollection)

		filter := bson.M{"serviceName": deletedService.Name}

		serviceFromDatabase := serviceCollection.FindOne(ctx, filter)

		if serviceFromDatabase.Err() != nil {
			if serviceFromDatabase.Err().Error() == "mongo: no documents in result" {
				log.Info().Msgf("Could not delete service %v from database. Does not exist ", deletedService.Name)
			} else {
				log.Fatal().Err(serviceFromDatabase.Err()).Msg("Failed to get service that should be deleted from database")
			}
		} else {

			var svc model.Service
			err := serviceFromDatabase.Decode(&svc)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed marshalling service that should be deleted from database")
			}

			if svc.IngressExists {

				update := bson.M{"$set": bson.M{"serviceExists": false}}

				_ = serviceCollection.FindOneAndUpdate(ctx, filter, update)

			} else {

				serviceCollection.FindOneAndDelete(ctx, filter)

			}
		}

	} else {
		log.Info().Msgf("onSvcDelete: Did not delete service %v from database, no label or set on false", deletedService.Name)
	}

}
