package router

import (
	"strings"
)

const (
	PlaceholderStart = '{'
	PlaceholderEnd   = '}'
)

type trieNode[T any] struct {
	children map[string]*trieNode[T]
	handler  T
}

func newTrieNode[T any]() *trieNode[T] {
	return &trieNode[T]{
		children: make(map[string]*trieNode[T]),
	}
}

func (t *trieNode[T]) Insert(key string, handler T) {
	parts := strings.Split(key, "/")

	node := t

	for _, p := range parts {
		_, ok := node.children[p]
		if !ok {
			node.children[p] = newTrieNode[T]()
		}

		node = node.children[p]
	}

	node.handler = handler
}

func (t *trieNode[T]) Search(key string) (T, map[string]string, bool) {
	parts := strings.Split(key, "/")

	var empty T

	node := t

	placeholders := make(map[string]string)

	for i, p := range parts {
		nextNode, ok := node.children[p]
		if ok {
			node = nextNode
			continue
		}

		if len(node.children) == 0 {
			if i != len(parts)-1 {
				return empty, placeholders, false
			}
		}

		hasPlaceholder := false
		for k, v := range node.children {
			if len(k) == 0 {
				continue
			}

			if k[0] == PlaceholderStart && k[len(k)-1] == PlaceholderEnd {
				node = v
				placeholders[k[1:len(k)-1]] = p
				hasPlaceholder = true
				break
			}
		}

		if !hasPlaceholder {
			return empty, nil, false
		}
	}

	return node.handler, placeholders, true
}
