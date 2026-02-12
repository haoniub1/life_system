package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ipRecord struct {
	LoginFailures int
	RegisterCount int
	Date          string // YYYY-MM-DD, reset on new day
}

type Limiter struct {
	mu                sync.Mutex
	records           map[string]*ipRecord
	maxLoginFailures  int
	maxDailyRegisters int
}

func NewLimiter(maxLoginFailures, maxDailyRegisters int) *Limiter {
	return &Limiter{
		records:           make(map[string]*ipRecord),
		maxLoginFailures:  maxLoginFailures,
		maxDailyRegisters: maxDailyRegisters,
	}
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func (l *Limiter) getRecord(ip string) *ipRecord {
	rec, ok := l.records[ip]
	if !ok || rec.Date != today() {
		rec = &ipRecord{Date: today()}
		l.records[ip] = rec
	}
	return rec
}

// CheckLogin returns true if the IP is allowed to attempt login.
func (l *Limiter) CheckLogin(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	rec := l.getRecord(ip)
	return rec.LoginFailures < l.maxLoginFailures
}

// RecordLoginFailure increments the failed login counter for the IP.
func (l *Limiter) RecordLoginFailure(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	rec := l.getRecord(ip)
	rec.LoginFailures++
}

// CheckRegister returns true if the IP is allowed to register.
func (l *Limiter) CheckRegister(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	rec := l.getRecord(ip)
	return rec.RegisterCount < l.maxDailyRegisters
}

// RecordRegister increments the register counter for the IP.
func (l *Limiter) RecordRegister(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	rec := l.getRecord(ip)
	rec.RegisterCount++
}

// LoginFailuresRemaining returns how many attempts are left today.
func (l *Limiter) LoginFailuresRemaining(ip string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	rec := l.getRecord(ip)
	remaining := l.maxLoginFailures - rec.LoginFailures
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetClientIP extracts the real client IP from a request.
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For first (reverse proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
