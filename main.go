package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"os/signal"
	"sync"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Application struct {
	client.Client
}

var (
	app Application
)

func main() {

	kubernetesHost := os.Getenv("KUBERNETES_HOST")   // https://192.168.99.110:8443
	kubernetesToken := os.Getenv("KUBERNETES_TOKEN") // ""

	if kubernetesHost == "" {
		panic("missing kubernetes host input")
	}

	if kubernetesToken == "" {
		panic("missing kubernetes token input")
	}

	config := &rest.Config{
		BearerToken: string(kubernetesToken),
		Host:        kubernetesHost,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	c, e := client.New(config, client.Options{})
	if e != nil {
		panic(e)
	}

	app.Client = c
	klog.Infoln("controller started")
	app.createAll()

	go func() {
		for {

			if e := c.Get(context.Background(), types.NamespacedName{
				Namespace: namespaced.Namespace,
				Name:      namespaced.Name,
			}, namespaced); e != nil {
				app.updateNamespace()
			}

			if e := c.Get(context.Background(), types.NamespacedName{
				Namespace: namespaced.Namespace,
				Name:      namespaced.Name,
			}, deployment); e != nil {
				app.updateDeployment()
			}

			if e := c.Get(context.Background(), types.NamespacedName{
				Namespace: namespaced.Namespace,
				Name:      namespaced.Name,
			}, service); e != nil {
				app.updateService()
			}

			time.Sleep(time.Second * 1)

		}
	}()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)

	task := sync.WaitGroup{}
	task.Add(1)
	select {
	case sig := <-channel:
		klog.Infof("received %s signal; now terminating\n", sig)
		app.cleanDeployment()
		app.cleanService()
		app.cleanNamespace()
		task.Done()
	}
}
