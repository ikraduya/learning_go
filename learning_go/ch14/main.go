package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

func middleware(timoutMs int) func(http.Handler) http.Handler {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timoutMs)*time.Millisecond)
	defer cancel()
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
		})
	}
}

func generate1234() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	var iterNum int
loop:
	for {
		iterNum += 1
		a := rand.Int() % 100_000_000
		b := rand.Int() % 100_000_000
		sum := a + b
		if sum == 1234 {
			fmt.Println("Found 1234 at iteration", iterNum, "last sum:", sum)
			break
		}
		select {
		case <-ctx.Done():
			fmt.Println("Time limit 2s exceed at iteration", iterNum, "last sum:", sum)
			break loop
		default:
		}
	}
}

type Level string

const (
	Debug = Level("debug")
	Info  = Level("info")
)

type logLevelKey int

const (
	_ logLevelKey = iota
	key
)

func contextWithLevel(ctx context.Context, level Level) context.Context {
	return context.WithValue(ctx, key, level)
}

func levelFromContext(ctx context.Context) (Level, bool) {
	l, ok := ctx.Value(key).(Level)
	return l, ok
}

func Log(ctx context.Context, level Level, message string) {
	var inLevel Level
	inLevel, ok := levelFromContext(ctx)
	if !ok {
		return
	}

	if level == Debug && inLevel == Debug {
		fmt.Println(message)
	}
	if level == Info && (inLevel == Debug || inLevel == Info) {
		fmt.Println(message)
	}
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		level := r.URL.Query().Get("log_level")
		ctx := contextWithLevel(r.Context(), Level(level))
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func main() {
	generate1234()
}
