package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const V4BITS int = net.IPv4len * 8

var bitsList = [27]int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
                       18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28 ,29, 30}


func getQuestion() (uint32, int) {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	ipint := rnd.Uint32()

	// 0.0.0.0 - 0.255.255.255 の場合は 16777217 を加算
	if ipint >= 0 && ipint <= 16777215 {
		ipint += 16777217
	}

	// 240.0.0.0 - 255.255.255.255 の場合は 268435457 を減算
	if ipint >= 4026531840 && ipint <= 4294967295 {
		ipint -= 268435457
	}

	var start int

	switch {
		case ipint >= 167772160 && ipint <= 184549375:
			// 10.0.0.0–10.255.255.255 -> /8 に制限
			start = 4

		case ipint >= 2130706432 && ipint <= 2147483647:
			// 127.0.0.0–127.255.255.255 -> /8 に制限
			start = 4

		case ipint >= 2851995648 && ipint <= 2852061183:
			// 169.254.0.0–169.254.255.255 -> /16 に制限
			start = 12

		case ipint >= 2886729728 && ipint <= 2887778303:
			// 172.16.0.0–172.31.255.255 -> /12 に制限
			start = 8

		case ipint >= 3232235520 && ipint <= 3232301055:
			// 192.168.0.0–192.168.255.255 -> /16 に制限
			start = 12

		default:
			start = rnd.Intn(21)
	}

	bits := bitsList[start:][rnd.Intn(27-start)]

	return ipint, bits
}

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

	return net.IPv4(i[0] | (m[0] ^ 255),
	                i[1] | (m[1] ^ 255),
	                i[2] | (m[2] ^ 255),
	                i[3] | (m[3] ^ 255))
}

func getSubnetMask(bits int) string {
	mask := getCIDRMask(bits)
	s, err := hex.DecodeString(mask.String())

	if err != nil {
		return "0.0.0.0"
	}

	return fmt.Sprintf("%d.%d.%d.%d", s[0], s[1], s[2], s[3])
}

func getCIDRMask(bits int) net.IPMask {
	if bits >= 0 && bits <= 32 {
		return net.CIDRMask(bits, V4BITS)
	} else {
		return net.IPv4Mask(0, 0, 0, 0)
	}
}

