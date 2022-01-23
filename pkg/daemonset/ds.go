package daemonset

import (
	"context"
	"fmt"

	"github.com/0xff-dev/k8s-basics/pkg/client"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDaemonSet(namespace, name, image string, labels map[string]string) (*v1.DaemonSet, error) {
	ds := v1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels:    labels,
		},
		Spec: v1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"ds-app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"ds-app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
						},
					},
				},
			},
		},
	}
	cli, err := client.K8sCli()
	if err != nil {
		return nil, err
	}

	return cli.AppsV1().DaemonSets(namespace).Create(context.TODO(), &ds, metav1.CreateOptions{})
}

func UpdateDaemonSet(ds *v1.DaemonSet, image string) (*v1.DaemonSet, error) {
	cli, err := client.K8sCli()
	if err != nil {
		return nil, err
	}

	if len(ds.Spec.Template.Spec.Containers) > 0 {
		for idx := 0; idx < len(ds.Spec.Template.Spec.Containers); idx++ {
			ds.Spec.Template.Spec.Containers[idx].Image = image
		}
	}

	return cli.AppsV1().DaemonSets(ds.GetNamespace()).Update(context.TODO(), ds, metav1.UpdateOptions{})
}

func ListDaemonSets(namespace string) error {
	cli, err := client.K8sCli()
	if err != nil {
		return err
	}

	dsList, err := cli.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ds := range dsList.Items {
		fmt.Printf("daemonset[%s], createtime: %s\n", ds.GetName(), ds.CreationTimestamp)
	}
	return nil
}

func GetDaemonSet(namespace, name string) (*v1.DaemonSet, error) {
	cli, err := client.K8sCli()
	if err != nil {
		return nil, err
	}

	return cli.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func DelDaemonSet(namespace, name string) error {
	cli, err := client.K8sCli()
	if err != nil {
		return err
	}

	return cli.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DSExample(namespace, name, image string) error {
	fmt.Println("Create DaemonSet ", name)
	ds, err := NewDaemonSet(namespace, name, image, nil)
	if err != nil {
		return err
	}
	fmt.Println("Create DaemonSet successfully")
	fmt.Println("Update DaemonSet")

	_, err = UpdateDaemonSet(ds, "tomcat")
	if err != nil {
		return err
	}
	fmt.Println("Update DaemonSet successfully")
	fmt.Println("List all DaemonSet")
	if err = ListDaemonSets(namespace); err != nil {
		return err
	}
	fmt.Println("List all DaemonSet successfully")

	fmt.Println("Get DaemonSet ", name)
	ds, err = GetDaemonSet(namespace, name)
	if err != nil {
		return err
	}
	fmt.Printf("Get DaemonSet %s successfully, createtime  is: %s ", name, ds.GetCreationTimestamp())
	fmt.Println("Del DaemonSet ", name)
	if err = DelDaemonSet(namespace, name); err != nil {
		return err
	}
	fmt.Printf("Del Daemonset %s successfully, Try Get it again...\n", name)
	ds, err = GetDaemonSet(namespace, name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		fmt.Println("not found DaemonSet ", name)
	}
	return nil
}
