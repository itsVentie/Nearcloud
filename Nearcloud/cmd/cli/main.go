package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/itsVentie/Nearcloud/pkg/discovery"
	"github.com/itsVentie/Nearcloud/pkg/network"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	localTCPPort := 9999
	server := network.NewTransferServer()
	go func() {
		err := server.Start(fmt.Sprintf(":%d", localTCPPort))
		if err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	defer server.Close()

	disc := discovery.NewDiscovery()
	err := disc.Register(localTCPPort)
	if err != nil {
		log.Fatalf("mDNS registration failed: %v", err)
	}
	defer disc.Stop()

	peerChan := make(chan discovery.Peer, 10)
	err = disc.StartScanning(ctx, peerChan)
	if err != nil {
		log.Fatalf("mDNS scanning failed: %v", err)
	}

	var peers []discovery.Peer
	go func() {
		for peer := range peerChan {
			alreadyExists := false
			for _, p := range peers {
				if p.ID == peer.ID {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				peers = append(peers, peer)
				fmt.Printf("\n[NEW PEER DETECTED] #%d - %s (%v:%d)\n", len(peers), peer.Hostname, peer.IPs[0], peer.Port)
			}
		}
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Println("\nMenu: [1] List peers & Send file | [2] Exit")
			fmt.Print("> ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			if choice == "2" {
				cancel()
				return
			}

			if choice == "1" {
				if len(peers) == 0 {
					fmt.Println("No peers discovered yet.")
					continue
				}

				fmt.Println("\nDiscovered Peers:")
				for i, p := range peers {
					fmt.Printf("[%d] %s (%s:%d)\n", i+1, p.Hostname, p.IPs[0], p.Port)
				}

				fmt.Print("Select peer index: ")
				indexStr, _ := reader.ReadString('\n')
				idx, err := strconv.Atoi(strings.TrimSpace(indexStr))
				if err != nil || idx < 1 || idx > len(peers) {
					fmt.Println("Invalid selection.")
					continue
				}

				selectedPeer := peers[idx-1]
				targetAddr := fmt.Sprintf("%s:%d", selectedPeer.IPs[0], selectedPeer.Port)

				fmt.Println("Sending mock transfer request (test.txt, 15.5 MB)...")
				accepted, err := network.SendFileRequest(targetAddr, "test.txt", 16252928)
				if err != nil {
					fmt.Printf("Request failed: %v\n", err)
				} else if accepted {
					fmt.Println("Result: Transfer APPROVED by peer!")
				} else {
					fmt.Println("Result: Transfer DENIED by peer.")
				}
			}
		}
	}()

	select {
	case <-sigChan:
		log.Println("Exiting...")
	case <-ctx.Done():
		log.Println("Context cancelled, exiting...")
	}
}
