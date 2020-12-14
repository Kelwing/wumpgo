package rest

import (
	"fmt"
	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
)

type GetAuditLogParams struct {
	UserID     objects.Snowflake     `json:"user_id,omitempty"`
	ActionType objects.AuditLogEvent `json:"action_type,omitempty"`
	Before     objects.Snowflake     `json:"before,omitempty"`
	Limit      int                   `json:"limit,omitempty"`
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

	res, err := c.request(http.MethodGet, u.String(), JsonContentType, nil)
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
