package grpcserv

import (
	"context"
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"github.com/gtgaleevtimur/reduction-url-service/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestNew(t *testing.T) {
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	serv := New(storage, conf)
	assert.IsType(t, &Shortener{}, serv)
}

func TestShortener_AddByText(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	links := []*proto.StringForm{
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
	}
	tokenBuf := make([]byte, aes.BlockSize)
	key := []byte("HdUeLk85Gp0i7pLh")
	nonce := []byte("userid")
	aesBlock, err := aes.NewCipher(key)
	require.NoError(t, err)
	b := make([]byte, 10)
	_, err = rand.Read(b)
	require.NoError(t, err)
	aesBlock.Encrypt(tokenBuf, append(b, nonce...))
	md := metadata.New(map[string]string{"token": hex.EncodeToString(tokenBuf)})
	fmt.Println("token:", hex.EncodeToString(tokenBuf))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	for _, v := range links {
		connData, err := client.AddByText(ctx, v)
		require.NoError(t, err)
		connResp := connData.Link
		require.NotNil(t, connResp)
	}
}

func TestShortener_GetByHashURL(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	links := []*proto.StringForm{
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
	}
	for i, v := range links {
		connAdd, err := client.AddByText(context.Background(), v)
		require.NoError(t, err)
		connResp := connAdd.Link
		connRespSlice := strings.Split(connResp, "/")
		connGet, err := client.GetByHashURL(context.Background(), &proto.StringForm{
			Link: connRespSlice[len(connRespSlice)-1],
		})
		require.NoError(t, err)
		require.Equal(t, links[i].Link, connGet.Link)
	}
}

func TestShortener_Ping(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	connPing, err := client.Ping(context.Background(), &proto.NoParam{})
	require.NoError(t, err)
	assert.Equal(t, int32(http.StatusOK), connPing.Value)
}

func TestShortener_Stats(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	connStats, err := client.Stats(context.Background(), &proto.NoParam{})
	require.NoError(t, err)
	assert.Equal(t, int32(0), connStats.Users)
	assert.Equal(t, int32(0), connStats.Urls)
}

func TestShortener_Delete(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	links := []*proto.StringForm{
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
	}
	tokenBuf := make([]byte, aes.BlockSize)
	key := []byte("HdUeLk85Gp0i7pLh")
	nonce := []byte("userid")
	aesBlock, err := aes.NewCipher(key)
	require.NoError(t, err)
	b := make([]byte, 10)
	_, err = rand.Read(b)
	require.NoError(t, err)
	aesBlock.Encrypt(tokenBuf, append(b, nonce...))
	md := metadata.New(map[string]string{"token": hex.EncodeToString(tokenBuf)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ids := make([]string, 0)
	for _, v := range links {
		connAdd, err := client.AddByText(ctx, v)
		require.NoError(t, err)
		require.NotNil(t, connAdd)
		connResp := connAdd.Link
		connRespSlice := strings.Split(connResp, "/")
		ids = append(ids, connRespSlice[len(connRespSlice)-1])
	}
	connDel, err := client.Delete(ctx, &proto.DeleteRequest{
		Id: ids,
	})
	require.NoError(t, err)
	require.Equal(t, int32(http.StatusAccepted), connDel.Value)
}

func TestShortener_GetUserURLs(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	tokenBuf := make([]byte, aes.BlockSize)
	key := []byte("HdUeLk85Gp0i7pLh")
	nonce := []byte("userid")
	aesBlock, err := aes.NewCipher(key)
	require.NoError(t, err)
	b := make([]byte, 10)
	_, err = rand.Read(b)
	require.NoError(t, err)
	aesBlock.Encrypt(tokenBuf, append(b, nonce...))
	md := metadata.New(map[string]string{"token": hex.EncodeToString(tokenBuf)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	links := []*proto.StringForm{
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
		{Link: `http://test.ru` + strconv.Itoa(rand.Intn(99))},
	}
	addSlice := make([]string, 0)
	for _, v := range links {
		connAdd, err := client.AddByText(ctx, v)
		require.NoError(t, err)
		require.NotNil(t, connAdd)
		addSlice = append(addSlice, connAdd.GetLink())
	}
	connGetAll, err := client.GetUserURLs(ctx, &proto.NoParam{})
	require.NoError(t, err)
	responseLinks := connGetAll.GetLinks()
	for j, k := range responseLinks {
		assert.Equal(t, k.Short, addSlice[j])
		assert.Equal(t, k.Full, links[j].Link)
	}
}

func TestShortener_PostJSON(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	full := &repository.FullURL{
		Full: `http://test.ru` + strconv.Itoa(rand.Intn(99)),
	}
	body, err := json.Marshal(full)
	require.NoError(t, err)
	connJSonPost, err := client.PostJSON(context.Background(), &proto.PostJSONRespReq{
		Json: body,
	})
	require.NoError(t, err)
	assert.NotNil(t, connJSonPost.GetJson())
}

func TestShortener_PostBatch(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	conf := config.NewConfig()
	storage, err := repository.NewDataSource(conf)
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(MyUnaryInterceptor))
	proto.RegisterShortenerServer(grpcServer, New(storage, conf))
	go grpcServer.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := proto.NewShortenerClient(conn)
	butch := make([]*proto.ButchLinks, 0)
	for i := 0; i < 3; i++ {
		but := &proto.ButchLinks{
			Id:   "test.ru",
			Link: `http://test.ru` + strconv.Itoa(rand.Intn(99)),
		}
		butch = append(butch, but)
	}
	connButch, err := client.PostBatch(context.Background(), &proto.PostBatchRequest{
		Links: butch,
	})
	require.NoError(t, err)
	assert.NotNil(t, connButch)
}
