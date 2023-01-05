package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"wumpgo.dev/wumpgo/gateway/receiver"
	"wumpgo.dev/wumpgo/router"
)

type ErrLoadDir struct {
	errors []error
}

func (e *ErrLoadDir) Error() string {
	sb := strings.Builder{}

	for i, err := range e.errors {
		sb.WriteString(err.Error())
		if i < len(e.errors)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (r *ErrLoadDir) Errors() []error {
	return r.errors
}

type Loader struct {
	rec    receiver.Receiver
	router *router.Router
}

func New(opts ...LoaderOption) *Loader {
	l := &Loader{}

	for _, o := range opts {
		o(l)
	}

	return l
}

func (l *Loader) LoadFile(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	setup, err := p.Lookup("Setup")
	if err != nil {
		return err
	}

	switch f := setup.(type) {
	case func(router *router.Router, rec receiver.Receiver) error:
		return f(l.router, l.rec)
	case func(router *router.Router) error:
		return f(l.router)
	case func(rec receiver.Receiver) error:
		return f(l.rec)
	default:
		return fmt.Errorf("setup function signature is invalid")
	}
}

func (l *Loader) LoadDir(p string) error {
	files, err := os.ReadDir(p)
	if err != nil {
		return err
	}

	errs := make([]error, 0)

	for _, fi := range files {
		filePath := filepath.Join(p, fi.Name())
		if err := l.LoadFile(filePath); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return &ErrLoadDir{
			errors: errs,
		}
	}

	return nil
}
