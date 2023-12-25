package middleware

import (
	"context"
	authservice "github.com/mamoru777/authservice2"
	"github.com/mamoru777/chatservice2/internal/jwttokens/jwtservice"
	"github.com/mamoru777/chatservice2/internal/mylogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type AuthInterceptor struct {
}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (i *AuthInterceptor) JWTInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	mylogger.Logger.Println("Использована прослойка jwt", info.FullMethod)
	if info.FullMethod == "/api.ChatService/SendMessage" || info.FullMethod == "/api.ChatService/GetChat" || info.FullMethod == "/api.ChatService/GetMessages" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "метаданные не предоставлены")
		}

		accessToken := i.getTokenFromMetadata(md)
		if accessToken == "" {
			return nil, status.Errorf(codes.Unauthenticated, "пустой access token")
		}
		refreshToken := i.getRefreshTokenFromMetadata(md)
		if refreshToken == "" {
			return nil, status.Errorf(codes.Unauthenticated, "пустой refresh token")
		}
		claims, err := jwtservice.VerifyToken(accessToken)
		if err != nil {
			mylogger.Logger.Println("Access Токен неверный")
			refreshClaims, refreshErr := jwtservice.VerifyToken(refreshToken)
			if refreshErr != nil {
				return nil, status.Errorf(codes.Unauthenticated, "неверный refresh token")
			}
			conn, err := grpc.Dial("localhost:13999", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				mylogger.Logger.Fatal("error connect to grpc server err:", err)
			}
			userId, ok := refreshClaims["user_id"].(string)
			if !ok {
				mylogger.Logger.Println("user_id не является строкой")
			}
			authservice2 := authservice.NewUsrServiceClient(conn)
			accessTokenNew, _, err := RefreshAccessToken(ctx, authservice2, &authservice.UpdateAccessTokenRequest{RefreshToken: refreshToken, Userid: userId})
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "не удалось обновить access token")
			}
			mylogger.Logger.Println(refreshClaims["user_id"])
			mylogger.Logger.Println("Новый access токен - ", accessTokenNew)
			ctx = context.WithValue(ctx, "user_id", userId)
			ctx = context.WithValue(ctx, "access_token", accessTokenNew)
			ctx = context.WithValue(ctx, "refresh_token", refreshToken)
			return handler(ctx, req)
		} else {
			userId, ok := claims["user_id"].(string)
			if !ok {
				mylogger.Logger.Println("user_id не является строкой")
			}
			mylogger.Logger.Println(userId)
			ctx = context.WithValue(ctx, "user_id", userId)
			ctx = context.WithValue(ctx, "access_token", accessToken)
			ctx = context.WithValue(ctx, "refresh_token", refreshToken)
			mylogger.Logger.Println(claims["user_id"])
			mylogger.Logger.Println("Токен верный")
			mylogger.Logger.Println("Завершение jwt middleware")
			return handler(ctx, req)
		}
	}
	return handler(ctx, req)
}

func (i *AuthInterceptor) getTokenFromMetadata(md metadata.MD) string {
	authorization := md.Get("authorization")
	if len(authorization) > 0 {
		parts := strings.Split(authorization[0], " ")
		if len(parts) == 3 {
			return parts[1]
		}
	}
	return ""
}

func (i *AuthInterceptor) getRefreshTokenFromMetadata(md metadata.MD) string {
	authorization := md.Get("authorization")
	if len(authorization) > 0 {
		parts := strings.Split(authorization[0], " ")
		if len(parts) == 3 {
			return parts[2]
		}
	}
	return ""
}

func RefreshAccessToken(ctx context.Context, client authservice.UsrServiceClient, req *authservice.UpdateAccessTokenRequest) (string, string, error) {
	res, err := client.UpdateAccessToken(ctx, req)
	if err != nil {
		mylogger.Logger.Println("Не удалось вызвать удаленную функцию обновления токена, либо функция возвратила ошибку")
		return "", "", err
	}
	accessToken := res.AccessToken
	refreshToken := res.RefreshToken
	return accessToken, refreshToken, nil
}
