package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mamoru777/chatservice2/internal/mylogger"
	"github.com/mamoru777/chatservice2/internal/repositories/chatrepository"
	"github.com/mamoru777/chatservice2/internal/repositories/chatusrrepository"
	"github.com/mamoru777/chatservice2/internal/repositories/messagerepository"
	"github.com/mamoru777/chatservice2/internal/repositories/usrrepository"

	gatewayapi "github.com/mamoru777/chatservice2/pkg/gateway-api"
	socketserviceapi "gitlab.com/mediasoft-internship/internship/mamoru777/socketservice/pkg/gateway-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ChatService struct {
	gatewayapi.UnimplementedChatServiceServer
	usrRep     usrrepository.IUsrRepository
	chatRep    chatrepository.IChatRepository
	chatUsrRep chatusrrepository.IChatUsrRepository
	messageRep messagerepository.IMessageRepository
}

func New(usrRep usrrepository.IUsrRepository, chatRep chatrepository.IChatRepository, chatUsrRep chatusrrepository.IChatUsrRepository, messageRep messagerepository.IMessageRepository) *ChatService {
	return &ChatService{
		usrRep:     usrRep,
		chatRep:    chatRep,
		chatUsrRep: chatUsrRep,
		messageRep: messageRep,
	}
}

func (cs *ChatService) CreateChats(ctx context.Context, request *gatewayapi.CreateChatsRequest) (*gatewayapi.CreateChatsResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		mylogger.Logger.Error(err)
		return &gatewayapi.CreateChatsResponse{}, err
	}

	conn, err := grpc.Dial("localhost:13993", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось подключиться к grpc серверу web socket", err)
		err = errors.New("Не удалось подключиться к grpc серверу web socket")
		return &gatewayapi.CreateChatsResponse{}, err
	}
	socketService := socketserviceapi.NewSocketServiceClient(conn)
	userUuid, err := StringToUuuid(request.Userid)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось конвертировать userId в uuid", err)
		err := errors.New("Не удалось конвертировать userId в uuid")
		return &gatewayapi.CreateChatsResponse{}, err
	}
	users, err := cs.usrRep.GetList(ctx)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось получить список всех пользователей", err)
		err := errors.New("Не удалось получить список всех пользователей")
		return &gatewayapi.CreateChatsResponse{}, err
	}
	u := &usrrepository.Usr{
		Id: userUuid,
	}
	err = cs.usrRep.Create(ctx, u)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать запись о пользователе в бд", err)
		err := errors.New("Не удалось создать запись о пользователе в бд")
		return &gatewayapi.CreateChatsResponse{}, err
	}
	for _, u := range users {
		id, err := cs.chatRep.Create(ctx)
		if err != nil {
			mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать запись о переписке в бд", err)
			err := errors.New("Не удалось создать запись о переписке в бд")
			return &gatewayapi.CreateChatsResponse{}, err
		}
		chatUsr := &chatusrrepository.ChatUsr{
			ChatId: id,
			UsrId:  userUuid,
		}
		err = cs.chatUsrRep.Create(ctx, chatUsr)
		if err != nil {
			mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать запись о переписке и пользователе в бд", err)
			err := errors.New("Не удалось создать запись о переписке и пользователе в бд")
			return &gatewayapi.CreateChatsResponse{}, err
		}
		chatUsr = &chatusrrepository.ChatUsr{
			ChatId: id,
			UsrId:  u.Id,
		}
		err = cs.chatUsrRep.Create(ctx, chatUsr)
		if err != nil {
			mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать запись о переписке и пользователе в бд", err)
			err := errors.New("Не удалось создать запись о переписке и пользователе в бд")
			return &gatewayapi.CreateChatsResponse{}, err
		}
		err = CreateHub(ctx, socketService, &socketserviceapi.CreateHubRequest{Chatid: id.String()})
		if err != nil {
			mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
			return &gatewayapi.CreateChatsResponse{}, err
		}
	}
	return &gatewayapi.CreateChatsResponse{}, nil
}

func (cs *ChatService) SendMessage(ctx context.Context, request *gatewayapi.SendMessageRequest) (*gatewayapi.SendMessageResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		mylogger.Logger.Error(err)
		return &gatewayapi.SendMessageResponse{}, err
	}

	if request.Chatid == "" {
		err := errors.New("Поле chatId не может быть пустым")
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.SendMessageResponse{}, err
	}
	if request.Text == "" {
		err := errors.New("Поле Text не может быть пустым")
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.SendMessageResponse{}, err
	}
	chatUuid, err := StringToUuuid(request.Chatid)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось конвертировать chatId в uuid", err)
		err := errors.New("Не удалось конвертировать chatId в uuid")
		return &gatewayapi.SendMessageResponse{}, err
	}
	userId := ctx.Value("user_id")
	userIdString, ok := userId.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx user id в string")
		err := errors.New("Не удалось преобразовать ctx user id в string")
		return &gatewayapi.SendMessageResponse{}, err
	}
	mylogger.Logger.Println("Запрос № ", xRequestIdString, " ", userIdString)
	userUuid, err := StringToUuuid(userIdString)
	if err != nil {
		mylogger.Logger.Println("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать строку в uuid", err)
		err := errors.New("Не удалось преобразовать строку в uuid")
		return &gatewayapi.SendMessageResponse{}, err
	}

	err = cs.messageRep.Create(ctx, &messagerepository.Message{
		ChatID: chatUuid,
		UsrID:  userUuid,
		Text:   request.Text,
		Data:   time.Now(),
	})
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось сохранить сообщение в бд", err)
		err := errors.New("Не удалось сохранить сообщение в бд")
		return &gatewayapi.SendMessageResponse{}, err
	}
	accessToken := ctx.Value("access_token")
	refreshToken := ctx.Value("refresh_token")
	accessTokenString, ok := accessToken.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx accessToken в string")
		err := errors.New("Не удалось преобразовать ctx accessToken в string")
		return &gatewayapi.SendMessageResponse{}, err
	}
	refreshTokenString, ok := refreshToken.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx refreshToken в string")
		err := errors.New("Не удалось преобразовать ctx refreshToken в string")
		return &gatewayapi.SendMessageResponse{}, err
	}
	u := "ws://localhost:13992/ws?id=" + request.Chatid
	mylogger.Logger.Println("Запрос № ", xRequestIdString, " ", "подключение к %s", u)
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "dial:", err)
	}
	defer c.Close()
	message := []byte(request.Text)
	err = c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось отправить сообщение", err)
		err := errors.New("Не удалось отправить сообщение")
		return &gatewayapi.SendMessageResponse{}, err
	}
	return &gatewayapi.SendMessageResponse{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}

func (cs *ChatService) GetChat(ctx context.Context, request *gatewayapi.GetChatRequest) (*gatewayapi.GetChatResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		mylogger.Logger.Error(err)
		return &gatewayapi.GetChatResponse{}, err
	}

	if request.Frinedid == "" {
		err := errors.New("Поле friendnId не может быть пустым")
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.GetChatResponse{}, err
	}
	friendUuid, err := StringToUuuid(request.Frinedid)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось конвертировать friendId в uuid", err)
		err := errors.New("Не удалось конвертировать friendId в uuid")
		return &gatewayapi.GetChatResponse{}, err
	}
	userId := ctx.Value("user_id")
	userIdString, ok := userId.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx user id в string")
		err := errors.New("Не удалось преобразовать ctx user id в string")
		return &gatewayapi.GetChatResponse{}, err
	}
	mylogger.Logger.Println("Запрос № ", xRequestIdString, " ", userIdString)
	userUuid, err := StringToUuuid(userIdString)
	if err != nil {
		mylogger.Logger.Println("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать строку в uuid", err)
		err := errors.New("Не удалось преобразовать строку в uuid")

		return &gatewayapi.GetChatResponse{}, err
	}
	chatId, err := cs.chatUsrRep.Get(ctx, userUuid, friendUuid)
	if err != nil {
		mylogger.Logger.Println("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать строку в uuid", err)
		err := errors.New("Не получить id переписки из бд")
		return &gatewayapi.GetChatResponse{}, err
	}
	accessToken := ctx.Value("access_token")
	refreshToken := ctx.Value("refresh_token")
	accessTokenString, ok := accessToken.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx accessToken в string")
		err := errors.New("Не удалось преобразовать ctx accessToken в string")
		return &gatewayapi.GetChatResponse{}, err
	}
	refreshTokenString, ok := refreshToken.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx refreshToken в string")
		err := errors.New("Не удалось преобразовать ctx refreshToken в string")
		return &gatewayapi.GetChatResponse{}, err
	}
	return &gatewayapi.GetChatResponse{Chatid: chatId.String(), AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}

func (cs *ChatService) GetMessages(ctx context.Context, request *gatewayapi.GetMessagesRequest) (*gatewayapi.GetMessagesResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		mylogger.Logger.Error(err)
		return &gatewayapi.GetMessagesResponse{}, err
	}

	if request.Chatid == "" {
		err := errors.New("Поле chatId не может быть пустым")
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.GetMessagesResponse{}, err
	}
	chatUuid, err := StringToUuuid(request.Chatid)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось конвертировать chatId в uuid", err)
		err := errors.New("Не удалось конвертировать chatId в uuid")
		return &gatewayapi.GetMessagesResponse{}, err
	}
	messages, err := cs.messageRep.GetList(ctx, chatUuid)
	if err != nil {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось получить список сообщений", err)
		err := errors.New("Не удалось получить список сообщений")
		return &gatewayapi.GetMessagesResponse{}, err
	}
	accessToken := ctx.Value("access_token")
	refreshToken := ctx.Value("refresh_token")
	accessTokenString, ok := accessToken.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx accessToken в string")
		err := errors.New("Не удалось преобразовать ctx accessToken в string")
		return &gatewayapi.GetMessagesResponse{}, err
	}
	refreshTokenString, ok := refreshToken.(string)
	if !ok {
		mylogger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось преобразовать ctx refreshToken в string")
		err := errors.New("Не удалось преобразовать ctx refreshToken в string")
		return &gatewayapi.GetMessagesResponse{}, err
	}
	var result []*gatewayapi.Message
	for _, m := range messages {
		result = append(result, &gatewayapi.Message{
			UsrId: m.UsrID.String(),
			Text:  m.Text,
			Data:  timestamppb.New(m.Data),
		})
	}
	return &gatewayapi.GetMessagesResponse{Result: result, AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}

func (cs *ChatService) GetAllChats(ctx context.Context, request *gatewayapi.GetAllChatsRequest) (*gatewayapi.GetAllChatsResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		mylogger.Logger.Error(err)
		return &gatewayapi.GetAllChatsResponse{}, err
	}
	_ = xRequestIdString
	chats, err := cs.chatRep.GetList(ctx)
	if err != nil {
		mylogger.Logger.Error("Не удалось получить список переписок", err)
		err := errors.New("Не удалось получить список переписок")
		return &gatewayapi.GetAllChatsResponse{}, err
	}
	var result []string
	for _, c := range chats {
		result = append(result, c.Id.String())
	}
	return &gatewayapi.GetAllChatsResponse{Result: result}, nil
}

func StringToUuuid(value string) (uuid.UUID, error) {
	emptyUUID := uuid.UUID{}
	uuid, err := uuid.Parse(value)
	if err != nil {
		return emptyUUID, err
	}
	return uuid, err
}

func CreateHub(ctx context.Context, client socketserviceapi.SocketServiceClient, req *socketserviceapi.CreateHubRequest) error {
	_, err := client.CreateHub(ctx, req)
	if err != nil {
		mylogger.Logger.Error("Не удалось вызвать удаленную функцию создания hub, либо функция возвратила ошибку")
		return err
	}
	return nil
}
