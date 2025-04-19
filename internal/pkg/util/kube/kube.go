package kube

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

func SetDefaultClientOptions(config *rest.Config) *rest.Config {
	config.DisableCompression = true
	config.QPS = float32(2000)
	config.Burst = 4000
	config.ContentType = runtime.ContentTypeJSON

	return config
}
