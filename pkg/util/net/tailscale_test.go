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
	"net"
	"testing"
)

func TestGetTailscaleIP(t *testing.T) {
	ip, err := GetTailscaleIP()

	if err != nil {
		t.Logf("No Tailscale IP found (expected if Tailscale is not configured): %v", err)
		return
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		t.Errorf("GetTailscaleIP returned invalid IP: %s", ip)
		return
	}

	_, tailscaleIPv4Net, _ := net.ParseCIDR("100.64.0.0/10")
	_, tailscaleIPv6Net, _ := net.ParseCIDR("fd7a:115c:a1e0::/48")

	if !tailscaleIPv4Net.Contains(parsedIP) && !tailscaleIPv6Net.Contains(parsedIP) {
		t.Errorf("GetTailscaleIP returned IP %s which is not in Tailscale range", ip)
	}

	t.Logf("Found Tailscale IP: %s", ip)
}
