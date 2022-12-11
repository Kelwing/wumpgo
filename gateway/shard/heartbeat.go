package shard

import (
	"time"

	"go.uber.org/atomic"
	"wumpgo.dev/wumpgo/objects"
)

type Heartbeat struct {
	interval      int64
	acked         *atomic.Bool
	gw            *Shard
	stop          chan bool
	lastHeartbeat time.Time
}

func NewHeartbeat(interval int64, shard *Shard) *Heartbeat {
	return &Heartbeat{
		interval: interval,
		acked:    atomic.NewBool(true),
		gw:       shard,
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
				h.gw.logger.Error().Msg("Discord missed last heartbeat")
				h.gw.conn.Close()
				return
			}
			h.acked.Store(false)
			err := h.SendHeartbeat()
			if err != nil {
				h.gw.logger.Err(err).Msg("Failed to send heartbeat")
			}
		case <-h.stop:
			return
		}
	}
}

func (h *Heartbeat) SendHeartbeat() error {
	h.gw.logger.Debug().Int64("op", int64(objects.OpHeartbeat)).Uint64("sequence", h.gw.seq.Load()).Msg("Sending heartbeat")
	h.lastHeartbeat = time.Now()
	return h.gw.Send(objects.OpHeartbeat, h.gw.seq.Load())
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
	h.gw.logger.Debug().Dur("latency", time.Since(h.lastHeartbeat)).Msg("Received heartbeat ACK")
}
