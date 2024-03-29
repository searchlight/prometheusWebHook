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
	"log"
	"path/filepath"
	_ "path/filepath"
)

var config *rest.Config
var clientSet *kubernetes.Clientset

const fetchThisManyLinesFromLog = 10

func init() {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	//uncomment this line, if you don't use helm
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		log.Fatal("Error building config from flags", err)
	}
	//cfg, _ := clientcmd.BuildConfigFromFlags("", "")
	config = cfg

	var err2 error
	clientSet, err2 = kubernetes.NewForConfig(config)
	if err2 != nil {
		log.Fatal("Error creating clientset", err)
	}
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
func GetPodLogs(namespace string) ([]byte, error) {
	pods, err := getPods(namespace)
	if err != nil {
		return []byte{}, err
	}

	logList := []byte{}

	for i, _ := range pods {
		a, err := getPodLog(pods[i], config)
		if err != nil {
			fmt.Println("Failed to get logs from pod", err)
			continue
		}

		logList = append(logList, a...)

		if i < len(pods) {
			logList = append(logList, []byte("\n")...)
		}
	}

	if len(logList) == 0 {
		logList = append(logList, []byte("Logs empty")...)
	}

	return logList, nil
}

func getPods(namespace string) ([]v1.Pod, error) {
	pods, err := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		fmt.Println("Error getting pods ", err.Error())
		return []v1.Pod{}, err
	}

	podList := []v1.Pod{}

	for _, pod := range pods.Items {
		podList = append(podList, pod)
	}

	return podList, nil
}

func toPtr(tailLines int64) *int64 {
	return &tailLines
}

// returns ["pod logs", error]
func getPodLog(pod v1.Pod, config *rest.Config) ([]byte, error) {
	podLogOpts := v1.PodLogOptions{TailLines: toPtr(fetchThisManyLinesFromLog)}

	req := clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return []byte{}, err

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
			return []byte{}, err
		}
		message = append(message, buf[:numBytes]...)
	}

	return message, nil
}
