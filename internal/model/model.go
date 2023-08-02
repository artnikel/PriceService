// Package model contains models of using entities
package model

import (
	"sync"
)

// Action is a struct for actions entity
type Action struct {
	Company string
	Price   float64
}

// SubscribersManager contains all subscribers by uuid and their subscribed shares in map subscribers
type SubscribersManager struct {
	Subscribers        sync.Map
	SubscribersActions sync.Map
}
