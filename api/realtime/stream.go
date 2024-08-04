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

// connectService defines the implementation of ConnectService interface
type connectService struct {
	redisClient  cache.Client
	streamClient stream.Client
	idGenerator  sonyflake.IDGenerator
}

// ConnectService defines the connect service interface
type ConnectService interface {
	HandleWebSocketStreamConnect(userID string, r *http.Request, w http.ResponseWriter)
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

// const wsTimeoutDuration = 300 * time.Second

func (impl *connectService) HandleWebSocketStreamConnect(userID string, r *http.Request, w http.ResponseWriter) {
	// Upgrade initial http request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading connection: %v", err)
	}
	defer ws.Close()

	// set websocket connection timeout
	// err = ws.SetReadDeadline(time.Now().Add(wsTimeoutDuration))
	// if err != nil {
	// 	log.Errorf("failed to set read deadline: %v", err)
	// }
	// err = ws.SetWriteDeadline(time.Now().Add(wsTimeoutDuration))
	// if err != nil {
	// 	log.Errorf("failed to set write deadline:: %v", err)
	// }
	ctx := r.Context()

	// parse chat room id and token
	queryParams := BindToStreamChatRoomQueryParams(r.URL.Query())
	chatroomID := queryParams.RoomID
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
	subID := fmt.Sprintf("streamSubscription:%s:%s", chatroomID, userID)
	subConfig := &stream.SubscriptionConfig{
		Topic:       impl.streamClient.Topic(topicID),
		TopicID:     topicID,
		ReadStartID: "0", // TODO
	}
	// TODO: if brrowser
	subscription := impl.streamClient.Subscription(subID, subConfig)
	if exist, err := subscription.Exists(ctx); err != nil {
		log.Errorf("failed to check if subscription already exist in topic: %v", err)
		return
	} else if !exist {
		if sub, err := impl.streamClient.CreateSubscription(ctx, subID, subConfig); err != nil {
			log.Errorf("failed to create subscription to stream topic: %v", err)
			return
		} else {
			subscription = sub
		}
	}
	// defer delete subscription for topic
	defer func() {
		if err := impl.streamClient.DeleteSubscription(context.Background(), topicID, subID); err != nil {
			log.Errorf("failed to delete subscription %s: %v", subID, err)
		} else {
			log.Infof("successfully delete subscription %s", subID)
		}
	}()

	// get message from sub and return to websocket client
	// channel to signal goroutine stop
	stopChan := make(chan struct{})
	go func() {
		if err := impl.listenToStream(ctx, ws, subscription, subID, stopChan); err != nil {
			log.Errorf("failed to listenToStream: %v", err)
		}
		log.Infof("listenToStream closed for client %s", ws.RemoteAddr().String())
	}()
	if err = impl.sendToStream(ctx, ws, topic, chatroomID, userID); err != nil {
		log.Errorf("failed to sendToStream: %v", err)
	}
	log.Infof("sendToStream closed for %s", ws.RemoteAddr().String())
	close(stopChan)
	log.Infof("HandleWebSocketStreamConnect closed %s", ws.RemoteAddr().String())
}

func (impl *connectService) sendToStream(ctx context.Context, conn *websocket.Conn, topic stream.Topic, chatroomID, userID string) error {
	// listen for new messages from the client
	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled")
			return nil
		default:
			// parse client chat room message
			var chatMsg StreamChatRoomMessage
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

			// TODO: check if userID and room ID matched
			// if !(chatMsg.UserID == userID && chatMsg.RoomID != chatroomID) {
			// 	return errors.Errorf("failed to handle msg from client ws due to mismatched userID orr chatroomID: %v", err)
			// }

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
	h := func(ctx context.Context, msg *stream.Message) {
		resp := &StreamChatRoomResponse{}
		if err := resp.convertRedisDataTo(msg.Attributes); err != nil {
			log.Errorf("failed to convert stream data to StreamChatRoomResponse: %v", err)
			return
		}
		resp.UserName = resp.UserID // TODO: replace with real user name
		err := conn.WriteJSON(resp)
		if err != nil {
			log.Errorf("failed to writing message to ws connection: %v", err)
			return
		}
		log.Infof("successfully sending message to ws client: %+v ", resp)
	}

	if err := sub.Receive(ctx, h, stopChan); err != nil {
		return errors.Errorf("failed to receive message from stream subscription: %v", err)
	}

	return nil
}
