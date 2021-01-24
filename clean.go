package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
)

func (client *Application) cleanNamespace() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespaced.Namespace,
		Name:      namespaced.Name,
	}, namespaced); e == nil {
		if err := client.Client.Delete(context.Background(), namespaced); err != nil {
			log(ns, "delete", "namespace", false, err)
			return err
		} else {
			log(ns, "delete", "namespace", true, nil)
		}
	}
	return nil
}

func (client *Application) cleanDeployment() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
	}, deployment); e == nil {
		if err := client.Client.Delete(context.Background(), deployment); err != nil {
			log(ns, "delete", "deployment", false, err)
			return err
		} else {
			log(ns, "delete", "deployment", true, nil)
		}
	}
	return nil
}

func (client *Application) cleanService() error {
	if e := client.Client.Get(context.Background(), types.NamespacedName{
		Namespace: service.Namespace,
		Name:      service.Name,
	}, deployment); e == nil {
		if err := client.Client.Delete(context.Background(), service); err != nil {
			log(ns, "delete", "service", false, err)
			return err
		} else {
			log(ns, "delete", "service", true, nil)
		}
	}
	return nil
}
