package app

import (
	"os"
	"strconv"

	"github.com/dgrr/http2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)
const (
	httpPort = 8081
	
)

var (//以下密钥配置注意不能有空格
	certBytes = []byte(`
-----BEGIN CERTIFICATE-----
MIIDdDCCAlygAwIBAgIBATANBgkqhkiG9w0BAQsFADAxMSAwHgYDVQQDDBdUTFNH
ZW5TZWxmU2lnbmVkdFJvb3RDQTENMAsGA1UEBwwEJCQkJDAeFw0yMTA5MTYwMjQ3
MzVaFw0zMTA5MTQwMjQ3MzVaMCUxEjAQBgNVBAMMCWxvY2FsaG9zdDEPMA0GA1UE
CgwGc2VydmVyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq6fZlQeA
Wk1ylmY24+tgaXYjqQMX5S2kgUExTjr5iKP3CRsHTOzpOedOdL+nm+oGlSMpscyR
eKoflFg5CZg0aMGKVPIcYIC9Cobh9BBksV1tN9SbDTIpB6tY2IHNJBaAt50wE2h3
/Lc0w1vNRPzmEjRCkIFkJtsQSjhRLPU1R1JlaLZq8q38rYfvpUSbJx0maqYu9ZXj
jyuMVPeyaUHV7vDTx0NJ3gtD9naypboY7Ze3lByu5E21KVqd5fTfpqDNjFkFHZop
57YRbJL450U/iqtIgKeaRoGgZVlkQnWQO6y2whukAhsi2E1TIO0UNm8/GeNVHhyd
feq52ElPvuPDFwIDAQABo4GiMIGfMAkGA1UdEwQCMAAwCwYDVR0PBAQDAgWgMBMG
A1UdJQQMMAoGCCsGAQUFBwMBMDAGA1UdEQQpMCeCCWxvY2FsaG9zdIIPREVTS1RP
UC0xODhRRjIygglsb2NhbGhvc3QwHQYDVR0OBBYEFBONVAnz36PTdsPoGK2ovZwg
BSioMB8GA1UdIwQYMBaAFPw+ngF0U1wvXPXeC+1DKLSJb4qbMA0GCSqGSIb3DQEB
CwUAA4IBAQCV6ursHHHSG5t1++bySzjmoGcf2K8zKyXGTA3xKdqOlo71w0h/H1Un
z+H9b2VtDZbZ25E7JPE0/IFDDWqnPlDK0xhpycf0gKwCxHluAgL4fCxk25TQPTk2
eYGZ/n6WLz1tb0MOwFHEoxp9GRd3Db6Q6xghO8gk2V/eC4c8akprRcAUk6BL16JF
fvgIaX0lv+QRmIIaKtKihut3i0cng+bbOYtdblVRZAW1f+OJrdq0CMpqXrwU/mMq
oFT+lxkNjEwpTl8yOO5J9elfFZH3E1kgWU3wT0xCh5C8JPe/JLevzdgQZJUXEv7J
WUuMr5sM/a9Cgkn+nqEYC262IdmE7cMt
-----END CERTIFICATE-----	
	`)
	keyBytes = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCrp9mVB4BaTXKW
Zjbj62BpdiOpAxflLaSBQTFOOvmIo/cJGwdM7Ok55050v6eb6gaVIymxzJF4qh+U
WDkJmDRowYpU8hxggL0KhuH0EGSxXW031JsNMikHq1jYgc0kFoC3nTATaHf8tzTD
W81E/OYSNEKQgWQm2xBKOFEs9TVHUmVotmryrfyth++lRJsnHSZqpi71leOPK4xU
97JpQdXu8NPHQ0neC0P2drKluhjtl7eUHK7kTbUpWp3l9N+moM2MWQUdminnthFs
kvjnRT+Kq0iAp5pGgaBlWWRCdZA7rLbCG6QCGyLYTVMg7RQ2bz8Z41UeHJ196rnY
SU++48MXAgMBAAECggEAFXmXT/yzQxjjWAuLnHILCsab6X4YlDRkm0MTrwzTwRN6
SWpXPHl7KCJW/2Ymyuu5TRksEzFblbP65W0wHZIsJFrqOnxbrnsMq296fzU507Kz
gkOX7kuzNGFsaRG8H2KtUctZg2QTdstYz4QBpzrYcbiWB0wYwn+vhwmKpkw7ESP6
SEjCj0gUtJEb0XB7E/yIC+5bYy+/zCmd8kpnYkwquvR8WfnPzcU+HJlmQJuO80R0
TYgCDKGb8ZccGWPrI4d4MbhGTrszILVA1gSHK3XDPtsdadAXxQDNEmkMO6rpNTth
jXUaEJQWDcCl/LGjvtZ2vkZVslWeuczmO/HmBBrdUQKBgQDbZ3SaFx721WaBbXrZ
k7vkA6mZU2Eyjqd4DVLFjxtwIzulxObooOoYOn8gwf3tVRZuGQ969wbzQrk6jMGf
Zx38xfuy4qqLQYLcCSSDfQjJj1ZZOTFE+7Vd3PGwajwCnp75zs36dpUhCGkoHX4z
OfCu+lVDPkvWgkayC21X2fZ8/wKBgQDISYmqY6QhVO+eCPGct7Rv8VMYKU4LlWgx
gf1cHDOKsqjuPJ5Ukz32YOdmnSUgLGWy9Bg83jTJO25o+/eBFb8Vzb/0b2gVbPSu
trlTib4b+hJTshL2XstSI0zckZfZ+zRuNdHgchJeDQEuv+CkleNVbnwuBE5DDSMz
tIA8RfEB6QKBgQDOMnSLPJ+FKxmjGdkTEpzKtgZ2ar42XYtWcG8R7GTFBtfP+zVn
+5MGIjPH/Yk/u2/RGQxLbE3D4TljpVVzEd5E6Wybuhq9tVven1kJmkDf7S4hvHZp
doYFKNicC7tKWvjdnVZHxZpx6Q2q/czVJ+bjC7GF+M4dU2JNgh/JKLdW0QKBgBiH
OP7O+RjD6Bx4h+5jaQuUiFKbLF2qzHnTq42OPpmry5hxgApnhd0YfP5KHHPWPBYw
Yo+BvwEt8BWXVfZPDXnEGs/6nMqS71w+MHAUnF2cwIXTdxMJBOloPU993RTq+L7O
hIdyMOGnwg9RnFdLq+2YfEi+aj836qm4X0QCZMORAoGAYWO53VKwlYoSIlYkyOO5
XjGOvhnXkV/eprKhKQaow61s9HAwYaAWcS4DPvk5AFj9Xidt4L0X2asis69Hxa/Y
Onu7o2PQiMcIkP8V1DYzGs3iLPErYMgallsE+YRyu3JBND9vgU8hyGQe71Eybm3b
AUMHaaXtBW2GWlO1sFqqlOI=
-----END PRIVATE KEY-----	
	`)
)
type HttpApp struct {
	Port uint32
}


func (app HttpApp) Start() (error){
	var address string
	if app.Port == 0 {
		address = ":" + strconv.Itoa(httpPort)
	} else {
		address = ":" + strconv.Itoa(int(app.Port))
	}
	if os.Getenv("HTTP2_ENABLED") == "true" {
		log.Info().Msg("使用http2")
		server := &fasthttp.Server{
			Handler: FastHTTPHandler(),
		}
		http2.ConfigureServer(server)
		return server.ListenAndServeTLSEmbed(address, certBytes, keyBytes)
	}
	log.Info().Msg("使用普通http")
	return fasthttp.ListenAndServe(address, FastHTTPHandler())
}
