// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/artnikel/PriceService/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PriceInterface is interface with method for reading prices
type PriceInterface interface {
	ReadPrices(ctx context.Context) (shares []*model.Share, err error)
	AddSubscriber(subscriberID uuid.UUID, selectedShare []string) error
	DeleteSubscriber(subscriberID uuid.UUID) error
	SendToSubscriber(ctx context.Context, subscriberID uuid.UUID) ([]*proto.Shares, error)
	SubscribeAll(ctx context.Context)
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
		logrus.Errorf("subscribe %v", err)
		return err
	}
	err = h.priceService.AddSubscriber(subscriberID, req.SelectedShares)
	if err != nil {
		logrus.Errorf("addSubscriber %v", err)
		return err
	}
	for {
		protoShares, errSend := h.priceService.SendToSubscriber(stream.Context(), subscriberID)
		if errSend != nil {
			logrus.Infof("sendToSubscriber: subscriber disconnected")
			errDelete := h.priceService.DeleteSubscriber(subscriberID)
			if errDelete != nil {
				logrus.Errorf("deleteSubscriber %v", errDelete)
			}
			return fmt.Errorf("sendToSubscriber %w", errSend)
		}

		err := stream.Send(&proto.SubscribeResponse{Shares: protoShares})
		if err != nil {
			logrus.Infof("stream.Send: %v", err)
			errDelete := h.priceService.DeleteSubscriber(subscriberID)
			if errDelete != nil {
				logrus.Errorf("deleteSubscriber %v", errDelete)
			}
			return fmt.Errorf("deleteSubscriber")
		}
	}
}
