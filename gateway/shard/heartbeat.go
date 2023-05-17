package shard

import (
	"math"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/atomic"
	"wumpgo.dev/wumpgo/objects"
)

type Heartbeat struct {
	interval      int64
	acked         *atomic.Bool
	gw            *Shard
	stop          chan bool
	lastHeartbeat time.Time
	logger        zerolog.Logger
	wg            sync.WaitGroup
	latency       *atomic.Duration
}

func NewHeartbeat(interval int64, shard *Shard, logger zerolog.Logger) *Heartbeat {
	return &Heartbeat{
		interval: interval,
		acked:    atomic.NewBool(true),
		gw:       shard,
		logger:   logger,
    latency: atomic.NewDuration(time.Duration(math.Inf(1))),
	}
}

func (h *Heartbeat) heartbeat(interval int64) {
	h.wg.Add(1)
	defer h.wg.Done()
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()
	h.gw.logger.Debug().Msg("Heartbeat loop running")
	for {
		select {
		case <-ticker.C:
			if !h.acked.Load() {
				h.logger.Error().Msg("Discord missed last heartbeat")
				h.gw.conn.Close()
				return
			}
			h.acked.Store(false)
			err := h.SendHeartbeat()
			if err != nil {
				h.logger.Err(err).Msg("Failed to send heartbeat")
				return
			}
		case <-h.stop:
			return
		}
	}
}

func (h *Heartbeat) SendHeartbeat() error {
	h.logger.Debug().Int64("op", int64(objects.OpHeartbeat)).
		Uint64("sequence", h.gw.seq.Load()).Msg("Sending heartbeat")
	h.lastHeartbeat = time.Now()
	return h.gw.Send(objects.OpHeartbeat, h.gw.seq.Load())
}

func (h *Heartbeat) Start() {
	h.logger.Info().Msg("Starting heartbeat loop")
	h.stop = make(chan bool)
	h.acked.Store(true)
	go h.heartbeat(h.interval)
}

func (h *Heartbeat) Stop() {
	h.logger.Info().Msg("Stopping heartbeat loop")
	h.stop <- true
	h.wg.Wait()
}

func (h *Heartbeat) ACK() {
	h.acked.Store(true)
	h.latency.Store(time.Since(h.lastHeartbeat))
	h.logger.Debug().Dur("latency", h.latency.Load()).Msg("Received heartbeat ACK")
}
