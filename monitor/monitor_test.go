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
		elbName        string
	)

	BeforeEach(func() {
		fileFH, err = ioutil.TempFile("", "config_writer_tests")
		configPath = fileFH.Name()
		elbName = "google-dns"
		interval = 500
		resolveMap := make(map[string][]string)
		resolveMap["google-dns"] = []string{"8.8.8.8", "8.8.4.4"}
		resolver := NewMockresolver(resolveMap)
		monitor = NewMonitor(configPath, dbusConnection, interval, elbName, 10080, "google-dns", resolver)
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
	}, 2)

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
		close(done)
	}, 2)
})
