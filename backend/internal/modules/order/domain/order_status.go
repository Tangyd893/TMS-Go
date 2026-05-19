package domain

import "github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"

const (
	OrderStatusDraft           = "draft"
	OrderStatusSubmitted       = "submitted"
	OrderStatusPendingDispatch = "pending_dispatch"
	OrderStatusDispatched      = "dispatched"
	OrderStatusInTransit       = "in_transit"
	OrderStatusSigned          = "signed"
	OrderStatusSettled         = "settled"
	OrderStatusCancelled       = "cancelled"
)

var (
	ErrInvalidStatusTransition = errors.New(610001, "订单状态不允许此操作")
	ErrOrderAlreadyCancelled   = errors.New(610002, "订单已取消")
	ErrOrderAlreadySigned      = errors.New(610003, "订单已签收")
)

type OrderStateMachine struct {
	status string
}

func NewOrderStateMachine(status string) *OrderStateMachine {
	return &OrderStateMachine{status: status}
}

func (m *OrderStateMachine) Status() string {
	return m.status
}

func (m *OrderStateMachine) CanSubmit() error {
	if m.status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if m.status != OrderStatusDraft {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (m *OrderStateMachine) Submit() string {
	m.status = OrderStatusSubmitted
	return m.status
}

func (m *OrderStateMachine) CanMoveToPendingDispatch() error {
	if m.status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if m.status != OrderStatusSubmitted {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (m *OrderStateMachine) MoveToPendingDispatch() string {
	m.status = OrderStatusPendingDispatch
	return m.status
}

func (m *OrderStateMachine) CanDispatch() error {
	if m.status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if m.status != OrderStatusPendingDispatch {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (m *OrderStateMachine) Dispatch() string {
	m.status = OrderStatusDispatched
	return m.status
}

func (m *OrderStateMachine) CanStartTransit() error {
	if m.status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if m.status != OrderStatusDispatched {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (m *OrderStateMachine) StartTransit() string {
	m.status = OrderStatusInTransit
	return m.status
}

func (m *OrderStateMachine) CanSign() error {
	if m.status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if m.status == OrderStatusSigned {
		return ErrOrderAlreadySigned
	}
	if m.status != OrderStatusInTransit {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (m *OrderStateMachine) Sign() string {
	m.status = OrderStatusSigned
	return m.status
}

func (m *OrderStateMachine) CanSettle() error {
	if m.status != OrderStatusSigned {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (m *OrderStateMachine) Settle() string {
	m.status = OrderStatusSettled
	return m.status
}

func (m *OrderStateMachine) CanCancel() error {
	if m.status == OrderStatusSigned || m.status == OrderStatusSettled {
		return ErrOrderAlreadySigned
	}
	return nil
}

func (m *OrderStateMachine) Cancel() string {
	m.status = OrderStatusCancelled
	return m.status
}
