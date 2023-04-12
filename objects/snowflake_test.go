package objects

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnowflakeMarshalJSON(t *testing.T) {
	snow := Snowflake(1234567890)
	b, err := json.Marshal(snow)
	require.NoError(t, err)
	require.Equal(t, "\"1234567890\"", string(b))
}

func TestSnowflakeUnmarshalJSON(t *testing.T) {
	var snow Snowflake
	err := json.Unmarshal([]byte("\"1234567890\""), &snow)
	require.NoError(t, err)
	require.Equal(t, Snowflake(1234567890), snow)
}

func TestSnowflakeUnmarshalJSON2(t *testing.T) {
	testObj := `{"id": "12345"}`
	obj := struct {
		ID Snowflake `json:"id"`
	}{}

	err := json.Unmarshal([]byte(testObj), &obj)
	require.NoError(t, err)
	require.Equal(t, Snowflake(12345), obj.ID)
}
