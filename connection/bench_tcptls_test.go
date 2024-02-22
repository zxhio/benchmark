package connection

import (
	"crypto/tls"
	"net"
	"testing"
)

func BenchmarkConnTCP(b *testing.B) {
	cs, ss, err := getTCPConnPair()
	if err != nil {
		b.Fatal(err)
	}
	defer cs.Close()
	defer ss.Close()
	bench(b, cs, ss)
}

func BenchmarkConnTLS(b *testing.B) {
	cs, ss, err := getTLSConnPair()
	if err != nil {
		b.Fatal(err)
	}
	defer cs.Close()
	defer ss.Close()
	bench(b, cs, ss)
}

func BenchmarkEchoTCP(b *testing.B) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	go serveEcho(b, lis)
	benchEcho(b, func(s string) (net.Conn, error) { return net.Dial("tcp", s) }, lis.Addr().String())
}

func BenchmarkEchoTLS(b *testing.B) {
	crt, err := tls.X509KeyPair([]byte(testCrtData), []byte(testKeyData))
	if err != nil {
		b.Fatal(err)
	}

	lis, err := tls.Listen("tcp", "localhost:0", &tls.Config{Certificates: []tls.Certificate{crt}})
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	go serveEcho(b, lis)
	benchEcho(b, func(s string) (net.Conn, error) { return tls.Dial("tcp", s, &tls.Config{InsecureSkipVerify: true}) }, lis.Addr().String())
}
