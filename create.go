package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createAll(client *kubernetes.Clientset) {
	createNamespace(client)
	createDeployment(client)
	createService(client)
}

func createNamespace(client *kubernetes.Clientset) error {
	if _, err := client.CoreV1().Namespaces().Create(context.Background(), namespaced, metav1.CreateOptions{}); err != nil {
		log(ns, "create", "namespace", false)
		return err
	}
	log(ns, "create", "namespace", true)
	return nil
}

func createDeployment(client *kubernetes.Clientset) error {
	if _, err := client.AppsV1().Deployments(namespaced.Name).Create(context.Background(), deployment, metav1.CreateOptions{}); err != nil {
		log(deployment.ObjectMeta.Name, "create", "deployment", false)
		return err
	}
	log(deployment.ObjectMeta.Name, "create", "deployment", true)
	return nil
}

func createService(client *kubernetes.Clientset) error {
	if _, err := client.CoreV1().Services(namespaced.Name).Create(context.Background(), service, metav1.CreateOptions{}); err != nil {
		log(service.ObjectMeta.Name, "create", "service", false)
		return err
	}
	log(service.ObjectMeta.Name, "create", "service", true)
	return nil
}
