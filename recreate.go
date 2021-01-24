package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
)

func (client *Application) RecreateAll() {
	client.updateNamespace()
	client.updateDeployment()
	client.updateService()
}

func (client *Application) recreateNamespace() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e != nil {
		client.cleanNamespace()
		namespaced.ResourceVersion = "0"
		client.createNamespace()
	}
	return nil
}

func (client *Application) recreateDeployment() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
	}, deployment); e != nil {
		client.cleanDeployment()
		deployment.ResourceVersion = "0"
		client.createDeployment()
	}
	return nil
}

func (client *Application) recreateService() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: service.Namespace,
		Name:      service.Name,
	}, deployment); e != nil {
		client.cleanService()
		service.ResourceVersion = "0"
		client.createService()
	}
	return nil
}
