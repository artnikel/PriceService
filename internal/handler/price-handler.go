// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/artnikel/PriceService/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PriceInterface is interface with method for reading prices
type PriceInterface interface {
	ReadPrices(ctx context.Context) (actions []*model.Action, err error)
	AddSubscriber(subscriberID uuid.UUID, selectedActions []string) error
	DeleteSubscriber(subscriberID uuid.UUID) error
	SendToSubscriber(ctx context.Context, subscriberID uuid.UUID) ([]*proto.Actions, error)
	SendToAllSubscribedChans(ctx context.Context)
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

// Subscribe takes message from redis stream through PriceServiceService and sends it to grpc stream.
func (h *PriceHandler) Subscribe(req *proto.SubscribeRequest, stream proto.PriceService_SubscribeServer) error {
	subscriberID, err := uuid.Parse(req.Uuid)
	if err != nil {
		logrus.Errorf("PriceHandler-Subscribe: error in method uuid.Parse: %v", err)
		return err
	}

	err = h.priceService.AddSubscriber(subscriberID, req.SelectedActions)
	if err != nil {
		logrus.Errorf("PriceHandler-Subscribe-AddSubscriber: error:%v", err)
		return err
	}

	for {
		protoActions, errSend := h.priceService.SendToSubscriber(stream.Context(), subscriberID)

		if errSend != nil {
			logrus.Errorf("PriceHandler-Subscribe-SendToSubscriber: error:%v", err)

			errDelete := h.priceService.DeleteSubscriber(subscriberID)
			if errDelete != nil {
				logrus.Errorf("PriceHandler-Subscribe-DeleteSubscriber: error:%v", err)
				return errDelete
			}

			return errSend
		}
		err := stream.Send(&proto.SubscribeResponse{Actions: protoActions})
		if err != nil {
			logrus.Errorf("PriceHandler-Subscribe: error in method stream.Send: %v", err)
			return err
		}
	}
}
