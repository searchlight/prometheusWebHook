package getLogs

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
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func GetLogs() {
	namespace := "default"
	kubeconfig := flag.String("kubeconfig", "/home/sohan/.kube/config", "location to your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	//for {
	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error %s, getting pods", err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range pods.Items {
		fmt.Println(getPodLogs(pod, config))
		fmt.Println("%s", pod.Name)
	}
	fmt.Printf("Deployment are: \n")
	deployments, err := clientset.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		fmt.Printf("error %s, getting pods", err.Error())
	}

	for _, dep := range deployments.Items {
		fmt.Println("%s", dep.Name)
	}

	//time.Sleep(10 * time.Second)
	//}
}

func intToPtr(x int64) *int64 {
	return &x
}

func getPodLogs(pod v1.Pod, config *rest.Config) string {
	fmt.Printf("hello world")
	podLogOpts := v1.PodLogOptions{TailLines: intToPtr(10)} // For last 10 lines
	//kubeconfig := flag.String("kubeconfig", "/home/sohan/.kube/config", "location to your kubeconfig file")
	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//if err != nil {
	//	return "error in getting config"
	//}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "error in getting access to K8S"
	}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

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
			return "err"
		}
		message := string(buf[:numBytes])
		fmt.Print(message)
	}
	return "nil"
}
