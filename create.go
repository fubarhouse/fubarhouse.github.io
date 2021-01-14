package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Application) createAll() {
	client.createNamespace()
	client.createDeployment()
	client.createService()
}

func (client *Application) createNamespace() error {
	if _, err := client.CoreV1().Namespaces().Create(context.Background(), namespaced, metav1.CreateOptions{}); err != nil {
		log(ns, "create", "namespace", false)
		return err
	}
	log(ns, "create", "namespace", true)
	return nil
}

func (client *Application) createDeployment() error {
	if _, err := client.AppsV1().Deployments(namespaced.Name).Create(context.Background(), deployment, metav1.CreateOptions{}); err != nil {
		log(deployment.ObjectMeta.Name, "create", "deployment", false)
		return err
	}
	log(deployment.ObjectMeta.Name, "create", "deployment", true)
	return nil
}

func (client *Application) createService() error {
	if _, err := client.CoreV1().Services(namespaced.Name).Create(context.Background(), service, metav1.CreateOptions{}); err != nil {
		log(service.ObjectMeta.Name, "create", "service", false)
		return err
	}
	log(service.ObjectMeta.Name, "create", "service", true)
	return nil
}
