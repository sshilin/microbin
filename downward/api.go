package downward

import "os"

type API interface {
	PodName() string
	PodNamespace() string
	NodeName() string
}

func New() API {
	return envDownwardAPI{}
}

type envDownwardAPI struct{}

func (e envDownwardAPI) PodName() string {
	return os.Getenv("K8S_POD_NAME")
}

func (e envDownwardAPI) PodNamespace() string {
	return os.Getenv("K8S_POD_NAMESPACE")
}

func (e envDownwardAPI) NodeName() string {
	return os.Getenv("K8S_NODE_NAME")
}
