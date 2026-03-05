package proxy

import (
"sync"
"time"
)

// State represents the state of the circuit breaker.
type State int

const (
StateClosed   State = iota // Normal operation
StateOpen                  // Rejecting requests
StateHalfOpen              // Testing if service recovered
)

func (s State) String() string {
switch s {
case StateClosed:
return "closed"
case StateOpen:
return "open"
case StateHalfOpen:
return "half-open"
default:
return "unknown"
}
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
mu               sync.Mutex
state            State
failureCount     int
successCount     int
failureThreshold int           // failures before opening
successThreshold int           // successes in half-open before closing
timeout          time.Duration // how long to stay open before half-open
lastFailure      time.Time
}

// NewCircuitBreaker creates a circuit breaker with sensible defaults.
func NewCircuitBreaker() *CircuitBreaker {
return &CircuitBreaker{
state:            StateClosed,
failureThreshold: 5,           // 5 consecutive failures → open
successThreshold: 2,           // 2 successes in half-open → close
timeout:          10 * time.Second, // wait 10s before trying again
}
}

// Allow checks if a request is allowed through. Returns false if circuit is open.
func (cb *CircuitBreaker) Allow() bool {
cb.mu.Lock()
defer cb.mu.Unlock()

switch cb.state {
case StateClosed:
return true
case StateOpen:
if time.Since(cb.lastFailure) > cb.timeout {
cb.state = StateHalfOpen
cb.successCount = 0
return true
}
return false
case StateHalfOpen:
return true
}
return true
}

// RecordSuccess records a successful request.
func (cb *CircuitBreaker) RecordSuccess() {
cb.mu.Lock()
defer cb.mu.Unlock()

switch cb.state {
case StateHalfOpen:
cb.successCount++
if cb.successCount >= cb.successThreshold {
cb.state = StateClosed
cb.failureCount = 0
}
case StateClosed:
cb.failureCount = 0
}
}

// RecordFailure records a failed request.
func (cb *CircuitBreaker) RecordFailure() {
cb.mu.Lock()
defer cb.mu.Unlock()

cb.lastFailure = time.Now()

switch cb.state {
case StateClosed:
cb.failureCount++
if cb.failureCount >= cb.failureThreshold {
cb.state = StateOpen
}
case StateHalfOpen:
cb.state = StateOpen
}
}

// State returns the current state.
func (cb *CircuitBreaker) GetState() State {
cb.mu.Lock()
defer cb.mu.Unlock()
return cb.state
}
