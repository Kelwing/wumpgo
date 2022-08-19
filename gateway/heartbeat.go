package gateway

import (
	"time"

	"go.uber.org/atomic"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

type Heartbeat struct {
	interval      int64
	acked         *atomic.Bool
	gw            *Gateway
	stop          chan bool
	lastHeartbeat time.Time
	hub           *sentry.Hub
}

func NewHeartbeat(interval int64, gw *Gateway, hub *sentry.Hub) *Heartbeat {
	return &Heartbeat{
		interval: interval,
		acked:    atomic.NewBool(true),
		gw:       gw,
		hub:      hub,
	}
}

func (h *Heartbeat) heartbeat(interval int64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()
	h.gw.logger.Debug().Msg("Heartbeat loop running")
	for {
		select {
		case <-ticker.C:
			if !h.acked.Load() {
				evt := sentry.NewEvent()
				evt.Message = "Heartbeat ACK not received"
				evt.Tags = map[string]string{
					"last_heartbeat": h.lastHeartbeat.Format(time.RFC3339),
				}
				h.hub.CaptureEvent(evt)
				h.gw.logger.Error().Msg("Discord missed last heartbeat")
				h.gw.conn.Close()
				return
			}
			h.acked.Store(false)
			err := h.SendHeartbeat()
			if err != nil {
				h.hub.CaptureException(err)
				h.gw.logger.Err(err).Msg("Failed to send heartbeat")
			}
		case <-h.stop:
			return
		}
	}
}

func (h *Heartbeat) SendHeartbeat() error {
	h.gw.logger.Debug().Int64("op", int64(OpHeartbeat)).Int64("sequence", h.gw.sequence.Load()).Msg("Sending heartbeat")
	h.lastHeartbeat = time.Now()
	return h.gw.Send(OpHeartbeat, h.gw.sequence.Load())
}

func (h *Heartbeat) Start() {
	h.gw.logger.Info().Msg("Starting heartbeat loop")
	h.stop = make(chan bool)
	h.acked.Store(true)
	go h.heartbeat(h.interval)
}

func (h *Heartbeat) Stop() {
	h.gw.logger.Info().Msg("Stopping heartbeat loop")
	h.stop <- true
}

func (h *Heartbeat) ACK() {
	h.acked.Store(true)
	log.Debug().Dur("latency", time.Since(h.lastHeartbeat)).Msg("Received heartbeat ACK")
}
