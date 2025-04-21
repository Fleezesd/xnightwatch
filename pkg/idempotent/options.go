package idempotent

import "github.com/redis/go-redis/v9"

type Options struct {
	redis  redis.UniversalClient
	prefix string
	expire int
}

func WithRedis(rd redis.UniversalClient) func(*Options) {
	return func(o *Options) {
		if rd == nil {
			return
		}
		getOptionsOrSetDefault(o).redis = rd
	}
}

func WithPrefix(prefix string) func(*Options) {
	return func(o *Options) {
		if prefix == "" {
			return
		}
		getOptionsOrSetDefault(o).prefix = prefix
	}
}

func WithExpire(expire int) func(*Options) {
	return func(o *Options) {
		if expire <= 0 {
			return
		}
		getOptionsOrSetDefault(o).expire = expire
	}
}

// getOptionsOrSetDefault returns the provided options if they are not nil,
// otherwise it returns a default set of options.
func getOptionsOrSetDefault(options *Options) *Options {
	if options != nil {
		return options
	}

	return &Options{
		prefix: "idempotent",
		expire: 60,
	}
}
