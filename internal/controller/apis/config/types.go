package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	componentbaseconfig "k8s.io/component-base/config"

	genericconfig "github.com/fleezesd/xnightwatch/pkg/config"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type XControllerManagerConfiguration struct {
	metav1.TypeMeta

	// Generic holds configuration for a generic controller-manager
	Generic GenericControllerManagerConfiguration

	// GarbageCollectorControllerConfiguration holds configuration for GarbageCollectorController related features.
	GarbageCollectorController GarbageCollectorControllerConfiguration

	// ChainControllerConfiguration holds configuration for ChainController related features.
	ChainController ChainControllerConfiguration
}

// GenericControllerManagerConfiguration holds configuration for a generic controller-manager.
type GenericControllerManagerConfiguration struct {
	// MySQL defines the configuration of mysql client.
	MySQL genericconfig.MySQLConfiguration

	// leaderElection defines the configuration of leader election client.
	LeaderElection componentbaseconfig.LeaderElectionConfiguration

	// Namespace that the controller watches to reconcile x-apiserver objects.
	Namespace string

	// bindAddress is the IP address for the proxy server to serve on (set to 0.0.0.0
	// for all interfaces)
	BindAddress string

	// MetricsBindAddress is the IP address and port for the metrics server to serve on,
	// defaulting to 127.0.0.1:20249 (set to 0.0.0.0 for all interfaces)
	MetricsBindAddress string

	// HealthzBindAddress is the IP address and port for the health check server to serve on,
	// defaulting to 0.0.0.0:20250
	HealthzBindAddress string

	// PprofBindAddress is the TCP address that the controller should bind to
	// for serving pprof.
	// It can be set to "" or "0" to disable the pprof serving.
	// Since pprof may contain sensitive information, make sure to protect it
	// before exposing it to public.
	PprofBindAddress string

	// Parallelism defines the amount of parallelism to process miners. Must be greater than 0. Defaults to 16.
	Parallelism int32

	// SyncPeriod determines the minimum frequency at which watched resources are
	// reconciled. A lower period will correct entropy more quickly, but reduce
	// responsiveness to change if there are many watched resources. Change this
	// value only if you know what you are doing. Defaults to 10 hours if unset.
	SyncPeriod metav1.Duration

	// Label value that the controller watches to reconcile cloud miner objects
	WatchFilterValue string

	// Controllers is the list of controllers to enable or disable
	// '*' means "all enabled by default controllers"
	// 'foo' means "enable 'foo'"
	// '-foo' means "disable 'foo'"
	// first item for a particular name wins
	Controllers []string
}

type ChainControllerConfiguration struct {
	// Image specify the blockchain node image.
	Image string
}

// GroupResource describes an group resource.
type GroupResource struct {
	// group is the group portion of the GroupResource.
	Group string
	// resource is the resource portion of the GroupResource.
	Resource string
}

// GarbageCollectorControllerConfiguration contains elements describing GarbageCollectorController.
type GarbageCollectorControllerConfiguration struct {
	// enables the generic garbage collector. MUST be synced with the
	// corresponding flag of the kube-apiserver. WARNING: the generic garbage
	// collector is an alpha feature.
	EnableGarbageCollector bool
	// concurrentGCSyncs is the number of garbage collector workers that are
	// allowed to sync concurrently.
	ConcurrentGCSyncs int32
	// gcIgnoredResources is the list of GroupResources that garbage collection should ignore.
	GCIgnoredResources []GroupResource
}
