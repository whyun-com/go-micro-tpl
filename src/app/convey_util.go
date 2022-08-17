package app 

import (
	"fmt"
	"time"
	"testing"
	"net/http"
	"github.com/gavv/httpexpect/v2"
	. "github.com/smartystreets/goconvey/convey"
)



type ConveyReporter struct {
	C
}

// Implements httpexpect.Reporter interface.
func (c ConveyReporter) Errorf(message string, args ...interface{}) {
	c.C.So(fmt.Sprintf(message, args...), assertionFails)
}

func assertionFails(actual interface{}, _ ...interface{}) string {
	return actual.(string)
}

type ConveyPrinter struct {
	C
}

// Request implements Printer.Request.
func (ConveyPrinter) Request(*http.Request) {
	// Does nothing.
}

// Response implements Printer.Response.
func (p ConveyPrinter) Response(*http.Response, time.Duration) {
	p.C.So(true, ShouldBeTrue)
}

func FastHTTPTester(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		// Pass requests directly to FastHTTPHandler.
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(FastHTTPHandler()),
			Jar:       httpexpect.NewJar(),
		},
		// Report errors using testify.
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

func ConveyHTTPTester(c C) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		// Pass requests directly to FastHTTPHandler.
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(FastHTTPHandler()),
			Jar:       httpexpect.NewJar(),
		},
		// Report errors using testify.
		Reporter: ConveyReporter{c},
		Printers: []httpexpect.Printer{ConveyPrinter{c}},
	})
}