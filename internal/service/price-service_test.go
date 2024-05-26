package service

import (
	"context"
	"testing"
	"time"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/artnikel/PriceService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testSubscriberID   = uuid.New()
	testSelectedShares = []string{
		"Porsche",
		"Facebook",
	}
	testShares = []*model.Share{
		{Company: "Porsche", Price: decimal.NewFromFloat(457.23)},
		{Company: "Facebook", Price: decimal.NewFromFloat(842.45)},
		{Company: "IKEA", Price: decimal.NewFromFloat(1743.88)},
	}
)

func TestAddSubscriber(t *testing.T) {
	rep := new(mocks.PriceRepository)
	srv := NewPriceService(rep)
	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
	require.NoError(t, err)
	selectedActions, ok := srv.manager.Subscribers[testSubscriberID]
	require.True(t, ok)
	require.Equal(t, len(selectedActions), len(testSelectedShares))
	rep.AssertExpectations(t)
}

func TestDeleteSubscriber(t *testing.T) {
	rep := new(mocks.PriceRepository)
	srv := NewPriceService(rep)
	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
	require.NoError(t, err)
	err = srv.DeleteSubscriber(testSubscriberID)
	require.NoError(t, err)
	err = srv.AddSubscriber(testSubscriberID, testSelectedShares)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestSubscribeAll(t *testing.T) {
	rep := new(mocks.PriceRepository)
	rep.On("ReadPrices", mock.Anything).
		Return(testShares, nil)

	srv := NewPriceService(rep)
	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		srv.SubscribeAll(ctx)
	}()
	time.Sleep(time.Second)
	cancel()
	for _, selectedShare := range testSelectedShares {
		expectedPrice := decimal.Zero
		for _, share := range testShares {
			if share.Company == selectedShare {
				expectedPrice = share.Price
				break
			}
		}
		receivedShare := <-srv.manager.SubscribersShare[testSubscriberID]
		require.Equal(t, expectedPrice, receivedShare.Price)
	}
	rep.AssertExpectations(t)
}

func TestSendToSubscriber(t *testing.T) {
	rep := new(mocks.PriceRepository)
	rep.On("ReadPrices", mock.Anything).
		Return(testShares, nil)

	srv := NewPriceService(rep)
	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
	require.NoError(t, err)
	ctxCanceled, cancel := context.WithCancel(context.Background())
	cancel()
	shares, err := srv.SendToSubscriber(ctxCanceled, testSubscriberID)
	require.Error(t, err)
	require.Nil(t, shares)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go srv.SubscribeAll(ctx)
	shares, err = srv.SendToSubscriber(ctx, testSubscriberID)
	require.NoError(t, err)
	require.NotNil(t, shares)
	require.Len(t, shares, len(testSelectedShares))
	for i, selectedShare := range testSelectedShares {
		expectedPrice := decimal.Zero
		for _, share := range testShares {
			if share.Company == selectedShare {
				expectedPrice = share.Price
				break
			}
		}
		require.Equal(t, expectedPrice.InexactFloat64(), shares[i].Price)
	}
	cancel()
	err = srv.DeleteSubscriber(testSubscriberID)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

