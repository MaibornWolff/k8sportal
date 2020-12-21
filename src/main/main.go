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

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	//initialize config from environment

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

	mongoClient, err := mongodb.Connect(ctx, "mongodb://root:dummypw@mongodb.default.svc.cluster.local:27017/?authSource=admin") //TODO Change host/pw to config
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	defer mongoClient.Disconnect(ctx)

	k8sclient.GetServices(kubeClient, mongoClient) //TODO parameterize mongodb

	log.Print("services successfully taken")
	//TODO Backend

	//TODO Services in DB are getting deleted, if they aren't in the List
	//TODO Services which aren't in the Database are added
	//TODO Get Ingresses of the services
	//TODO ADD FQDN of the Ingress to the DB

	//start the informer factory, to react to changes of services in the cluster
	//TODO
	go k8sclient.ServiceInform(ctx, kubeClient, mongoClient)

	//TODO Server
	//TODO Get list of running services from mongodb, if a request comes in

	web.StartWebserver(mongoClient)

}
