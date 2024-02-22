package connection

import (
	"net"
	"testing"

	"github.com/xtaci/smux"
)

type SmuxSession struct{ *smux.Session }

func (s *SmuxSession) Addr() net.Addr            { return s.Session.LocalAddr() }
func (s *SmuxSession) Accept() (net.Conn, error) { return s.Session.AcceptStream() }
func (s *SmuxSession) Close() error              { return s.Session.Close() }

func BenchmarkConnSmux(b *testing.B) {
	b.Run("OverTCP", func(b *testing.B) {
		cs, ss, err := getTCPSmuxStreamPair()
		if err != nil {
			b.Fatal(err)
		}
		defer cs.Close()
		defer ss.Close()
		bench(b, cs, ss)
	})

	b.Run("OverTLS", func(b *testing.B) {
		cs, ss, err := getTLSSmuxStreamPair()
		if err != nil {
			b.Fatal(err)
		}
		defer cs.Close()
		defer ss.Close()
		bench(b, cs, ss)
	})
}

func BenchmarkEchoSmux(b *testing.B) {
	b.Run("OverTCP", func(b *testing.B) {
		cs, ss, err := getTCPConnPair()
		if err != nil {
			b.Fatal(err)
		}
		benchEchoSmux(b, cs, ss)
	})

	b.Run("OverTLS", func(b *testing.B) {
		conn0, conn1, err := getTLSConnPair()
		if err != nil {
			b.Fatal(err)
		}
		benchEchoSmux(b, conn0, conn1)
	})
}

func benchEchoSmux(b *testing.B, conn0, conn1 net.Conn) {
	defer conn1.Close()
	defer conn0.Close()

	cs, _ := smux.Client(conn0, nil)
	ss, _ := smux.Server(conn1, nil)

	go serveEcho(b, &SmuxSession{Session: ss})
	benchEcho(b, func(s string) (net.Conn, error) { return cs.OpenStream() }, "")
}

// Code from https://github.com/xtaci/smux/blob/master/session_test.go#L1005
func getTCPSmuxStreamPair() (*smux.Stream, *smux.Stream, error) {
	c1, c2, err := getTCPConnPair()
	if err != nil {
		return nil, nil, err
	}
	return getSmuxStreamPair(c1, c2)
}

func getTLSSmuxStreamPair() (*smux.Stream, *smux.Stream, error) {
	c1, c2, err := getTLSConnPair()
	if err != nil {
		return nil, nil, err
	}
	return getSmuxStreamPair(c1, c2)
}

func getSmuxStreamPair(c1, c2 net.Conn) (*smux.Stream, *smux.Stream, error) {
	s, err := smux.Server(c2, nil)
	if err != nil {
		return nil, nil, err
	}
	c, err := smux.Client(c1, nil)
	if err != nil {
		return nil, nil, err
	}
	var ss *smux.Stream
	done := make(chan error)
	go func() {
		var rerr error
		ss, rerr = s.AcceptStream()
		done <- rerr
		close(done)
	}()
	cs, err := c.OpenStream()
	if err != nil {
		return nil, nil, err
	}
	err = <-done
	if err != nil {
		return nil, nil, err
	}

	return cs, ss, nil
}
