package service

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/artnikel/PriceService/internal/model"
// 	"github.com/artnikel/PriceService/internal/service/mocks"
// 	"github.com/google/uuid"
// 	"github.com/shopspring/decimal"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// var (
// 	testSubscriberID   = uuid.New()
// 	testSelectedShares = []string{
// 		"Porsche",
// 		"Facebook",
// 	}
// 	testShares = []*model.Share{
// 		{Company: "Porsche", Price: decimal.NewFromFloat(457.23)},
// 		{Company: "Facebook", Price: decimal.NewFromFloat(842.45)},
// 		{Company: "IKEA", Price: decimal.NewFromFloat(1743.88)},
// 	}
// )

// func TestAddSubscriber(t *testing.T) {
// 	rep := new(mocks.PriceRepository)
// 	srv := NewPriceService(rep)

// 	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
// 	require.NoError(t, err)

// 	selectedActions, ok := srv.manager.Subscribers[testSubscriberID]
// 	require.True(t, ok)
// 	require.Equal(t, len(selectedActions), len(testSelectedShares))
// 	rep.AssertExpectations(t)
// }

// func TestDeleteSubscriber(t *testing.T) {
// 	rep := new(mocks.PriceRepository)
// 	srv := NewPriceService(rep)

// 	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
// 	require.NoError(t, err)

// 	err = srv.DeleteSubscriber(testSubscriberID)
// 	require.NoError(t, err)

// 	err = srv.AddSubscriber(testSubscriberID, testSelectedShares)
// 	require.NoError(t, err)
// 	rep.AssertExpectations(t)
// }

// func TestSendToAllSubscribedChans(t *testing.T) {
// 	rep := new(mocks.PriceRepository)
// 	rep.On("ReadPrices", mock.Anything).
// 		Return(testShares, nil)

// 	srv := NewPriceService(rep)

// 	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
// 	require.NoError(t, err)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	srv.SendToAllSubscribedChans(ctx)

// 	close(srv.manager.SubscribersShare[testSubscriberID])
// 	shares := <-srv.manager.SubscribersShare[testSubscriberID]
// 	cancel()
// 	require.Equal(t, len(testSelectedShares), len(shares))

// 	rep.AssertExpectations(t)
// }

// func TestSendToSubscriber(t *testing.T) {
// 	rep := new(mocks.PriceRepository)
// 	srv := NewPriceService(rep)

// 	err := srv.AddSubscriber(testSubscriberID, testSelectedShares)
// 	require.NoError(t, err)

// 	srv.manager.SubscribersShare[testSubscriberID] <- testShares

// 	close(srv.manager.SubscribersShare[testSubscriberID])

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	actions, err := srv.SendToSubscriber(ctx, testSubscriberID)
// 	cancel()
// 	require.NoError(t, err)

// 	require.Equal(t, len(testShares), len(actions))

// 	rep.AssertExpectations(t)
// }

// func TestReadPrices(t *testing.T) {
// 	rep := new(mocks.PriceRepository)
// 	rep.On("ReadPrices", mock.Anything).
// 		Return(testShares, nil).
// 		Once()
// 	srv := NewPriceService(rep)
// 	actions, err := srv.ReadPrices(context.Background())
// 	require.NoError(t, err)
// 	require.Equal(t, len(actions), len(testShares))
// 	rep.AssertExpectations(t)
// }
