package connection

import (
	"crypto/tls"
	"encoding/binary"
	"io"
	"net"
	"testing"

	"github.com/soheilhy/cmux"
)

const (
	PacketMagic = 0x00114514
	PacketToken = "xyz_token"
)

func PacketMagicMatcher(r io.Reader) bool {
	buf := make([]byte, 4)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return false
	}
	return binary.BigEndian.Uint32(buf[:n]) == 0x00114514
}

func PacketTokenMatcher(r io.Reader) bool {
	buf := make([]byte, len(PacketToken))
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return false
	}
	return string(buf[:n]) == PacketToken
}

func wrapConnWithHeader(conn net.Conn, err error, header []byte) (net.Conn, error) {
	if err != nil {
		return nil, err
	}
	conn.Write(header)
	return conn, nil
}

func wrapTLSListener(lis net.Listener) (net.Listener, error) {
	crt, err := tls.X509KeyPair([]byte(testCrtData), []byte(testKeyData))
	if err != nil {
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{crt}}
	return tls.NewListener(lis, config), nil
}

func BenchmarkConnCmux(b *testing.B) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	tcpMux := cmux.New(lis)
	magicLis := tcpMux.Match(PacketMagicMatcher)
	tokenLis := tcpMux.Match(PacketTokenMatcher)
	tlsLis := tcpMux.Match(cmux.TLS())
	anyLis := tcpMux.Match(cmux.Any())
	go tcpMux.Serve()

	b.Run("magic", func(b *testing.B) {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, PacketMagic)

		conn0, conn1, err := getConnPair(
			func() (net.Listener, error) { return magicLis, nil },
			func(s string) (net.Conn, error) {
				conn, err := net.Dial("tcp", s)
				return wrapConnWithHeader(conn, err, buf)
			},
		)
		if err != nil {
			b.Fatal(err)
		}
		defer conn0.Close()
		defer conn1.Close()

		io.ReadAtLeast(conn0, buf, len(buf))
		bench(b, conn0, conn1)
	})

	b.Run("token", func(b *testing.B) {
		conn0, conn1, err := getConnPair(
			func() (net.Listener, error) { return tokenLis, nil },
			func(s string) (net.Conn, error) {
				conn, err := net.Dial("tcp", s)
				return wrapConnWithHeader(conn, err, []byte(PacketToken))
			},
		)
		if err != nil {
			b.Fatal(err)
		}
		defer conn0.Close()
		defer conn1.Close()

		buf := make([]byte, len(PacketToken))
		io.ReadAtLeast(conn0, buf, len(buf))
		bench(b, conn0, conn1)
	})

	b.Run("tls", func(b *testing.B) {
		b.Run("bare", func(b *testing.B) {
			conn0, conn1, err := getConnPair(
				func() (net.Listener, error) { return wrapTLSListener(tlsLis) },
				func(s string) (net.Conn, error) { return tls.Dial("tcp", s, &tls.Config{InsecureSkipVerify: true}) },
			)
			if err != nil {
				b.Fatal(err)
			}
			defer conn0.Close()
			defer conn1.Close()

			bench(b, conn0, conn1)
		})

		tlsMuxLis, err := wrapTLSListener(tlsLis)
		if err != nil {
			b.Fatal(err)
		}

		tlsMux := cmux.New(tlsMuxLis)
		tlsMagicLis := tlsMux.Match(PacketMagicMatcher)
		go tlsMux.Serve()

		b.Run("magic", func(b *testing.B) {
			buf := make([]byte, 4)
			binary.BigEndian.PutUint32(buf, PacketMagic)

			conn0, conn1, err := getConnPair(
				func() (net.Listener, error) { return tlsMagicLis, nil },
				func(s string) (net.Conn, error) {
					conn, err := tls.Dial("tcp", s, &tls.Config{InsecureSkipVerify: true})
					return wrapConnWithHeader(conn, err, buf)
				},
			)
			if err != nil {
				b.Fatal(err)
			}
			defer conn0.Close()
			defer conn1.Close()

			io.ReadAtLeast(conn0, buf, len(buf))
			bench(b, conn0, conn1)
		})

	})

	b.Run("any", func(b *testing.B) {
		buf := []byte("use_max_len_test_data")
		conn0, conn1, err := getConnPair(
			func() (net.Listener, error) { return anyLis, nil },
			func(s string) (net.Conn, error) {
				conn, err := net.Dial("tcp", s)
				return wrapConnWithHeader(conn, err, []byte(buf))
			},
		)
		if err != nil {
			b.Fatal(err)
		}
		defer conn0.Close()
		defer conn1.Close()

		io.ReadAtLeast(conn0, buf, len(buf))
		bench(b, conn0, conn1)
	})
}

func BenchmarkEchoCmux(b *testing.B) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	m := cmux.New(lis)
	l := m.Match(cmux.Any())
	go m.Serve()

	go serveEcho(b, l)
	benchEcho(b, func(s string) (net.Conn, error) { return net.Dial("tcp", s) }, lis.Addr().String())
}
