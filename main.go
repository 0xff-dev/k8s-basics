package main

import (
	"flag"
	"fmt"
	"github.com/0xff-dev/k8s-basics/pkg/crontab"
	"github.com/0xff-dev/k8s-basics/pkg/daemonset"
	"github.com/0xff-dev/k8s-basics/pkg/deployment"
	"github.com/0xff-dev/k8s-basics/pkg/statefulset"

	"github.com/0xff-dev/k8s-basics/pkg/client"
	"github.com/0xff-dev/k8s-basics/pkg/nodes"
)

var (
	// command args
	nodeExample = flag.Bool("node-example", false, "get all nodes")
	deploy  = flag.Bool("deploy", false, "usage of deployment")
	ds   = flag.Bool("ds", false, "usage of daemonset")
	sts = flag.Bool("sts", false, "usgae of daemonset")
	crd = flag.Bool("crd", false, "usage of crd")
	namespace = flag.String("namespace", "default", "the namespace which you want to use")
	kubeConfPath = flag.String("kubeconf", "", "the path of kubeconfig")
)

func init() {
	flag.Parse()
	client.InitClient(*kubeConfPath)
	client.InitDynamicCli(*kubeConfPath)
}

func main() {
	if *nodeExample {
		fmt.Println("list nodes")
		nodes.Nodes()
	}

	if *deploy {
		if err := deployment.DeployExamples(*namespace, "test", "nginx", 1); err != nil {
			fmt.Println(err)
		}
	}

	if *ds {
		if err := daemonset.DSExample(*namespace, "test", "nginx"); err != nil {
			fmt.Println(err)
		}
	}

	if *sts {
		if err := statefulset.StsExample(*namespace, "test", "nginx", 1); err != nil {
			fmt.Println(err)
		}
	}

	if *crd {
		if err := crontab.CRDExample(*namespace, "test", "my-crontab"); err != nil {
			fmt.Println(err)
		}
	}
}