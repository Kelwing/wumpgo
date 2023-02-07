package shard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/atomic"
	"golang.org/x/time/rate"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	"wumpgo.dev/wumpgo/objects"
)

type IdentifyLocker interface {
	Lock()
	Unlock()
}

// Shard represents a single Shard connection
type Shard struct {
	conn         *Websocket
	seq          *atomic.Uint64
	identify     objects.Identify
	dispatcher   dispatcher.Dispatcher
	session_id   string
	resume_url   string
	resume       *atomic.Bool
	gateway_url  string
	hello        *atomic.Bool
	limiter      *rate.Limiter
	identified   *atomic.Bool
	stopping     *atomic.Bool
	done         chan error
	processors   map[objects.OpCode]packetProcessor
	identifyLock IdentifyLocker

	heartbeat *Heartbeat

	logger zerolog.Logger
}

func New(token string, opts ...ShardOption) *Shard {
	s := &Shard{
		seq: atomic.NewUint64(0),
		identify: objects.Identify{
			Token:          token,
			Intents:        objects.IntentsNone,
			Shard:          []int{0, 1},
			LargeThreshold: 50,
			Compress:       true,
			Properties: objects.Properties{
				OS:      runtime.GOOS,
				Browser: "wumpgo",
				Device:  "wumpgo",
			},
		},
		resume:       atomic.NewBool(false),
		hello:        atomic.NewBool(false),
		dispatcher:   dispatcher.NewNOOPDispatcher(),
		gateway_url:  fmt.Sprintf(GatewayAddressFmt, GatewayDefaultURL, GatewayVersion, GatewayEncoding),
		limiter:      rate.NewLimiter(2, 120),
		identified:   atomic.NewBool(false),
		stopping:     atomic.NewBool(false),
		logger:       zerolog.Nop(), // By default log nothing
		processors:   make(map[objects.OpCode]packetProcessor),
		identifyLock: nil,
	}

	for _, o := range opts {
		o(s)
	}

	s.addProcessors(
		&dispatchProcessor{},
		&helloProcessor{},
		&heartbeatProcessor{},
		&reconnectProcessor{},
		&invalidSessionProcessor{},
		&heartbeatAckProcessor{},
	)

	s.conn = NewWebsocket(&s.logger)

	return s
}

func (s *Shard) MarshalZerologObject(e *zerolog.Event) {
	e.Int("shard_id", s.identify.Shard[0]).Bool("identified", s.IsIdentified())
}

func (s *Shard) String() string {
	return "Shard " + strconv.Itoa(s.identify.Shard[0])
}

func (s *Shard) addProcessors(processors ...packetProcessor) {
	for _, p := range processors {
		s.processors[p.op()] = p
	}
}

func (s *Shard) Send(op objects.OpCode, data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	p := objects.Payload{
		Op:   op,
		Data: d,
	}
	reservation := s.limiter.Reserve()
	if !reservation.OK() {
		s.logger.Warn().Msg("Ratelimiter cannot provide a reservation in the maximum wait time")
		return nil
	}
	delay := reservation.Delay()
	s.logger.Debug().Dur("delay", delay).Msg("Ratelimiter reservation")
	time.Sleep(delay)
	s.logger.Debug().Int("op", int(op)).Msg("Sending payload")
	return s.conn.WriteJSON(&p)
}

func (s *Shard) IsIdentified() bool {
	return s.identified.Load()
}

func (s *Shard) Close() {
	s.done <- errGeneric("connection closed")
}

func (s *Shard) close() error {
	return s.conn.Close()
}

func (s *Shard) setResume() {
	s.resume.Store(true)
}

func (s *Shard) sendIdentify() error {
	if s.identifyLock != nil {
		s.identifyLock.Lock()
		defer func() {
			time.Sleep(time.Second * 5)
			s.identifyLock.Unlock()
		}()
	}
	err := s.Send(objects.OpIdentify, s.identify)
	if err != nil {
		s.logger.Err(err).Msg("failed to send identify payload")
		return err
	}
	s.identified.Store(true)

	return nil
}

func (s *Shard) sendResume() error {
	resume := objects.Resume{
		Token:     s.identify.Token,
		Sequence:  s.seq.Load(),
		SessionID: s.session_id,
	}

	log.Info().Uint64("sequence", resume.Sequence).
		Str("session_id", resume.SessionID).
		Msg("Sending resume")

	err := s.Send(objects.OpResume, resume)
	if err != nil {
		s.logger.Err(err).Msg("failed to send resume payload")
		return err
	}
	s.identified.Store(true)
	return nil
}

func (s *Shard) connect() error {
	s.identified.Store(false)
	if !s.resume.Load() {
		s.seq.Store(0)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	url := s.gateway_url
	if s.resume.Load() {
		url = s.resume_url
	}

	log.Debug().Str("url", url).Msg("opening websocket connection")

	return s.conn.Open(ctx, url, header)
}

func (s *Shard) read(out chan objects.Payload) error {
	for {
		packet, err := s.conn.Read(context.Background())
		if err != nil {
			return err
		}

		var p objects.Payload
		err = json.Unmarshal(packet, &p)
		if err != nil {
			return err
		}

		out <- p
	}
}

func (s *Shard) receive() error {
	s.done = make(chan error)
	msgs := make(chan objects.Payload, 10)

	go func() {
		s.logger.Debug().Msg("starting read loop")
		if err := s.read(msgs); err != nil {
			s.logger.Err(err).Msg("read failed")
			s.done <- err
		}
		s.logger.Debug().Msg("read loop stopped")
	}()

	for {
		select {
		case err := <-s.done:
			s.close()
			return err
		case p := <-msgs:
			if err := s.process(p); err != nil {
				s.close()
				return err
			}
		}
	}
}

func (s *Shard) process(p objects.Payload) error {
	if p.Sequence > s.seq.Load() {
		s.seq.Store(p.Sequence)
	}

	s.logger.Debug().Uint64("sequence", p.Sequence).Int64("op", int64(p.Op)).Msg("received packet")

	if !s.hello.Load() && p.Op != objects.OpHello {
		s.logger.Error().Int64("op", int64(p.Op)).Msg("expected Hello Op")
		return errGeneric("no hello")
	}

	if processor, ok := s.processors[p.Op]; ok {
		if err := processor.process(s, p); err != nil {
			return err
		}
	} else {
		log.Warn().Int64("op", int64(p.Op)).Msg("no process found for op")
	}

	return nil
}

func (s *Shard) Run() error {
	var err error
	for {
		err = s.connect()
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to connect")
			time.Sleep(time.Second * 3)
			continue
		}

		err = s.receive()
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to receive")
			if errors.Is(err, ErrReconnect) {
				time.Sleep(time.Second * 3)
				s.setResume()
				continue
			}
			return err
		}
	}
}
