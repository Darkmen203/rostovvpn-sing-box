//go:build go1.20 && windows

package inbound

import (
    "context"
    "net"
)

const go120Available = true

// Windows: no TFO support; fall back to normal Listen silently.
func listenTFO(listenConfig net.ListenConfig, ctx context.Context, network string, address string) (net.Listener, error) {
    return listenConfig.Listen(ctx, network, address)
}
