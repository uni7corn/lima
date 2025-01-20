//go:build !darwin

package portfwd

import (
	"context"
	"net"
)

func Listen(ctx context.Context, listenConfig net.ListenConfig, hostAddress string) (net.Listener, error) {
	return listenConfig.Listen(ctx, "tcp", hostAddress)
}

func ListenPacket(ctx context.Context, listenConfig net.ListenConfig, hostAddress string) (net.PacketConn, error) {
	return listenConfig.ListenPacket(ctx, "udp", hostAddress)
}
