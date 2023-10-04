package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jaysonzhao/gotest/vo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/labstack/echo"
)

// http://localhost:8000/deploys/
// json: ns=mytest&deploy=testdeploy&image=docker.io/nginx:latest
func AddDeploy(c echo.Context) error {
	deploy := vo.Deploy{}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&deploy)
	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	var clientFactory K8sDynamicClientFactory
	client := clientFactory.GetInClusterClient()
	client.CreateDeployment(deploy.Name, deploy.Namespace, deploy.Image, 80, 2)
	log.Printf("this is yout deploy request %#v", deploy)
	return c.String(http.StatusOK, "We got your deploy request!!!")

}

// http://localhost:8000/pods
func GetPods(c echo.Context) error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := clientset.CoreV1().Pods("haotest").List(context.TODO(), metav1.ListOptions{})
	//clientset.AppsV1().Deployments("haotest").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("internal error: %s", err.Error()),
		})
	}
	podsname := make(map[string]string)
	for i, pod := range pods.Items {
		podsname[fmt.Sprintf("pod%d", i)] = pod.Name
		i++
	}
	return c.JSON(http.StatusBadRequest, podsname)

}

// http://localhost:8000/deploys
func GetDeploys(c echo.Context) error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	//pods, err := clientset.CoreV1().Pods("haotest").List(context.TODO(), metav1.ListOptions{})
	deploys, err := clientset.AppsV1().Deployments("haotest").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("internal error: %s", err.Error()),
		})
	}
	deploysname := make(map[string]string)
	for i, deploy := range deploys.Items {
		deploysname[fmt.Sprintf("Deploy%d", i)] = deploy.Name
		i++
	}
	return c.JSON(http.StatusBadRequest, deploysname)

}

func AddCat(c echo.Context) error {
	cat := vo.Cat{}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&cat)
	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	log.Printf("this is yout cat %#v", cat)
	return c.String(http.StatusOK, "We got your Cat!!!")
}
