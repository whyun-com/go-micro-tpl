package app

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dgrr/http2"
	"github.com/valyala/fasthttp"
)
//由于要使用 ssl 证书，所以运行测试之前需要先将测试用的 CA 证书拷贝到系统目录中，参见命令 make copy-ca
//这个用例只能手动测试，启动完项目后，再手动运行这个测试用例
func TestHttp2(t *testing.T) {
	if os.Getenv("HTTP2_ENABLED") != "true" {
		fmt.Println("HTTP2未启动")
		return
	}
	addr := fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT"))
	hc := &fasthttp.HostClient{
		Addr: addr,
	}

	if err := http2.ConfigureClient(hc, http2.ClientOpts{}); err != nil {
		log.Printf("%s doesn't support http/2, %s\n", hc.Addr, err)
		t.Fatal()
		return
	}

	statusCode, body, err := hc.Get(nil, fmt.Sprintf("https://%s/healthz", addr))
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Printf("%d: %s\n", statusCode, body)
}