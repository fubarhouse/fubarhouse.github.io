package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type Application struct {
	kubernetes.Clientset
}

var (
	app Application
)

func main() {

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client, _ := kubernetes.NewForConfig(config)

	app.Clientset = *client
	klog.Infoln("controller started")
	app.createAll()

	go func() {
		for {

			if _, e := client.CoreV1().Namespaces().Get(context.Background(), namespaced.Name, metav1.GetOptions{}); e != nil {
				app.updateNamespace()
			}
			if _, e := client.AppsV1().Deployments(namespaced.Name).Get(context.Background(), deployment.ObjectMeta.Name, metav1.GetOptions{}); e != nil {
				app.updateDeployment()
			}
			if _, e := client.CoreV1().Services(namespaced.Name).Get(context.Background(), service.ObjectMeta.Name, metav1.GetOptions{}); e != nil {
				app.updateService()
			}

			time.Sleep(time.Second * 1)

		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task := sync.WaitGroup{}
	task.Add(1)
	select {
	case sig := <-c:
		klog.Infof("received %s signal; now terminating\n", sig)
		app.cleanDeployment()
		app.cleanService()
		app.cleanNamespace()
		task.Done()
	}
}
