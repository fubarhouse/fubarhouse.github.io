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
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e == nil {
		if err := client.Client.Update(context.Background(), namespaced); err != nil {
			log(ns, "update", "namespace", false, err)
			return err
		} else {
			log(ns, "update", "namespace", true, nil)
		}
	} else {
		client.cleanNamespace()
		client.createNamespace()
	}
	return nil
}

func (client *Application) updateDeployment() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
	}, deployment); e == nil {
		if err := client.Client.Update(context.Background(), deployment); err != nil {
			log(ns, "update", "deployment", false, err)
			return err
		} else {
			log(ns, "update", "deployment", true, nil)
		}
	} else {
		client.cleanDeployment()
		client.createDeployment()
	}
	return nil
}

func (client *Application) updateService() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: service.Namespace,
		Name:      service.Name,
	}, deployment); e == nil {
		if err := client.Client.Update(context.Background(), service); err != nil {
			log(ns, "update", "service", false, err)
			return err
		} else {
			log(ns, "update", "service", true, nil)
		}
	} else {
		client.cleanService()
		client.createService()
	}
	return nil
}
