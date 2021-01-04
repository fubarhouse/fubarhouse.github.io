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

const (
	NAMESPACE = "awesome-o"
)

var (
	KUBECONFIG = os.Getenv("HOME") + "/.kube/config"

	POD = &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-awesome-pod",
			Namespace: NAMESPACE,
			Labels: map[string]string{
				"run": "nginx-pod",
			},
		},
		Spec: v1.PodSpec{
			Volumes: nil,
			Containers: []v1.Container{
				{
					Name:  "nginx-pod-awesome",
					Image: "nginx",
				},
			},
			DNSPolicy:     "ClusterFirst",
			RestartPolicy: "Always",
		},
	}
)

func clean(client *kubernetes.Clientset) error {
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
		} else {
			fmt.Println("pod was removed")
		}
	}
	err := client.CoreV1().Namespaces().Delete(context.TODO(), NAMESPACE, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("error deleting namespace", err)
	} else {
		fmt.Println("namespace was removed")
	}
	return nil
}

func create(client *kubernetes.Clientset) *v1.Pod {
	// Check if namespace exists.
	namespace, err := client.CoreV1().Namespaces().Get(context.TODO(), NAMESPACE, metav1.GetOptions{})
	if err != nil {
		namespace.Name = NAMESPACE
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		if err != nil {
			fmt.Println("cannot create namespace", err)
			return &v1.Pod{}
		} else {
			fmt.Println("created namespace")
		}
	}
	// Check if pod exists.
	old, _ := client.CoreV1().Pods(NAMESPACE).Get(context.TODO(), POD.Name, metav1.GetOptions{})
	if old.Name != POD.Name {
		// If pod does not exist create it.
		pod, err := client.CoreV1().Pods(NAMESPACE).Create(context.TODO(), POD, metav1.CreateOptions{})
		if err != nil {
			// Report error and return empty object.
			fmt.Println("Error creating pod", err)
			return &v1.Pod{}
		} else {
			fmt.Println("Pod was created")
		}
		// return pod object.
		return pod
	}
	return &v1.Pod{}
}

func update(client *kubernetes.Clientset) *v1.Pod {
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

	client, _ := kubernetes.NewForConfig(config)

	fmt.Println("listening")

	go func() {
		for {

			create(client)

		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task := sync.WaitGroup{}
	task.Add(1)
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		clean(client)
		task.Done()
	}
}
