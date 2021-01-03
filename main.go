package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	client = &kubernetes.Clientset{}

	//KUBECONFIG = "/home/karl/.kube/config"
	KUBECONFIG = "/Users/karl/.kube/config"

	POD = &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-pod",
			Namespace: "default",
			Labels:    map[string]string{
				"run": "nginx-pod",
			},
		},
		Spec: v1.PodSpec{
			Volumes: nil,
			Containers: []v1.Container{
				{
					Name: "nginx-pod",
					Image: "nginx",
				},
			},
			DNSPolicy: "ClusterFirst",
			RestartPolicy: "Always",
		},
	}

	NAMESPACE = "default"
)

func clean() error {
	fmt.Println("clean invoked.")
	// Check if pod exists.
	old, _ := client.CoreV1().Pods(NAMESPACE).Get(context.TODO(), POD.Name, metav1.GetOptions{})
	if old.Name == POD.Name {
		// If pod does exist, delete it.
		e := client.CoreV1().Pods(NAMESPACE).Delete(context.TODO(), POD.Name, metav1.DeleteOptions{})
		if e != nil {
			// Report error.
			fmt.Println("Error deleting pod")
			return e
		}
		// return pod object.
		return e
	}
	return nil
}

func create() *v1.Pod{
	// Check if pod exists.
	old, _ := client.CoreV1().Pods(NAMESPACE).Get(context.TODO(), POD.Name, metav1.GetOptions{})
	if old.Name != POD.Name {
		// If pod does not exist create it.
		pod, err := client.CoreV1().Pods(NAMESPACE).Create(context.TODO(), POD, metav1.CreateOptions{})
		if err != nil {
			// Report error and return empty object.
			fmt.Println("Error creating pod")
			return &v1.Pod{}
		}
		// return pod object.
		return pod
	}
	return &v1.Pod{}
}

func update() *v1.Pod {
	// Check if pod exists.
	old, _ := client.CoreV1().Pods(NAMESPACE).Get(context.TODO(), POD.Name, metav1.GetOptions{})
	if old.Name == POD.Name {
		// If pod does exist, update it.
		pod, err := client.CoreV1().Pods(NAMESPACE).Update(context.TODO(), POD, metav1.UpdateOptions{})
		if err != nil {
			// Report error and return empty object.
			fmt.Println("Error updating pod")
			return &v1.Pod{}
		}
		// return pod object.
		return pod
	}
	return &v1.Pod{}
}

func main() {

	config, _ := clientcmd.BuildConfigFromFlags("", KUBECONFIG)

	client, _ = kubernetes.NewForConfig(config)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task := sync.WaitGroup{}
	task.Add(1)
	go func() {
		defer clean()
	}()

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		task.Done()
	}

	fmt.Println("listening")

	for {

		create()

	}
	

}
