package peer

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
)

func CallCommand(cmd string, stream quic.Stream) error {

	_, err := stream.Write([]byte(cmd + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func Connect(addr string) (quic.Stream, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-p2p"},
	}
	conn, err := quic.DialAddr(context.Background(), addr, tlsConfig, nil)
	if err != nil {
		return nil, err
	}

	return conn.OpenStreamSync(context.Background())
}
