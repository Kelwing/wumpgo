package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"go.uber.org/atomic"
	"golang.org/x/time/rate"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	gatewayerrors "wumpgo.dev/wumpgo/gateway/internal/gateway_errors"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type Gateway struct {
	sync.Mutex
	conn        *Websocket
	sequence    *atomic.Int64
	identify    *Identify
	dispatcher  dispatcher.Dispatcher
	session_id  string
	resume      bool
	logger      zerolog.Logger
	gatewayInfo *objects.Gateway
	limiter     *rate.Limiter
	client      *rest.Client
	heartbeat   *Heartbeat
	hello       bool
	identified  *atomic.Bool
	hub         *sentry.Hub
	stopping    *atomic.Bool
	serverCount int
}

type GatewayConfig struct {
	Intents     []string
	Dispatcher  dispatcher.Dispatcher
	Token       string
	Shard       int
	ShardCount  int
	Client      *rest.Client
	GatewayInfo *objects.Gateway
	Hub         *sentry.Hub
}

func New(config *GatewayConfig) *Gateway {
	intents, err := ParseIntents(config.Intents)
	if err != nil {
		config.Hub.CaptureException(err)
		log.Panic().Err(err).Msg("failed to parse intents")
	}

	return &Gateway{
		conn:     NewWebsocket(log.Logger),
		sequence: atomic.NewInt64(0),
		identify: &Identify{
			Token:          config.Token,
			Intents:        intents,
			Shard:          []int{config.Shard, config.ShardCount},
			LargeThreshold: 50,
			Compress:       true,
			Properties: Properties{
				OS:      runtime.GOOS,
				Browser: "wumpgo",
				Device:  "wumpgo",
			},
		},
		dispatcher:  config.Dispatcher,
		logger:      log.With().Int("shard", config.Shard).Logger(),
		gatewayInfo: config.GatewayInfo,
		limiter:     rate.NewLimiter(2, 120),
		client:      config.Client,
		identified:  atomic.NewBool(false),
		hub:         config.Hub,
		stopping:    atomic.NewBool(false),
	}
}

func (g *Gateway) SetDispatcher(d dispatcher.Dispatcher) {
	g.dispatcher = d
}

func (g *Gateway) Servers() int {
	return g.serverCount
}

func (g *Gateway) connect() error {
	g.identified.Store(false)
	if !g.resume {
		g.sequence.Store(0)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	gatewayAddr := fmt.Sprintf(GatewayAddressFmt, g.gatewayInfo.URL, GatewayVersion, GatewayEncoding)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	return g.conn.Open(ctx, gatewayAddr, header)
}

func (g *Gateway) Close() error {
	return g.conn.Close()
}

func (g *Gateway) Receive() error {
	for {
		packet, err := g.conn.Read(context.Background())
		if err != nil {
			g.logger.Debug().Err(err).Msg("failed to read packet")
			return err
		}

		var p Payload
		err = json.Unmarshal(packet, &p)
		if err != nil {
			g.hub.CaptureException(err)
			g.logger.Error().Err(err).Msg("failed to unmarshal packet")
			return err
		}
		if p.Sequence > g.sequence.Load() {
			g.sequence.Store(p.Sequence)
		}

		log.Debug().Int64("sequence", p.Sequence).Int64("op", int64(p.Op)).Msg("received packet")

		if !g.hello && p.Op != OpHello {
			g.logger.Error().Int64("op", int64(p.Op)).Msg("expected Hello Op")
			return errGeneric("no hello")
		}

		switch p.Op {
		case OpDispatch:
			g.logger.Debug().Str("event", p.EventName).Msg("Received dispatch")
			if p.EventName == "READY" {
				ready := &objects.Ready{}
				err = json.Unmarshal(p.Data, ready)
				if err != nil {
					g.hub.CaptureException(err)
					g.logger.Err(err).Msg("Failed to unmarshal ready")
					return err
				}
				g.session_id = ready.SessionID
				g.logger.Info().Str("session_id", g.session_id).Str("user", ready.User.Username).Msg("We are ready!")
			}
			go func(event string, data json.RawMessage) {
				start := time.Now()
				err := g.dispatcher.Dispatch(p.EventName, p.Data)
				log.Debug().Dur("duration", time.Since(start)).Str("event", p.EventName).Msg("Dispatch finished")
				if err != nil {
					g.hub.CaptureException(err)
					g.logger.Err(err).Msg("Failed to dispatch")
				}
			}(p.EventName, p.Data)
		case OpHeartbeat:
			if err := g.heartbeat.SendHeartbeat(); err != nil {
				g.hub.CaptureException(err)
				g.logger.Err(err).Msg("Failed to send heartbeat")
				return err
			}
		case OpReconnect:
			g.logger.Info().Msg("Received reconnect")
			g.Close()
			return errReconnect()
		case OpInvalidSession:
			g.logger.Info().Msg("Invalid session")
			g.Close()
			return errInvalidSession()
		case OpHeartbeatACK:
			g.heartbeat.ACK()
		case OpHello:
			g.logger.Info().Msgf("Received Hello Op")
			g.hello = true
			hello := Hello{}
			err = json.Unmarshal(p.Data, &hello)
			if err != nil {
				g.hub.CaptureException(err)
				g.logger.Err(err).Msg("Failed to unmarshal hello")
				return err
			}
			hbHub := g.hub.Clone()
			hbHub.Scope().SetTag("component", "heartbeat")
			g.heartbeat = NewHeartbeat(hello.HeartbeatInterval, g, hbHub)
			defer g.heartbeat.Stop()
			g.heartbeat.Start()

			if g.resume && g.session_id != "" {
				g.logger.Info().Msg("Resuming session")
				g.resume = false
				err = g.sendResume()
				if err != nil {
					g.hub.CaptureException(err)
					g.logger.Err(err).Msg("Failed to resume session")
					return err
				}
			} else {
				g.logger.Info().Msg("Identifying")
				err = g.sendIdentify()
				if err != nil {
					g.hub.CaptureException(err)
					g.logger.Err(err).Msg("Failed to identify")
					return err
				}
			}
		}
	}
}

func (g *Gateway) Send(op OpCode, data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	p := Payload{
		Op:   op,
		Data: d,
	}
	reservation := g.limiter.Reserve()
	if !reservation.OK() {
		g.logger.Warn().Msg("Ratelimiter cannot provide a reservation in the maximum wait time")
		return nil
	}
	delay := reservation.Delay()
	g.logger.Debug().Dur("delay", delay).Msg("Ratelimiter reservation")
	time.Sleep(delay)
	g.logger.Debug().Int("op", int(op)).Msg("Sending payload")
	return g.conn.WriteJSON(&p)
}

func (g *Gateway) setResume() {
	g.resume = true
}

func (g *Gateway) sendIdentify() error {
	err := g.Send(OpIdentify, g.identify)
	if err != nil {
		g.hub.CaptureException(err)
		g.logger.Err(err).Msg("failed to send identify payload")
		return err
	}
	g.identified.Store(true)
	return nil
}

func (g *Gateway) sendResume() error {
	resume := Resume{
		Token:     g.identify.Token,
		Sequence:  g.sequence.Load(),
		SessionID: g.session_id,
	}

	log.Info().Int64("sequence", resume.Sequence).
		Str("session_id", resume.SessionID).
		Msg("Sending resume")

	err := g.Send(OpResume, resume)
	if err != nil {
		g.hub.CaptureException(err)
		g.logger.Err(err).Msg("failed to send resume payload")
		return err
	}
	g.identified.Store(true)
	return nil
}

func (g *Gateway) Run() error {
	var err error
	for {
		err = g.connect()
		if err != nil {
			g.hub.CaptureException(err)
			g.logger.Error().Err(err).Msg("failed to connect")
			time.Sleep(time.Second * 3)
			continue
		}

		err = g.Receive()
		if err != nil {
			g.logger.Err(err).Msg("failed to receive")
			cont, wait := g.ReconnectHandler(err)
			if cont {
				time.Sleep(wait)
				continue
			} else {
				return err
			}
		}
	}
}

func (g *Gateway) IsIdentified() bool {
	return g.identified.Load()
}

func (g *Gateway) ReconnectHandler(err error) (cont bool, wait time.Duration) {
	cont = true
	wait = time.Second * 3
	closeCode := websocket.CloseStatus(err)
	if g.conn.IsConnected() {
		g.conn.Close()
	}
	switch closeCode {
	case websocket.StatusNormalClosure:
		g.logger.Info().Msg("normal closure")
		if g.stopping.Load() {
			g.logger.Info().Msg("shutting down")
			cont = false
			return
		}
		return
	}

	if DiscordCloseCode(closeCode) >= CloseUnknownError {
		switch code := DiscordCloseCode(closeCode); code {
		case CloseAuthenticationFailed:
			fallthrough
		case CloseInvalidShard:
			fallthrough
		case CloseShardingRequired:
			fallthrough
		case CloseInvalidIntents:
			fallthrough
		case CloseDisallowedIntents:
			g.logger.Fatal().Msg(code.String())
		case CloseInvalidSeq:
			fallthrough
		case CloseSessionTimeout:
			g.logger.Warn().Int("code", int(code)).Msg(code.String())
		default:
			g.logger.Warn().Int("code", int(code)).Msg(code.String())
			g.setResume()
		}
		return
	}

	if errors.Is(err, ErrInvalidSession) {
		g.logger.Info().Msg("session invalidated")
		return
	}
	if errors.Is(err, ErrSessionStartLimitReached) {
		limitErr := err.(*gatewayerrors.ErrSessionStartLimitReached)
		g.logger.Warn().Msgf("session start limit reached, retrying in %sms", limitErr.ResetAfter)
		wait = time.Duration(limitErr.ResetAfter) * time.Millisecond
		return
	}
	if closeCode != -1 || errors.Is(err, ErrReconnect) {
		g.logger.Info().Msg("connection closed")
		g.setResume()
		return
	}
	g.hub.WithScope(func(scope *sentry.Scope) {
		scope.SetExtra("close_code", closeCode)
		g.hub.CaptureException(err)
	})
	g.logger.Error().Err(err).Msg("unknown or unhandled error")
	return
}
