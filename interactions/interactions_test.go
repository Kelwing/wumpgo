package interactions

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func generateValidKeys() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}

	return priv, pub, nil
}

func PrepareTest() (*App, ed25519.PrivateKey, ed25519.PublicKey, error) {
	// Generate a keypair for use in testing
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, nil, err
	}

	app, err := New(hex.EncodeToString(pub))
	if err != nil {
		return nil, nil, nil, err
	}

	return app, priv, pub, nil
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
	app, priv, _, err := PrepareTest()
	require.NoError(t, err)

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
	app, _, _, err := PrepareTest()
	require.NoError(t, err)

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err = enc.Encode(objects.Interaction{
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
	app, priv, _, err := PrepareTest()
	require.NoError(t, err)

	app.CommandHandler(func(context.Context, *objects.Interaction) *objects.InteractionResponse {
		return &objects.InteractionResponse{
			Type: objects.ResponseChannelMessageWithSource,
			Data: &objects.InteractionMessagesCallbackData{
				Content: "Success",
				Files: []*objects.DiscordFile{
					{
						Buffer:      bytes.NewBufferString("testing"),
						Filename:    "test.txt",
						ContentType: "text/plain",
					},
				},
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

	require.NoError(t, err)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Header().Get("Content-Type"), "multipart/form-data")

	resp := w.Result()
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	require.NoError(t, err)
	multipart.NewReader(resp.Body, params["boundary"])
	require.NoError(t, err)
}

func TestNew(t *testing.T) {
	_, pub, err := generateValidKeys()
	require.NoError(t, err)

	client := rest.New()
	tests := []struct {
		PublicKey string
		Options   []InteractionOption
		Error     bool
		Require   func(t *testing.T, a *App)
	}{
		{
			PublicKey: "invalid",
			Options:   []InteractionOption{},
			Error:     true,
		},
		{
			PublicKey: hex.EncodeToString(pub),
			Options: []InteractionOption{
				WithLogger(log.Logger),
			},
			Error: false,
			Require: func(t *testing.T, a *App) {
				require.Equal(t, log.Logger, a.logger)
			},
		},
		{
			PublicKey: hex.EncodeToString(pub),
			Options: []InteractionOption{
				WithClient(client),
			},
			Error: false,
			Require: func(t *testing.T, a *App) {
				require.Equal(t, client, a.restClient)
			},
		},
	}

	for _, tc := range tests {
		app, err := New(tc.PublicKey, tc.Options...)
		if tc.Error {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		if tc.Require != nil {
			tc.Require(t, app)
		}
	}
}

func TestComponentHandler(t *testing.T) {
	app, priv, _, err := PrepareTest()
	require.NoError(t, err)

	app.ComponentHandler(func(context.Context, *objects.Interaction) *objects.InteractionResponse {
		return &objects.InteractionResponse{
			Type: objects.ResponseModal,
			Data: &objects.InteractionModalCallbackData{
				CustomID: "mycustomID",
				Title:    "My Modal",
				Components: []*objects.Component{
					{
						Type: objects.ComponentTypeActionRow,
						Components: []*objects.Component{
							{
								Type:  objects.ComponentTypeButton,
								Style: objects.ButtonStyleLink,
								URL:   "https://wumpgo.dev",
							},
						},
					},
				},
			},
		}
	})

	data, err := json.Marshal(&objects.MessageComponentData{
		CustomID:      "test",
		ComponentType: objects.ComponentTypeButton,
	})
	require.NoError(t, err)

	req, err := generateValid(&objects.Interaction{
		ID:            objects.Snowflake(1234),
		Type:          objects.InteractionComponent,
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
	require.NoError(t, err)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
