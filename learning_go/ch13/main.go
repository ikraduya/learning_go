package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"slices"
	"time"
)

// {
// 	"day_of_week": "Monday",
// 	"day_of_month": 10,
// 	"month": "April",
// 	"year": 2023,
// 	"hour": 20,
// 	"minute": 15,
// 	"second": 20
// }

type TimeResponse struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth int    `json:"day_of_month"`
	Month      string `json:"month"`
	Year       int    `json:"year"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	Second     int    `json:"second"`
}

func CurrentTimeHandle(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	if slices.Contains(r.Header["Accept"], "application/json") {
		timeResponse := TimeResponse{
			DayOfWeek:  currentTime.Weekday().String(),
			DayOfMonth: currentTime.Day(),
			Month:      currentTime.Month().String(),
			Year:       currentTime.Year(),
			Hour:       currentTime.Hour(),
			Minute:     currentTime.Minute(),
			Second:     currentTime.Second(),
		}
		timeByte, err := json.Marshal(timeResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			w.Write([]byte("\n"))
			return
		}
		w.Write(timeByte)
		w.Write([]byte("\n"))
	} else {
		w.Write([]byte(currentTime.Format(time.RFC3339)))
		w.Write([]byte("\n"))
	}

}

func main() {
	mux := http.NewServeMux()

	mSlog := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	loggerMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Can't get client IP")
			}
			mSlog.Info("", "client_ip", clientIP)
			h.ServeHTTP(w, r)
		})
	}

	mux.Handle("GET /", loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CurrentTimeHandle(w, r)
		})))

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
