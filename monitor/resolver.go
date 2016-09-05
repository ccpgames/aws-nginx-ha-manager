package monitor

import (
	"net"

	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
)

// Resolver requires Resolve
type Resolver interface {
	Resolve(string) ([]string, error)
}

// AWSResolver is a very simple resolver
type AWSResolver struct {
	elb *elb.ELB
	ec2 *ec2.EC2
}

// NewAWSResolver returns a new instance of AWSResolver
func NewAWSResolver() *AWSResolver {
	sess := session.New()
	conf := aws.Config{Region: aws.String("eu-west-1")}
	r := AWSResolver{
		elb: elb.New(sess, &conf),
		ec2: ec2.New(sess, &conf),
	}
	return &r
}

// Resolve resolves a single host
func (r *AWSResolver) Resolve(host string) ([]string, error) {
	input := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(host),
	}
	ips := []string{}
	ret, err := r.elb.DescribeInstanceHealth(input)
	if err != nil {
		log.Errorf("Error resolving %s: %s", host, err)
	} else {
		instanceIds := make([]*string, 0)
		for _, state := range ret.InstanceStates {
			if *state.State == "InService" {
				instanceIds = append(instanceIds, state.InstanceId)
			}
		}
		var instances *ec2.DescribeInstancesOutput
		instances, err = r.ec2.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: instanceIds,
		})
		ips = make([]string, len(instances.Reservations))
		for _, res := range instances.Reservations {
			for _, inst := range res.Instances {
				ips = append(ips, *inst.PrivateIpAddress)
			}
		}
	}
	// Ensure we get only ipv4
	return Filter(ips, func(s string) bool {
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
