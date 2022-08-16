package objects

import (
	"encoding/json"
	"strings"
)

//go:generate stringer -type StickerType,StickerFormatType -output sticker_string.go

var _ SnowflakeObject = (*StickerItem)(nil)
var _ SnowflakeObject = (*Sticker)(nil)
var _ SnowflakeObject = (*StickerPack)(nil)

type StickerType int

const (
	StickerTypeStandard StickerType = iota + 1
	StickerTypeGuild
)

type StickerFormatType int

const (
	StickerFormatTypePNG StickerFormatType = iota + 1
	StickerFormatTypeAPNG
	StickerFormatTypeLOTTIE
)

type StickerTags []string

func (t StickerTags) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		strings.Join(t, ","),
	)
}

func (t StickerTags) UnmarshalJSON(bytes []byte) error {
	var tags string
	err := json.Unmarshal(bytes, &tags)
	if err != nil {
		return err
	}

	t = strings.Split(tags, ",")
	return nil
}

type StickerItem struct {
	DiscordBaseObject
	Name       string            `json:"name"`
	FormatType StickerFormatType `json:"format_type"`
}

type Sticker struct {
	StickerItem
	PackID      Snowflake   `json:"pack_id"`
	Description string      `json:"description"`
	Tags        string      `json:"tags"`
	Type        StickerType `json:"type"`
	Available   bool        `json:"available,omitempty"`
	GuildID     Snowflake   `json:"guild_id,omitempty"`
	User        *User       `json:"user,omitempty"`
	SortValue   *int        `json:"sort_value,omitempty"`
}

type StickerPack struct {
	DiscordBaseObject
	Stickers       []*Sticker `json:"stickers"`
	Name           string     `json:"name"`
	SKU            Snowflake  `json:"sku_id"`
	CoverStickerID Snowflake  `json:"cover_sticker_id"`
	Description    string     `json:"description"`
	BannerAssetID  Snowflake  `json:"banner_asset_id"`
}
