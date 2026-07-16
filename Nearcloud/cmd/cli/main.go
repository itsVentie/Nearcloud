package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsVentie/Nearcloud/pkg/discovery"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	disc := discovery.NewDiscovery()

	localPort := 8888
	log.Printf("Registering service on port %d...", localPort)
	err := disc.Register(localPort)
	if err != nil {
		log.Fatalf("Failed to register mDNS service: %v", err)
	}
	defer disc.Stop()

	peerChan := make(chan discovery.Peer, 10)

	log.Println("Scanning for local peers...")
	err = disc.StartScanning(ctx, peerChan)
	if err != nil {
		log.Fatalf("Failed to start scanning: %v", err)
	}

	go func() {
		for peer := range peerChan {
			fmt.Printf("\n[PEER FOUND]\n")
			fmt.Printf("ID:       %s\n", peer.ID)
			fmt.Printf("IPs:      %v\n", peer.IPs)
			fmt.Printf("Port:     %d\n", peer.Port)
			fmt.Printf("Hostname: %s\n\n", peer.Hostname)
		}
	}()

	<-sigChan
	log.Println("Shutting down discovery...")
}
