package client

import (
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/pflag"
)

var (
	UserAgent  = "xnightwatch"
	Debug      = false
	RetryCount = 3
	Timeout    = 30 * time.Second
)

func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&UserAgent, "client.user-agent", UserAgent, ""+
		"Used to specify the Resty client User-Agent.")

	fs.BoolVar(&Debug, "client.debug", Debug, ""+
		"Enables the debug mode on Resty client.")

	fs.IntVar(&RetryCount, "client.retry-count", RetryCount, ""+
		"Enables retry on Resty client and allows you to set no. of retry count. Resty uses a Backoff mechanism.")

	fs.DurationVar(&Timeout, "client.timeout", Timeout, ""+
		"Request timeout for client.")
}

func NewRequest(url string) *resty.Request {
	return resty.New().
		SetBaseURL(url).
		SetRetryCount(RetryCount).
		SetDebug(Debug).
		R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   UserAgent,
		})
}

// IsDiscoveryEndpoint used to determine if the given endpoint is a service discovery endpoint.
func IsDiscoveryEndpoint(server string) bool {
	return strings.HasPrefix(server, "discovery:///")
}
