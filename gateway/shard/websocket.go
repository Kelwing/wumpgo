package shard

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/atomic"
	"nhooyr.io/websocket"
)

func NewWebsocket(logger *zerolog.Logger) *Websocket {
	if logger == nil {
		logger = &log.Logger
	}
	return &Websocket{
		logger:      *logger,
		isConnected: atomic.NewBool(false),
	}
}

type Websocket struct {
	c           *websocket.Conn
	isConnected *atomic.Bool
	logger      zerolog.Logger
}

func (w *Websocket) Open(ctx context.Context, endpoint string, requestHeader http.Header) (err error) {
	w.c, _, err = websocket.Dial(ctx, endpoint, &websocket.DialOptions{
		HTTPHeader: requestHeader,
	})
	if err != nil {
		if w.c != nil {
			_ = w.Close()
		}
		return err
	}
	w.isConnected.Store(true)

	w.c.SetReadLimit(32768 * 10000)
	return nil
}

func (w *Websocket) WriteJSON(v interface{}) error {
	wr, err := w.c.Writer(context.Background(), websocket.MessageText)
	if err != nil {
		return err
	}
	defer wr.Close()

	err = json.NewEncoder(wr).Encode(v)
	if err != nil {
		return err
	}
	return nil
}

func (w *Websocket) Close() error {
	err := w.c.Close(websocket.StatusNormalClosure, "Shutting down")
	if !w.isConnected.Load() {
		return nil
	}
	w.isConnected.Store(false)
	return err
}

func (w *Websocket) Read(ctx context.Context) (data []byte, err error) {
	var mt websocket.MessageType
	mt, data, err = w.c.Read(ctx)
	if err != nil {
		if ctx.Err() != nil && errors.Is(err, context.Canceled) {
			w.isConnected.Store(false)
			return nil, context.Canceled
		}
		var closeErr websocket.CloseError
		if errors.As(err, &closeErr) {
			w.isConnected.Store(false)
			return nil, closeErr
		}
		return nil, err
	}

	if mt == websocket.MessageBinary {
		return w.decompressPacket(data)
	}

	return data, nil
}

func (w *Websocket) decompressPacket(b []byte) ([]byte, error) {
	rdr := bytes.NewReader(b)
	z, err := zlib.NewReader(rdr)
	if err != nil {
		w.logger.Err(err).Msg("Failed to create zlib reader")
		return nil, err
	}

	defer func() {
		err := z.Close()
		if err != nil {
			w.logger.Warn().Err(err).Msg("Failed to close zlib reader")
		}
	}()

	return io.ReadAll(z)
}

func (w *Websocket) IsConnected() bool {
	return w.isConnected.Load()
}
