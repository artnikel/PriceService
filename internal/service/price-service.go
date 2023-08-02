// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/artnikel/PriceService/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PriceRepository is interface with method for reading prices
type PriceRepository interface {
	ReadPrices(ctx context.Context) ([]*model.Action, error)
}

// PriceService contains PriceRepository interface
type PriceService struct {
	manager  model.SubscribersManager
	priceRep PriceRepository
}

// NewPriceService accepts PriceRepository object and returnes an object of type *PriceService
func NewPriceService(priceRep PriceRepository) *PriceService {
	return &PriceService{
		priceRep: priceRep,
		manager: model.SubscribersManager{
			Subscribers:        sync.Map{},
			SubscribersActions: sync.Map{},
		},
	}
}

// AddSubscriber adds new subscriber to subscribe map in SubscriberManager
func (p *PriceService) AddSubscriber(subscriberID uuid.UUID, selectedActions []string) error {
	const messages = 1
	if subscriberID == uuid.Nil {
		return fmt.Errorf("PriceService-AddSubscriber: error: subscriber has nil uuid")
	}
	if len(selectedActions) == 0 {
		return fmt.Errorf("PriceService-AddSubscriber: error: subscriber hasn't subscribed to any shares")
	}
	if _, loaded := p.manager.Subscribers.LoadOrStore(subscriberID, selectedActions); loaded {
		return fmt.Errorf("PriceService-AddSubscriber: error: subscriber with such ID already exists")
	}
	if _, loaded := p.manager.SubscribersActions.LoadOrStore(subscriberID, make(chan []*model.Action, messages)); loaded {
		p.manager.Subscribers.Delete(subscriberID)
		return fmt.Errorf("PriceService-AddSubscriber: error: subscriber with such ID already exists")
	}
	return nil
}

// DeleteSubscriber delete subscriber from subscribe map in SubscriberManager by uuid
func (p *PriceService) DeleteSubscriber(subscriberID uuid.UUID) error {
	if _, found := p.manager.Subscribers.LoadAndDelete(subscriberID); !found {
		return fmt.Errorf("PriceService-DeleteSubscriber: error: subscriber with such ID doesn't exist")
	}
	if ch, found := p.manager.SubscribersActions.LoadAndDelete(subscriberID); found {
		close(ch.(chan []*model.Action))
	}
	return nil
}

// ReadPrices is a method of GeneratorService that calls method of Repository
func (p *PriceService) ReadPrices(ctx context.Context) (actions []*model.Action, e error) {
	actions, err := p.priceRep.ReadPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("PriceService-ReadPrices: error: %w", err)
	}
	return actions, nil
}

// SendToAllSubscribedChans sends in loop actual info about subscribed shares to subscribers chans
func (p *PriceService) SendToAllSubscribedChans(ctx context.Context) {
	for {
		subscriberExists := false
		p.manager.Subscribers.Range(func(subID, selectedActions interface{}) bool {
			subscriberExists = true
			subscriberID := subID.(uuid.UUID)
			subscriberActionsChan, ok := p.manager.SubscribersActions.Load(subscriberID)
			if !ok {
				return true
			}
			actions, err := p.ReadPrices(ctx)
			if err != nil {
				logrus.Errorf("PriceServiceService-SendToAllSubscribedChans: error:%v", err)
				return false
			}
			tempActions := make([]*model.Action, 0)
			for _, action := range actions {
				selected := selectedActions.([]string)
				for _, selectedCompany := range selected {
					if action.Company == selectedCompany {
						tempActions = append(tempActions, action)
						break
					}
				}
			}
			select {
			case <-ctx.Done():
				return false
			case subscriberActionsChan.(chan []*model.Action) <- tempActions:
			}
			return true
		})
		if !subscriberExists {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}
}

// SendToSubscriber calls SendToSubscriber method of repository
func (p *PriceService) SendToSubscriber(ctx context.Context, subscriberID uuid.UUID) (protoActions []*proto.Actions, err error) {
	actionsChan, found := p.manager.SubscribersActions.Load(subscriberID)
	if !found {
		return nil, fmt.Errorf("PriceService-SendYoSubscriber: subscriber with ID %s not found", subscriberID)
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case actions := <-actionsChan.(chan []*model.Action):
		for _, action := range actions {
			protoActions = append(protoActions, &proto.Actions{
				Company: action.Company,
				Price:   action.Price,
			})
		}
		return protoActions, nil
	}
}
