// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"time"

	"github.com/artnikel/PriceService/proto"
	"github.com/sirupsen/logrus"
)

// PriceInterface is interface with method for reading prices
type PriceInterface interface {
	ReadPrices(ctx context.Context, company string) (float64, error)
}

// PriceHandler contains PriceInterface interface and UnimplementedPriceServiceServer
type PriceHandler struct {
	priceService PriceInterface
	proto.UnimplementedPriceServiceServer
}

// NewPriceHandler accepts PriceInterface object and returnes an object of type *PriceHandler
func NewPriceHandler(priceService PriceInterface) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
	}
}

// ReadPrices is a method of PriceHandler that calls method of Service
func (s *PriceHandler) ReadPrices(req *proto.ReadPricesRequest, stream proto.PriceService_ReadPricesServer) error {
	for {
		price, err := s.priceService.ReadPrices(stream.Context(), req.Company)
		if err != nil {
			logrus.Errorf("PriceHandler-ReadPrices: error getting price: %v", err)
			return err
		}

		action := &proto.Actions{
			Company: req.Company,
			Price:   price,
		}

		if err := stream.Send(&proto.ReadPricesResponse{Actions: []*proto.Actions{action}}); err != nil {
			logrus.Errorf("PriceHandler-ReadPrices: error sending response: %v", err)
			return err
		}
		//nolint gonmd
		time.Sleep(time.Second / 2)
	}
}
