package backend

import (
	"testing"

	"net"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

//
// Component test suite entry point
//
func TestComponent(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("server_test.xml")
	RunSpecsWithDefaultAndCustomReporters(
		t,
		"coding-challenge server unit test",
		[]Reporter{junitReporter},
	)
}

func init() {
}

var (
	err      error
	listener net.Listener
	opts     []Option
)

var _ = BeforeSuite(func() {
	listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		listener, err = net.Listen("tcp", "[::1]:0")
		Expect(err).NotTo(HaveOccurred())
	}
	Expect(err).NotTo(HaveOccurred())
	Expect(listener).NotTo(BeNil())
})

var _ = AfterSuite(func() {
})
