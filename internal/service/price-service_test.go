package service

import (
	"context"
	"testing"
	"time"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/artnikel/PriceService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testSubscriberID    = uuid.New()
	testSelectedActions = []string{
		"Porsche",
		"Facebook",
	}
	testActions = []*model.Action{
		{Company: "Porsche", Price: 457.23},
		{Company: "Facebook", Price: 842.45},
		{Company: "IKEA", Price: 1743.88},
	}
)

func TestAddSubscriber(t *testing.T) {
	rep := new(mocks.PriceRepository)
	srv := NewPriceService(rep)

	err := srv.AddSubscriber(testSubscriberID, testSelectedActions)
	require.NoError(t, err)

	selectedSharesInterface, ok := srv.manager.Subscribers.Load(testSubscriberID)
	require.True(t, ok)

	selectedShares, ok := selectedSharesInterface.([]string)
	require.True(t, ok)
	require.Equal(t, len(selectedShares), len(testSelectedActions))
}

func TestDeleteSubscriber(t *testing.T) {
	rep := new(mocks.PriceRepository)
	srv := NewPriceService(rep)

	err := srv.AddSubscriber(testSubscriberID, testSelectedActions)
	require.NoError(t, err)

	err = srv.DeleteSubscriber(testSubscriberID)
	require.NoError(t, err)
}

func TestReadPrices(t *testing.T) {
	rep := new(mocks.PriceRepository)
	rep.On("ReadPrices", mock.Anything).
		Return(testActions, nil).
		Once()
	srv := NewPriceService(rep)
	actions, err := srv.ReadPrices(context.Background())
	require.NoError(t, err)
	require.Equal(t, len(actions), len(testActions))
	rep.AssertExpectations(t)
}

func TestSendToAllSubscribedChans(t *testing.T) {
	rep := new(mocks.PriceRepository)
	rep.On("ReadPrices", mock.Anything).
		Return(testActions, nil)

	srv := NewPriceService(rep)

	err := srv.AddSubscriber(testSubscriberID, testSelectedActions)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	go srv.SendToAllSubscribedChans(ctx)

	defer func() {
		channel, ok := srv.manager.SubscribersActions.Load(testSubscriberID)
		if !ok {
			t.Errorf("Channel not found for testSubscriberID: %v", testSubscriberID)
			return
		}

		close(channel.(chan []*model.Action))
		cancel()
	}()

	channel, ok := srv.manager.SubscribersActions.Load(testSubscriberID)
	if !ok {
		t.Errorf("Channel not found for testSubID: %v", testSubscriberID)
		return
	}

	actions := <-channel.(chan []*model.Action)
	require.Equal(t, len(testSelectedActions), len(actions))

	rep.AssertExpectations(t)
}

func TestSendToSubscriber(t *testing.T) {
	rep := new(mocks.PriceRepository)
	srv := NewPriceService(rep)

	err := srv.AddSubscriber(testSubscriberID, testSelectedActions)
	require.NoError(t, err)

	channel := make(chan []*model.Action, len(testActions))
	srv.manager.SubscribersActions.Store(testSubscriberID, channel)
	defer close(channel)

	go func() {
		channel <- testActions
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	actions, err := srv.SendToSubscriber(ctx, testSubscriberID)
	require.NoError(t, err)

	require.Equal(t, len(testActions), len(actions))

	rep.AssertExpectations(t)
}
