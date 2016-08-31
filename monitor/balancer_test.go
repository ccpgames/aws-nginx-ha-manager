package monitor_test

import (
	. "github.com/ccpgames/aws-nginx-ha-manager/monitor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ec2Mock struct {
}

var _ = Describe("Balancer", func() {
	var (
		balancer     *Balancer
		deadBalancer *Balancer
		resolver     MockResolver
		liveHost     string
		deadHost     string
		ipList       []string
		resolveMap   map[string][]string
	)

	BeforeEach(func() {
		liveHost = "internal.balancer.test"
		ipList = []string{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"}
		resolveMap = make(map[string][]string)
		resolveMap[liveHost] = ipList
		resolveMap[deadHost] = make([]string, 0)
		resolver = MockResolver{resolveMap}
		balancer = NewBalancer(resolver, liveHost)
		deadBalancer = NewBalancer(resolver, deadHost)
	})

	AfterEach(func() {
	})

	Describe("Running health check", func() {
		Context("With a healthy balancer", func() {
			It("should return true", func() {
				Expect(balancer.IsHealthy()).To(BeTrue())
			})
			It("should equal expected list", func() {
				Expect(balancer.GetIPList()).To(ConsistOf(ipList))
			})
		})
		Context("With an unhealthy balancer", func() {
			It("should return unhealthy", func() {
				Expect(deadBalancer.IsHealthy()).To(BeFalse())
			})
			It("should return emtpy list", func() {
				Expect(deadBalancer.GetIPList()).To(ConsistOf([]string{}))
			})
		})
	})

})
