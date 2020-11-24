package interactions

type Ratelimiter interface {
	Request(method, url, contentType string, body []byte) ([]byte, error)
}
