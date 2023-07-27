package service

import (
	"context"
	"fmt"
)

type PriceRepository interface {
	ReadPrices(ctx context.Context, company string) (float64, error)
}

type PriceService struct {
	priceRep PriceRepository
}

func NewPriceService(priceRep PriceRepository) *PriceService {
	return &PriceService{priceRep: priceRep}
}

func (p *PriceService) ReadPrices(ctx context.Context, company string) (float64, error) {
	actions, err := p.priceRep.ReadPrices(ctx, company)
	if err != nil {
		return 0, fmt.Errorf("Price-ReadPrices: error: %w", err)
	}
	return actions, nil
}
