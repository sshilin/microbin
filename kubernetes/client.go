package kubernetes

import "os"

type Client interface {
	PodName() string
	PodNamespace() string
	NodeName() string
}

func NewClient() Client {
	return k8sClient{}
}

type k8sClient struct{}

func (e k8sClient) PodName() string {
	return os.Getenv("K8S_POD_NAME")
}

func (e k8sClient) PodNamespace() string {
	return os.Getenv("K8S_POD_NAMESPACE")
}

func (e k8sClient) NodeName() string {
	return os.Getenv("K8S_NODE_NAME")
}
