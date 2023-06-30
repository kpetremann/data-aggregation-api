package routingpolicy

type Decision string

const (
	Permit Decision = "permit"
	Deny   Decision = "deny"
)
