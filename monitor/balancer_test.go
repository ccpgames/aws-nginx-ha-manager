package monitor_test

import (
	. "github.com/ccpgames/aws-nginx-ha-manager/monitor"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ec2Mock struct {
}

var _ = Describe("Balancer", func() {
	var (
		balancer  *Balancer
		ipList    []string
		ipListStr string
	)

	BeforeEach(func() {
		balancer = NewBalancer("internal.balancer.test")
		ipList = []string{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"}
		ipListStr = strings.Join(ipList, ", ")
	})

	AfterEach(func() {
	})

	Describe("Running health check", func() {
		Context("With a healthy balancer", func() {
			It("should return true", func() {
				Expect(balancer.IsHealthy()).To(BeTrue())
			})
		})
	})

	Describe("Calling for a list of IPs", func() {
		Context("With ips "+ipListStr, func() {
			It("should be of length 4", func() {
				Expect(balancer.GetIPList()).To(HaveLen(len(ipList)))
			})
			It("should equal "+ipListStr, func() {
				Expect(balancer.GetIPList()).To(ConsistOf(ipList))
			})
		})
	})
})
