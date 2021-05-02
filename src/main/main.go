package main

import (
	"k8sportal/k8sclient"
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

	factory := informers.NewSharedInformerFactory(kubeClient, 0)

	//start the informer to react to changes of services in the cluster
	go k8sclient.ServiceInform(factory)
	go k8sclient.IngressInform(factory)

	web.StartWebserver()

}
