// Copyright 2023 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"fmt"
	"net"
)

// GetTailscaleIP returns the Tailscale IP address of the current machine.
// Tailscale IPs are in the CGNAT range 100.64.0.0/10 (100.64.0.0 - 100.127.255.255)
// and for IPv6, they use the fd7a:115c:a1e0::/48 prefix.
// This works in Docker containers when Tailscale is configured via:
// - Sidecar container sharing the network namespace
// - Host network mode
// - Userspace networking mode
func GetTailscaleIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Define Tailscale IP ranges
	_, tailscaleIPv4Net, _ := net.ParseCIDR("100.64.0.0/10")
	_, tailscaleIPv6Net, _ := net.ParseCIDR("fd7a:115c:a1e0::/48")

	for _, iface := range interfaces {
		// Skip interfaces that are down
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			// Prefer IPv4 Tailscale addresses
			if ip.To4() != nil && tailscaleIPv4Net.Contains(ip) {
				return ip.String(), nil
			}
		}
	}

	// Fall back to IPv6 if no IPv4 Tailscale address found
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			if tailscaleIPv6Net.Contains(ip) {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no Tailscale IP address found on any network interface")
}
