package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
)

type GetAuditLogParams struct {
	UserID     objects.Snowflake     `url:"user_id,omitempty"`
	ActionType objects.AuditLogEvent `url:"action_type,omitempty"`
	Before     objects.Snowflake     `url:"before,omitempty"`
	Limit      int                   `url:"limit,omitempty"`
}

func (c *Client) GetAuditLogs(guild objects.Snowflake, params *GetAuditLogParams) (*objects.AuditLog, error) {
	u, err := url.Parse(fmt.Sprintf(GuildAuditLogsFmt, guild))
	if err != nil {
		return nil, err
	}

	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = q.Encode()

	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        u.String(),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	entries := &objects.AuditLog{}

	err = res.JSON(entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
