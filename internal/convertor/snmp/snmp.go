package snmp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/snmp"
	"github.com/criteo/data-aggregation-api/internal/model/ietf"
)

func SNMPtoIETFfSystem(snmp *snmp.SNMP) ietf.IETFSystem_System {
	system := ietf.IETFSystem_System{}

	// Set Contact if available
	if snmp.Contact != "" {
		system.Contact = &snmp.Contact
	}

	// Set Location if available
	if snmp.Location != "" {
		system.Location = &snmp.Location
	}

	return system
}

func SNMPtoIETFsnmp(snmp *snmp.SNMP) ietf.IETFSnmp_Snmp {
	IETFsnmp := ietf.IETFSnmp_Snmp{
		Community: make(map[string]*ietf.IETFSnmp_Snmp_Community),
	}
	for _, community := range snmp.CommunityList {
		communitynew := community // create a new variable to avoid assigning the address of the range-loop variable
		snmpcommunity := ietf.IETFSnmp_Snmp_Community{
			Index:        &communitynew.Name,
			SecurityName: &communitynew.Type,
			TextName:     &communitynew.Community,
		}
		IETFsnmp.Community[community.Name] = &snmpcommunity
	}

	return IETFsnmp
}
