package rest

import (
	"errors"
	"regexp"
	"strings"
)

type Ratelimiter interface {
	Request(method, url, contentType string, body []byte) (*DiscordResponse, error)
}

var userIdRe = regexp.MustCompile(`\d{16,19}`)

var majorParams = [...]string{"channels", "guilds", "webhooks"}

var MaxRetriesExceeded = errors.New("max retried exceeded")

func isBucketedParam(comp string) bool {
	for _, param := range majorParams {
		if comp == param {
			return true
		}
	}
	return false
}

func getBucketID(urlStr string) string {
	comps := strings.Split(urlStr, "/")
	bucketComps := make([]string, 0)
	for i, comp := range comps[5:] {
		if comp == "reactions" {
			bucketComps = append(bucketComps, comp)
			break
		}

		if userIdRe.MatchString(comp) && !isBucketedParam(comps[5:][i]) {
			bucketComps = append(bucketComps, "id")
		} else {
			bucketComps = append(bucketComps, comp)
		}
	}

	return strings.Join(bucketComps, ":")
}
