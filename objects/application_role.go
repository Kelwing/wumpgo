package objects

type ApplicationRoleConnectionMetadataType uint

const (
	TypeIntegerLessThanOrEqual = iota + 1
	TypeIntegerGreaterThanOrEqual
	TypeIntegerEqual
	TypeIntegerNotEqual
	TypeDatetimeLessThanOrEqual
	TypeDatetimeGreaterThanOrEqual
	TypeBooleanEqual
	TypeBooleanNotEqual
)

type ApplicationRoleConnectionMetadata struct {
	Type                     ApplicationRoleConnectionMetadataType `json:"type"`
	Key                      string                                `json:"key"`
	Name                     string                                `json:"name"`
	NameLocalizations        map[string]string                     `json:"name_localizations"`
	Description              string                                `json:"description"`
	DescriptionLocalizations map[string]string                     `json:"description_localizations"`
}

type ApplicationRoleConnection struct {
	PlatformName     string                 `json:"platform_name"`
	PlatformUsername string                 `json:"platform_username"`
	Metadata         map[string]interface{} `json:"metadata"`
}
