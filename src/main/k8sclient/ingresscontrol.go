package k8sclient

import (
	"context"
	"fmt"
	"k8sportal/model"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	networkingv1 "k8s.io/api/networking/v1"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

//InitIngress Returns all services with the label showOnCLusterPortal: true
/*func InitIngress(kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

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

}*/

//IngressInform reacts to changed services  TODO Add mongodb client, so changes can be made
func IngressInform(ctx context.Context, kubeClient kubernetes.Interface, mongoClient *mongo.Client) {

	factory := informers.NewSharedInformerFactory(kubeClient, 0)
	informer := factory.Networking().V1().Ingresses().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onIngAdd(ctx, obj, mongoClient)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			onIngUpdate(ctx, old, new, mongoClient)
		},
		DeleteFunc: func(obj interface{}) {
			onIngDelete(ctx, obj, mongoClient)
		},
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onIngAdd(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	newIngress := obj.(*networkingv1.Ingress)

	newIngressRule := newIngress.Spec.Rules[0]

	newIngressHostName := newIngressRule.Host
	newIngressPath := newIngressRule.HTTP.Paths[0].Path

	newIngressServiceName := newIngressRule.HTTP.Paths[0].Backend.Service.Name

	log.Info().Msgf("Ingress detected")
	log.Info().Msgf("Ingress Host Name:  %v", newIngressHostName)
	log.Info().Msgf("Ingress Path:  %v", newIngressPath)
	log.Info().Msgf("Ingress Service Name:  %v", newIngressServiceName)
	log.Info().Msgf("Let's see if it works")

	newIngressAnnotations := newIngress.GetLabels()
	log.Info().Msgf("onAdd received ingress %v", newIngress.Name)
	log.Info().Msgf("ingresss tags  %v", newIngressAnnotations)

	if val, ok := newIngressAnnotations["showOnClusterPortal"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection)

		filter := bson.M{"serviceName": newIngressServiceName}
		update := bson.M{
			"$set": bson.M{
				"ingressHost":   newIngressHostName,
				"ingressPath":   newIngressPath,
				"ingressOnline": true,
			},
		}

		after := options.After
		upsert := false
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
			Upsert:         &upsert,
		}

		result := serviceCollection.FindOneAndUpdate(ctx, filter, update, &opt)
		if result.Err() != nil {
			if result.Err().Error() == "ErrNoDocuments" {

				ing := model.Service{
					ServiceName:   newIngressServiceName,
					Category:      "",
					ServiceOnline: false,
					IngressHost:   newIngressHostName,
					IngressPath:   newIngressPath,
					IngressOnline: true,
				}

				_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).InsertOne(ctx, ing) //TODO Parameterize
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to insert new added service into mongodb")
				}

			} else {
				log.Fatal().Err(result.Err()).Msg("Failed to insert new added ingress into mongodb")
			}
		}

		log.Info().Msgf("Added the service %v to the database\n", newIngress.Name)
	} else {
		log.Info().Msgf("Did not add service %v to the database, no annotation or set on false\n", newIngress.Name)
	}

}

func onIngUpdate(ctx context.Context, old interface{}, new interface{}, mongoClient *mongo.Client) {
	// Cast the obj as Service
	//service := obj.(*corev1.Service)
	log.Print("Ingress Changed")

}

func onIngDelete(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	deletedIngress := obj.(*networkingv1.Ingress)

	deletedIngressServiceName := deletedIngress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name

	deletedIngressAnnotations := deletedIngress.GetLabels()
	log.Info().Msgf("onDelete Ingress service %v", deletedIngress.Name)
	log.Info().Msgf("Services tags  %v", deletedIngressAnnotations)

	if val, ok := deletedIngressAnnotations["showOnClusterPortal"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection)

		filter := bson.M{"serviceName": deletedIngressServiceName}

		svcFromDatabase := serviceCollection.FindOne(ctx, filter)

		if svcFromDatabase.Err() != nil {
			if svcFromDatabase.Err().Error() == "ErrNoDocuments" {
				log.Info().Msgf("Could not delete service %v from database. Does not exist ", deletedIngressServiceName)
			} else {
				log.Fatal().Err(svcFromDatabase.Err()).Msg("Failed to get service that should be deleted from database")
			}
		} else {

			var svc model.Service
			err := svcFromDatabase.Decode(&svc)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed marshalling service that should be deleted from database")
			}

			if svc.ServiceOnline {

				update := bson.M{
					"$set": bson.M{
						"ingressHost":   "",
						"ingressPath":   "",
						"ingressOnline": false,
					},
				}

				_ = serviceCollection.FindOneAndUpdate(ctx, filter, update)

			} else {

				serviceCollection.FindOneAndDelete(ctx, filter)

			}
		}

	} else {
		log.Info().Msgf("Did not delete ingress to service %v from database, no annotation or set on false\n", deletedIngressServiceName)
	}

}
