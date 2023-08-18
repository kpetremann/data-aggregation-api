package app

type info struct {
	Version   string
	BuildTime string
	BuildUser string
	Commit    string
}

var Info = info{}
