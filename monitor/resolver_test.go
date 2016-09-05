package monitor_test

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/ccpgames/aws-nginx-ha-manager/monitor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockResolver struct {
	ResolveMap map[string][]string
}

func NewMockresolver(resolveMap map[string][]string) *MockResolver {
	r := MockResolver{resolveMap}
	return &r
}

func (r MockResolver) Resolve(host string) (ipList []string, err error) {
	return r.ResolveMap[host], nil
}

var _ = Describe("Monitor/Resolver", func() {
	var (
		resolver *AWSResolver
	)
	BeforeEach(func() {
		resolver = NewAWSResolver()
	})

	It("Should fail to resolve", func() {
		log.Error("The next AWS resolve error is expected :)")
		actual, err := resolver.Resolve("dummy.host.random")
		Expect(err).NotTo(BeNil())
		Expect(actual).To(BeEmpty())
	})

	It("Should freak out", func() {
		log.Error("Don't to this at home kids")
		actual, err := resolver.Resolve("nginx-ha-testing")
		if err != nil {
			log.Errorln(err)
		}
		log.Println(actual)
	})
})
