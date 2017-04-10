package api

import (
	_AC "jojopoper/NBi/StressTest/appconf"
	"strings"

	"github.com/stellar/go/build"
)

var (
	NBiPublicNetwork = build.Network{"Public Global NBiLe Network 201612"}
	NBiTestNetwork   = build.Network{"Public Global NBiLe Network ; 201611"}
)

// GetNetwork 获取当前使用的network
func GetNetwork() build.Network {
	if strings.ToLower(_AC.ConfigInstance.GetString("Network", "current_network")) == "live" {
		return NBiPublicNetwork
	}
	return NBiTestNetwork
}
