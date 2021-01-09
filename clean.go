package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func cleanNamespace(client *kubernetes.Clientset) error {
	if er := client.CoreV1().Namespaces().Delete(context.Background(), namespaced.Name, metav1.DeleteOptions{}); er != nil {
		log(namespaced.Name, "delete", "namespace", false)
		return er
	}
	log(namespaced.Name, "delete", "namespace", true)
	return nil
}

func cleanDeployment(client *kubernetes.Clientset) error {
	if er := client.AppsV1().Deployments(ns).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{}); er != nil {
		log(deployment.ObjectMeta.Name, "delete", "deployment", false)
		return er
	}
	log(deployment.ObjectMeta.Name, "delete", "deployment", true)
	return nil
}

func cleanService(client *kubernetes.Clientset) error {
	if err := client.CoreV1().Services(ns).Delete(context.Background(), service.Name, metav1.DeleteOptions{}); err != nil {
		log(service.Name, "delete", "service", false)
		return err
	}
	log(service.Name, "delete", "service", true)
	return nil
}
