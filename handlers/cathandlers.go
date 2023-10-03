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

// http://localhost:8000/cats/json?name=arnold&type=fluffy
func GetCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is : %s\nand cat type is : %s\n", catName, catType))
	} else if dataType == "json" {
		cat := vo.Cat{
			Name: catName,
			Type: catType,
		}
		return c.JSON(http.StatusOK, cat)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the data type as Sting or JSON",
		})
	}

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
