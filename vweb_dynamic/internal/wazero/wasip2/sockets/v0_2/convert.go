package v0_2

import (
	"encoding/binary"
	"errors"
	"io/fs"
	"net"
	"syscall"
)

// fromIPAddressFamily 将 WIT 的 IPAddressFamily 转换为内部的通用类型。
func fromIPAddressFamily(family IPAddressFamily) (IPAddressFamily, error) {
	switch family {
	case IPAddressFamilyIPV4:
		return IPAddressFamilyIPV4, nil
	case IPAddressFamilyIPV6:
		return IPAddressFamilyIPV6, nil
	default:
		return 0, errors.New("invalid ip-address-family")
	}
}

// toIPAddressFamily 将内部的通用类型转换为 WIT 的 IPAddressFamily。
func toIPAddressFamily(family IPAddressFamily) (IPAddressFamily, error) {
	switch family {
	case IPAddressFamilyIPV4:
		return IPAddressFamilyIPV4, nil
	case IPAddressFamilyIPV6:
		return IPAddressFamilyIPV6, nil
	default:
		return 0, errors.New("invalid ip-address-family")
	}
}

func zoneFromScopeID(id uint32) string {
	if id == 0 {
		return ""
	}
	ifi, err := net.InterfaceByIndex(int(id))
	if err != nil {
		return ""
	}
	return ifi.Name
}

func scopeIDFromZone(zone string) uint32 {
	if zone == "" {
		return 0
	}
	ifi, err := net.InterfaceByName(zone)
	if err != nil || ifi.Index < 0 {
		return 0
	}
	return uint32(ifi.Index)
}

// fromIPSocketAddressToTCPAddr 将 WIT 的 IPSocketAddress 转换为 Go 的 *net.TCPAddr。
func fromIPSocketAddressToTCPAddr(addr IPSocketAddress) (*net.TCPAddr, error) {
	if addr.IPV4 != nil {
		ip := make(net.IP, 4)
		copy(ip, addr.IPV4.Address[:])
		return &net.TCPAddr{IP: ip, Port: int(addr.IPV4.Port)}, nil
	}
	if addr.IPV6 != nil {
		ip := make(net.IP, 16)
		for i, part := range addr.IPV6.Address {
			binary.BigEndian.PutUint16(ip[i*2:], part)
		}
		// Zone 恒空则 link-local 无法绑定/连接
		return &net.TCPAddr{IP: ip, Port: int(addr.IPV6.Port), Zone: zoneFromScopeID(addr.IPV6.ScopeID)}, nil
	}
	return nil, errors.New("invalid ip-socket-address")
}

// fromIPSocketAddressToUDPAddr 将 WIT 的 IPSocketAddress 转换为 Go 的 *net.UDPAddr。
func fromIPSocketAddressToUDPAddr(addr IPSocketAddress) (*net.UDPAddr, error) {
	if addr.IPV4 != nil {
		ip := make(net.IP, 4)
		copy(ip, addr.IPV4.Address[:])
		return &net.UDPAddr{IP: ip, Port: int(addr.IPV4.Port)}, nil
	}
	if addr.IPV6 != nil {
		ip := make(net.IP, 16)
		for i, part := range addr.IPV6.Address {
			binary.BigEndian.PutUint16(ip[i*2:], part)
		}
		return &net.UDPAddr{IP: ip, Port: int(addr.IPV6.Port), Zone: zoneFromScopeID(addr.IPV6.ScopeID)}, nil
	}
	return nil, errors.New("invalid ip-socket-address")
}

// mapDnsError 将 Go 的 net.DNSError 映射到 wasi:sockets 的 ErrorCode。
func mapDnsError(err error) ErrorCode {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		if dnsErr.IsTemporary {
			return ErrorCodeTemporaryResolverFailure
		}
		if dnsErr.IsNotFound {
			return ErrorCodeNameUnresolvable
		}
		// 如果不是临时或未找到，则认为是永久性故障
		return ErrorCodePermanentResolverFailure
	}
	// 对于其他类型的网络错误，返回一个通用的不可解析错误
	return ErrorCodeNameUnresolvable
}

// mapOsError 将 Go 的 os/syscall 网络错误映射到 wasi:sockets 的 ErrorCode。
func mapOsError(err error) ErrorCode {
	if err == nil {
		return 0 // Not an error
	}
	if errors.Is(err, fs.ErrPermission) {
		return ErrorCodeAccessDenied
	}
	if errors.Is(err, fs.ErrInvalid) {
		return ErrorCodeInvalidArgument
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		err = opErr.Err // 深入到根本的 syscall 错误
	}

	var errno syscall.Errno
	if errors.As(err, &errno) {
		if code, ok := mapNetErrno(errno); ok {
			return code
		}
	}
	return ErrorCodeUnknown
}

// toIPSocketAddress 将 Go 的 net.Addr 转换为 WIT 的 IPSocketAddress。
func toIPSocketAddress(addr net.Addr) (IPSocketAddress, error) {
	switch tcpAddr := addr.(type) {
	case *net.TCPAddr:

		if ipv4 := tcpAddr.IP.To4(); ipv4 != nil {
			var wasiAddr IPv4Address
			copy(wasiAddr[:], ipv4)
			return IPSocketAddress{
				IPV4: &IPv4SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
				},
			}, nil
		}

		if ipv6 := tcpAddr.IP.To16(); ipv6 != nil {
			var wasiAddr IPv6Address
			for i := 0; i < 8; i++ {
				wasiAddr[i] = binary.BigEndian.Uint16(ipv6[i*2:])
			}
			return IPSocketAddress{
				IPV6: &IPv6SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
					ScopeID: scopeIDFromZone(tcpAddr.Zone),
				},
			}, nil
		}
	case *net.UDPAddr:
		if ipv4 := tcpAddr.IP.To4(); ipv4 != nil {
			var wasiAddr IPv4Address
			copy(wasiAddr[:], ipv4)
			return IPSocketAddress{
				IPV4: &IPv4SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
				},
			}, nil
		}

		if ipv6 := tcpAddr.IP.To16(); ipv6 != nil {
			var wasiAddr IPv6Address
			for i := 0; i < 8; i++ {
				wasiAddr[i] = binary.BigEndian.Uint16(ipv6[i*2:])
			}
			return IPSocketAddress{
				IPV6: &IPv6SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
					ScopeID: scopeIDFromZone(tcpAddr.Zone),
				},
			}, nil
		}
	default:
	}
	return IPSocketAddress{}, errors.New("address is not TCPAddr")
}

func fromIPSocketAddressToSockaddr(addr IPSocketAddress) (syscall.Sockaddr, error) {
	if addr.IPV4 != nil {
		return &syscall.SockaddrInet4{
			Port: int(addr.IPV4.Port),
			Addr: addr.IPV4.Address,
		}, nil
	}
	if addr.IPV6 != nil {
		var ip [16]byte
		for i, part := range addr.IPV6.Address {
			binary.BigEndian.PutUint16(ip[i*2:], part)
		}
		return &syscall.SockaddrInet6{
			Port:   int(addr.IPV6.Port),
			ZoneId: addr.IPV6.ScopeID,
			Addr:   ip,
		}, nil
	}
	return nil, errors.New("invalid ip-socket-address")
}

// toIPAddress 将 Go 的 net.IP 转换为 WIT 的 IPAddress。
func toIPAddress(ip net.IP) (IPAddress, error) {
	if ipv4 := ip.To4(); ipv4 != nil {
		var wasiAddr IPv4Address
		copy(wasiAddr[:], ipv4)
		return IPAddress{IPV4: &wasiAddr}, nil
	}
	if ipv6 := ip.To16(); ipv6 != nil {
		var wasiAddr IPv6Address
		for i := 0; i < 8; i++ {
			wasiAddr[i] = binary.BigEndian.Uint16(ipv6[i*2:])
		}
		return IPAddress{IPV6: &wasiAddr}, nil
	}
	return IPAddress{}, errors.New("unsupported IP address format")
}
