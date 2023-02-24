// Package grpcserv описывает все методы и структуру grpc сервера.
package grpcserv

import (
	"context"
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"github.com/gtgaleevtimur/reduction-url-service/proto"
)

// Shortener реализует методы grpc сервера.
type Shortener struct {
	proto.UnimplementedShortenerServer

	conf       *config.Config
	repository repository.Storager
}

// New - конструктор grpc Shortener.
func New(s repository.Storager, conf *config.Config) *Shortener {
	return &Shortener{
		proto.UnimplementedShortenerServer{},
		conf,
		s,
	}
}

// AddByText - сокращает полный URL, добавляя в БД.
func (s *Shortener) AddByText(ctx context.Context, r *proto.StringForm) (*proto.CommonResponse, error) {
	var response proto.CommonResponse
	var token string
	url := r.GetLink()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}
		res, err := s.repository.InsertURL(ctx, url, token)
		if err != nil {
			return &response, status.Error(codes.Internal, "method AddByText not realise")
		}
		exShortURL := s.conf.ExpShortURL(res)
		response.Link = exShortURL
	}
	return &response, nil
}

// GetByHashURL - возвращает оригинальный URL по хэшу.
func (s *Shortener) GetByHashURL(ctx context.Context, r *proto.StringForm) (*proto.CommonResponse, error) {
	var response proto.CommonResponse
	hash := r.GetLink()
	res, err := s.repository.GetFullURL(ctx, hash)
	if err != nil {
		return &response, status.Error(codes.Internal, "method GetByHashURL not realise")
	}
	if !strings.HasPrefix(res, config.HTTP) {
		res = config.HTTP + strings.TrimPrefix(res, "//")
	}
	response.Link = res
	return &response, nil
}

// Ping - возвращает 200 в случае успешного Ping, возвращает 500 , если БД не доступна.
func (s *Shortener) Ping(ctx context.Context, no *proto.NoParam) (*proto.IntForm, error) {
	var response proto.IntForm
	err := s.repository.Ping(ctx)
	if err != nil {
		response.Value = http.StatusInternalServerError
		return &response, status.Error(codes.Internal, "method Ping not realise")
	}
	response.Value = http.StatusOK
	return &response, nil
}

// Stats - возвращает количество пользователей и сохраненных url в БД.
func (s *Shortener) Stats(ctx context.Context, no *proto.NoParam) (*proto.StatsResponse, error) {
	var response proto.StatsResponse
	urls, err := s.repository.GetCountURL(ctx)
	if err != nil {
		return &response, status.Errorf(codes.Internal, "method Stats not realise count urls")
	}
	users, err := s.repository.GetCountUsers(ctx)
	if err != nil {
		return &response, status.Errorf(codes.Internal, "method Stats not realise count users")
	}
	response.Urls = int32(urls)
	response.Users = int32(users)
	return &response, nil
}

// Delete - удаляет все url пользователя, возвращает http.StatusAccepted.
func (s *Shortener) Delete(ctx context.Context, r *proto.DeleteRequest) (*proto.IntForm, error) {
	var response proto.IntForm
	var token string
	ids := r.GetId()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}
		go s.repository.Delete(ctx, ids, token)
		response.Value = http.StatusAccepted
	}
	return &response, nil
}

// GetUserURLs - возвращает сохраненные пользователем url.
func (s *Shortener) GetUserURLs(ctx context.Context, r *proto.NoParam) (*proto.GetUserURLsResponse, error) {
	var response proto.GetUserURLsResponse
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}
		res, err := s.repository.GetAllUserURLs(ctx, token)
		if err != nil {
			return nil, status.Errorf(codes.Unimplemented, "method GetUserURLs not implemented")
		}
		result := make([]*proto.Links, len(res))
		for i, v := range res {
			result[i] = &proto.Links{
				Short: s.conf.ExpShortURL(v.Short),
				Full:  v.Full,
			}
		}
		response.Links = result
	}
	return &response, nil
}

// PostJSON - метод принимающий в теле json с полным url,возвращает сокращенный URL в json.
func (s *Shortener) PostJSON(ctx context.Context, r *proto.PostJSONRespReq) (*proto.PostJSONRespReq, error) {
	var response proto.PostJSONRespReq
	var token string
	body := r.GetJson()
	var full repository.FullURL
	err := json.Unmarshal(body, &full)
	if err != nil {
		return &response, status.Errorf(codes.Internal, "method PostJSON->json not realise")
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
		if len(token) == 0 {
			return &response, status.Error(codes.Unauthenticated, "missing token")
		}
		var sURL repository.ShortURL
		sURL.Short, err = s.repository.InsertURL(ctx, full.Full, token)
		if err != nil && !errors.Is(err, repository.ErrConflictInsert) {
			return &response, status.Errorf(codes.Internal, "method PostJSON->InsertURL not realise")
		}
		sURL.Short = s.conf.ExpShortURL(sURL.Short)
		respBody, err := json.Marshal(sURL)
		if err != nil {
			return &response, status.Errorf(codes.Internal, "method PostJSON->json.Marshal not realise")
		}
		response.Json = respBody
	}

	return &response, status.Errorf(codes.Unimplemented, "method PostJSON not implemented")
}

// PostBatch - метод реализующий загрузку массива с url.
func (s *Shortener) PostBatch(ctx context.Context, r *proto.PostBatchRequest) (*proto.PostBatchResponse, error) {
	var response proto.PostBatchResponse
	var token string
	batch := r.GetLinks()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
		if len(token) == 0 {
			return &response, status.Error(codes.Unauthenticated, "missing token")
		}
		result := make([]*proto.ButchLinks, 0)
		for _, v := range batch {
			short, err := s.repository.InsertURL(ctx, v.Link, token)
			if err != nil && !errors.Is(err, repository.ErrConflictInsert) {
				return &response, status.Errorf(codes.Internal, "method PostBatch not realise Insert")
			}
			result = append(result, &proto.ButchLinks{
				Link: short,
				Id:   v.Id,
			})
		}
		response.Links = result
	}
	return &response, nil
}

// MyUnaryInterceptor - перехватчик-аутентификатор, проверяет заголовок метаданных userid,
// если он пуст или проверка токен не удалась, то он выдает новый userid.
func MyUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
	}
	tokenBuf := make([]byte, aes.BlockSize)
	key := []byte("HdUeLk85Gp0i7pLh")
	nonce := []byte("userid")
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("aesBlock inizialise error:%s", err.Error())
	}
	// Проверка токена на подлинность.
	if len(token) != 0 {
		userIDByte, err := hex.DecodeString(token)
		if err != nil {
			log.Printf("Cookie decoding: %v\n", err)
		}
		aesBlock.Decrypt(tokenBuf, userIDByte)
		if string(tokenBuf[len(tokenBuf)-len(nonce):]) == string(nonce) {
			return handler(ctx, req)
		}
	}
	// Если токена нет или значение-пустая строка, то создаем новый токен и пишем в метаданные.
	userID, err := generateRandom(10)
	if err != nil {
		log.Printf("UserID generate: %v\n", err)
	}
	aesBlock.Encrypt(tokenBuf, append(userID, nonce...))
	md := metadata.New(map[string]string{"token": hex.EncodeToString(tokenBuf)})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	return handler(ctx, req)
}

// generateRandom - генератор случайных байт длинной size.
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
