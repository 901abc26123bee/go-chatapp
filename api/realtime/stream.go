package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"gsm/pkg/realtime"
	"gsm/pkg/stream"
	"gsm/pkg/util/sonyflake"
)

// streamChatService defines the implementation of StreamChatService interface
type streamChatService struct {
	redisClient  cache.Client
	streamClient stream.Client
	idGenerator  sonyflake.IDGenerator
}

// StreamChatService defines the connect service interface
type StreamChatService interface {
	HandleWebSocketStreamConnect(userID string, r *http.Request, w http.ResponseWriter)
	JoinChatRoom(ctx context.Context, userID, chatRoomID, topicID, subID string) (stream.Subscription, error)
	LeaveChatRoom(ctx context.Context, topicID string, subID string) error
	PushMessage(ctx context.Context, userID string, msgReq *ChatRoomMessageRequest) error
}

// NewStreamChatService init the connect service
func NewStreamChatService(redisClient cache.Client, streamClient stream.Client, idGenerator sonyflake.IDGenerator) StreamChatService {
	return &streamChatService{
		redisClient:  redisClient,
		streamClient: streamClient,
		idGenerator:  idGenerator,
	}
}

type changeRoomSub struct {
	roomID       string
	subscription stream.Subscription
	action       ChatRoomAction
}

func (impl *streamChatService) HandleWebSocketStreamConnect(userID string, r *http.Request, w http.ResponseWriter) {
	// Upgrade initial http request to a WebSocket
	ws, err := realtime.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading connection: %v", err)
	}
	defer ws.Close()

	ctx := r.Context()

	// handle receiving message from chat room
	changeRoomCh := make(chan changeRoomSub)
	go impl.receiveMessages(ctx, ws, changeRoomCh, userID)

	// handle message request
	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled")
			return
		default:
			// parse client message request.
			var msgReq StreamWSMessageRequest
			err := ws.ReadJSON(&msgReq)
			log.Infof("successfully receive message from ws client %+v ", msgReq)
			if err != nil {
				// connection closed normally.
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived) {
					log.Info("ws connection is closed by client")
				}
				log.Errorf("failed to read message from ws client: %v", err)
				return
			}

			// handle messages from the WebSocket client and route them to the appropriate handler.
			if handleErr := impl.handleWebSocketStreamConnectMessage(ctx, msgReq, userID, changeRoomCh); handleErr != nil {
				err := ws.WriteJSON(&StreamWSMessageResponse{
					Type:    WSMessageACK,
					Status:  MessageStatusSuccess,
					ErrCode: errors.MessageHandleError,
					ErrMsg:  fmt.Sprintf("%s", handleErr),
				})
				if err != nil {
					log.Errorf("failed to writing ack error message to ws connection: %v", err)
					return
				}
			} else {
				err := ws.WriteJSON(&StreamWSMessageResponse{
					Type:   WSMessageACK,
					Status: MessageStatusSuccess,
				})
				if err != nil {
					log.Errorf("failed to writing ack success message to ws connection: %v", err)
					return
				}
			}
		}
	}
}

// handleWebSocketStreamConnectMessage manage messages from the WebSocket client and route them to the appropriate handler.
func (impl *streamChatService) handleWebSocketStreamConnectMessage(ctx context.Context, msgReq StreamWSMessageRequest, userID string, changeRoomCh chan changeRoomSub) error {
	switch msgReq.Type {
	case WSMessageTypeChatRoomAction:
		var actionPayload ChatRoomActionRequest
		if err := json.Unmarshal(msgReq.Payload, &actionPayload); err != nil {
			return err
		}
		if actionPayload.RoomID == "" {
			return errors.NewErrorf(errors.InvalidArgument, "empty chat room id")
		}

		topicID := stream.ConstructChatRoomTopicID(actionPayload.RoomID)
		subID := stream.ConstructChatRoomSubscriptionID(actionPayload.RoomID, userID)

		switch actionPayload.Action {
		case ActionJoinChatRoom:
			if sub, err := impl.JoinChatRoom(ctx, userID, actionPayload.RoomID, topicID, subID); err != nil {
				return errors.Errorf("failed to join chat room %s", actionPayload.RoomID)
			} else {
				changeMsg := changeRoomSub{
					roomID:       actionPayload.RoomID,
					subscription: sub,
					action:       ActionJoinChatRoom,
				}
				changeRoomCh <- changeMsg
			}
		case ActionLeaveChatRoom:
			if err := impl.LeaveChatRoom(ctx, topicID, subID); err != nil {
				return errors.Errorf("failed to leave chat room %s", actionPayload.RoomID)
			}
		}

	case WSMessageTypeChat:
		var chatMsgPayload ChatRoomMessageRequest
		if err := json.Unmarshal(msgReq.Payload, &chatMsgPayload); err != nil {
			return err
		}
		if chatMsgPayload.RoomID == "" {
			return errors.NewErrorf(errors.InvalidArgument, "empty chat room id")
		}

		if err := impl.PushMessage(ctx, userID, &chatMsgPayload); err != nil {
			return errors.Errorf("failed to push chat message to chat room %s", chatMsgPayload.RoomID)
		}
	}
	return nil
}

func (impl *streamChatService) JoinChatRoom(ctx context.Context, userID, chatRoomID, topicID, subID string) (stream.Subscription, error) {
	// TODO: remove when CreateChatRoom is implemented
	topic := impl.streamClient.Topic(topicID)
	exist, err := topic.Exists(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to check if stream topic exist: %v", err)
	}
	if exist {
		_, err := impl.streamClient.CreateTopic(ctx, topicID)
		if err != nil {
			return nil, errors.Errorf("failed to create stream topic: %v", err)
		}
	}

	// create subscription of topic channel if not exist
	subConfig := &stream.SubscriptionConfig{
		Topic:       impl.streamClient.Topic(topicID),
		TopicID:     topicID,
		ReadStartID: "0", // TODO
	}

	subscription := impl.streamClient.Subscription(subID, subConfig)
	if exist, err := subscription.Exists(ctx); err != nil {
		return nil, errors.Errorf("failed to check if subscription already exist in topic: %v", err)
	} else if !exist {
		if sub, err := impl.streamClient.CreateSubscription(ctx, subID, subConfig); err != nil {
			return nil, errors.Errorf("failed to create subscription to stream topic: %v", err)
		} else {
			subscription = sub
		}
	}
	log.Infof("successfully create subscription %s for topic %s", subID, topicID)

	// TODO: add online chatroom member in redis

	return subscription, nil
}

func (impl *streamChatService) LeaveChatRoom(ctx context.Context, topicID, subID string) error {
	if err := impl.streamClient.DeleteSubscription(ctx, topicID, subID); err != nil {
		return errors.Errorf("failed to delete subscription %s: %v", subID, err)
	}
	log.Infof("successfully delete subscription %s", subID)

	// TODO: delete online chatroom member in redis

	return nil
}

func (impl *streamChatService) PushMessage(ctx context.Context, userID string, msgReq *ChatRoomMessageRequest) error {
	// send the new message to redis stream
	id, err := impl.idGenerator.NextID()
	if err != nil {
		return errors.Errorf("failed to generate id: %v", err)
	}
	sID := fmt.Sprintf("%d", id)
	// TODO: save message in db

	msgReq.ID = sID
	streamMsg := &stream.Message{
		ID:         sID,
		Attributes: msgReq.convertToKeyValuePairs(userID),
	}
	topicID := stream.ConstructChatRoomTopicID(msgReq.RoomID)
	if res := impl.streamClient.Topic(topicID).Send(ctx, streamMsg); res.Get(ctx) != nil {
		return errors.Errorf("failed to send message %+v to stream topic: %v", streamMsg, err)
	}
	return nil
}

func (impl *streamChatService) receiveMessages(ctx context.Context, conn *websocket.Conn, subCh chan changeRoomSub, userID string) error {
	h := func(ctx context.Context, msg *stream.Message) {
		payload := &ChatRoomMessageResponse{}
		if err := payload.convertFromKeyValuePairs(msg.Attributes); err != nil {
			log.Errorf("failed to convert stream data to StreamChatRoomResponse: %v", err)
			return
		}
		content, err := json.Marshal(payload)
		if err != nil {
			log.Errorf("failed to marshal message payload: %v", err)
			return
		}
		resp := &StreamWSMessageResponse{
			Type:    WSMessageTypeChatStream,
			Payload: content,
		}
		err = conn.WriteJSON(resp)
		if err != nil {
			log.Errorf("failed to writing message to ws connection: %v", err)
			return
		}
		log.Infof("successfully sending message to ws client: %+v ", resp)
	}

	subscriptions := make(map[string]map[*stream.Subscription]chan struct{}) // [room_id]:[stream subscription]:[struct{}]

	var curSub stream.Subscription
	for {
		select {
		case <-ctx.Done():
			// handle connection closure, e.g., close subscriptions
			for roomID, sub := range subscriptions {
				for _, ch := range sub {
					// signal all subscription to stop
					ch <- struct{}{}
					topicID := stream.ConstructChatRoomTopicID(roomID)
					subID := stream.ConstructChatRoomSubscriptionID(roomID, userID)
					impl.LeaveChatRoom(context.Background(), topicID, subID)
				}
			}
			return nil
		case subMsg := <-subCh:
			switch subMsg.action {
			case ActionJoinChatRoom:
				curSub = subMsg.subscription
				stopCh := make(chan struct{})
				if curSub != nil {
					if err := curSub.Receive(ctx, h, stopCh); err != nil {
						return errors.Errorf("failed to receive message from stream subscription: %v", err)
					}
				}
			case ActionLeaveChatRoom:
				subID := stream.ConstructChatRoomSubscriptionID(subMsg.roomID, userID)
				if subID != curSub.GetSubID() {
					log.Errorf("message with action leave chat room doesn't match current chat room")
				} else {
					if subInfo, ok := subscriptions[subMsg.roomID]; ok {
						for _, ch := range subInfo {
							// signal subscription to stop
							ch <- struct{}{}
						}
					}
				}
			}
		}
	}
}
