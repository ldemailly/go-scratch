// Modified from https://trstringer.com/connect-to-kubernetes-from-go/
// to get the users's default namespace from their context

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}, nil)

	rCfg, err := cfg.RawConfig()
	if err != nil {
		fmt.Printf("error getting raw config: %v\n", err)
		os.Exit(1)
	}
	namespace := rCfg.Contexts[rCfg.CurrentContext].Namespace
	fmt.Printf("Using (default) namespace: %q\n", namespace)
	kubeConfig, err := cfg.ClientConfig()
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error getting Kubernetes clientset: %v\n", err)
		os.Exit(1)
	}
	printedNS := namespace
	if namespace == "" {
		printedNS = "<all namespaces>"
	}
	fmt.Println("Get Kubernetes pods for namespace:", printedNS)
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Printf("error getting pods: %v\n", err)
		os.Exit(1)
	}
	for _, pod := range pods.Items {
		fmt.Printf("Pod name: %s %s\n", pod.Namespace, pod.Name)
	}
}
