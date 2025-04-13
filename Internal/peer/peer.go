package peer

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"strconv"
	"strings"
)

func StartServer(port int) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-p2p"},
	}
	listener, err := quic.ListenAddr("0.0.0.0:"+strconv.Itoa(port), tlsConfig, nil)
	if err != nil {
		return err
	}
	s := fmt.Sprint("running on: ", listener.Addr())
	fmt.Println(s)
	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	err = HandleStream(stream)
	if err != nil {
		return err
	}

	return nil
}

func Connect(addr string, command string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-p2p"},
	}
	conn, err := quic.DialAddr(context.Background(), addr, tlsConfig, nil)
	if err != nil {
		return err
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	_, err = stream.Write([]byte(command + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func HandleStream(stream quic.Stream) error {
	reader := bufio.NewReader(stream)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	cmd := strings.TrimSpace(line)

	fmt.Println(cmd)
	return nil
}
