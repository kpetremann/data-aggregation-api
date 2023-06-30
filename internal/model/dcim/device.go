package dcim

type NetworkDevice struct {
	Hostname     string `json:"name" validate:"required"`
	SerialNumber string `json:"serial" validate:"omitempty"`
	Tags         []struct {
		Name string `json:"name" validate:"required"`
	} `json:"tags" validate:"omitempty"`
}
