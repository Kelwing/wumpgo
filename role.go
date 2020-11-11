package objects

type Role struct {
	ID          Snowflake     `json:"id"`
	Name        string        `json:"name"`
	Color       int           `json:"color"`
	Hoist       bool          `json:"hoist"`
	Position    int           `json:"position"`
	Permissions PermissionBit `json:"permissions"`
	Managed     bool          `json:"managed"`
	Mentionable bool          `json:"mentionable"`
}
