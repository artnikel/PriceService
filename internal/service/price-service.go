package service

import (
	"context"
	"fmt"
)

type PriceRepository interface {
	ReadPrices(ctx context.Context) (map[string]float64, error)
}

type PriceService struct {
	priceRep PriceRepository
}

func NewPriceService(priceRep PriceRepository) *PriceService {
	return &PriceService{priceRep: priceRep}
}

func (p *PriceService) ReadPrices(ctx context.Context) (map[string]float64, error) {
	actions, err := p.priceRep.ReadPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("Price-ReadPrices: error: %w", err)
	}
	return actions, nil
}
