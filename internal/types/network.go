package types

import (
	"encoding/json"
	"fmt"
	"net"
)

type CIDR struct {
	IP      net.IP
	Netmask int
}

func (c *CIDR) String() string {
	return fmt.Sprintf("%s/%d", c.IP.String(), c.Netmask)
}

func (c *CIDR) UnmarshalJSON(cidr []byte) error {
	var cidrStr string
	if err := json.Unmarshal(cidr, &cidrStr); err != nil {
		return err
	}

	if ip, network, err := net.ParseCIDR(cidrStr); err != nil {
		return err
	} else {
		c.IP = ip
		c.Netmask, _ = network.Mask.Size()
		return nil
	}
}

func (c *CIDR) MarshalJSON() ([]byte, error) {
	return []byte(c.String()), nil
}
