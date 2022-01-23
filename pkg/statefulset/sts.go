package statefulset

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/0xff-dev/k8s-basics/pkg/client"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewStatefulSet(namespace, name, image string, replicas int32) (*v1.StatefulSet, error) {
	sts := v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name: name,
		},
		Spec: v1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"sts-app": name,
				},
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"sts-app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: image,
							Name: name,
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
	return cli.AppsV1().StatefulSets(namespace).Create(context.TODO(), &sts, metav1.CreateOptions{})
}

func UpdateStatefulSet(sts *v1.StatefulSet, replicas int32) (*v1.StatefulSet, error) {
	cli, err := client.K8sCli()
	if err != nil {
		return nil, err
	}

	sts.Spec.Replicas = &replicas
	return cli.AppsV1().StatefulSets(sts.GetNamespace()).Update(context.TODO(), sts, metav1.UpdateOptions{})
}

func ListStatefulSets(namespace string) error {
	cli, err := client.K8sCli()
	if err != nil {
		return err
	}

	stsList, err := cli.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, sts := range stsList.Items {
		fmt.Printf("statefulset[%s], creattime: %s", sts.GetName(), sts.CreationTimestamp)
	}

	return nil
}

func GetStatefulSet(namespace, name string) (*v1.StatefulSet, error) {
	cli, err := client.K8sCli()
	if err != nil {
		return nil, err
	}

	return cli.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func DelStatefulSet(namespace, name string) error {
	cli, err := client.K8sCli()
	if err != nil {
		return err
	}

	return cli.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func StsExample(namespace, name, image string, replicas int32) error {
	fmt.Println("Create StatefulSet ", name)
	sts, err := NewStatefulSet(namespace, name, image, replicas)
	if err != nil {
		return err
	}
	fmt.Println("Create StatefulSet successfully")
	fmt.Printf("Update StatefulSet replicas to 2...")
	sts, err = UpdateStatefulSet(sts, 2)
	if err != nil {
		return err
	}
	fmt.Println("Update StatefulSet successfully, replicas is: ", *sts.Spec.Replicas)

	fmt.Println("List StatefulSets...")
	if err = ListStatefulSets(namespace); err != nil {
		return err
	}

	fmt.Println("Get StatefulSet ", name)
	sts, err = GetStatefulSet(namespace,name)
	if err != nil {
		return err
	}
	fmt.Printf("Get StatefulSet %s successfully, craetetime is: %s", sts.GetName(), sts.CreationTimestamp)
	fmt.Println("Del StatefulSet ", name)
	if err = DelStatefulSet(namespace, name); err != nil {
		return err
	}
	fmt.Println("Del StatefulSet successfully, try get it again...")
	sts, err = GetStatefulSet(namespace, name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		fmt.Println("not found statefulset ", name)
	}

	return nil
}
