package deployment

import (
	"context"
	"fmt"
	"github.com/0xff-dev/k8s-basics/pkg/client"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// NewDeploy
func NewDeploy(namespace, name, image string, labels map[string]string, replicas int32) (*v1.Deployment, error) {
	dep := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
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
	return cli.AppsV1().Deployments(namespace).Create(context.TODO(), &dep, metav1.CreateOptions{})
}

func UpdateDeploy(deploy *v1.Deployment, replicas int32) (*v1.Deployment, error) {
	deploy.Spec.Replicas = &replicas
	cli, err := client.K8sCli()
	if err != nil {
		return deploy, err
	}

	return cli.AppsV1().Deployments(deploy.GetNamespace()).Update(context.TODO(), deploy, metav1.UpdateOptions{})
}


func ListDeployments(namespace string) error {
	cli, err := client.K8sCli()
	if err != nil {
		return err
	}

	deployList, err := cli.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, dep := range deployList.Items {
		fmt.Printf("deployment[%s], replicas: %d\n", dep.GetName(), *dep.Spec.Replicas)
	}

	return nil
}

func GetDeploy(namespace, name string) (*v1.Deployment, error) {
	cli, err := client.K8sCli()
	if err != nil {
		return nil, err
	}
	return cli.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func DelDeployment(namespace, name string) error {
	cli, err := client.K8sCli()
	if err != nil {
		return err
	}

	return cli.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
func DeployExamples(namespace, name, image string, replicas int32) error {
	fmt.Println("Crate Deployment ", name)
	_, err := NewDeploy(namespace, name, image, nil, replicas)
	if err != nil {
		return err
	}
	fmt.Println("create deploy successfully")

	//fmt.Println("Update deploy, set replicas to zero!!")
	//<- time.NewTimer(5*time.Second).C
	//dep, err = UpdateDeploy(dep, 0)
	//if err != nil {
	//	return err
	//}

	fmt.Println("List Deployments")
	<- time.NewTimer(5*time.Second).C
	if err = ListDeployments(namespace); err != nil {
		return err
	}
	fmt.Println("Delete deployment")
	if err = DelDeployment(namespace, name); err != nil {
		return err
	}
	fmt.Println("Get Deployment")
	dep, err := GetDeploy(namespace, name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		fmt.Println("not found deployment ", name)
		return nil
	}
	fmt.Println("dep status: ", dep.Status)
	return nil
}
