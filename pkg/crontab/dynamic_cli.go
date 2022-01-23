package crontab

import (
	"context"
	"fmt"

	"github.com/0xff-dev/k8s-basics/pkg/client"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

/*
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: my-new-cron-object
spec:
  cronSpec: "* * * * *"
image: my-awesome-cron-image
*/
const (
	Group = "stable.example.com"
	Resource = "crontabs"
	Version = "v1"
)
func NewCrontab(namespace, name, image string) (*unstructured.Unstructured, error) {
	cr := map[string]interface{}{
		"apiVersion": fmt.Sprintf("%s/%s", Group, Version),
		"kind": "CronTab",
		"metadata": map[string]interface{}{
			"name": name,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"cronSpec": "* * * * */5",
			"image": "my-crontab-image:v1",
		},
	}
	dcli, err := client.DynamicCli()
	if err != nil {
		return nil, err
	}

	return dcli.Resource(schema.GroupVersionResource{Group: Group, Resource: Resource, Version: Version}).
		Namespace(namespace).Create(context.TODO(), &unstructured.Unstructured{Object: cr}, metav1.CreateOptions{})

}

func UpdateCrontab(cr *unstructured.Unstructured, newImage string) (*unstructured.Unstructured, error) {
	dcli, err := client.DynamicCli()
	if err != nil {
		return nil, err
	}

	if err = unstructured.SetNestedField(cr.Object, newImage, "spec", "image"); err != nil {
		return nil, err
	}

	return dcli.Resource(schema.GroupVersionResource{Group: Group, Resource: Resource, Version: Version}).
		Namespace(cr.GetNamespace()).Update(context.TODO(), cr, metav1.UpdateOptions{})
}

func ListCRs(namespace string) error {
	dcli, err := client.DynamicCli()
	if err != nil {
		return err
	}

	crList, err := dcli.Resource(schema.GroupVersionResource{Group: Group, Resource: Resource, Version: Version}).
		Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cr := range crList.Items {
		fmt.Printf("cr[%s], cratetime: %s", cr.GetName(), cr.GetCreationTimestamp())
	}
	return nil
}

func GetCR(namespace, name string) (*unstructured.Unstructured, error) {
	dcli, err := client.DynamicCli()
	if err != nil {
		return nil, err
	}
	return dcli.Resource(schema.GroupVersionResource{Group: Group, Version: Version, Resource: Resource}).
		Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func DelCR(namespace, name string) error {
	dcli, err := client.DynamicCli()
	if err != nil {
		return err
	}

	return dcli.Resource(schema.GroupVersionResource{Group: Group, Resource: Resource, Version: Version}).
		Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
func CRDExample(namespace, name, image string) error {
	fmt.Println("Create CR ", name)
	cr, err := NewCrontab(namespace, name, image)
	if err != nil {
		return err
	}
	fmt.Printf("Create CR successfully, name: %s, createtime: %s\n", cr.GetName(), cr.GetCreationTimestamp())
	fmt.Println("Update Cr image ...")
	cr, err = UpdateCrontab(cr, "new-image-abc:v2")
	if err != nil {
		return err
	}
	newImage, _, _ := unstructured.NestedString(cr.Object, "spec", "image")
	fmt.Printf("Update Cr successfully, new image is: %s", newImage)
	fmt.Println("List CRs")
	if err = ListCRs(namespace); err != nil {
		return err
	}

	fmt.Println("Get CR")
	cr, err = GetCR(namespace, name)
	if err != nil {
		return err
	}
	fmt.Println("Get Cr successfully name: ", cr.GetName())
	fmt.Println("Delete CR ")
	if err = DelCR(namespace, name); err != nil {
		return err
	}
	fmt.Println("Del CR successfully, try get it again...")
	cr, err = GetCR(namespace, name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		fmt.Println("not found cr ", name)
	}
	return nil
}