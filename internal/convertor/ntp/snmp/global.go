package snmp

type Community struct {
	Name      string `json:"name"`
	Community string `json:"community"`
	Type      string `json:"type"`
}

type SNMP struct {
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
	Location      string      `json:"location"                   validate:"omitempty"`
	Contact       string      `json:"contact"                    validate:"omitempty"`
	CommunityList []Community `json:"community_list" validate:"omitempty,dive"`
}
