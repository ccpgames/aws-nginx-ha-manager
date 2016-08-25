package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAwsNginxHaManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AwsNginxHaManager Suite")
}
