package realtime

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"gsm/pkg/cache"
	"gsm/pkg/errors"
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
	PushMessage(ctx context.Context, chatroomID, userID string, msgReq *StreamChatRoomMessageRequest) error
}

// NewStreamChatService init the connect service
func NewStreamChatService(redisClient cache.Client, streamClient stream.Client, idGenerator sonyflake.IDGenerator) StreamChatService {
	return &streamChatService{
		redisClient:  redisClient,
		streamClient: streamClient,
		idGenerator:  idGenerator,
	}
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for simplicity, adjust as needed
		},
	}
)

type changeRoomSub struct {
	roomID       string
	subscription stream.Subscription
	action       ChatRoomAction
}

func (impl *streamChatService) HandleWebSocketStreamConnect(userID string, r *http.Request, w http.ResponseWriter) {
	// Upgrade initial http request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading connection: %v", err)
	}
	defer ws.Close()

	ctx := r.Context()

	impl.handleWebSocketStreamConnectMessage(ctx, ws, userID)
}

// handleWebSocketStreamConnectMessage manage messages from the WebSocket client and route them to the appropriate handler.
func (impl *streamChatService) handleWebSocketStreamConnectMessage(ctx context.Context, conn *websocket.Conn, userID string) error {
	// handle receiving message from chat room
	changeRoomCh := make(chan changeRoomSub)
	go impl.receiveMessages(ctx, conn, changeRoomCh, userID)

	// handle message request
	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled")
			return nil
		default:
			// parse client message request.
			var msgReq StreamChatRoomMessageRequest
			err := conn.ReadJSON(&msgReq)
			log.Infof("successfully receive message from ws client %+v ", msgReq)
			if err != nil {
				// connection closed normally.
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived) {
					log.Info("ws connection is closed by client")
					return nil
				}
				return errors.Errorf("failed to read message from ws client: %v", err)
			}

			// check params
			if msgReq.RoomID == "" {
				return errors.NewError(errors.InvalidArgument, "empty room_id")
			}

			topicID := stream.ConstructChatRoomTopicID(msgReq.RoomID)
			subID := stream.ConstructChatRoomSubscriptionID(msgReq.RoomID, userID)
			// manage messages from the WebSocket client and route them to the appropriate handler.
			switch msgReq.Action {
			case ActionJoinChatRoom:
				if sub, err := impl.JoinChatRoom(ctx, userID, msgReq.RoomID, topicID, subID); err != nil {
					return errors.Errorf("failed to join chat room %s", msgReq.RoomID)
				} else {
					changeMsg := changeRoomSub{
						roomID:       msgReq.RoomID,
						subscription: sub,
						action:       ActionJoinChatRoom,
					}
					changeRoomCh <- changeMsg
				}
			case ActionLeaveChatRoom:
				if err := impl.LeaveChatRoom(ctx, topicID, subID); err != nil {
					return errors.Errorf("failed to leave chat room %s", msgReq.RoomID)
				}
			case ActionChatRoomMessage:
				if err := impl.PushMessage(ctx, msgReq.RoomID, userID, &msgReq); err != nil {
					return errors.Errorf("failed to push chat message to chat room %s", msgReq.RoomID)
				}
			}
		}
	}
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

func (impl *streamChatService) PushMessage(ctx context.Context, chatroomID, userID string, msgReq *StreamChatRoomMessageRequest) error {
	// send the new message to redis stream
	id, err := impl.idGenerator.NextID()
	if err != nil {
		return errors.Errorf("failed to generate id: %v", err)
	}
	sID := fmt.Sprintf("%d", id)
	// TODO: save message in db

	msgReq.SenderID = userID
	msgReq.ID = sID
	streamMsg := &stream.Message{
		ID:         sID,
		Attributes: msgReq.convertToKeyValuePairs(),
	}
	topicID := stream.ConstructChatRoomTopicID(chatroomID)
	if res := impl.streamClient.Topic(topicID).Send(ctx, streamMsg); res.Get(ctx) != nil {
		return errors.Errorf("failed to send message %+v to stream topic: %v", streamMsg, err)
	}
	return nil
}

func (impl *streamChatService) receiveMessages(ctx context.Context, conn *websocket.Conn, subCh chan changeRoomSub, userID string) error {
	h := func(ctx context.Context, msg *stream.Message) {
		resp := &StreamChatRoomMessageResponse{}
		if err := resp.convertFromKeyValuePairs(msg.Attributes); err != nil {
			log.Errorf("failed to convert stream data to StreamChatRoomResponse: %v", err)
			return
		}
		resp.SenderName = resp.SenderID // TODO: replace with real user name
		err := conn.WriteJSON(resp)
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
