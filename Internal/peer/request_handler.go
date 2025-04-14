package peer

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"net"
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

	fmt.Println("Listening on: ")
	err = PrintAddr(port)
	if err != nil {
		return err
	}

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

func PrintAddr(port int) error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && ipnet.IP.To4() != nil {
				fmt.Printf("  - %s:%d\n", ipnet.IP.String(), port)
			}
		}
	}

	return nil
}
