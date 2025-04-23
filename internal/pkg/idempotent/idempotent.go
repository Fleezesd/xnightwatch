package idempotent

import (
	"github.com/fleezesd/xnightwatch/pkg/idempotent"
	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(NewIdempotent)

type Idempotent struct {
	idempotent *idempotent.Idempotent
}

func NewIdempotent(redis redis.UniversalClient) (idt *Idempotent, err error) {
	ins := idempotent.New(idempotent.WithRedis(redis))
	idt = &Idempotent{
		idempotent: ins,
	}

	log.Infow("Initialize idempotent success")
	return idt, nil
}
