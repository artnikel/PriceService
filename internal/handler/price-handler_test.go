package handler

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/artnikel/TradingSystem/internal/handler/mocks"
	"github.com/artnikel/TradingSystem/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testActions struct {
	Company string
	Price   float64
}

func server(ctx context.Context, s *mocks.PriceInterface) (psClient proto.PriceServiceClient, clientCloser func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	proto.RegisterPriceServiceServer(baseServer, NewPriceHandler(s))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := proto.NewPriceServiceClient(conn)

	return client, closer
}

func TestReadPrices(t *testing.T) {
	action := testActions{"Microsoft", 725.52}

	prc := new(mocks.PriceInterface)
	prc.On("ReadPrices", mock.Anything, mock.AnythingOfType("string")).Return(action.Price, nil)

	client, closer := server(context.Background(), prc)
	defer closer()

	type expectation struct {
		rpResponses []*proto.ReadPricesResponse
		err         error
	}

	respProtoActions := make([]*proto.Actions, 0)
	respProtoActions = append(respProtoActions,
		&proto.Actions{Company: action.Company, Price: action.Price})
	reqSelectedActions := "Microsoft"

	testReqResp := struct {
		rpRequest    *proto.ReadPricesRequest
		expectedResp expectation
	}{rpRequest: &proto.ReadPricesRequest{Company: reqSelectedActions},
		expectedResp: expectation{
			rpResponses: []*proto.ReadPricesResponse{
				{Actions: respProtoActions},
			},
			err: nil,
		},
	}
	out, err := client.ReadPrices(context.Background(), testReqResp.rpRequest)
	require.NoError(t, err)
	var outs []*proto.ReadPricesResponse
	o, err := out.Recv()
	require.NoError(t, err)
	outs = append(outs, o)
	require.Equal(t, len(testReqResp.expectedResp.rpResponses), len(outs))
	for i, actions := range outs[0].Actions {
		require.Equal(t, actions.Company, testReqResp.expectedResp.rpResponses[0].Actions[i].Company)
		require.Equal(t, actions.Price, testReqResp.expectedResp.rpResponses[0].Actions[i].Price)
	}
	prc.AssertExpectations(t)
}
