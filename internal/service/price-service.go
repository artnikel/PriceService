// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"
	"log"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/artnikel/PriceService/proto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
		manager: &model.SubscribersManager{SubscribersShare: make(map[uuid.UUID]chan model.Share),
			Subscribers: make(map[uuid.UUID][]string)}}
}

// AddSubscriber adds new subscriber to subscribe map in SubscriberManager
func (p *PriceService) AddSubscriber(subscriberID uuid.UUID, selectedShares []string) error {
	p.manager.Mu.Lock()
	defer p.manager.Mu.Unlock()
	if _, ok := p.manager.Subscribers[subscriberID]; !ok {
		p.manager.Subscribers[subscriberID] = selectedShares
		p.manager.SubscribersShare[subscriberID] = make(chan model.Share, len(selectedShares))
		return nil
	}
	return fmt.Errorf("addSubscriber ")
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
	return fmt.Errorf("priceService: subscriber with ID %s doesn't exist", subscriberID)
}

// ReadPrices is a method of GeneratorService that calls method of Repository
func (p *PriceService) ReadPrices(ctx context.Context) (shares []*model.Share, e error) {
	shares, err := p.priceRep.ReadPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("priceService %w", err)
	}
	return shares, nil
}

// SubscribeAll is a method that send to all subscribes chanells
func (p *PriceService) SubscribeAll(ctx context.Context) {
	subShares := make(map[string]decimal.Decimal)
	for {
		if len(p.manager.Subscribers) == 0 {
			continue
		}
		shares, err := p.ReadPrices(ctx)
		if err != nil {
			log.Fatalf("subscribeAll %v", err)
			return
		}
		for _, share := range shares {
			subShares[share.Company] = share.Price
		}
		p.manager.Mu.Lock()
		for subscriberID, selectedShares := range p.manager.Subscribers {
			if len(p.manager.SubscribersShare[subscriberID]) != 0 {
				continue
			}
			for _, selectedShare := range selectedShares {
				select {
				case <-ctx.Done():
					p.manager.Mu.Unlock()
					return
				case p.manager.SubscribersShare[subscriberID] <- model.Share{Company: selectedShare, Price: subShares[selectedShare]}:
				}
			}
		}
		p.manager.Mu.Unlock()
	}
}

// SendToSubscriber send shares in proto format to subscriber
func (p *PriceService) SendToSubscriber(ctx context.Context, subscriberID uuid.UUID) (protoShares []*proto.Shares, e error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case share := <-p.manager.SubscribersShare[subscriberID]:
		protoShares = append(protoShares, &proto.Shares{
			Company: share.Company,
			Price:   share.Price.InexactFloat64(),
		})
		for i := 1; i < len(p.manager.Subscribers[subscriberID]); i++ {
			share = <-p.manager.SubscribersShare[subscriberID]
			protoShares = append(protoShares, &proto.Shares{
				Company: share.Company,
				Price:   share.Price.InexactFloat64(),
			})
		}
		return protoShares, nil
	}
}
