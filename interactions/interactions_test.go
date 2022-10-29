package interactions

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wumpgo.dev/wumpgo/objects"
)

func PrepareTest() (*App, ed25519.PrivateKey, ed25519.PublicKey) {
	// Generate a keypair for use in testing
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	app, err := New(&Config{
		PublicKey: hex.EncodeToString(pub),
	})
	if err != nil {
		panic(err)
	}

	return app, priv, pub
}

// Generates a valid request for the given interaction
func generateValid(i *objects.Interaction, priv ed25519.PrivateKey) (*http.Request, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/", &buf)
	timestamp := time.Now().Format(time.RFC3339)
	req.Header.Set("X-Signature-Timestamp", timestamp)
	signature := ed25519.Sign(priv, append([]byte(timestamp), buf.Bytes()...))
	req.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature))

	return req, nil
}

func Test_HTTPHandler(t *testing.T) {
	// Generate a keypair for use in testing
	app, priv, _ := PrepareTest()

	req, err := generateValid(&objects.Interaction{
		ID:      objects.Snowflake(1234),
		Type:    objects.InteractionRequestPing,
		Version: 1,
	}, priv)

	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	app.HTTPHandler()(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func Test_HTTPHandler_InvalidSignature(t *testing.T) {
	// Generate a keypair for use in testing
	app, _, _ := PrepareTest()

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(objects.Interaction{
		ID:      objects.Snowflake(1234),
		Type:    objects.InteractionRequestPing,
		Version: 1,
	})
	if err != nil {
		t.Fail()
	}

	req := httptest.NewRequest("POST", "/", &buf)
	timestamp := time.Now().Format(time.RFC3339)
	req.Header.Set("X-Signature-Timestamp", timestamp)
	req.Header.Set("X-Signature-Ed25519", "invalid")

	w := httptest.NewRecorder()
	app.HTTPHandler()(w, req)

	if w.Code != 401 {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

// Test a full interaction and ensure it is properly handled, and returns a valid response
func Test_HTTPHandler_FullEvent(t *testing.T) {
	// Generate a keypair for use in testing
	app, priv, _ := PrepareTest()

	app.CommandHandler(func(context.Context, *objects.Interaction) *objects.InteractionResponse {
		return &objects.InteractionResponse{
			Type: objects.ResponseChannelMessageWithSource,
			Data: &objects.InteractionApplicationCommandCallbackData{
				Content: "Success",
			},
		}
	})

	data, err := json.Marshal(&objects.ApplicationCommandData{
		ID:   objects.Snowflake(1234),
		Name: "test",
		Type: objects.CommandTypeChatInput,
	})
	if err != nil {
		t.Errorf("Failed to marshal test interaction: %s", err)
	}

	req, err := generateValid(&objects.Interaction{
		ID:            objects.Snowflake(1234),
		Type:          objects.InteractionApplicationCommand,
		ApplicationID: objects.Snowflake(1234),
		Data:          data,
		GuildID:       objects.Snowflake(1234),
		ChannelID:     objects.Snowflake(1234),
		Member: &objects.GuildMember{
			User: &objects.User{
				ID:            objects.Snowflake(1234),
				Username:      "Test",
				Discriminator: "1234",
			},
		},
		Version: 1,
	}, priv)

	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	app.HTTPHandler()(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected application/json, got %s", w.Header().Get("Content-Type"))
	}

	var resp objects.InteractionResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Expected valid response, got %v", err)
	}

	if resp.Data.Content != "Success" {
		t.Errorf("Expected 'Success', got '%s'", resp.Data.Content)
	}
}
