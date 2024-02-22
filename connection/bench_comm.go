package connection

import (
	"crypto/tls"
	"io"
	"net"
	"sync"
	"testing"
)

const (
	testCrtData = `-----BEGIN CERTIFICATE-----
MIICbzCCAdgCCQCSkmSsatQrejANBgkqhkiG9w0BAQsFADB8MQswCQYDVQQGEwJV
UzENMAsGA1UECAwETWFyczETMBEGA1UEBwwKaVRyYW5zd2FycDETMBEGA1UECgwK
aVRyYW5zd2FycDETMBEGA1UECwwKaVRyYW5zd2FycDEfMB0GA1UEAwwWd3d3Lnpo
YW5nbWVuZ190ZXN0LmNvbTAeFw0yMTAzMTUwOTE3MzJaFw0zMTAzMTMwOTE3MzJa
MHwxCzAJBgNVBAYTAlVTMQ0wCwYDVQQIDARNYXJzMRMwEQYDVQQHDAppVHJhbnN3
YXJwMRMwEQYDVQQKDAppVHJhbnN3YXJwMRMwEQYDVQQLDAppVHJhbnN3YXJwMR8w
HQYDVQQDDBZ3d3cuemhhbmdtZW5nX3Rlc3QuY29tMIGfMA0GCSqGSIb3DQEBAQUA
A4GNADCBiQKBgQC/XvQcJlZgE5nIeRbixgHLgmWvkQRF3ZGJ8MAi1oApQtctyaX2
KPWJizqvCHwg7TS7pVbIcXP+CFB5dGnT1B66KoOlz3y5aQK5VWGStzUN3uIEjajV
BuDNfRpI4zk6dWZW8N6lnPIQExDAsJ+kbJW2Arm+NHwIuMXG22VQZcORVwIDAQAB
MA0GCSqGSIb3DQEBCwUAA4GBAGiBn96QciYhOjAJoW8BMNJ3AzExj4rpSk54JOT7
oFkBxnFCvVvTfwo+FK+VJIIe9ALVx/exT+QBjtfdfnVvqv6+N2t6ZYWEYsmDt884
+df5mnmSQ6FOQ1Oqng2qs1fi8AZCTMfI+5ek1dNZVVYwFHkgTEOy6X5Hm4HMztvs
W1aL
-----END CERTIFICATE-----
`

	testKeyData = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQC/XvQcJlZgE5nIeRbixgHLgmWvkQRF3ZGJ8MAi1oApQtctyaX2
KPWJizqvCHwg7TS7pVbIcXP+CFB5dGnT1B66KoOlz3y5aQK5VWGStzUN3uIEjajV
BuDNfRpI4zk6dWZW8N6lnPIQExDAsJ+kbJW2Arm+NHwIuMXG22VQZcORVwIDAQAB
AoGAagtYAfFMk9jIsspG4EsQ25DagDs/vudUqrd6ANQUGMktK/Y9vPZdeWZpkmyF
PEm1mvW37ULRH8fDsEnOCs/UZieXUvq1Id3XUIcBRXcmtkJePY3VKG17K0tDxD3u
NCOaM//RgUL9CAVIlheo3nykM5zncPt71hOb4HDUL0kGooECQQDj/fhJG/0UcxD8
hlbhRPUMuf7aT/VsR8y1B8hR5gsU6uu/2stdHQbSDR1mAcXrJad4MDy8M04EWA7P
QyqTWfERAkEA1uFMvRt63UBoUt2u/S/P5WIUj1IP4vtnp648xTCoTQJ8H6KX+6n8
hVTDzQGzCyHzEc30YUUGRsj9+Ig4SEpb5wJASv8iCzqPt4haUBcIwTVjvnn4YWvn
+WRs7CfRN0+K2ailQAkC2HBR7AqwXvu6VS2ftyN29xmRUlB9HqSjfrEZYQJAczMA
UBXubbV8+IgOq4A5hbFqcle9WqQLszLPM6xdXkPpxZAGyQ4d6mFCQ6MYmOxPgwkW
bhtyPPq+ZcKp4d+zmwJAeKNksU4mDkcWox7/rLaWep9zPGRKMZcGnzzi9EqDrx9+
MwNhHleOXDhRgibypADDZwe9SYHCIGcENru1izTY4Q==
-----END RSA PRIVATE KEY-----
`
)

func getTCPConnPair() (net.Conn, net.Conn, error) {
	return getConnPair(
		func() (net.Listener, error) { return net.Listen("tcp", "localhost:0") },
		func(s string) (net.Conn, error) { return net.Dial("tcp", s) },
	)
}

func getTLSConnPair() (net.Conn, net.Conn, error) {
	crt, err := tls.X509KeyPair([]byte(testCrtData), []byte(testKeyData))
	if err != nil {
		return nil, nil, err
	}

	return getConnPair(
		func() (net.Listener, error) {
			return tls.Listen("tcp", "localhost:0", &tls.Config{Certificates: []tls.Certificate{crt}})
		},
		func(s string) (net.Conn, error) { return tls.Dial("tcp", s, &tls.Config{InsecureSkipVerify: true}) },
	)
}

func getConnPair(makeListen func() (net.Listener, error), makeConn func(string) (net.Conn, error)) (net.Conn, net.Conn, error) {
	done := make(chan struct{})
	var (
		sErr  error
		conn0 net.Conn
		conn1 net.Conn
	)

	lis, err := makeListen()
	if err != nil {
		return nil, nil, err
	}

	go func() {
		conn0, sErr = lis.Accept()
		if _, ok := conn0.(*tls.Conn); ok && sErr == nil {
			conn0.Read([]byte{}) // for TLS handeshake
		}
		close(done)
	}()

	conn1, err = makeConn(lis.Addr().String())
	if err != nil {
		return nil, nil, err
	}
	<-done
	if sErr != nil {
		return nil, nil, sErr
	}

	return conn0, conn1, nil
}

func bench(b *testing.B, rd io.Reader, wr io.Writer) {
	buf := make([]byte, 128*1024)
	buf2 := make([]byte, 128*1024)
	b.SetBytes(128 * 1024)
	b.ResetTimer()
	b.ReportAllocs()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for {
			n, _ := rd.Read(buf2)
			count += n
			if count == 128*1024*b.N {
				return
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		wr.Write(buf)
	}
	wg.Wait()
}

const BufSize = 128 * 1024

func serveEcho(b *testing.B, lis net.Listener) {
	b.SetBytes(BufSize)
	b.ResetTimer()
	b.ReportAllocs()

	copy := func(conn net.Conn) {
		defer conn.Close()
		io.Copy(conn, conn)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		go copy(conn)
	}
}

func benchEcho(b *testing.B, dialer func(string) (net.Conn, error), addr string) {
	buf := make([]byte, BufSize)
	buf2 := make([]byte, BufSize)
	b.SetBytes(BufSize)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		conn, _ := dialer(addr)
		conn.Write(buf)
		conn.Read(buf2)
		conn.Close()
	}
}
