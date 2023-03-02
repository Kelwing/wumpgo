package shard

import "fmt"

type ShardError struct {
	Message string
}

func (s ShardError) Error() string {
	return fmt.Sprintf("ShardError: %s", s.Message)
}

func shardError(m string) ShardError {
	return ShardError{Message: m}
}
