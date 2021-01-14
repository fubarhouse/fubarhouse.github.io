package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Application) UpdateAll() {
	client.updateNamespace()
	client.updateDeployment()
	client.updateService()
}

func (client *Application) updateNamespace() error {
	if _, err := client.CoreV1().Namespaces().Update(context.Background(), namespaced, metav1.UpdateOptions{}); err == nil {
		log(ns, "update", "namespace", false)
		return err
	} else {
		return client.createNamespace()
	}
	log(ns, "update", "namespace", true)
	return nil
}

func (client *Application) updateDeployment() error {
	if _, err := client.AppsV1().Deployments(namespaced.Name).Update(context.Background(), deployment, metav1.UpdateOptions{}); err == nil {
		log(deployment.ObjectMeta.Name, "update", "deployment", false)
		return err
	} else {
		return client.createDeployment()
	}
	log(deployment.ObjectMeta.Name, "update", "deployment", true)
	return nil
}

func (client *Application) updateService() error {
	if _, err := client.CoreV1().Services(namespaced.Name).Update(context.Background(), service, metav1.UpdateOptions{}); err == nil {
		log(service.ObjectMeta.Name, "update", "service", false)
		return err
	} else {
		return client.createService()
	}
	log(service.ObjectMeta.Name, "update", "service", true)
	return nil
}
