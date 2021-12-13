package srv

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// withMiddlewares creates a new HTTP handler with the chain of middlewares received as parameter
func withMiddlewares(h http.Handler, middlewares ...(func(next http.Handler) http.Handler)) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// PanicRecoveryMiddleware treats any panic error that happens after this middleware
// and writes correct error to HTTP response and log
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defering panic recover logic (i.e. catching it after next handler(s) execution)
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

// ETagMiddleware adds ETag header for proper client caching based on a version received as parameter (e.g. the users data version)
// and returns 304 status code (not modified) if client requests the same version
func ETagMiddleware(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("ETag", version)

			if r.Header.Get("If-None-Match") == version {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware sets proper CORS headers and handles the preflight OPTIONS request
func CORSMiddleware(allowOrigin string, allowMethods []string, allowHeaders []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))

			// just returns if it's a prefligh request
			if r.Method == http.MethodOptions {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiterMiddleware blocks the user from making a big amount of requests in a small amount of time,
// receives some configuration:
// - maxFrequency: maximum allowed frequency (requests per second)
// - burstSize: maximum bursts permitted
// - memoryDuration: duration of users rate limiter memory before it's cleaned
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
				// returns status code 429 ("too many requests") if rate limit is reached
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
