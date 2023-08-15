package routingpolicy

type CommunityListTerm struct {
	Community string `json:"community" validate:"required"`
}

type CommunityList struct {
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
	Name  string               `json:"name"  validate:"required"`
	Terms []*CommunityListTerm `json:"terms" validate:"required"`
}
