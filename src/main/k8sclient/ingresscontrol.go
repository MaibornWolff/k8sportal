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

//IngressInform reacts to changed services  TODO Add mongodb client, so changes can be made
func IngressInform(ctx context.Context, kubeClient kubernetes.Interface, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	factory := informers.NewSharedInformerFactory(kubeClient, 0)
	informer := factory.Networking().V1().Ingresses().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onIngAdd(ctx, obj, mongoClient, mongodbDatabase, mongodbCollection)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			onIngUpdate(ctx, old, new, mongoClient, mongodbDatabase, mongodbCollection)
		},
		DeleteFunc: func(obj interface{}) {
			onIngDelete(ctx, obj, mongoClient, mongodbDatabase, mongodbCollection)
		},
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onIngAdd(ctx context.Context, obj interface{}, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	newIngress := obj.(*networkingv1.Ingress)

	log.Info().Msgf("onIngAdd: Received ingress to add: %v", newIngress.Name)

	newIngressLabels := newIngress.GetLabels()

	if val, ok := newIngressLabels["clusterPortalShow"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbDatabase).Collection(mongodbCollection)

		newIngressRules := newIngress.Spec.Rules

		for _, ingressRule := range newIngressRules {

			if ingressRule.HTTP.Paths[0].Backend.Resource == nil {

				ingressRuleHostName := ingressRule.Host
				ingressRulePath := ingressRule.HTTP.Paths[0].Path
				ingressRuleServiceName := ingressRule.HTTP.Paths[0].Backend.Service.Name

				filter := bson.M{"serviceName": ingressRuleServiceName}

				svcFromDatabase := serviceCollection.FindOne(ctx, filter)

				if svcFromDatabase.Err() != nil {
					if svcFromDatabase.Err().Error() == "mongo: no documents in result" {

						ing := model.IngressRule{
							IngressHost: ingressRuleHostName,
							IngressPath: ingressRulePath,
						}

						newSvc := model.Service{
							ServiceName:   ingressRuleServiceName,
							IngressRules:  []model.IngressRule{ing},
							ServiceOnline: false,
							IngressOnline: true,
						}

						_, err := serviceCollection.InsertOne(ctx, newSvc)
						if err != nil {
							log.Fatal().Err(err).Msg("Failed to insert new added service into database")
						}

					} else {
						log.Fatal().Err(svcFromDatabase.Err()).Msg("Failed to insert new added ingress into database")
					}
				} else {

					var svc model.Service
					err := svcFromDatabase.Decode(&svc)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed marshalling service that should be modified in database")
					}

					ing := model.IngressRule{
						IngressHost: ingressRuleHostName,
						IngressPath: ingressRulePath,
					}

					update := bson.M{
						"$set": bson.M{
							"ingressRules":  append(svc.IngressRules, ing),
							"ingressOnline": true,
						}}

					after := options.After
					upsert := false
					opt := options.FindOneAndUpdateOptions{
						ReturnDocument: &after,
						Upsert:         &upsert,
					}

					result := serviceCollection.FindOneAndUpdate(ctx, filter, update, &opt)
					if result.Err() != nil {
						log.Fatal().Err(result.Err()).Msg("Failed to insert Rules in new added ingress into database")
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

func onIngUpdate(ctx context.Context, old interface{}, new interface{}, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	log.Info().Msgf("onIngUpdate: Received ingress to update")
	log.Info().Msgf("onIngUpdate: Delete outdated Ingress")
	onIngDelete(ctx, old, mongoClient, mongodbDatabase, mongodbCollection)
	log.Info().Msgf("onIngUpdate: Add updated Ingress")
	onIngAdd(ctx, new, mongoClient, mongodbDatabase, mongodbCollection)

}

func onIngDelete(ctx context.Context, obj interface{}, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	deletedIngress := obj.(*networkingv1.Ingress)

	log.Info().Msgf("onIngDelete: Received ingress to delete: %v", deletedIngress.Name)

	deletedIngressLabels := deletedIngress.GetLabels()

	if val, ok := deletedIngressLabels["clusterPortalShow"]; ok && val == "true" {

		serviceCollection := mongoClient.Database(mongodbDatabase).Collection(mongodbCollection)

		deletedIngressRules := deletedIngress.Spec.Rules

		for _, ingressRule := range deletedIngressRules {

			if ingressRule.HTTP.Paths[0].Backend.Resource == nil {

				ingressRuleHostName := ingressRule.Host
				ingressRulePath := ingressRule.HTTP.Paths[0].Path
				ingressRuleServiceName := ingressRule.HTTP.Paths[0].Backend.Service.Name

				filter := bson.M{"serviceName": ingressRuleServiceName}

				svcFromDatabase := serviceCollection.FindOne(ctx, filter)

				if svcFromDatabase.Err() != nil {
					if svcFromDatabase.Err().Error() == "mongo: no documents in result" {
						log.Info().Msgf("Could not delete Ingress Rule from database. Backend Service does not exist ")
					} else {
						log.Fatal().Err(svcFromDatabase.Err()).Msg("Failed to get service that should be deleted from database")
					}
				} else {

					var svc model.Service
					err := svcFromDatabase.Decode(&svc)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed marshalling service that should be deleted from database")
					}

					ruleToBeRemoved := model.IngressRule{
						IngressHost: ingressRuleHostName,
						IngressPath: ingressRulePath,
					}

					modifiedRulesSlice := removeRule(svc.IngressRules, ruleToBeRemoved)

					if len(modifiedRulesSlice) == 0 {

						if svc.ServiceOnline {

							update := bson.M{
								"$set": bson.M{
									"ingressRules":  modifiedRulesSlice,
									"ingressOnline": false,
								}}

							_ = serviceCollection.FindOneAndUpdate(ctx, filter, update)

						} else {
							serviceCollection.FindOneAndDelete(ctx, filter)
						}

					} else {

						update := bson.M{
							"$set": bson.M{
								"ingressRules": modifiedRulesSlice,
							}}

						_ = serviceCollection.FindOneAndUpdate(ctx, filter, update)

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

func removeRule(l []model.IngressRule, rule model.IngressRule) []model.IngressRule {
	index := linearSearch(l, rule)
	if index != -1 {
		return append(l[:index], l[index+1:]...)
	}
	return l

}

func linearSearch(l []model.IngressRule, rule model.IngressRule) int {
	for i, n := range l {
		if n == rule {
			return i
		}
	}
	return -1
}

func testEq(a, b []model.IngressRule) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
