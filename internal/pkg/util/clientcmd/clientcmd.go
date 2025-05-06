package clientcmd

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

const (
	RecommendedConfigPathFlag   = "kubeconfig"
	RecommendedConfigPathEnvVar = "KUBECONFIG"
	RecommendedHomeDir          = ".kube"
	RecommendedFileName         = "config"
)

var (
	RecommendedConfigDir = filepath.Join(homedir.HomeDir(), RecommendedHomeDir)
	RecommendedHomeFile  = filepath.Join(RecommendedConfigDir, RecommendedFileName)
)

func DefaultKubeconfig() string {
	defaultKubeconfig := os.Getenv(RecommendedConfigPathEnvVar)
	if defaultKubeconfig == "" {
		defaultKubeconfig = RecommendedHomeFile
	}

	return defaultKubeconfig
}
