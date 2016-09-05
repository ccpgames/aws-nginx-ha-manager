package monitor

import (
	"net"

	log "github.com/Sirupsen/logrus"
)

// Resolver requires Resolve
type Resolver interface {
	Resolve(string) ([]string, error)
}

// AWSResolver is a very simple resolver
type AWSResolver struct{}

// NewAWSResolver returns a new instance of AWSResolver
func NewAWSResolver() *AWSResolver {
	r := AWSResolver{}
	return &r
}

// Resolve resolves a single host
func (r *AWSResolver) Resolve(host string) (ipList []string, err error) {
	ret, err := net.LookupHost(host)
	if err != nil {
		log.Errorf("Error resolving %s: %s", host, err)
	} else {
		log.Infof("Resolved %s to %s", host, ret)
	}
	// Ensure we get only ipv4

	return Filter(ret, func(s string) bool {
		return net.ParseIP(s).To4() != nil
	}), err
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
