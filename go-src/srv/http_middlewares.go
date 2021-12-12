package srv

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func withMiddlewares(h http.Handler, middlewares ...(func(next http.Handler) http.Handler)) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec != nil {
				switch recValue := rec.(type) {
				case error:
					writeError(w, recValue)
				case string:
					writeError(w, fmt.Errorf(recValue))
				default:
					writeError(w, fmt.Errorf("internal error"))
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RateLimiterMiddleware(maxFrequency int, burstSize int, memoryDuration time.Duration) func(next http.Handler) http.Handler {
	// Using rate limiter from package "golang.org/x/time/rate"
	rateLimiterPerUser := make(map[string]*rate.Limiter)
	var mutex sync.Mutex

	// Cleaning the map from time to time to release the memory
	go func(m *sync.Mutex) {
		for {
			time.Sleep(memoryDuration)
			m.Lock()
			for user := range rateLimiterPerUser {
				delete(rateLimiterPerUser, user)
			}
			m.Unlock()
		}
	}(&mutex)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Getting user IP address (first checking if the server is under a reverse proxy by
			// trying to get it from the headers "X-Real-Ip" and "X-Forwarded-For")
			userIPAddress := r.Header.Get("X-Real-Ip")
			if userIPAddress == "" {
				userIPAddress = r.Header.Get("X-Forwarded-For")
			}
			if userIPAddress == "" {
				var err error
				userIPAddress, _, err = net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					writeError(w, err)
					return
				}
			}

			// Creating rate limiter for the user (if not created yet)
			mutex.Lock()
			userRateLimiter, alreadyCreated := rateLimiterPerUser[userIPAddress]
			if !alreadyCreated {
				userRateLimiter = rate.NewLimiter(rate.Limit(maxFrequency), burstSize)
				rateLimiterPerUser[userIPAddress] = userRateLimiter
			}
			mutex.Unlock()

			// Checking user rate limiter
			if !userRateLimiter.Allow() {
				writeError(w, &httpError{
					StatusCode: http.StatusTooManyRequests,
					Message:    "Too many requests",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
