package main

import (
	"context"
	"k8sportal/k8sclient"
	"k8sportal/mongodb"
	"k8sportal/web"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	err := envconfig.Init(&config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	//set loglevel
	logLevel, err := zerolog.ParseLevel(strings.ToLower(config.LogLevel))

	zerolog.SetGlobalLevel(logLevel)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse log level")
	}

	//create kubernetes client
	kubeconfig := os.Getenv("KUBECONFIG")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build kubeConfig")
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build kubeClient")
	}

	//create new mongodb client
	ctx := context.Background()

	mongoClient, err := mongodb.Connect(ctx, config.Mongodb.Host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	defer mongoClient.Disconnect(ctx)

	err = mongoClient.Database(config.Mongodb.Database).Collection(config.Mongodb.Collection).Drop(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to drop up k8sportal collection in mongodb")
	}
	err = mongoClient.Database(config.Mongodb.Database).CreateCollection(ctx, config.Mongodb.Collection)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create k8sportal collection in mongodb")
	}

	factory := informers.NewSharedInformerFactory(kubeClient, 0)

	//start the informer to react to changes of services in the cluster
	go k8sclient.ServiceInform(ctx, factory, mongoClient, config.Mongodb.Database, config.Mongodb.Collection)
	go k8sclient.IngressInform(ctx, factory, mongoClient, config.Mongodb.Database, config.Mongodb.Collection)

	web.StartWebserver(ctx, mongoClient, config.Mongodb.Database, config.Mongodb.Collection)

}
