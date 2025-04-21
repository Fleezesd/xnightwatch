package bootstrap

import (
	"github.com/fleezesd/xnightwatch/pkg/log"
	krtlog "github.com/go-kratos/kratos/v2/log"
)

func NewLogger(info AppInfo) krtlog.Logger {
	return krtlog.With(
		log.Default(),
		"ts", krtlog.DefaultTimestamp,
		"caller", krtlog.DefaultCaller,
		"service.id", info.ID,
		"service.name", info.Name,
		"service.version", info.Version,
	)
}
