package e2e

import (
	"context"
	"errors"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/caarlos0/env/v6"

	"github.com/anoriar/shortener/internal/e2e/config"
	pb "github.com/anoriar/shortener/proto/generated/shortener/proto"
)

const testURL = "https://github.com/"

type ShortenerGRPCSuite struct {
	suite.Suite

	conf *config.TestConfig

	urlServiceClient   pb.URLServiceClient
	statsServiceClient pb.StatsServiceClient
}

func (suite *ShortenerGRPCSuite) SetupSuite() {
	conf := config.NewTestConfig()

	err := env.Parse(conf)

	suite.NoError(err)

	suite.conf = conf
	if conf.BaseURL == "" {
		suite.T().Skip()
	}

	suite.Require().NotEmpty(conf.BaseURL)
	suite.Require().NotEmpty(conf.ServerAddr)

	conn, err := grpc.Dial(conf.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	suite.urlServiceClient = pb.NewURLServiceClient(conn)
	suite.statsServiceClient = pb.NewStatsServiceClient(conn)

}

func (suite *ShortenerGRPCSuite) TestShortenerV2() {

	addResponse, err := suite.urlServiceClient.AddURL(context.Background(), &pb.AddURLRequest{Url: testURL})

	suite.Assert().NoError(err)
	suite.Assert().True(strings.HasPrefix(addResponse.Result, suite.conf.BaseURL))

	splittedURL := strings.Split(addResponse.Result, "/")
	key := splittedURL[len(splittedURL)-1]
	defer func() {
		_, err = suite.urlServiceClient.DeleteURLBatch(context.Background(),
			&pb.DeleteURLBatchRequest{
				ShortUrls: []string{key},
			},
		)
		suite.Assert().NoError(err)
	}()

	getResponse, err := suite.urlServiceClient.GetURL(context.Background(),
		&pb.GetURLRequest{ShortUrl: key},
	)
	suite.Assert().NoError(err)
	suite.Assert().Equal(testURL, getResponse.GetOriginalUrl())
}

const (
	originalURL1   = "https://practicum.yandex.ru"
	correlationID1 = "g0fsdf9fj"
	originalURL2   = "https://practicum2.yandex.ru"
	correlationID2 = "ngfdsf3"
	originalURL3   = "https://practicum3.yandex.ru"
	correlationID3 = "by4564trg"
)

func (suite *ShortenerGRPCSuite) Test_ShortenerAddURlBatch() {

	correlationIDSHortKeyMap := map[string]*pb.AddURLBatchRequest_Item{
		correlationID1: {
			CorrelationId: correlationID1,
			OriginalUrl:   originalURL1,
		},
		correlationID2: {
			CorrelationId: correlationID2,
			OriginalUrl:   originalURL2,
		},
		correlationID3: {
			CorrelationId: correlationID3,
			OriginalUrl:   originalURL3,
		},
	}
	batchRequestItems := make([]*pb.AddURLBatchRequest_Item, len(correlationIDSHortKeyMap))
	i := 0
	for _, mapItem := range correlationIDSHortKeyMap {
		batchRequestItems[i] = mapItem
		i++
	}

	addResponse, err := suite.urlServiceClient.AddURLBatch(context.Background(),
		&pb.AddURLBatchRequest{
			Items: batchRequestItems,
		},
	)
	suite.Assert().NoError(err)

	suite.Assert().True(len(addResponse.Items) > 0)

	correlationIDShortKeyMap := make(map[string]string)
	for _, item := range addResponse.Items {
		suite.Assert().True(strings.HasPrefix(item.ShortUrl, suite.conf.BaseURL))
		splittedURL := strings.Split(item.ShortUrl, "/")
		correlationIDShortKeyMap[item.CorrelationId] = splittedURL[len(splittedURL)-1]
	}
	defer func() {
		keysForDelete := make([]string, len(correlationIDShortKeyMap))
		i := 0
		for _, key := range correlationIDShortKeyMap {
			keysForDelete[i] = key
			i++
		}
		_, err = suite.urlServiceClient.DeleteURLBatch(
			context.Background(),
			&pb.DeleteURLBatchRequest{
				ShortUrls: keysForDelete,
			},
		)
		suite.Assert().NoError(err)
	}()

	for correlationID, shortKey := range correlationIDShortKeyMap {
		getResponse, err := suite.urlServiceClient.GetURL(context.Background(), &pb.GetURLRequest{ShortUrl: shortKey})
		suite.Assert().NoError(err)

		mapItem, existed := correlationIDSHortKeyMap[correlationID]
		suite.Assert().True(existed)
		suite.Assert().Equal(mapItem.GetOriginalUrl(), getResponse.GetOriginalUrl())
	}
}

func (suite *ShortenerGRPCSuite) Test_ShortenerGetUserURLs() {

	var expectedURLs []*pb.GetUserURLsResponse_URL
	var keysForDelete []string
	originalURLs := []string{originalURL1, originalURL2, originalURL3}

	token := ""
	var header metadata.MD
	ctx := context.Background()
	for _, url := range originalURLs {
		if token != "" {
			md := metadata.Pairs("token", token)
			ctx = metadata.NewOutgoingContext(context.Background(), md)
		}

		addResponse, err := suite.urlServiceClient.AddURL(ctx, &pb.AddURLRequest{Url: url}, grpc.Header(&header))
		suite.Assert().NoError(err)

		values := header.Get("token")
		if len(values) > 0 {
			token = values[0]
		}

		splittedURL := strings.Split(addResponse.Result, "/")
		keysForDelete = append(keysForDelete, splittedURL[len(splittedURL)-1])
		expectedURLs = append(expectedURLs, &pb.GetUserURLsResponse_URL{
			ShortUrl:    addResponse.Result,
			OriginalUrl: url,
		})
	}

	defer func() {
		_, err := suite.urlServiceClient.DeleteURLBatch(context.Background(), &pb.DeleteURLBatchRequest{ShortUrls: keysForDelete})
		suite.Assert().NoError(err)
	}()

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	response, err := suite.urlServiceClient.GetUserURLs(ctx, &pb.Empty{})
	suite.Assert().NoError(err)

	suite.Assert().True(len(expectedURLs) == len(response.Items))
	suite.Assert().Equal(expectedURLs, response.Items)
}

func (suite *ShortenerGRPCSuite) Test_ShortenerDeleteUserURLs() {

	var keysForDelete []string
	originalURLs := []string{originalURL1, originalURL2, originalURL3}

	token := ""
	var header metadata.MD
	ctx := context.Background()
	for _, url := range originalURLs {
		if token != "" {
			md := metadata.Pairs("token", token)
			ctx = metadata.NewOutgoingContext(context.Background(), md)
		}

		addResponse, err := suite.urlServiceClient.AddURL(ctx, &pb.AddURLRequest{Url: url}, grpc.Header(&header))
		suite.Assert().NoError(err)

		values := header.Get("token")
		if len(values) > 0 {
			token = values[0]
		}

		splittedURL := strings.Split(addResponse.Result, "/")
		keysForDelete = append(keysForDelete, splittedURL[len(splittedURL)-1])
	}

	defer func() {
		_, err := suite.urlServiceClient.DeleteURLBatch(context.Background(), &pb.DeleteURLBatchRequest{ShortUrls: keysForDelete})
		suite.Assert().NoError(err)
	}()

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	_, err := suite.urlServiceClient.DeleteUserURLs(ctx, &pb.DeleteUserURLsRequest{ShortUrls: keysForDelete})
	suite.Assert().NoError(err)

	//Операция асинхронная
	time.Sleep(1 * time.Second)

	for _, key := range keysForDelete {
		_, err := suite.urlServiceClient.GetURL(context.Background(), &pb.GetURLRequest{ShortUrl: key})
		if e, ok := status.FromError(err); ok {
			suite.Assert().Equal(codes.FailedPrecondition.String(), e.Code().String())
		} else {
			suite.Error(errors.New("not expected url existing"))
		}
	}
}
func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(ShortenerGRPCSuite))
}
