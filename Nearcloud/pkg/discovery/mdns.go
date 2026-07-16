package discovery

import (
	"context"
	"net"
	"os"

	"github.com/grandcat/zeroconf"
)

type Peer struct {
	ID       string
	IPs      []net.IP
	Port     int
	Hostname string
}

type Discovery struct {
	server *zeroconf.Server
}

func NewDiscovery() *Discovery {
	return &Discovery{}
}

func (d *Discovery) Register(port int) error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "nearcloud-peer"
	}

	server, err := zeroconf.Register(
		hostname,
		"_nearcloud._tcp",
		"local.",
		port,
		[]string{"txtv=1", "app=nearcloud"},
		nil,
	)
	if err != nil {
		return err
	}

	d.server = server
	return nil
}

func (d *Discovery) StartScanning(ctx context.Context, peerChan chan<- Peer) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func() {
		for entry := range entries {
			var ips []net.IP
			ips = append(ips, entry.AddrIPv4...)
			ips = append(ips, entry.AddrIPv6...)

			if len(ips) == 0 {
				continue
			}

			peerChan <- Peer{
				ID:       entry.Instance,
				IPs:      ips,
				Port:     entry.Port,
				Hostname: entry.HostName,
			}
		}
	}()

	err = resolver.Browse(ctx, "_nearcloud._tcp", "local.", entries)
	if err != nil {
		return err
	}

	return nil
}

func (d *Discovery) Stop() {
	if d.server != nil {
		d.server.Shutdown()
	}
}
