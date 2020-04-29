package main

import (
	"fmt"
	"net/http"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	http.HandleFunc("/healthz", healthHandler)
	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			panic(err.Error())
		}
	}()

	removeTaint := "NodeWithImpairedVolumes"

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

	for {
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		currentTime := time.Now().Format("01-02-2006 15:04:05")
		for _, node := range nodes.Items {
			fmt.Printf("[%s]: Found %d taints on node %s.\n", currentTime, len(node.Spec.Taints), node.Name)
			for _, taint := range node.Spec.Taints {
				if taint.Key == removeTaint {
					fmt.Printf("[%s]: Taint %s found on node %s.\n", currentTime, removeTaint, node.Name)
					node.Spec.Taints = make([]apiv1.Taint, 0)
					_, err := clientset.CoreV1().Nodes().Update(&node)
					if err != nil {
						panic(err.Error())
					}
					fmt.Printf("[%s]: Successfully removed taint on node %s.\n", currentTime, node.Name)
				} else {
					fmt.Printf("[%s]: Taint %s does not match taint %s on node %s.\n", currentTime, removeTaint, taint.Key, node.Name)
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}
