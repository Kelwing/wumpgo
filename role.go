package objects

import (
	"encoding/json"
	"strconv"

	"github.com/Postcord/objects/permissions"
)

var _ Mentionable = (*Role)(nil)
var _ SnowflakeObject = (*Role)(nil)

type Role struct {
	DiscordBaseObject
	Name         string                    `json:"name"`
	Color        int                       `json:"color"`
	Hoist        bool                      `json:"hoist"`
	Icon         string                    `json:"icon"`
	UnicodeEmoji string                    `json:"unicode_emoji"`
	Position     int                       `json:"position"`
	Permissions  permissions.PermissionBit `json:"permissions"`
	Managed      bool                      `json:"managed"`
	Mentionable  bool                      `json:"mentionable"`
	Tags         RoleTags                  `json:"tags"`
}

type RoleTags struct {
	BotID             Snowflake `json:"bot_id"`
	IntegrationID     Snowflake `json:"integration_id"`
	PremiumSubscriber bool
}

func (r *RoleTags) UnmarshalJSON(in []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(in, &m)
	if err != nil {
		return err
	}

	if _, exists := m["premium_subscriber"]; exists {
		r.PremiumSubscriber = true
	}

	if botID, exists := m["bot_id"]; exists {
		intSnow, err := strconv.ParseUint(botID.(string), 10, 64)
		if err != nil {
			return err
		}

		r.BotID = Snowflake(intSnow)
	}

	if integrationID, exists := m["integration_id"]; exists {
		intSnow, err := strconv.ParseUint(integrationID.(string), 10, 64)
		if err != nil {
			return err
		}

		r.IntegrationID = Snowflake(intSnow)
	}

	return nil
}

func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}
