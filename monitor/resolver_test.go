package monitor_test

import (
	"fmt"

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
		testList []string
		testHost string
	)
	BeforeEach(func() {
		testList = []string{"8.8.8.8"}
		testHost = "google-public-dns-a.google.com"
		resolver = NewAWSResolver()
	})

	It("Should resolve address", func() {
		actual, err := resolver.Resolve(testHost)
		if err != nil {
			Fail(fmt.Sprintf("Error resolving %s: %s", testHost, err))
		}
		Expect(actual).To(ConsistOf(testList))
	})

	It("Should fail to resolve", func() {
		_, err := resolver.Resolve("dummy.host.random")
		Expect(err).NotTo(BeNil())
	})
})
