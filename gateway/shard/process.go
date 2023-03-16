package shard

import (
	"encoding/json"
	"time"

	"wumpgo.dev/wumpgo/objects"
)

var (
	_ packetProcessor = (*dispatchProcessor)(nil)
	_ packetProcessor = (*helloProcessor)(nil)
	_ packetProcessor = (*heartbeatProcessor)(nil)
	_ packetProcessor = (*reconnectProcessor)(nil)
	_ packetProcessor = (*invalidSessionProcessor)(nil)
	_ packetProcessor = (*heartbeatAckProcessor)(nil)
)

type packetProcessor interface {
	op() objects.OpCode
	process(*Shard, objects.Payload) error
}

type dispatchProcessor struct{}

func (d *dispatchProcessor) op() objects.OpCode {
	return objects.OpDispatch
}

func (d *dispatchProcessor) process(s *Shard, p objects.Payload) error {
	s.logger.Debug().Str("event", p.EventName).Msg("Received dispatch")
	if p.EventName == "READY" {
		ready := &objects.Ready{}
		err := json.Unmarshal(p.Data, ready)
		if err != nil {
			s.logger.Err(err).Msg("Failed to unmarshal ready")
			return err
		}
		s.session_id = ready.SessionID
		s.resume_url = ready.ResumeGatewayURL
		s.logger.Info().Str("session_id", s.session_id).Str("user", ready.User.Username).Msg("We are ready!")
	}
	go func(event string, data json.RawMessage) {
		start := time.Now()
		err := s.dispatcher.Dispatch(p.EventName, p.Data)
		s.logger.Debug().Dur("duration", time.Since(start)).Str("event", p.EventName).Msg("Dispatch finished")
		if err != nil {
			s.logger.Err(err).Msg("Failed to dispatch")
		}
	}(p.EventName, p.Data)
	return nil
}

type helloProcessor struct{}

func (h *helloProcessor) op() objects.OpCode {
	return objects.OpHello
}

func (h *helloProcessor) process(s *Shard, p objects.Payload) error {
	s.logger.Info().Msgf("Received Hello Op")
	s.hello.Store(true)
	hello := objects.Hello{}
	err := json.Unmarshal(p.Data, &hello)
	if err != nil {
		s.logger.Err(err).Msg("Failed to unmarshal hello")
		return err
	}
	s.heartbeat = NewHeartbeat(hello.HeartbeatInterval, s, s.logger)
	s.heartbeat.Start()

	if s.resume.Load() && s.session_id != "" {
		s.logger.Info().Msg("Resuming session")
		s.resume.Store(true)
		err = s.sendResume()
		if err != nil {
			s.logger.Err(err).Msg("Failed to resume session")
			return err
		}
	} else {
		s.logger.Info().Msg("Identifying")
		err = s.sendIdentify()
		if err != nil {
			s.logger.Err(err).Msg("Failed to identify")
			return err
		}
	}
	return nil
}

type heartbeatProcessor struct{}

func (h *heartbeatProcessor) op() objects.OpCode {
	return objects.OpHeartbeat
}

func (h *heartbeatProcessor) process(s *Shard, p objects.Payload) error {
	if err := s.heartbeat.SendHeartbeat(); err != nil {
		s.logger.Err(err).Msg("Failed to send heartbeat")
		return err
	}
	return nil
}

type reconnectProcessor struct{}

func (r *reconnectProcessor) op() objects.OpCode {
	return objects.OpReconnect
}

func (r *reconnectProcessor) process(s *Shard, p objects.Payload) error {
	s.logger.Info().Msg("Received reconnect")
	s.resume.Store(true)
	s.Close()
	return shardError("reconnect")
}

type invalidSessionProcessor struct{}

func (i *invalidSessionProcessor) op() objects.OpCode {
	return objects.OpInvalidSession
}

func (i *invalidSessionProcessor) process(s *Shard, p objects.Payload) error {
	s.logger.Info().Msg("Invalid session")
	// Ensure we're not resuming
	s.resume.Store(false)
	s.Close()
	return shardError("invalid session")
}

type heartbeatAckProcessor struct{}

func (h *heartbeatAckProcessor) op() objects.OpCode {
	return objects.OpHeartbeatACK
}

func (h *heartbeatAckProcessor) process(s *Shard, p objects.Payload) error {
	s.heartbeat.ACK()
	return nil
}
