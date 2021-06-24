package k8s

import "os"

type Info struct {
	Pod       string
	Namespace string
	Node      string
}

func GetSelfInfo() Info {
	return Info{
		Pod:       os.Getenv("K8S_POD_NAME"),
		Namespace: os.Getenv("K8S_POD_NAMESPACE"),
		Node:      os.Getenv("K8S_NODE_NAME"),
	}
}
