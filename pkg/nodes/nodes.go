package nodes

import (
	"context"
	"fmt"

	"github.com/0xff-dev/k8s-basics/pkg/client"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Nodes() {
	cli, err := client.K8sCli()
	if err != nil {
		glog.Error(err)
		return
	}

	nodeList, err := cli.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Error(err)
		return
	}
	for _, node := range nodeList.Items {
		fmt.Printf("get node[%s], arch[%s]\n", node.GetName(), node.Status.NodeInfo.Architecture)
	}
}
