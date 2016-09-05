package monitor

import (
	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

// Resolver requires Resolve
type Resolver interface {
	Resolve(string) ([]string, error)
}

// AWSResolver is a very simple resolver
type AWSResolver struct {
	svc *elb.ELB
}

// NewAWSResolver returns a new instance of AWSResolver
func NewAWSResolver() *AWSResolver {
	r := AWSResolver{
		svc: elb.New(session.New(), &aws.Config{Region: aws.String("eu-west-1")}),
	}
	return &r
}

// Resolve resolves a single host
func (r *AWSResolver) Resolve(host string) (ipList []string, err error) {
	input := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(host),
	}
	ret, err := r.svc.DescribeInstanceHealth(input)
	if err != nil {
		log.Errorf("Error resolving %s: %s", host, err)
	} else {
		log.Infof("Resolved %s to %s", host, ret)
	}
	// Ensure we get only ipv4

	return []string{}, err

	// return Filter(ret, func(s string) bool {
	// 	return net.ParseIP(s).To4() != nil
	// }), err
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
