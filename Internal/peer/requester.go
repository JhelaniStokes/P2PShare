package peer

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
)

func CallCommand(cmd string, conn quic.Connection) (quic.Stream, error) {
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	_, err = stream.Write([]byte(cmd + "\n"))
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func Connect(addr string) (quic.Connection, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-p2p"},
		ServerName:         "dummy",
	}
	conn, err := quic.DialAddr(context.Background(), addr, tlsConfig, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil

}
