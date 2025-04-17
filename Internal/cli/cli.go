package cli

import (
	"P2PShare/Internal/peer"
	"bufio"
	"fmt"
	"github.com/google/shlex"
	"github.com/quic-go/quic-go"
	"log"
	"os"
	"strings"
)

func StartCli() {

	ready := make(chan struct{}, 1)
	go func() {
		err := peer.StartServer(7134, ready)
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-ready

	buf := bufio.NewReader(os.Stdin)
	var conn quic.Connection
	for {
		fmt.Print("> ")

		line, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		line = strings.TrimSpace(line)

		args, err := shlex.Split(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch args[0] {
		case "connect":
			if len(args) < 2 {
				fmt.Println("Usage: connect <address:port>")
				continue
			}
			conn, err = peer.Connect(args[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "exit":
			os.Exit(0)
		default:
			if conn != nil {
				stream, err := peer.CallCommand(args[0], conn)

				if err != nil {
					fmt.Println(err)

					continue
				}
				reader := bufio.NewReader(stream)
				line, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(line)
				err = stream.Close()
				if err != nil {
					fmt.Println(err)
					continue
				}
			} else {
				fmt.Println("Unknown command or not connected.")
			}
		}

	}
}
