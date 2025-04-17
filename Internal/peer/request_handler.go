package peer

import (
	"P2PShare/Internal/p2ptls"
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/quic-go/quic-go"
	"net"
	"strconv"
	"strings"
)

func StartServer(port int, ready chan<- struct{}) error {
	cert, err := p2ptls.GenerateSelfCert()
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-p2p"},
		Certificates:       []tls.Certificate{cert},
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

	close(ready)
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}

		go HandleConnection(conn)
	}

}

func HandleConnection(conn quic.Connection) error {
	fmt.Println("Connection established: " + conn.RemoteAddr().String())
	fmt.Print("> ")
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			fmt.Println(err)
			return err
		}
		go func(s quic.Stream) {
			err := HandleStream(s)
			if err != nil {
				fmt.Println("streamhandler error: ", err)
				fmt.Print("> ")
			}
		}(stream)
	}
}

func HandleStream(stream quic.Stream) error {
	if stream == nil {
		return errors.New("stream is nil")
	}
	reader := bufio.NewReader(stream)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	cmd := strings.TrimSpace(line)

	HandleCommand(cmd, stream)
	return nil
}

func HandleCommand(cmd string, stream quic.Stream) {

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
