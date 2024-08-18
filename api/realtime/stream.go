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
	"gsm/pkg/util/convert"
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
	PushMessage(ctx context.Context, userID string, msgReq *ChatRoomMessageRequest, clientRespCh chan *StreamWSMessageResponse) error
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
	subscription stream.Subscription
	action       ChatRoomAction
}

func (impl *streamChatService) HandleWebSocketStreamConnect(userID string, r *http.Request, w http.ResponseWriter) {
	// Upgrade initial http request to a WebSocket
	ws, err := realtime.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("failed to upgrade connection to websocket: %v", err)
	}
	remoteIP := ws.RemoteAddr().String()
	defer func() {
		ws.Close()
		log.Infof("HandleWebSocketStreamConnect fo connection %s is stopped", remoteIP)
	}()
	log.Infof("HandleWebSocketStreamConnect fo connection %s is started", remoteIP)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel() // Ensure that the context is cancelled when main exits

	// handle receiving message from chat room
	changeRoomCh := make(chan changeRoomSub)
	clientRespCh := make(chan *StreamWSMessageResponse)
	go func() {
		log.Infof("collectMsgFromSubs for connection %s is started", remoteIP)
		impl.collectMsgFromSubs(ctx, changeRoomCh, userID, clientRespCh)
		log.Infof("collectMsgFromSubs for connection %s is stopped", remoteIP)
	}()
	go func() {
		log.Infof("writeRespToWS for connection %s is started", remoteIP)
		impl.writeRespToWS(ctx, ws, clientRespCh)
		log.Infof("writeRespToWS for connection %s is stopped", remoteIP)
	}()

	// handle message request
	for {
		select {
		case <-ctx.Done():
			log.Infof("context canceled in HandleWebSocketStreamConnect, reason: %v", ctx.Err())
			return
		default:
			// parse client message request.
			var msgReq StreamWSMessageRequest
			err := ws.ReadJSON(&msgReq)
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
			if handleErr := impl.handleWebSocketStreamConnectMessage(ctx, msgReq, userID, changeRoomCh, remoteIP, clientRespCh); handleErr != nil {
				return
			}
		}
	}
}

// handleWebSocketStreamConnectMessage manage messages from the WebSocket client and route them to the appropriate handler.
func (impl *streamChatService) handleWebSocketStreamConnectMessage(ctx context.Context, msgReq StreamWSMessageRequest, userID string, changeRoomCh chan changeRoomSub, remoteIP string, clientRespCh chan *StreamWSMessageResponse) error {
	switch msgReq.Type {
	case WSMessageTypeChatRoomAction:
		var actionPayload ChatRoomActionRequest
		if err := json.Unmarshal(msgReq.Payload, &actionPayload); err != nil {
			return errors.Errorf("failed to Unmarshal payload to ChatRoomActionRequest: %v", err)
		}
		if actionPayload.RoomID == "" {
			return errors.NewErrorf(errors.InvalidArgument, "empty chat room id")
		}
		log.Infof("receive message from ws client with message type %s, payload %+v from ip %s", WSMessageTypeChatRoomAction, actionPayload, remoteIP)

		topicID := stream.ConstructChatRoomTopicID(actionPayload.RoomID)
		subID := stream.ConstructChatRoomSubscriptionID(actionPayload.RoomID, userID, remoteIP)

		switch actionPayload.Action {
		case ActionJoinChatRoom:
			if sub, err := impl.JoinChatRoom(ctx, userID, actionPayload.RoomID, topicID, subID); err != nil {
				return errors.Errorf("failed to join chat room %s", actionPayload.RoomID)
			} else {
				changeMsg := changeRoomSub{
					subscription: sub,
					action:       ActionJoinChatRoom,
				}
				changeRoomCh <- changeMsg
			}
		case ActionLeaveChatRoom:
			subConfig := &stream.SubscriptionConfig{
				Topic:       impl.streamClient.Topic(topicID),
				TopicID:     topicID,
				ReadStartID: "0", // TODO
			}
			subscription := impl.streamClient.Subscription(subID, subConfig)
			changeMsg := changeRoomSub{
				subscription: subscription,
				action:       ActionLeaveChatRoom,
			}
			changeRoomCh <- changeMsg
		}

	case WSMessageTypeChat:
		var chatMsgPayload ChatRoomMessageRequest
		if err := json.Unmarshal(msgReq.Payload, &chatMsgPayload); err != nil {
			return err
		}
		if chatMsgPayload.RoomID == "" {
			return errors.NewErrorf(errors.InvalidArgument, "empty chat room id")
		}
		log.Infof("receive message from ws client with message type %s, payload %+v from ip %s", WSMessageTypeChat, chatMsgPayload, remoteIP)

		if err := impl.PushMessage(ctx, userID, &chatMsgPayload, clientRespCh); err != nil {
			return errors.Errorf("failed to push chat message to chat room %s", chatMsgPayload.RoomID)
		}
	default:
		return nil
	}
	return nil
}

func (impl *streamChatService) JoinChatRoom(ctx context.Context, userID, chatRoomID, topicID, subID string) (stream.Subscription, error) {
	// TODO: remove when CreateChatRoom is implemented
	// create stream topic if not exist
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

func (impl *streamChatService) PushMessage(ctx context.Context, userID string, msgReq *ChatRoomMessageRequest, clientRespCh chan *StreamWSMessageResponse) error {
	// TODO: get id from id producer center
	// send the new message to redis stream
	id, err := impl.idGenerator.NextID()
	if err != nil {
		return errors.Errorf("failed to generate id: %v", err)
	}
	sID := fmt.Sprintf("%d", id)

	// TODO: save message in db

	// push message into topic stream
	msgReq.ID = sID
	streamMsg := &stream.Message{
		ID:         sID,
		Attributes: msgReq.convertToKeyValuePairs(userID),
	}
	topicID := stream.ConstructChatRoomTopicID(msgReq.RoomID)
	if res := impl.streamClient.Topic(topicID).Send(ctx, streamMsg); res.Err(ctx) != nil {
		resp := &StreamWSMessageResponse{
			Type:    WSMessageACK,
			Status:  MessageStatusError,
			ErrCode: errors.MessageHandleError,
			ErrMsg:  fmt.Sprintf("%s", res.Err(ctx)),
		}
		clientRespCh <- resp
		return errors.Errorf("failed to send message %+v to stream topic: %v", streamMsg, res.Err(ctx))
	} else {
		resp := &StreamWSMessageResponse{
			Type:   WSMessageACK,
			Status: MessageStatusSuccess,
		}
		clientRespCh <- resp
	}
	return nil
}

func (impl *streamChatService) writeRespToWS(ctx context.Context, conn *websocket.Conn, clientRespCh chan *StreamWSMessageResponse) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("context canceled in writeRespToWS, reason: %v", ctx.Err())
			return
		case resp := <-clientRespCh:
			if err := conn.WriteJSON(resp); err != nil {
				log.Errorf("failed to writing message to ws client with address %s: %v", conn.RemoteAddr().String(), err)
				return
			}
			log.Infof("successfully sending message to ws client with address %s , type %s, status %s, payload %+v", conn.RemoteAddr().String(), resp.Type, resp.Status, convert.FormatJsonString(string(resp.Payload)))
		}
	}
}

func (impl *streamChatService) collectMsgFromSubs(ctx context.Context, changeRoomCh chan changeRoomSub, userID string, clientRespCh chan *StreamWSMessageResponse) error {
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
		clientRespCh <- resp
	}

	connSubHistory := make(map[string]map[stream.Subscription]chan struct{}) // [subscription id]:[stream subscription]:[stop chan struct{} for stopping the subscription]

	var curSub stream.Subscription
	for {
		select {
		case <-ctx.Done():
			// handle connection closure, e.g., close subscriptions
			for _, subMap := range connSubHistory {
				for sub, ch := range subMap {
					// signal all subscription to stop
					ch <- struct{}{}
					// delete subscription
					if err := impl.streamClient.DeleteSubscription(ctx, sub.GetTopicID(), sub.GetSubID()); err != nil {
						return errors.Errorf("failed to delete subscription %s: %v", sub.GetSubID(), err)
					}
					log.Infof("successfully delete subscription %s", sub.GetSubID())
				}
			}
			log.Infof("context canceled in collectMsgFromSubs, reason: %v", ctx.Err())
			return nil
		case subMsg := <-changeRoomCh:
			// skip invalid subMsg
			if subMsg.subscription == nil {
				continue
			}
			msgSubID := subMsg.subscription.GetSubID()

			switch subMsg.action {
			case ActionJoinChatRoom:
				// record message subID into connSubscriptions history
				var stopCh chan struct{}
				if _, ok := connSubHistory[msgSubID]; !ok {
					connSubHistory[msgSubID] = make(map[stream.Subscription]chan struct{})
					stopCh = make(chan struct{})
					connSubHistory[msgSubID][subMsg.subscription] = stopCh
				} else {
					stopCh = connSubHistory[msgSubID][subMsg.subscription]
				}

				// stop curSub before subscribe to other chat room to avoid panic due to occurrence write to ws connect
				if curSub != nil {
					for _, sub := range connSubHistory {
						for _, stopCh := range sub {
							// signal all other subscription to stop
							stopCh <- struct{}{}
						}
					}
				}
				// send join room ack to client
				resp := &StreamWSMessageResponse{
					Type:   WSMessageTypeChatRoomAction,
					Status: MessageStatusSuccess,
				}
				clientRespCh <- resp

				// receiving message from selected chat room
				curSub = subMsg.subscription
				if err := curSub.Receive(ctx, h, stopCh); err != nil {
					return errors.Errorf("failed to receive message from stream subscription: %v", err)
				}
			case ActionLeaveChatRoom:
				if curSub == nil || curSub.GetSubID() != msgSubID {
					log.Errorf("message with action leave chat room doesn't match current chat room")
				} else {
					if subInfo, ok := connSubHistory[msgSubID]; ok {
						for _, ch := range subInfo {
							// signal subscription to stop
							ch <- struct{}{}
						}
					}
				}
				// send leave room ack to client
				resp := &StreamWSMessageResponse{
					Type:   WSMessageTypeChatRoomAction,
					Status: MessageStatusSuccess,
				}
				clientRespCh <- resp
			}
		}
	}
}
