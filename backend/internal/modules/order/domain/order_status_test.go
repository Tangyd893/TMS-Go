package domain

import (
	"testing"
)

func TestStateMachineDraftToSubmitted(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusDraft)
	if err := sm.CanSubmit(); err != nil {
		t.Fatalf("CanSubmit() from draft should succeed, got: %v", err)
	}
	status := sm.Submit()
	if status != OrderStatusSubmitted {
		t.Fatalf("Submit() = %s, want %s", status, OrderStatusSubmitted)
	}
}

func TestStateMachineSubmittedToPendingDispatch(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSubmitted)
	if err := sm.CanMoveToPendingDispatch(); err != nil {
		t.Fatalf("CanMoveToPendingDispatch() from submitted should succeed")
	}
	status := sm.MoveToPendingDispatch()
	if status != OrderStatusPendingDispatch {
		t.Fatalf("MoveToPendingDispatch() = %s, want %s", status, OrderStatusPendingDispatch)
	}
}

func TestStateMachinePendingDispatchToDispatched(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusPendingDispatch)
	if err := sm.CanDispatch(); err != nil {
		t.Fatalf("CanDispatch() from pending_dispatch should succeed")
	}
	status := sm.Dispatch()
	if status != OrderStatusDispatched {
		t.Fatalf("Dispatch() = %s, want %s", status, OrderStatusDispatched)
	}
}

func TestStateMachineDispatchedToInTransit(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusDispatched)
	if err := sm.CanStartTransit(); err != nil {
		t.Fatalf("CanStartTransit() from dispatched should succeed")
	}
	status := sm.StartTransit()
	if status != OrderStatusInTransit {
		t.Fatalf("StartTransit() = %s, want %s", status, OrderStatusInTransit)
	}
}

func TestStateMachineInTransitToSigned(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusInTransit)
	if err := sm.CanSign(); err != nil {
		t.Fatalf("CanSign() from in_transit should succeed")
	}
	status := sm.Sign()
	if status != OrderStatusSigned {
		t.Fatalf("Sign() = %s, want %s", status, OrderStatusSigned)
	}
}

func TestStateMachineSignedToSettled(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSigned)
	if err := sm.CanSettle(); err != nil {
		t.Fatalf("CanSettle() from signed should succeed")
	}
	status := sm.Settle()
	if status != OrderStatusSettled {
		t.Fatalf("Settle() = %s, want %s", status, OrderStatusSettled)
	}
}

func TestStateMachineInvalidTransition(t *testing.T) {
	tests := []struct {
		name   string
		status string
		action func(*OrderStateMachine) error
	}{
		{"Submit from submitted", OrderStatusSubmitted, (*OrderStateMachine).CanSubmit},
		{"Dispatch from draft", OrderStatusDraft, (*OrderStateMachine).CanDispatch},
		{"Sign from draft", OrderStatusDraft, (*OrderStateMachine).CanSign},
		{"StartTransit from submitted", OrderStatusSubmitted, (*OrderStateMachine).CanStartTransit},
		{"Settle from draft", OrderStatusDraft, (*OrderStateMachine).CanSettle},
		{"MoveToPendingDispatch from draft", OrderStatusDraft, (*OrderStateMachine).CanMoveToPendingDispatch},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewOrderStateMachine(tt.status)
			if err := tt.action(sm); err == nil {
				t.Fatalf("%s should fail with ErrInvalidStatusTransition", tt.name)
			}
		})
	}
}

func TestStateMachineCancel(t *testing.T) {
	cancelableStatuses := []string{
		OrderStatusDraft, OrderStatusSubmitted, OrderStatusPendingDispatch,
		OrderStatusDispatched, OrderStatusInTransit,
	}
	for _, status := range cancelableStatuses {
		t.Run("Cancel from "+status, func(t *testing.T) {
			sm := NewOrderStateMachine(status)
			if err := sm.CanCancel(); err != nil {
				t.Fatalf("CanCancel() from %s should succeed, got: %v", status, err)
			}
			if s := sm.Cancel(); s != OrderStatusCancelled {
				t.Fatalf("Cancel() = %s from %s, want %s", s, status, OrderStatusCancelled)
			}
		})
	}
}

func TestStateMachineCannotCancelSigned(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSigned)
	if err := sm.CanCancel(); err != ErrOrderAlreadySigned {
		t.Fatalf("CanCancel() from signed should return ErrOrderAlreadySigned, got: %v", err)
	}
}

func TestStateMachineCannotCancelSettled(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSettled)
	if err := sm.CanCancel(); err != ErrOrderAlreadySigned {
		t.Fatalf("CanCancel() from settled should return ErrOrderAlreadySigned, got: %v", err)
	}
}

func TestStateMachineCannotSubmitFromCancelled(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusDraft)
	sm.Cancel()
	if err := sm.CanSubmit(); err != ErrOrderAlreadyCancelled {
		t.Fatalf("CanSubmit() from cancelled should return ErrOrderAlreadyCancelled, got: %v", err)
	}
}

func TestStateMachineCannotSignTwice(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusInTransit)
	sm.Sign()
	if err := sm.CanSign(); err != ErrOrderAlreadySigned {
		t.Fatalf("CanSign() on already signed should return ErrOrderAlreadySigned, got: %v", err)
	}
}
