package objects

import (
	"encoding/json"
	"testing"
)

func TestRoleTagsUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected RoleTags
	}{
		{
			input: `{"bot_id": "123456789012345678", "integration_id": "123456789012345678", "premium_subscriber": true}`,
			expected: RoleTags{
				BotID:             Snowflake(123456789012345678),
				IntegrationID:     Snowflake(123456789012345678),
				PremiumSubscriber: true,
			},
		},
	}

	for _, test := range tests {
		var actual RoleTags
		err := json.Unmarshal([]byte(test.input), &actual)
		if err != nil {
			t.Errorf("Failed to unmarshal: %s", err)
		}

		if actual != test.expected {
			t.Errorf("Expected %+v, got %+v", test.expected, actual)
		}
	}
}
