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
		interval = 1
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
		sig := make(chan os.Signal, 1)
		msgOut := make(chan string, 1)
		go monitor.Loop(sig, msgOut)
		time.Sleep(time.Millisecond * 100)
		Expect(<-msgOut).To(Equal("Updated and reloaded configuration"))
		sig <- os.Interrupt
		time.Sleep(time.Millisecond * 100)
		Expect(<-msgOut).To(Equal("Exit"))
		Expect(monitor.IsStopped()).To(BeTrue())
		close(done)
	}, 5)

	It("Should send a reload signal", func(done Done) {
		sig := make(chan os.Signal, 1)
		msgOut := make(chan string, 1)
		go monitor.Loop(sig, msgOut)
		time.Sleep(time.Millisecond * 100)
		Expect(<-msgOut).To(Equal("Updated and reloaded configuration"))
		sig <- syscall.SIGHUP
		time.Sleep(time.Millisecond * 1000)
		Expect(<-msgOut).To(Equal("Reloaded configuration"))
		sig <- os.Interrupt
		time.Sleep(time.Millisecond * 100)
		close(done)
	}, 5)
})
