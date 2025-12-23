package test

import (
	"context"
	"fmt"
	"strings"

	"github.com/foxie-io/ng"
)

type TraceMiddleware struct {
	ng.DefaultID[TraceMiddleware]
}

func (u TraceMiddleware) Use(ctx context.Context, next ng.Handler) {
	t := &Tracer{}
	ng.Store(ctx, t)
	next(ctx)
}

type Tracer struct {
	forwards  []string
	backwards []string
}

func (t *Tracer) traceForward(pairs ...string) {
	if len(pairs)%2 != 0 {
		panic("pairs must be even")
	}
	t.forwards = append(t.forwards, pairs...)
}

func (t *Tracer) traceBackward(pairs ...string) {
	if len(pairs)%2 != 0 {
		panic("pairs must be even")
	}
	t.backwards = append(t.backwards, pairs...)
}

func (t *Tracer) Tree() string {
	indent := "   "
	whitespace := ""
	str := ""
	for i := 0; i < len(t.forwards); i += 2 {
		k, v := t.forwards[i], t.forwards[i+1]
		str += fmt.Sprintf("%s%s %s\n", whitespace, k, v)
		isGuard := strings.Contains(v, "guard")
		isHandler := strings.Contains(k, "Handler")

		switch {
		case isHandler:
			whitespace = strings.Replace(whitespace, indent, "", 1)
		case isGuard:
		default:
			whitespace = whitespace + indent
		}
	}

	for i := 0; i < len(t.backwards); i += 2 {
		k, v := t.backwards[i], t.backwards[i+1]

		str += fmt.Sprintf("%s%s %s\n", whitespace, k, v)
		whitespace = strings.Replace(whitespace, indent, "", 1)
	}
	return str
}

func (t *Tracer) PrintTree() {
	fmt.Println("Trace Tree:\n", t.Tree())
}
