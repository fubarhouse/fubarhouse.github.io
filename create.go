package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
)

func (client *Application) createAll() {
	client.createNamespace()
	client.createDeployment()
	client.createService()
}

func (client *Application) createNamespace() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e != nil {
		err := client.Create(context.Background(), namespaced)
		log(ns, "create", "namespace", false)
		return err
	}
	log(ns, "create", "namespace", true)
	return nil
}

func (client *Application) createDeployment() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		err := client.Create(context.Background(), deployment)
		log(ns, "create", "deployment", false)
		return err
	}
	log(ns, "create", "deployment", true)
	return nil
}

func (client *Application) createService() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		err := client.Create(context.Background(), service)
		log(ns, "create", "service", false)
		return err
	}
	log(ns, "create", "service", true)
	return nil
}
