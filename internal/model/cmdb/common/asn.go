package common

type ASN struct {
	Number       *uint32 `json:"number"            validate:"required"`
	Organization string  `json:"organization_name" validate:"required"`
}
