package cmdb_test

import (
	"net"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/ygot/ygot"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/interfaces"
	"github.com/criteo/data-aggregation-api/internal/types"
)

const resolveHostname = "tor01-01"

func newResolveDeviceInterface(name string) *interfaces.DeviceInterface {
	deviceInterface := interfaces.DeviceInterface{Name: name, Enabled: ygot.Bool(true)}
	deviceInterface.Device.Name = resolveHostname
	return &deviceInterface
}

func newResolveLogicalInterface(parent string, index uint32) *interfaces.LogicalInterface {
	logicalInterface := interfaces.LogicalInterface{Index: ygot.Uint32(index), Enabled: ygot.Bool(true), Type: interfaces.TypeL3}
	logicalInterface.ParentInterface.Name = parent
	logicalInterface.ParentInterface.Device.Name = resolveHostname
	return &logicalInterface
}

func newResolveManagementRoute(parent string, index uint32, kind string) *interfaces.ManagementRoute {
	route := interfaces.ManagementRoute{Kind: kind, Prefix: types.CIDR{IP: net.IPv4zero, Netmask: 0}}
	route.LogicalInterface.Index = ygot.Uint32(index)
	route.LogicalInterface.ParentInterface.Name = parent
	route.LogicalInterface.ParentInterface.Device.Name = resolveHostname
	return &route
}

func TestResolveInterfaces(t *testing.T) {
	uplink := newResolveDeviceInterface("to_spine_01")
	mgmt := newResolveDeviceInterface("eth0")
	loopback := newResolveDeviceInterface("Loopback0")

	uplinkUnit0 := newResolveLogicalInterface("to_spine_01", 0)
	uplinkUnit0.UseIPv6LinkLocalOnly = ygot.Bool(true)
	uplinkUnit100 := newResolveLogicalInterface("to_spine_01", 100)
	uplinkUnit100.UseIPv6LinkLocalOnly = ygot.Bool(false) // the default: accepted on any unit
	mgmtUnit0 := newResolveLogicalInterface("eth0", 0)
	mgmtUnit0.IPv4Address = &interfaces.IPAddress{Address: types.CIDR{IP: net.ParseIP("10.0.3.69"), Netmask: 21}, Family: 4}

	static := newResolveManagementRoute("eth0", 0, interfaces.ManagementRouteStatic)
	forced := newResolveManagementRoute("eth0", 0, interfaces.ManagementRouteForced)

	port := &interfaces.PortLayout{Name: "etp1", LogicalName: "to_spine_01", VendorName: "Ethernet0"}
	portLayout := interfaces.PortLayoutTable{"to_spine_01": port}
	neighbor := &interfaces.Neighbor{Device: "spine01-01", Port: "Ethernet16"}
	neighbors := interfaces.NeighborTable{"to_spine_01": neighbor}

	got, err := cmdb.ResolveInterfaces(resolveHostname,
		[]*interfaces.DeviceInterface{uplink, mgmt, loopback},
		[]*interfaces.LogicalInterface{uplinkUnit0, mgmtUnit0, uplinkUnit100},
		[]*interfaces.ManagementRoute{static, forced},
		portLayout, neighbors)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	want := []*interfaces.Interface{
		{
			DeviceInterface:   uplink,
			Name:              "Ethernet0",
			Port:              port,
			Neighbor:          neighbor,
			LogicalInterfaces: []*interfaces.LogicalInterface{uplinkUnit0, uplinkUnit100},
		},
		{
			DeviceInterface:   mgmt,
			Name:              "eth0",
			LogicalInterfaces: []*interfaces.LogicalInterface{mgmtUnit0},
			ManagementRoutes:  map[uint32][]*interfaces.ManagementRoute{0: {static, forced}},
		},
		{
			DeviceInterface: loopback,
			Name:            "Loopback0",
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected diff: %s", diff)
	}
}

func TestResolveInterfacesEmpty(t *testing.T) {
	got, err := cmdb.ResolveInterfaces(resolveHostname, nil, nil, nil, nil, nil)
	if err != nil || len(got) != 0 {
		t.Errorf("expected no interface, got %v, %v", got, err)
	}
}

func TestResolveInterfacesErrors(t *testing.T) {
	ipv4 := &interfaces.IPAddress{Address: types.CIDR{IP: net.ParseIP("10.0.3.69"), Netmask: 21}, Family: 4}
	mgmtUnit0 := newResolveLogicalInterface("eth0", 0)
	mgmtUnit0.IPv4Address = ipv4

	tests := []struct {
		name              string
		deviceInterfaces  []string
		logicalInterfaces func() []*interfaces.LogicalInterface
		routes            []*interfaces.ManagementRoute
		want              string
	}{
		{
			name:             "duplicate device interface",
			deviceInterfaces: []string{"Ethernet0", "Ethernet0"},
			want:             "duplicate device interface Ethernet0",
		},
		{
			name:             "two device interfaces on the same port",
			deviceInterfaces: []string{"to_spine_01", "Ethernet0"},
			want:             "device interfaces Ethernet0 and to_spine_01 of tor01-01 both map to port Ethernet0",
		},
		{
			name:             "device interface named after an ambiguous port",
			deviceInterfaces: []string{"eth0", "UNUSED"},
			want:             "device interface UNUSED of tor01-01 matches several ports of the port layout",
		},
		{
			name: "missing logical interface index",
			logicalInterfaces: func() []*interfaces.LogicalInterface {
				logicalInterface := newResolveLogicalInterface("eth0", 0)
				logicalInterface.Index = nil
				return []*interfaces.LogicalInterface{logicalInterface}
			},
			want: "missing index on logical interface of eth0",
		},
		{
			name: "logical interface of an unknown device interface",
			logicalInterfaces: func() []*interfaces.LogicalInterface {
				return []*interfaces.LogicalInterface{newResolveLogicalInterface("Ethernet4", 0)}
			},
			want: "logical interface Ethernet4.0 of tor01-01 references unknown device interface Ethernet4",
		},
		{
			name: "duplicate logical interface",
			logicalInterfaces: func() []*interfaces.LogicalInterface {
				return []*interfaces.LogicalInterface{newResolveLogicalInterface("eth0", 0), newResolveLogicalInterface("eth0", 0)}
			},
			want: "duplicate logical interface eth0.0",
		},
		{
			name: "use_ipv6_link_local_only on a non-zero unit",
			logicalInterfaces: func() []*interfaces.LogicalInterface {
				logicalInterface := newResolveLogicalInterface("eth0", 100)
				logicalInterface.UseIPv6LinkLocalOnly = ygot.Bool(true)
				return []*interfaces.LogicalInterface{logicalInterface}
			},
			want: "use_ipv6_link_local_only is set on logical interface eth0.100",
		},
		{
			name:   "management route of an unknown device interface",
			routes: []*interfaces.ManagementRoute{newResolveManagementRoute("eth1", 0, interfaces.ManagementRouteStatic)},
			want:   "references unknown device interface eth1",
		},
		{
			name:   "management route on another interface",
			routes: []*interfaces.ManagementRoute{newResolveManagementRoute("to_spine_01", 0, interfaces.ManagementRouteStatic)},
			want:   "attached to to_spine_01.0, which is not the management interface",
		},
		{
			name:   "management route of an unknown logical interface",
			routes: []*interfaces.ManagementRoute{newResolveManagementRoute("eth0", 1, interfaces.ManagementRouteStatic)},
			want:   "references unknown logical interface eth0.1",
		},
		{
			name: "management route without IPv4 address",
			logicalInterfaces: func() []*interfaces.LogicalInterface {
				return []*interfaces.LogicalInterface{newResolveLogicalInterface("eth0", 0)}
			},
			routes: []*interfaces.ManagementRoute{newResolveManagementRoute("eth0", 0, interfaces.ManagementRouteStatic)},
			want:   "no IPv4 address on eth0.0",
		},
		{
			name:   "unknown management route kind",
			routes: []*interfaces.ManagementRoute{newResolveManagementRoute("eth0", 0, "connected")},
			want:   `unsupported management route kind "connected"`,
		},
	}

	portLayout := interfaces.PortLayoutTable{
		"to_spine_01": {Name: "etp1", LogicalName: "to_spine_01", VendorName: "Ethernet0"},
		"UNUSED":      nil,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			names := test.deviceInterfaces
			if names == nil {
				names = []string{"eth0", "to_spine_01"}
			}
			deviceInterfaces := make([]*interfaces.DeviceInterface, 0, len(names))
			for _, name := range names {
				deviceInterfaces = append(deviceInterfaces, newResolveDeviceInterface(name))
			}
			logicalInterfaces := []*interfaces.LogicalInterface{mgmtUnit0}
			if test.logicalInterfaces != nil {
				logicalInterfaces = test.logicalInterfaces()
			}

			_, err := cmdb.ResolveInterfaces(resolveHostname, deviceInterfaces, logicalInterfaces, test.routes, portLayout, nil)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Errorf("got error %v, want %q", err, test.want)
			}
		})
	}
}
