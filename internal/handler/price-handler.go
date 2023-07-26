package handler

import (
	"context"

	"github.com/artnikel/TradingSystem/proto"
)

type PriceInterface interface {
	ReadPrices(ctx context.Context) (map[string]float64, error)
	GetPriceByCompany(ctx context.Context, company string) (map[string]float64, error)
}

type PriceHandler struct {
	priceService PriceInterface
	proto.UnimplementedPriceServiceServer
}

func NewPriceHandler(priceService PriceInterface) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
	}
}

func (s *PriceHandler) ReadPrices(ctx context.Context, req *proto.ReadPricesRequest, stream proto.PriceService_ReadPricesServer) error {
	for {
		actions, err := s.priceService.ReadPrices(ctx)
	}
}
