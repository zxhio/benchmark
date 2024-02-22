package connection

import (
	"net"
	"net/http"
	"testing"

	"github.com/valyala/fasthttp"
)

func BenchmarkEchoFastHTTP(b *testing.B) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	buf := make([]byte, BufSize)
	go fasthttp.Serve(lis, func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBody(buf) })
	benchHTTPEcho(b, func(s string) (net.Conn, error) { return net.Dial("tcp", s) }, lis.Addr().String())
}

func BenchmarkEchoNetHTTP(b *testing.B) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	buf := make([]byte, BufSize)

	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write(buf) })
	go http.Serve(lis, h)
	benchHTTPEcho(b, func(s string) (net.Conn, error) { return net.Dial("tcp", s) }, lis.Addr().String())
}

func benchHTTPEcho(b *testing.B, dialer func(string) (net.Conn, error), addr string) {
	buf := make([]byte, BufSize)
	buf2 := make([]byte, BufSize)
	b.SetBytes(BufSize)
	b.ResetTimer()
	b.ReportAllocs()

	rawRequest := `GET / HTTP/1.1
Host: 192.168.119.200:443

hello`

	copy(buf, []byte(rawRequest))

	for i := 0; i < b.N; i++ {
		conn, _ := dialer(addr)
		conn.Write(buf)
		conn.Read(buf2)
		conn.Close()
	}
}
