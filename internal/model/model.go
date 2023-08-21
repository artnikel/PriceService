// Package model contains models of using entities
package model

import (
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Share is a struct for shares entity
type Share struct {
	Company string
	Price   decimal.Decimal
}

// SubscribersManager contains all subscribers by uuid and their subscribed shares in map subscribers
type SubscribersManager struct {
	Mu               sync.RWMutex
	Subscribers      map[uuid.UUID][]string
	SubscribersShare map[uuid.UUID]chan Share
}
