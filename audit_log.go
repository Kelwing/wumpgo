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
	entries := &objects.AuditLog{}
	err = NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(entries).
		Expect(http.StatusOK).
		Send(c)

	return entries, err
}
