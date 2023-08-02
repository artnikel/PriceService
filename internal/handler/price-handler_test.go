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

var testActions = []*proto.Actions{
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
		Return(testActions, nil)
	client, closer := Server(context.Background(), s)

	type expectation struct {
		subResponses []*proto.SubscribeResponse
		err          error
	}
	respProtoActions := make([]*proto.Actions, 0)
	respProtoActions = append(respProtoActions,
		&proto.Actions{Company: "Starbucks", Price: 178.46},
		&proto.Actions{Company: "McDonalds", Price: 872.96})

	reqSelectedActions := make([]string, 0)
	reqSelectedActions = append(reqSelectedActions, "Starbucks", "McDonalds")

	testReqResp := struct {
		subReq       *proto.SubscribeRequest
		expectedResp expectation
	}{subReq: &proto.SubscribeRequest{Uuid: "747b6b85-9441-48cd-aee5-932f386ba381", SelectedActions: reqSelectedActions},
		expectedResp: expectation{
			subResponses: []*proto.SubscribeResponse{
				{Actions: respProtoActions},
				{Actions: respProtoActions},
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

	for i, action := range outs[0].Actions {
		require.Equal(t, action.Company, testReqResp.expectedResp.subResponses[0].Actions[i].Company)
		require.Equal(t, action.Price, testReqResp.expectedResp.subResponses[0].Actions[i].Price)
	}
	for i, action := range outs[1].Actions {
		require.Equal(t, action.Company, testReqResp.expectedResp.subResponses[0].Actions[i].Company)
		require.Equal(t, action.Price, testReqResp.expectedResp.subResponses[0].Actions[i].Price)
	}

	closer()
}
