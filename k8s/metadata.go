package k8s

import "os"

type Pod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Node      string `json:"node"`
}

func GetPodInfo() Pod {
	return Pod{
		Name:      os.Getenv("K8S_POD_NAME"),
		Namespace: os.Getenv("K8S_POD_NAMESPACE"),
		Node:      os.Getenv("K8S_NODE_NAME"),
	}
}
