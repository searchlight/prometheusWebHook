package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var config *rest.Config

func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	config = cfg

	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}
}

func GetPodLogs(namespace string) []string {
	pods := getPods(namespace)
	logList := []string{}
	for i, _ := range pods {
		_, a := getPodLog(pods[i], config)
		logList = append(logList, a)
	}

	return logList
}

func getPods(namespace string) []v1.Pod {

	clientset, err := kubernetes.NewForConfig(config)
	//for {
	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	podList := []v1.Pod{}

	for _, pod := range pods.Items {
		/*fmt.Println(getPodLogs(pod, config))
		fmt.Println("%s", pod.Name)*/
		podList = append(podList, pod)
	}

	return podList
	/*	fmt.Printf("Deployment are: \n")
		deployments, _ := clientset.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})

		for _, dep := range deployments.Items {
			fmt.Println("%s", dep.Name)
		}*/

	//time.Sleep(10 * time.Second)
	//}
}

func toPtr(xd int64) *int64 {
	return &xd
}

func getPodLog(pod v1.Pod, config *rest.Config) (string, string) {
	podLogOpts := v1.PodLogOptions{TailLines: toPtr(10)}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "error in getting access to K8S", ""
	}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening stream", ""

	}
	defer podLogs.Close()

	var message []byte

	for {
		buf := make([]byte, 2000)
		numBytes, err := podLogs.Read(buf)
		if err == io.EOF {
			break
		}
		if numBytes == 0 {
			continue
		}

		if err != nil {
			return "err", ""
		}
		message = append(message, buf[:numBytes]...)
	}

	return "", string(message)
}
