// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/artnikel/PriceService/proto"
	"github.com/google/uuid"
)

// PriceRepository is interface with method for reading prices
type PriceRepository interface {
	ReadPrices(ctx context.Context) ([]*model.Share, error)
}

// PriceService contains PriceRepository interface
type PriceService struct {
	manager  *model.SubscribersManager
	priceRep PriceRepository
}

// NewPriceService accepts PriceRepository object and returnes an object of type *PriceService
func NewPriceService(priceRep PriceRepository) *PriceService {
	return &PriceService{
		priceRep: priceRep,
		manager: &model.SubscribersManager{SubscribersShare: make(map[uuid.UUID]chan []*model.Share),
			Subscribers: make(map[uuid.UUID][]string)}}
}

// AddSubscriber adds new subscriber to subscribe map in SubscriberManager
func (p *PriceService) AddSubscriber(subscriberID uuid.UUID, selectedActions []string) error {
	const msgs = 1
	p.manager.Mu.Lock()
	defer p.manager.Mu.Unlock()
	if _, ok := p.manager.Subscribers[subscriberID]; !ok {
		p.manager.Subscribers[subscriberID] = selectedActions
		p.manager.SubscribersShare[subscriberID] = make(chan []*model.Share, msgs)
		return nil
	}
	return fmt.Errorf("PriceService-AddSubscriber: error: subscriber with such ID already exists")
}

// DeleteSubscriber deletes a subscriber from the subscribe map in SubscriberManager by uuid
func (p *PriceService) DeleteSubscriber(subscriberID uuid.UUID) error {
	p.manager.Mu.Lock()
	defer p.manager.Mu.Unlock()
	if _, ok := p.manager.Subscribers[subscriberID]; ok {
		delete(p.manager.Subscribers, subscriberID)
		close(p.manager.SubscribersShare[subscriberID])
		delete(p.manager.SubscribersShare, subscriberID)
		return nil
	}
	return fmt.Errorf("PriceService-DeleteSubscriber: error: subscriber with ID %s doesn't exist", subscriberID)
}

// ReadPrices is a method of GeneratorService that calls method of Repository
func (p *PriceService) ReadPrices(ctx context.Context) (shares []*model.Share, e error) {
	shares, err := p.priceRep.ReadPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("PriceService-ReadPrices: error: %w", err)
	}
	return shares, nil
}

// SendToAllSubscribedChans is a method that send to all subscribes chanells
func (p *PriceService) SendToAllSubscribedChans(ctx context.Context) {
	for {
		if len(p.manager.Subscribers) > 0 {
			actions, err := p.ReadPrices(ctx)
			if err != nil {
				log.Fatalf("PriceServiceService-SendToAllSubscribedChans-ReadPrices: error %v", err)
				return
			}
			for subID, selectedActions := range p.manager.Subscribers {
				tempActions := make([]*model.Share, 0)
				for _, action := range actions {
					if strings.Contains(strings.Join(selectedActions, ","), action.Company) {
						tempActions = append(tempActions, action)
					}
				}
				select {
				case <-ctx.Done():
					return
				case p.manager.SubscribersShare[subID] <- tempActions:
				}
			}
		}
	}
}

// SendToSubscriber calls SendToSubscriber method of repository
func (p *PriceService) SendToSubscriber(ctx context.Context, subscriberID uuid.UUID) (protoShares []*proto.Shares, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case shares := <-p.manager.SubscribersShare[subscriberID]:
		for _, share := range shares {
			protoShares = append(protoShares, &proto.Shares{
				Company: share.Company,
				Price:   share.Price.InexactFloat64(),
			})
		}
		return protoShares, nil
	}
}
