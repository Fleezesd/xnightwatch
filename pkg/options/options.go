package options

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
	netutils "k8s.io/utils/net"
)

type IOptions interface {
	Validate() []error

	AddFlags(fs *pflag.FlagSet, prefixes ...string)
}

// ValidateAddress takes an address as a string and validates it.
// If the input address is not in a valid :port or IP:port format, it returns an error.
// It also checks if the host part of the address is a valid IP address and if the port number is valid.
func ValidateAddress(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("%q is not in a valid format (:port or ip:port): %w", addr, err)
	}
	if host != "" && netutils.ParseIPSloppy(host) == nil {
		return fmt.Errorf("%q is not a valid IP address", host)
	}
	if _, err := netutils.ParsePort(port, true); err != nil {
		return fmt.Errorf("%q is not a valid number", port)
	}

	return nil
}
