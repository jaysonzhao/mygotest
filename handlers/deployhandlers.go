package handlers

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sDynamicClientFactory struct {
}

type K8sDynamicClient struct {
	Client dynamic.Interface
}

func (c *K8sDynamicClientFactory) GetInClusterClient() *K8sDynamicClient {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		err := fmt.Errorf("error getting kubernetes config: %v\n", err)
		log.Fatal(err.Error())
	}

	fmt.Printf("%T\n", client)
	return &K8sDynamicClient{
		Client: client,
	}
}

func (c *K8sDynamicClientFactory) GetK8sClient() *K8sDynamicClient {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		err := fmt.Errorf("Error getting kubernetes config: %v\n", err)
		log.Fatal(err.Error())
	}
	client, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		err := fmt.Errorf("error getting kubernetes config: %v\n", err)
		log.Fatal(err.Error())
	}

	fmt.Printf("%T\n", client)
	return &K8sDynamicClient{
		Client: client,
	}
}

func (c *K8sDynamicClient) CreateDeployment(name, namespace, image string, port, replicas int32) {
	deploymentRes := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	deploymentObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"replicas": replicas,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": name,
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": name,
						},
					},
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  name,
								"image": image,
								"ports": []map[string]interface{}{
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": port,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating deployment using the unstructured object")
	deployment, err := c.Client.Resource(deploymentRes).Namespace(namespace).Create(context.TODO(), deploymentObject, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create deploymetn: %v\n", err)
	}
	fmt.Printf("Created deployment: %s", deployment.GetName())
}

func (c *K8sDynamicClient) CreateService(name, selector, namespace string, port, targetPort int32) {
	serviceRes := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}

	serviceObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"app": selector,
				},
				"ports": []map[string]interface{}{
					{
						"name":       "http",
						"protocol":   "TCP",
						"port":       port,
						"targetPort": targetPort,
					},
				},
				"type": "ClusterIP",
			},
		},
	}

	service, err := c.Client.Resource(serviceRes).Namespace(namespace).Create(context.TODO(), serviceObject, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create service: %v\n", err)
	}
	fmt.Printf("Service created: %v\n", service.GetName())
}

func (c *K8sDynamicClient) DeleteDeployment(name, namespace string) {
	deploymentRes := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	deploymentObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}

	err := c.Client.Resource(deploymentRes).Namespace(namespace).Delete(context.TODO(), deploymentObject.GetName(), metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Failed to delete deployment: %v\n", err)
	}
	fmt.Printf("Deployment deleted: %v\n", deploymentObject.GetName())
}

func (c *K8sDynamicClient) DeleteService(name, namespace string) {
	serviceRes := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}

	serviceObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}

	err := c.Client.Resource(serviceRes).Namespace(namespace).Delete(context.TODO(), serviceObject.GetName(), metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Failed to delete service: %v\n", err)
	}
	fmt.Printf("Service deleted: %v\n", serviceObject.GetName())
}
