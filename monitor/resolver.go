package monitor

import "net"

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
	return net.LookupHost(host)
}
