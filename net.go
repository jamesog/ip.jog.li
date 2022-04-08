package ipjogli

import (
	"net/netip"
)

// Helpers for evaluating an IP address

func init() {
	buildNonGlobalNets()
}

// The list of CIDRs net.IP.IsGlobaUnicast() incorrectly reports as global addresses
var nonGlobalCIDRs = []string{
	"2001:db8::/32",   // RFC3849 Documentation
	"fc00::/7",        // RFC4193 Unique Local Addressing
	"10.0.0.0/8",      // RFC1918 Private
	"100.64.0.0/10",   // RFC6598 Private (CGNAT)
	"172.16.0.0/12",   // RFC1918 Private
	"192.0.0.0/24",    // RFC6890 Private
	"192.0.2.0/24",    // RFC5737 Documentation
	"192.168.0.0/16",  // RFC1918 Private
	"198.18.0.0/15",   // RFC6815 Private
	"198.51.100.0/24", // RFC5737 Documentation
	"203.0.113.0/24",  // RFC5737 Documentation
}

var nonGlobalNets []netip.Prefix

func buildNonGlobalNets() {
	nonGlobalNets = make([]netip.Prefix, 0, len(nonGlobalCIDRs))
	for _, cidr := range nonGlobalCIDRs {
		net := netip.MustParsePrefix(cidr)
		nonGlobalNets = append(nonGlobalNets, net)
	}
}

func isRoutableAddr(ip netip.Addr) bool {
	for _, net := range nonGlobalNets {
		if net.Contains(ip) {
			return false
		}
	}
	return true
}
