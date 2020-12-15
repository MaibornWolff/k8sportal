package main

import (
	"k8sportal/k8sclient"
	"k8sportal/web"
)

func main() {

	web.StartWebserver()

	k8sclient.Inform()
}
