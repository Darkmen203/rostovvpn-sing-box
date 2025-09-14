//go:build go1.20 && windows

package dialer

import (
    "net"
)

type tcpDialer = ExtendedTCPDialer

func newTCPDialer(dialer net.Dialer, tfoEnabled bool, tlsFragment *TLSFragment) (tcpDialer, error) {
    // Force-disable TFO on Windows; tfo-go is unsupported here.
    return tcpDialer{Dialer: dialer, DisableTFO: true, TLSFragment: tlsFragment}, nil
}

