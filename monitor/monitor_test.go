package monitor_test

import (
	"io/ioutil"
	"os"
	"syscall"
	"time"

	. "github.com/ccpgames/aws-nginx-ha-manager/monitor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockDbusConnection struct{}

var reloadReceived bool

func (m MockDbusConnection) ReloadOrRestartUnit(name string, mode string, ch chan<- string) (int, error) {
	reloadReceived = true
	return 0, nil
}

var _ = Describe("Monitor/Monitor", func() {
	var (
		err            error
		fileFH         *os.File
		monitor        *Monitor
		configPath     string
		dbusConnection MockDbusConnection
		interval       int
		fqdn           string
	)

	BeforeEach(func() {
		fileFH, err = ioutil.TempFile("", "config_writer_tests")
		configPath = fileFH.Name()
		fqdn = "google-public-dns-a.google.com"
		interval = 500
		monitor = NewMonitor(configPath, dbusConnection, interval, fqdn)
		reloadReceived = false
	})

	AfterEach(func() {
		fileFH.Close()
	})

	It("Should return working monitor", func() {
		Expect(monitor).ToNot(BeNil(), "monitor instance of Monitor should not be nil")
	})

	It("Should exit gracefully", func(done Done) {
		ch := make(chan syscall.Signal)
		go monitor.Loop(ch)
		time.Sleep(time.Millisecond * 100)
		sig := <-ch
		Expect(sig).To(Equal(syscall.Signal(0)))
		ch <- syscall.SIGABRT
		time.Sleep(time.Millisecond * 100)
		sig = <-ch
		Expect(sig).To(Equal(syscall.SIGABRT))
		Expect(monitor.IsStopped()).To(BeTrue())
		close(done)
	})

	It("Should send a reload signal", func(done Done) {
		ch := make(chan syscall.Signal)
		go monitor.Loop(ch)
		time.Sleep(time.Millisecond * 100)
		Expect(<-ch).To(Equal(syscall.Signal(0)))
		ch <- syscall.SIGHUP
		time.Sleep(time.Millisecond * 1000)
		Expect(<-ch).To(Equal(syscall.SIGHUP))
		ch <- syscall.SIGABRT
		time.Sleep(time.Millisecond * 100)
		Expect(<-ch).To(Equal(syscall.SIGABRT))
		close(done)
	}, 2)
})
