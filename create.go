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
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e != nil {
		if err := client.Client.Create(context.Background(), namespaced); err != nil {
			log(ns, "create", "namespace", false, err)
			return err
		} else {
			log(ns, "create", "namespace", true, nil)
		}
	}
	return nil
}

func (client *Application) createDeployment() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
	}, deployment); e != nil {
		if err := client.Client.Create(context.Background(), deployment); err != nil {
			log(ns, "create", "deployment", false, err)
			return err
		} else {
			log(ns, "create", "deployment", true, nil)
		}
	}
	return nil
}

func (client *Application) createService() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: service.Namespace,
		Name:      service.Name,
	}, deployment); e != nil {
		if err := client.Client.Create(context.Background(), service); err != nil {
			log(ns, "create", "service", false, err)
			return err
		} else {
			log(ns, "create", "service", true, nil)
		}
	}
	return nil
}
