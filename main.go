package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	ns = "awesome-o"
)

var (
	kubeconfig = os.Getenv("HOME") + "/.kube/config"

	replicas   = int32(1)
	deployment = &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-awesome-deployment",
			Namespace: ns,
			Labels: map[string]string{
				"app": "nginx-awesome",
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "nginx",
							Image:           "nginx",
							ImagePullPolicy: "Always",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-service",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Type: "LoadBalancer",
			Ports: []corev1.ServicePort{
				{
					Port: 80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
		},
	}
)

func log(name string, action string, category string, state bool) {
	if !state {
		fmt.Printf("could not %v %v %v\n", action, category, name)
	} else {
		fmt.Printf("%v %v was %vd\n", name, category, action)
	}
}

func clean(client *kubernetes.Clientset) error {
	if _, e := client.AppsV1().Deployments(ns).Get(context.TODO(), deployment.ObjectMeta.Name, metav1.GetOptions{}); e == nil {
		if er := client.AppsV1().Deployments(ns).Delete(context.TODO(), deployment.ObjectMeta.Name, metav1.DeleteOptions{}); er != nil {
			log(deployment.ObjectMeta.Name, "delete", "deployment", false)
		} else {
			log(deployment.ObjectMeta.Name, "delete", "deployment", true)
		}
	}

	if _, e := client.CoreV1().Services(ns).Get(context.TODO(), service.ObjectMeta.Name, metav1.GetOptions{}); e == nil {
		if er := client.CoreV1().Services(ns).Delete(context.TODO(), service.ObjectMeta.Name, metav1.DeleteOptions{}); er != nil {
			log(service.ObjectMeta.Name, "delete", "service", false)
		} else {
			log(service.ObjectMeta.Name, "delete", "service", true)
		}
	}

	if err := client.CoreV1().Namespaces().Delete(context.TODO(), ns, metav1.DeleteOptions{}); err != nil {
		log(ns, "delete", "namespace", false)
	} else {
		log(ns, "delete", "namespace", true)
	}
	return nil
}

func create(client *kubernetes.Clientset) {
	if namespace, err := client.CoreV1().Namespaces().Get(context.TODO(), ns, metav1.GetOptions{}); err != nil {
		namespace.Name = ns
		if _, err := client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{}); err != nil {
			log(ns, "create", "namespace", false)
		} else {
			log(ns, "create", "namespace", true)
		}
	}

	if _, e := client.AppsV1().Deployments(ns).Get(context.TODO(), deployment.ObjectMeta.Name, metav1.GetOptions{}); e != nil {
		if _, er := client.AppsV1().Deployments(ns).Create(context.TODO(), deployment, metav1.CreateOptions{}); er != nil {
			log(deployment.ObjectMeta.Name, "create", "deployment", false)
		} else {
			log(deployment.ObjectMeta.Name, "create", "deployment", true)
		}
	}

	if _, e := client.CoreV1().Services(ns).Get(context.TODO(), service.ObjectMeta.Name, metav1.GetOptions{}); e != nil {
		if _, er := client.CoreV1().Services(ns).Create(context.TODO(), service, metav1.CreateOptions{}); er != nil {
			log(service.ObjectMeta.Name, "create", "service", false)
		} else {
			log(service.ObjectMeta.Name, "create", "service", true)
		}
	}

}

func main() {

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client, _ := kubernetes.NewForConfig(config)
	fmt.Println("Listening...")

	go func() {
		for {

			create(client)

		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task := sync.WaitGroup{}
	task.Add(1)
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Gracefully closing the controller...\n", sig)
		clean(client)
		task.Done()
	}
}
