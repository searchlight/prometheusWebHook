package podLogs

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
	_ "k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/client-go/util/homedir"
	"path/filepath"
	_ "path/filepath"
)

var config *rest.Config

func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	fmt.Println(*kubeconfig)
	_ = kubeconfig
	//uncomment this line, if you don't use helm
	cfg, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//cfg, _ := clientcmd.BuildConfigFromFlags("", "")
	config = cfg

	//uncomment these lines, if you don't use helm
	//if err != nil {
	//	// handle error
	//	fmt.Printf("erorr %s building config from flags\n", err.Error())
	//	config, err = rest.InClusterConfig()
	//	if err != nil {
	//		fmt.Printf("error %s, getting inclusterconfig", err.Error())
	//	}
	//}
}
func GetPodLogs(namespace string) []byte {
	pods := getPods(namespace)

	logList := []byte{}

	for i, _ := range pods {
		_, a := getPodLog(pods[i], config)
		logList = append(logList, a...)

		if i < len(pods) {
			logList = append(logList, []byte("\n")...)
		}
	}

	if len(logList) == 0 {
		logList = append(logList, []byte("Logs empty")...)
	}

	return logList
}

func getPods(namespace string) []v1.Pod {

	clientset, err := kubernetes.NewForConfig(config)
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error getting pods ", err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	podList := []v1.Pod{}

	for _, pod := range pods.Items {
		podList = append(podList, pod)
	}

	return podList
}

func toPtr(xd int64) *int64 {
	return &xd
}

// returns ["error message", "pod logs"]
func getPodLog(pod v1.Pod, config *rest.Config) (string, []byte) {
	podLogOpts := v1.PodLogOptions{TailLines: toPtr(10)}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "error in getting access to K8S", []byte{}
	}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening stream", []byte{}

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
			return "err", []byte{}
		}
		message = append(message, buf[:numBytes]...)
	}

	return "", message
}
