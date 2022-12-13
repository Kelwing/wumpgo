package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"wumpgo.dev/wumpgo/objects"
)

func (c *Client) GetApplicationRoleConnectionMetadataRecords(ctx context.Context, application objects.SnowflakeObject) ([]*objects.ApplicationRoleConnectionMetadata, error) {
	metadata := []*objects.ApplicationRoleConnectionMetadata{}

	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ApplicationRoleConnection, application.GetID())).
		ContentType(JsonContentType).
		Bind(&metadata).
		Send(c)

	return metadata, err
}

func (c *Client) UpdateApplicationRoleConnectionMetadataRecords(
	ctx context.Context, application objects.SnowflakeObject, params []*objects.ApplicationRoleConnectionMetadata,
) ([]*objects.ApplicationRoleConnectionMetadata, error) {
	metadata := []*objects.ApplicationRoleConnectionMetadata{}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	err = NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(ApplicationRoleConnection, application.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Bind(&metadata).
		Send(c)

	return metadata, err
}

func (c *Client) GetUserApplicationRoleConnection(ctx context.Context, application objects.SnowflakeObject) (*objects.ApplicationRoleConnection, error) {
	var connection objects.ApplicationRoleConnection

	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(UserApplicationRoleConnection, application.GetID())).
		ContentType(JsonContentType).
		Bind(&connection).
		Send(c)

	return &connection, err
}

func (c *Client) UpdateUserApplicationRoleConnection(
	ctx context.Context, application objects.SnowflakeObject, params *objects.ApplicationRoleConnection,
) (*objects.ApplicationRoleConnection, error) {
	var connection objects.ApplicationRoleConnection

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	err = NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(UserApplicationRoleConnection, application.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Bind(&connection).
		Send(c)

	return &connection, err
}
