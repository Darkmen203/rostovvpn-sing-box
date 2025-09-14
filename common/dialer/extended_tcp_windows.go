//go:build go1.20 && windows

package dialer

import (
    "context"
    "net"

    M "github.com/sagernet/sing/common/metadata"
    N "github.com/sagernet/sing/common/network"
)

// Windows variant without TFO dependency. Supports optional TLS fragmentation only.
type ExtendedTCPDialer struct {
    net.Dialer
    DisableTFO  bool
    TLSFragment *TLSFragment
}

func (d *ExtendedTCPDialer) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
    // If not TCP or TLS fragmentation is disabled, fall back to standard dial.
    if d.TLSFragment == nil || !d.TLSFragment.Enabled || N.NetworkName(network) != N.NetworkTCP {
        switch N.NetworkName(network) {
        case N.NetworkTCP, N.NetworkUDP:
            return d.Dialer.DialContext(ctx, network, destination.String())
        default:
            return d.Dialer.DialContext(ctx, network, destination.AddrString())
        }
    }

    // TLS Fragmentation path (no TFO on Windows)
    fragmentConn := &fragmentConn{
        dialer:      d.Dialer,
        fragment:    *d.TLSFragment,
        network:     network,
        destination: destination,
    }
    conn, err := d.Dialer.DialContext(ctx, network, destination.String())
    if err != nil {
        fragmentConn.err = err
        return nil, err
    }
    fragmentConn.conn = conn
    return fragmentConn, nil
}

