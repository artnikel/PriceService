// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"
)

// PriceRepository is interface with method for reading prices
type PriceRepository interface {
	ReadPrices(ctx context.Context, company string) (float64, error)
}

// PriceService contains PriceRepository interface
type PriceService struct {
	priceRep PriceRepository
}

// NewPriceService accepts PriceRepository object and returnes an object of type *PriceService
func NewPriceService(priceRep PriceRepository) *PriceService {
	return &PriceService{priceRep: priceRep}
}

// ReadPrices is a method of GeneratorService that calls method of Repository
func (p *PriceService) ReadPrices(ctx context.Context, company string) (float64, error) {
	actions, err := p.priceRep.ReadPrices(ctx, company)
	if err != nil {
		return 0, fmt.Errorf("Price-ReadPrices: error: %w", err)
	}
	return actions, nil
}
