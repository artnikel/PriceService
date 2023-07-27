package handler

import (
	"context"
	"io"

	"github.com/artnikel/TradingSystem/proto"
	"github.com/sirupsen/logrus"
)

type PriceInterface interface {
	ReadPrices(ctx context.Context, company string) (float64, error)
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

func (s *PriceHandler) ReadPrices(stream proto.PriceService_ReadPricesServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			logrus.Errorf("PriceHandler-ReadPrices: error receiving request: %v", err)
			return err
		}

		price, err := s.priceService.ReadPrices(stream.Context(), req.Company)
		if err != nil {
			logrus.Errorf("PriceHandler-ReadPrices: error getting price: %v", err)
			return err
		}

		action := &proto.Actions{
			Company: req.Company,
			Price:   price,
		}

		err = stream.Send(&proto.ReadPricesResponse{Actions: []*proto.Actions{action}})
		if err != nil {
			logrus.Errorf("PriceHandler-ReadPrices: error sending response: %v", err)
			return err
		}
	}
}
