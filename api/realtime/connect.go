package realtime

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
	log "github.com/sirupsen/logrus"

	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"gsm/pkg/stream"
	"gsm/pkg/util/sonyflake"
)

// connectService defines the implementation of ConnectService interface
type connectService struct {
	redisClient  cache.Client
	streamClient stream.Client
	idGenerator  sonyflake.IDGenerator
}

// ConnectService defines the connect service interface
type ConnectService interface {
	HandleWebSocketStreamConnect(w http.ResponseWriter, r *http.Request)
}

// NewConnectService init the connect service
func NewConnectService(redisClient cache.Client, streamClient stream.Client, idGenerator sonyflake.IDGenerator) ConnectService {
	return &connectService{
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

const wsTimeoutDuration = 300 * time.Second

// ChatRoomMessage Define a chat room message object
type ChatRoomMessage struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
	Chat   string `json:"chat"`
}

func (msg *ChatRoomMessage) convertToKeyValuePair() map[string]interface{} {
	return map[string]interface{}{
		"id":      msg.ID,
		"room_id": msg.RoomID,
		"user_id": msg.UserID,
		"chat":    msg.Chat,
	}
}

// Chat Define a chat object
type Chat struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"message"`
	MsgType   string `json:"msg_type"`
	Timestamp int64  `json:"timestamp"`
}

func (impl *connectService) HandleWebSocketStreamConnect(w http.ResponseWriter, r *http.Request) {
	// channel to signal goroutine stop
	stopChan := make(chan struct{})
	defer close(stopChan)

	// Upgrade initial http request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading connection: %v", err)
	}
	defer ws.Close()

	// set websocket connection timeout
	err = ws.SetReadDeadline(time.Now().Add(wsTimeoutDuration))
	if err != nil {
		log.Errorf("failed to set read deadline: %v", err)
	}
	err = ws.SetWriteDeadline(time.Now().Add(wsTimeoutDuration))
	if err != nil {
		log.Errorf("failed to set write deadline:: %v", err)
	}
	ctx := r.Context()

	// parse use_id from jwt token
	userID := ulid.Make().String()

	// set user online in redis
	// key := fmt.Sprintf("chatroom:%s:%s", chatroomID, userID)
	// if err = impl.redisClient.Set(ctx, key, true, 600*time.Second); err != nil {
	// 	log.Errorf("failed to set user online in redis: %v", err)
	// 	return
	// }

	// parse chat room id
	queryParams := r.URL.Query()
	chatroomID := queryParams.Get("room-id")
	if chatroomID == "" {
		log.Errorf("empty room-id parameter")
		return
	}

	// create topic if not exist
	topicID := fmt.Sprintf("streamTopic:%s", chatroomID)
	topic, err := impl.streamClient.CreateTopic(ctx, topicID)
	if err != nil {
		log.Errorf("failed to create stream topic: %v", err)
		return
	}

	// create subscription of topic channel if not exist
	subscriptionID := fmt.Sprintf("streamSubscription:%s:%s", chatroomID, userID)
	sub, err := impl.streamClient.CreateSubscription(ctx, subscriptionID, &stream.SubscriptionConfig{
		Topic:       impl.streamClient.Topic(topicID),
		TopicID:     topicID,
		ReadStartID: "0", // TODO
	})
	if err != nil {
		log.Errorf("failed to create subscription to stream topic: %v", err)
		return
	}
	defer func() {
		if err := impl.streamClient.DeleteSubscription(ctx, subscriptionID); err != nil {
			log.Errorf("failed to create subscription to stream topic: %v", err)
		}
	}()

	// get message from sub and return to websocket client
	go func(stopChan chan struct{}) {
		if err := impl.listenToStream(ctx, ws, sub, subscriptionID, stopChan); err != nil {
			log.Errorf("failed to listenToStream: %v", err)
		}
	}(stopChan)
	err = impl.sendToStream(ctx, ws, topic)
	if err != nil {
		log.Errorf("failed to sendToStream: %v", err)
	}
	log.Info("sendToStream closed")
	log.Info("HandleWebSocketStreamConnect closed")
}

func (impl *connectService) sendToStream(ctx context.Context, conn *websocket.Conn, topic stream.Topic) error {
	// listen for new messages from the client
	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled")
			return nil
		default:
			// parse client chat room message
			var chatMsg ChatRoomMessage
			err := conn.ReadJSON(&chatMsg)
			log.Infof("successfully receive message from ws client %+v ", chatMsg)
			if err != nil {
				// connection closed normally
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived) {
					log.Info("ws connection closed by client")
					return nil
				}
				return errors.Errorf("failed to read chat room msg: %v", err)
			}

			// send the new message to redis stream
			id, err := impl.idGenerator.NextID()
			if err != nil {
				return errors.Errorf("failed to generate id: %v", err)
			}
			chatMsg.ID = fmt.Sprintf("%d", id)
			streamMsg := &stream.Message{
				ID:         fmt.Sprintf("%d", id),
				Attributes: chatMsg.convertToKeyValuePair(),
			}
			if res := topic.Send(ctx, streamMsg); res.Get(ctx) != nil {
				return errors.Errorf("failed to send message %+v to stream topic: %v", streamMsg, err)
			}
		}
	}
}

func (impl *connectService) listenToStream(ctx context.Context, conn *websocket.Conn, sub stream.Subscription, subID string, stopChan chan struct{}) error {
	defer func() {
		if err := impl.streamClient.DeleteSubscription(context.Background(), subID); err != nil {
			log.Errorf("failed to delete subscription %s: %v", subID, err)
		} else {
			log.Infof("successfully delete subscription %s: %v", subID, err)
		}
	}()

	h := func(ctx context.Context, msg *stream.Message) {
		resp := &StreamChatRoomResponse{}
		if err := resp.convertRedisDataTo(msg.Attributes); err != nil {
			log.Errorf("failed to convert stream data to StreamChatRoomResponse: %v", err)
			return
		}
		err := conn.WriteJSON(resp)
		if err != nil {
			log.Errorf("failed to writing message to ws connection: %v", err)
			return
		}
		log.Infof("successfully sending message to ws client: %+v ", resp)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled")
			return nil
		case <-stopChan:
			log.Info("receive stop signal in listenToStream")
			return nil
		default:
			if err := sub.Receive(ctx, h); err != nil {
				return errors.Errorf("failed to receive message from stream subscription: %v", err)
			}
		}
	}
}
