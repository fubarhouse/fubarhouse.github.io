package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
)

func (client *Application) cleanNamespace() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e != nil {
		err := client.Delete(context.Background(), namespaced)
		log(ns, "delete", "namespace", false)
		return err
	}
	log(ns, "delete", "namespace", true)
	return nil
}

func (client *Application) cleanDeployment() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		err := client.Delete(context.Background(), deployment)
		log(ns, "delete", "deployment", false)
		return err
	}
	log(ns, "delete", "deployment", true)
	return nil
}

func (client *Application) cleanService() error {
	if e := client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, deployment); e != nil {
		err := client.Delete(context.Background(), service)
		log(ns, "delete", "service", false)
		return err
	}
	log(ns, "delete", "service", true)
	return nil
}
