package client

import (
	"k8s.io/client-go/dynamic"
	"path/filepath"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var K8sCli func() (*kubernetes.Clientset, error)
var DynamicCli func() (dynamic.Interface, error)

func InitClient(kubeConfPath string) func() (*kubernetes.Clientset, error) {
	funcName := "InitClient"
	K8sCli = func() (*kubernetes.Clientset, error) {
		if kubeConfPath == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeConfPath = filepath.Join(home, ".kube", "config")
			} else {
				kubeConfPath = "./config"
			}
		}
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
		if err != nil {
			glog.Errorf("%s use kubeConfPath[%s] init k8s cfg error: %s", funcName, kubeConfPath, err)
			return nil, err
		}
		client , err := kubernetes.NewForConfig(cfg)
		if err != nil {
			glog.Errorf("%s init k8s client error: %s", funcName, err)
			return nil, err
		}
		return client, nil
	}
	return K8sCli
}

func InitDynamicCli(kubeConfPath string) func() (dynamic.Interface, error) {
	funcName := "InitDynamicCli"
	DynamicCli = func() (dynamic.Interface, error) {
		if kubeConfPath == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeConfPath = filepath.Join(home, ".kube", "config")
			} else {
				kubeConfPath = "./config"
			}
		}
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
		if err != nil {
			glog.Errorf("%s use kubeConfPath[%s] init k8s cfg error: %s", funcName, kubeConfPath, err)
			return nil, err
		}
		client, err := dynamic.NewForConfig(cfg)
		if err != nil {
			glog.Errorf("%s init k8s client error: %s", funcName, err)
			return nil, err
		}
		return client, nil
	}

	return DynamicCli
}