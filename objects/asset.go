package objects

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const baseURL = "https://cdn.discordapp.com"

type assetOptions struct {
	Size   int
	Format string
}

func newOptions(opts ...AssetOption) *assetOptions {
	opt := &assetOptions{
		Size:   0,
		Format: "png",
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

type AssetOption func(o *assetOptions)

func WithSize(size int) AssetOption {
	return func(o *assetOptions) {
		o.Size = size
	}
}

func WithExtension(ext string) AssetOption {
	return func(o *assetOptions) {
		o.Format = ext
	}
}

type CDNObject interface {
	Avatar | CustomEmoji | GuildIcon
}

type Asset[T CDNObject] struct {
	asset T
}

func (a *Asset[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.asset)
}

func (a *Asset[T]) UnmarshalJSON(data []byte) error {
	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	a.asset = v
	return nil
}

func (a *Asset[T]) Asset() T {
	return a.asset
}

func buildURL(in string, opt *assetOptions) string {
	u, _ := url.Parse(in)
	if opt.Size > 0 {
		q := u.Query()
		q.Set("size", strconv.FormatInt(int64(opt.Size), 10))
		u.RawQuery = q.Encode()
	}

	return u.String()
}

type CustomEmoji Snowflake

func (c CustomEmoji) URL(opts ...AssetOption) string {
	opt := newOptions(opts...)

	return buildURL(fmt.Sprintf("%s/emojis/%d.%s", baseURL, c, opt.Format), opt)
}

type GuildIcon string

func (g GuildIcon) URL(guildID Snowflake, opts ...AssetOption) string {
	opt := newOptions(opts...)

	return buildURL(fmt.Sprintf("%s/guilds/%d/%s.%s", baseURL, guildID, g, opt.Format), opt)
}

type Avatar string

func (a Avatar) URL(userID Snowflake, opts ...AssetOption) string {
	opt := newOptions(opts...)
	return buildURL(fmt.Sprintf("%s/avatars/%d/%s.%s", baseURL, userID, a, opt.Format), opt)
}
