package handler

import (
	"testing"

	"github.com/artnikel/TradingSystem/internal/handler/mocks"
	"github.com/artnikel/TradingSystem/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testActions struct {
	Company string
	Price   float64
}

func TestReadPrices(t *testing.T) {
	action := testActions{"Amazon", 158.45}

	prc := new(mocks.PriceInterface)
	stream := new(proto.PriceService_ReadPricesServer)

	handl := NewPriceHandler(prc)

	prc.On("Recv").Return(&proto.ReadPricesRequest{Company: action.Company}, nil)
	prc.On("Send", mock.AnythingOfType("*proto.ReadPricesResponse")).Return(nil)

	err := handl.ReadPrices(*stream)
	require.NoError(t, err)

	prc.AssertExpectations(t)
}
