package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
)

const V4BITS int = net.IPv4len * 8


func ip2uint(ip net.IP) uint32 {
	return binary.BigEndian.Uint32([]byte(ip.To4()))
}

func uint2ip(ip uint32) net.IP {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, ip)

	return net.IP(b)
}

func getNetworkAddress(ip net.IP, bits int) net.IP {
	return ip.Mask(getCIDRMask(bits))
}

func getBroadcastAddress(ip net.IP, bits int) net.IP {
	i := []byte(ip.To4())
	m := []byte(getCIDRMask(bits))

	r := []byte{i[0] | (m[0] ^ 255),
	            i[1] | (m[1] ^ 255),
	            i[2] | (m[2] ^ 255),
	            i[3] | (m[3] ^ 255)}

	return net.IP(r)
}

func getCIDRMask(bits int) net.IPMask {
	if bits >= 0 && bits <= 32 {
		return net.CIDRMask(bits, V4BITS)
	} else {
		return net.IPv4Mask(0, 0, 0, 0)
	}
}

func getSubnetMask(bits int) string {
	mask := getCIDRMask(bits)
	s, err := hex.DecodeString(mask.String())

	if err != nil {
		return "0.0.0.0"
	}

	return fmt.Sprintf("%d.%d.%d.%d", s[0], s[1], s[2], s[3])
}

func getCIDRbits(mask net.IPMask) int {
	s, _ := mask.Size()

	if s >= 0 && s <= 32 {
		return s
	} else {
		return 0
	}
}

