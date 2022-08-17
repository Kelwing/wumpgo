package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Kelwing/wumpgo/objects"
)

type BaseStickerParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"-"`
	RawTags     string   `json:"tags"`
	Reason      string   `json:"-"`
}

type CreateGuildStickerParams struct {
	BaseStickerParams
	File io.Reader
}

func (c *Client) GetSticker(ctx context.Context, id objects.SnowflakeObject) (*objects.Sticker, error) {
	sticker := &objects.Sticker{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(StickerFmt, id.GetID())).
		ContentType(JsonContentType).
		Bind(sticker).
		Send(c)

	return sticker, err
}

func (c *Client) ListNitroStickerPacks(ctx context.Context) ([]*objects.StickerPack, error) {
	var packs []*objects.StickerPack
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(NitroStickerFmt).
		ContentType(JsonContentType).
		Bind(&packs).
		Send(c)

	return packs, err
}

func (c *Client) ListGuildStickers(ctx context.Context, guildID objects.SnowflakeObject) ([]*objects.Sticker, error) {
	var stickers []*objects.Sticker
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildStickersFmt, guildID.GetID())).
		ContentType(JsonContentType).
		Bind(&stickers).
		Send(c)

	return stickers, err
}

func (c *Client) GetGuildSticker(ctx context.Context, guildID objects.SnowflakeObject, stickerID objects.SnowflakeObject) (*objects.Sticker, error) {
	sticker := &objects.Sticker{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildStickerFmt, guildID.GetID(), stickerID.GetID())).
		ContentType(JsonContentType).
		Bind(sticker).
		Send(c)

	return sticker, err
}

func (c *Client) CreateGuildSticker(ctx context.Context, guildID objects.SnowflakeObject, params *CreateGuildStickerParams) (*objects.Sticker, error) {
	buffer := new(bytes.Buffer)
	m := multipart.NewWriter(buffer)

	var err error

	err = m.WriteField("name", params.Name)
	if err != nil {
		return nil, err
	}
	err = m.WriteField("description", params.Description)
	if err != nil {
		return nil, err
	}

	if params.RawTags == "" {
		params.RawTags = strings.Join(params.Tags, ", ")
	}

	err = m.WriteField("tags", params.RawTags)
	if err != nil {
		return nil, err
	}

	f, err := m.CreateFormField("file")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, params.File)
	if err != nil {
		return nil, err
	}

	sticker := &objects.Sticker{}

	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildStickersFmt, guildID.GetID())).
		ContentType(JsonContentType).
		Body(buffer.Bytes()).
		Bind(sticker).
		Reason(params.Reason).
		Send(c)

	return sticker, err
}

func (c *Client) ModifyGuildSticker(ctx context.Context, guildID, id objects.SnowflakeObject, params *BaseStickerParams) (*objects.Sticker, error) {
	if params.RawTags == "" {
		params.RawTags = strings.Join(params.Tags, ", ")
	}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	sticker := &objects.Sticker{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildStickerFmt, guildID.GetID(), id.GetID())).
		ContentType(JsonContentType).
		Bind(sticker).
		Body(data).
		Reason(params.Reason).
		Send(c)

	return sticker, err
}

func (c *Client) DeleteGuildSticker(ctx context.Context, guildID, id objects.SnowflakeObject, reason ...string) error {
	req := NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildStickerFmt, guildID.GetID(), id.GetID())).
		ContentType(JsonContentType)

	if len(reason) > 0 {
		req.Reason(reason[0])
	}

	return req.Send(c)
}
