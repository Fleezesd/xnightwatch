package options

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/fleezesd/xnightwatch/internal/controller/apis/config"
	"github.com/fleezesd/xnightwatch/pkg/apis/apps/v1beta1"
	"github.com/fleezesd/xnightwatch/pkg/generated/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

// LogOrWriteConfig logs the completed component config and writes it into the given file name as YAML, if either is enabled.
func LogOrWriteConfig(fileName string, cfg *config.XControllerManagerConfiguration) error {
	if !(klog.V(2).Enabled() || len(fileName) > 0) {
		return nil
	}

	buf, err := encodeConfig(cfg)
	if err != nil {
		return err
	}

	if klog.V(2).Enabled() {
		klog.Info("Using component config", "config", buf.String())
	}

	if len(fileName) > 0 {
		configFile, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer configFile.Close()
		if _, err := io.Copy(configFile, buf); err != nil {
			return err
		}
		klog.InfoS("Wrote configuration", "file", fileName)
		os.Exit(0)
	}
	return nil
}

func encodeConfig(cfg *config.XControllerManagerConfiguration) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	const mediaType = runtime.ContentTypeYAML
	info, ok := runtime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return buf, fmt.Errorf("unable to locate encoder -- %q is not a supported media type", mediaType)
	}

	var encoder runtime.Encoder
	switch cfg.TypeMeta.APIVersion {
	case v1beta1.SchemeGroupVersion.String():
		encoder = scheme.Codecs.EncoderForVersion(info.Serializer, v1beta1.SchemeGroupVersion)
	default:
		encoder = scheme.Codecs.EncoderForVersion(info.Serializer, v1beta1.SchemeGroupVersion)
	}
	if err := encoder.Encode(cfg, buf); err != nil {
		return buf, err
	}
	return buf, nil
}
