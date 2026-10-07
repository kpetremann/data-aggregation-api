package syslog

type Server struct {
	ServerAddress string `json:"server_address" validate:"required"`
}

type Syslog struct {
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
	ServerList []Server `json:"server_list" validate:"omitempty,dive"`
}
