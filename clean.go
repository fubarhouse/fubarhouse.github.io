package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
)

func (client *Application) cleanNamespace() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e != nil {
		if err := client.Client.Delete(context.Background(), namespaced); err != nil {
			log(ns, "delete", "namespace", false)
			return err
		}
	}
	log(ns, "delete", "namespace", true)
	return nil
}

func (client *Application) cleanDeployment() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		if err := client.Client.Delete(context.Background(), deployment); err != nil {
			log(ns, "delete", "deployment", false)
			return err
		}
	}
	log(ns, "delete", "deployment", true)
	return nil
}

func (client *Application) cleanService() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		if err := client.Client.Delete(context.Background(), service); err != nil {
			log(ns, "delete", "service", false)
			return err
		}
	}
	log(ns, "delete", "service", true)
	return nil
}
