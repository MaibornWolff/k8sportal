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

	log.Info().Msgf("onIngAdd: Received ingress to add: %v", newIngress.Name)

	newIngressLabels := newIngress.GetLabels()

	if val, ok := newIngressLabels["showOnClusterPortal"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection)

		newIngressRules := newIngress.Spec.Rules

		for _, ingressRule := range newIngressRules {

			if ingressRule.HTTP.Paths[0].Backend.Resource == nil {

				ingressRuleHostName := ingressRule.Host
				ingressRulePath := ingressRule.HTTP.Paths[0].Path
				ingressRuleServiceName := ingressRule.HTTP.Paths[0].Backend.Service.Name

				filter := bson.M{"serviceName": ingressRuleServiceName}
				update := bson.M{
					"$set": bson.M{
						"ingressHost":   ingressRuleHostName,
						"ingressPath":   ingressRulePath,
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
							ServiceName:   ingressRuleServiceName,
							Category:      "",
							ServiceOnline: false,
							IngressHost:   ingressRuleHostName,
							IngressPath:   ingressRulePath,
							IngressOnline: true,
						}

						_, err := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection).InsertOne(ctx, ing) //TODO Parameterize
						if err != nil {
							log.Fatal().Err(err).Msg("Failed to insert new added service into database")
						}

					} else {
						log.Fatal().Err(result.Err()).Msg("Failed to insert new added ingress into database")
					}
				}

				log.Info().Msgf("OnIngAdd: Added ingress rule for service %v to database", ingressRuleServiceName)

			} else {
				log.Info().Msgf("OnIngAdd: Did not add ingress rule to database, backend is resource")
			}

		}

	} else {
		log.Info().Msgf("onIngAdd: Did not add rules of ingress %v to database, no label or set on false", newIngress.Name)
	}
}

func onIngUpdate(ctx context.Context, old interface{}, new interface{}, mongoClient *mongo.Client) {

	log.Info().Msgf("onIngUpdate: Received ingress to update")
	log.Info().Msgf("onIngUpdate: Delete outdated Ingress")
	onIngDelete(ctx, old, mongoClient)
	log.Info().Msgf("onIngUpdate: Add updated Ingress")
	onIngAdd(ctx, new, mongoClient)

}

func onIngDelete(ctx context.Context, obj interface{}, mongoClient *mongo.Client) {

	deletedIngress := obj.(*networkingv1.Ingress)

	log.Info().Msgf("onIngDelete: Received ingress to delete: %v", deletedIngress.Name)

	deletedIngressLabels := deletedIngress.GetLabels()

	if val, ok := deletedIngressLabels["showOnClusterPortal"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbdatabase).Collection(mongodbcollection)

		deletedIngressRules := deletedIngress.Spec.Rules

		for _, ingressRule := range deletedIngressRules {

			if ingressRule.HTTP.Paths[0].Backend.Resource == nil {

				ingressRuleServiceName := ingressRule.HTTP.Paths[0].Backend.Service.Name

				filter := bson.M{"serviceName": ingressRuleServiceName}

				svcFromDatabase := serviceCollection.FindOne(ctx, filter)

				if svcFromDatabase.Err() != nil {
					if svcFromDatabase.Err().Error() == "ErrNoDocuments" {
						log.Info().Msgf("Could not delete service %v from database. Does not exist ", ingressRuleServiceName)
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
				log.Info().Msgf("OnIngDelete: Did not delete ingress rule from database, backend is resource")
			}
		}

	} else {
		log.Info().Msgf("onIngDelete: Did not delete ingress %v from database, no label or set on false", deletedIngress.Name)
	}

}
