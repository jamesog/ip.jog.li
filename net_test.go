package ipjogli

func init() {
	// Override nonGlobalCIDRs in tests so we can use the documentation prefixes
	nonGlobalCIDRs = []string{
		"2001:db8::/32",  // RFC3849 Documentation
		"fc00::/7",       // RFC4193 Unique Local Addressing
		"10.0.0.0/8",     // RFC1918 Private
		"100.64.0.0/10",  // RFC6598 Private (CGNAT)
		"172.16.0.0/12",  // RFC1918 Private
		"192.0.0.0/24",   // RFC6890 Private
		"192.168.0.0/16", // RFC1918 Private
		"198.18.0.0/15",  // RFC6815 Private
	}
	buildNonGlobalNets()
}
