package main

import (
	"os"

	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	fmt.Printf("Fuck")
	fmt.Printf("My")
	fmt.Printf("Life")
	fmt.Printf("Please")
	/*r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/lastChange", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":80")

	*/
	log.Print("Shared Informer app started")
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Services().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: onUpdate,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper

}

//func onAdd(obj interface{}) {
// Cast the obj as node
//service := obj.(*corev1.Node)
//	_, ok := node.GetLabels()[]
//if ok {
//	fmt.Printf("It has the label!")
//	}
//}
func onUpdate(obj interface{}, obj2 interface{}) {
	// Cast the obj as node
	//service := obj.(*corev1.Service)
	fmt.Printf("Service Changed")

}
