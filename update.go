package main

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
)

func (client *Application) UpdateAll() {
	client.updateNamespace()
	client.updateDeployment()
	client.updateService()
}

func (client *Application) updateNamespace() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e != nil {
		err := client.Update(context.Background(), namespaced)
		log(ns, "update", "namespace", false)
		return err
	}
	log(ns, "update", "namespace", true)
	return nil
}

func (client *Application) updateDeployment() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		err := client.Update(context.Background(), deployment)
		log(ns, "update", "deployment", false)
		return err
	}
	log(ns, "update", "deployment", true)
	return nil
}

func (client *Application) updateService() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		err := client.Update(context.Background(), service)
		log(ns, "update", "service", false)
		return err
	}
	log(ns, "update", "service", true)
	return nil
}
