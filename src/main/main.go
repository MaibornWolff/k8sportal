package main

import (
	"k8sportal/k8sclient"
	"k8sportal/web"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	kubeconfig := os.Getenv("KUBECONFIG")

	//initialize kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}

	//create new client with the given config
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	//Get all running services in the Cluster, which have the Label "showOnClusterPortal: true"
	k8sclient.GetServices(kubeClient)

	//TODO Backend

	//TODO Services in DB are getting deleted, if they aren't in the List
	//TODO Services which aren't in the Database are added
	//TODO Get Ingresses of the services
	//TODO ADD FQDN of the Ingress to the DB

	//start the informer factory, to react to changes of services in the cluster
	//TODO
	k8sclient.Inform(kubeClient)

	//TODO Server
	//TODO Get list of running services from mongodb, if a request comes in

	web.StartWebserver()

}
