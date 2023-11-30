package handler

import (
	"context"
	"net"
	"testing"

	"github.com/artnikel/PriceService/internal/handler/mocks"
	"github.com/artnikel/PriceService/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var testShares = []*proto.Shares{
	{Company: "Starbucks", Price: 178.46},
	{Company: "McDonalds", Price: 872.96},
}

func Server(ctx context.Context, s *mocks.PriceInterface) (psClient proto.PriceServiceClient, clientCloser func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	baseServer := grpc.NewServer()
	proto.RegisterPriceServiceServer(baseServer, NewPriceHandler(s))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			logrus.Printf("error serving server: %v", err)
		}
	}()
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("error connecting to server: %v", err)
	}
	closer := func() {
		err := lis.Close()
		if err != nil {
			logrus.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}
	client := proto.NewPriceServiceClient(conn)

	return client, closer
}

func TestSubscribe(t *testing.T) {
	s := new(mocks.PriceInterface)

	s.On("AddSubscriber", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("[]string")).
		Return(nil)
	s.On("DeleteSubscriber", mock.AnythingOfType("uuid.UUID")).
		Return(nil)
	s.On("SendToSubscriber", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(testShares, nil)
	client, closer := Server(context.Background(), s)

	type expectation struct {
		subResponses []*proto.SubscribeResponse
		err          error
	}
	respProtoShares := make([]*proto.Shares, 0)
	respProtoShares = append(respProtoShares,
		&proto.Shares{Company: "Starbucks", Price: 178.46},
		&proto.Shares{Company: "McDonalds", Price: 872.96})

	reqSelectedShares := make([]string, 0)
	reqSelectedShares = append(reqSelectedShares, "Starbucks", "McDonalds")

	testReqResp := struct {
		subReq       *proto.SubscribeRequest
		expectedResp expectation
	}{subReq: &proto.SubscribeRequest{Uuid: "747b6b85-9441-48cd-aee5-932f386ba381", SelectedShares: reqSelectedShares},
		expectedResp: expectation{
			subResponses: []*proto.SubscribeResponse{
				{Shares: respProtoShares},
				{Shares: respProtoShares},
			},
			err: nil,
		},
	}
	out, err := client.Subscribe(context.Background(), testReqResp.subReq)
	require.NoError(t, err)
	var outs []*proto.SubscribeResponse
	for i := 0; i < 2; i++ {
		o, err := out.Recv()
		require.NoError(t, err)

		outs = append(outs, o)
	}
	require.Equal(t, len(testReqResp.expectedResp.subResponses), len(outs))
	for i, share := range outs[0].Shares {
		require.Equal(t, share.Company, testReqResp.expectedResp.subResponses[0].Shares[i].Company)
		require.Equal(t, share.Price, testReqResp.expectedResp.subResponses[0].Shares[i].Price)
	}
	for i, share := range outs[1].Shares {
		require.Equal(t, share.Company, testReqResp.expectedResp.subResponses[0].Shares[i].Company)
		require.Equal(t, share.Price, testReqResp.expectedResp.subResponses[0].Shares[i].Price)
	}
	closer()
}
