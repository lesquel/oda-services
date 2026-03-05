package proxy

import (
	"testing"
	"time"
)

func TestCircuitBreaker_StartsClosedAndAllows(t *testing.T) {
	cb := NewCircuitBreaker()
	if cb.GetState() != StateClosed {
		t.Fatalf("expected closed, got %s", cb.GetState())
	}
	if !cb.Allow() {
		t.Fatal("expected allow in closed state")
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker()
	for i := 0; i < cb.failureThreshold; i++ {
		cb.RecordFailure()
	}
	if cb.GetState() != StateOpen {
		t.Fatalf("expected open after %d failures, got %s", cb.failureThreshold, cb.GetState())
	}
	if cb.Allow() {
		t.Fatal("expected deny in open state")
	}
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker()
	for i := 0; i < cb.failureThreshold-1; i++ {
		cb.RecordFailure()
	}
	cb.RecordSuccess()
	for i := 0; i < cb.failureThreshold-1; i++ {
		cb.RecordFailure()
	}
	if cb.GetState() != StateClosed {
		t.Fatalf("expected closed (success should have reset count), got %s", cb.GetState())
	}
}

func TestCircuitBreaker_TransitionsToHalfOpen(t *testing.T) {
	cb := &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: 2,
		successThreshold: 1,
		timeout:          50 * time.Millisecond,
	}
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.GetState() != StateOpen {
		t.Fatalf("expected open, got %s", cb.GetState())
	}
	time.Sleep(60 * time.Millisecond)
	if !cb.Allow() {
		t.Fatal("expected allow after timeout (half-open)")
	}
	if cb.GetState() != StateHalfOpen {
		t.Fatalf("expected half-open, got %s", cb.GetState())
	}
}

func TestCircuitBreaker_HalfOpenClosesOnSuccess(t *testing.T) {
	cb := &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: 2,
		successThreshold: 1,
		timeout:          50 * time.Millisecond,
	}
	cb.RecordFailure()
	cb.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	cb.Allow()
	cb.RecordSuccess()
	if cb.GetState() != StateClosed {
		t.Fatalf("expected closed after success in half-open, got %s", cb.GetState())
	}
}

func TestCircuitBreaker_HalfOpenReopensOnFailure(t *testing.T) {
	cb := &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: 2,
		successThreshold: 2,
		timeout:          50 * time.Millisecond,
	}
	cb.RecordFailure()
	cb.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	cb.Allow()
	cb.RecordFailure()
	if cb.GetState() != StateOpen {
		t.Fatalf("expected open after failure in half-open, got %s", cb.GetState())
	}
}

func TestState_String(t *testing.T) {
	cases := map[State]string{
		StateClosed:   "closed",
		StateOpen:     "open",
		StateHalfOpen: "half-open",
		State(99):     "unknown",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("State(%d).String() = %q, want %q", int(s), got, want)
		}
	}
}
