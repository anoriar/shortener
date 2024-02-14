package utilip

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// ErrIPNotFound missing godoc.
var ErrIPNotFound = errors.New("IP not found")

// IPService missing godoc.
type IPService struct {
}

// IsIPBelongToSubnet missing godoc.
func (s *IPService) IsIPBelongToSubnet(ip net.IP, cidr string) (bool, error) {
	_, subnet, err := net.ParseCIDR(cidr)

	if err != nil {
		return false, fmt.Errorf("parse CIDR error %v", err)
	}

	return subnet.Contains(ip), nil
}

// GetIPFromRequest missing godoc.
func (s *IPService) GetIPFromRequest(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get("X-Real-IP")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		ips := r.Header.Get("X-Forwarded-For")
		ipStrs := strings.Split(ips, ",")
		ipStr = ipStrs[0]
		ip = net.ParseIP(ipStr)
	}
	if ip == nil {
		return nil, ErrIPNotFound
	}
	return ip, nil
}
