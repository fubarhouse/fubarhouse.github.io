package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

const (
	ns = "awesome-o"
)

var (
	kubeconfig = os.Getenv("HOME") + "/.kube/config"

	replicas   = int32(1)
	namespaced = &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}
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
					"app": "nginx-awesome",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx-awesome",
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
			Name: "nginx-awesome-service",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx-awesome",
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
		klog.Warningf("could not %v %v %v\n", action, category, name)
		return
	}
	klog.Infof("%v %v was %vd\n", name, category, action)
}

func cleanNamespace(client *kubernetes.Clientset) error {
	if er := client.CoreV1().Namespaces().Delete(context.Background(), namespaced.Name, metav1.DeleteOptions{}); er != nil {
		log(namespaced.Name, "delete", "namespace", false)
		return er
	}
	log(namespaced.Name, "delete", "namespace", true)
	return nil
}

func cleanDeployment(client *kubernetes.Clientset) error {
	if er := client.AppsV1().Deployments(ns).Delete(context.Background(), deployment.ObjectMeta.Name, metav1.DeleteOptions{}); er != nil {
		log(deployment.ObjectMeta.Name, "delete", "deployment", false)
		return er
	}
	log(deployment.ObjectMeta.Name, "delete", "deployment", true)
	return nil
}

func cleanService(client *kubernetes.Clientset) error {
	if err := client.CoreV1().Services(ns).Delete(context.Background(), service.Name, metav1.DeleteOptions{}); err != nil {
		log(service.Name, "delete", "service", false)
		return err
	}
	log(service.Name, "delete", "service", true)
	return nil
}

func createAll(client *kubernetes.Clientset) {
	createNamespace(client)
	createDeployment(client)
	createService(client)
}

func createNamespace(client *kubernetes.Clientset) error {
	if _, err := client.CoreV1().Namespaces().Create(context.Background(), namespaced, metav1.CreateOptions{}); err != nil {
		log(ns, "create", "namespace", false)
		return err
	}
	log(ns, "create", "namespace", true)
	return nil
}

func createDeployment(client *kubernetes.Clientset) error {
	if _, err := client.AppsV1().Deployments(namespaced.Name).Create(context.Background(), deployment, metav1.CreateOptions{}); err != nil {
		log(deployment.ObjectMeta.Name, "create", "deployment", false)
		return err
	}
	log(deployment.ObjectMeta.Name, "create", "deployment", true)
	return nil
}

func createService(client *kubernetes.Clientset) error {
	if _, err := client.CoreV1().Services(namespaced.Name).Create(context.Background(), service, metav1.CreateOptions{}); err != nil {
		log(service.ObjectMeta.Name, "create", "service", false)
		return err
	}
	log(service.ObjectMeta.Name, "create", "service", true)
	return nil
}

func main() {

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client, _ := kubernetes.NewForConfig(config)
	klog.Infoln("controller started")

	go func() {
		for {

			if _, e := client.CoreV1().Namespaces().Get(context.Background(), namespaced.Name, metav1.GetOptions{}); e != nil {
				createNamespace(client)
			}
			if _, e := client.AppsV1().Deployments(namespaced.Name).Get(context.Background(), deployment.ObjectMeta.Name, metav1.GetOptions{}); e != nil {
				createDeployment(client)
			}
			if _, e := client.CoreV1().Services(namespaced.Name).Get(context.Background(), service.ObjectMeta.Name, metav1.GetOptions{}); e != nil {
				createService(client)
			}

			time.Sleep(time.Second * 1)

		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task := sync.WaitGroup{}
	task.Add(1)
	select {
	case sig := <-c:
		klog.Infof("received %s signal; now terminating\n", sig)
		cleanDeployment(client)
		cleanService(client)
		cleanNamespace(client)
		task.Done()
	}
}
