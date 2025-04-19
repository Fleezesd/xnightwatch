package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	MachinePool featuregate.Feature = "MachinePool"
)

func init() {
	runtime.Must(DefaultMutableFeatureGate.Add(defaultXFeatureGates))
}

var defaultXFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	MachinePool: {Default: false, PreRelease: featuregate.Alpha},
}
