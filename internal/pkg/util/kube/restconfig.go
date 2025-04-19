package kube

import (
	"fmt"
	restclient "k8s.io/client-go/rest"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fleezesd/xnightwatch/pkg/version"
)

const (
	unknowString = "unknown"
)

// buildUserAgent builds a User-Agent string from given args.
func buildUserAgent(command, version, os, arch, commit string) string {
	return fmt.Sprintf(
		"%s/%s (%s/%s) fleezesd.io/%s", command, version, os, arch, commit)
}

func DefaultXUserAgent() string {
	return buildUserAgent(
		adjustCommand(os.Args[0]),
		adjustVersion(version.Get().GitVersion),
		runtime.GOOS,
		runtime.GOARCH,
		adjustCommit(version.Get().GitCommit))
}

func SetXDefaults(config *restclient.Config) {
	if len(config.UserAgent) == 0 {
		config.UserAgent = DefaultXUserAgent()
	}
}

// adjustSourceName returns the name of the source calling the client.
func adjustSourceName(c string) string {
	if c == "" {
		return unknowString
	}
	return c
}

// adjustCommit returns sufficient significant figures of the commit's git hash.
func adjustCommit(c string) string {
	if c == "" {
		return unknowString
	}
	if len(c) > 7 {
		return c[:7]
	}
	return c
}

// adjustVersion strips "alpha", "beta", etc. from version in form
// major.minor.patch-[alpha|beta|etc].
func adjustVersion(v string) string {
	if v == "" {
		return unknowString
	}
	seg := strings.SplitN(v, "-", 2)
	return seg[0]
}

// adjustCommand returns the last component of the
// OS-specific command path for use in User-Agent.
func adjustCommand(p string) string {
	// Unlikely, but better than returning "".
	if p == "" {
		return unknowString
	}
	return filepath.Base(p)
}

func GetUserAgent(userAgent string) string {
	return DefaultXUserAgent() + "/" + adjustSourceName(userAgent)
}

func AddUserAgent(config *restclient.Config, userAgent string) *restclient.Config {
	config.UserAgent = GetUserAgent(userAgent)
	return config
}
