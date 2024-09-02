package realtime

import (
	"encoding/json"
	"fmt"
	"time"

	"gsm/pkg/errors"
	"gsm/pkg/util/convert"
)

type ChatRoomAction string

const (
	ActionJoinChatRoom    ChatRoomAction = "JOIN_CHAT_ROOM"
	ActionLeaveChatRoom   ChatRoomAction = "LEAVE_CHAT_ROOM"
	ActionChatRoomMessage ChatRoomAction = "CHAT_ROOM_MESSAGE"
)

type MessageStatus string

const (
	MessageStatusSuccess MessageStatus = "SUCCESS"
	MessageStatusError   MessageStatus = "ERROR"
)

type WSMessageType string

const (
	// req resp
	WSMessageTypeChat           WSMessageType = "CHAT"
	WSMessageTypeChatRoomAction WSMessageType = "CHAT_ROOM_ACTION"
	WSMessageACK                WSMessageType = "ACK"

	// stream receiving
	WSMessageTypeChatStream WSMessageType = "CHAT_STREAM"
)

// StreamWSMessageResponse defines a stream chat room message response of api HandleWebSocketStreamConnect
type StreamWSMessageResponse struct {
	Type    WSMessageType   `json:"type"`
	Status  MessageStatus   `json:"status"`
	ErrCode errors.ErrCode  `json:"err_code"`
	ErrMsg  string          `json:"err_msg"`
	Payload json.RawMessage `json:"payload"` // Allows for flexible payloads
}

// ChatRoomMessageResponsePayload defines a chat room message response
type ChatRoomMessageResponsePayload struct {
	ID       uint64 `json:"id"`
	RoomID   string `json:"room_id"`
	Chat     string `json:"chat"`
	SenderID string `json:"sender_id"`
}

// ChatRoomActionResponsePayload defines a chat room action response
type ChatRoomActionResponsePayload struct {
}

func (resp *ChatRoomMessageResponsePayload) convertFromKeyValuePairs(attr map[string]interface{}) error {
	if resp == nil {
		return fmt.Errorf("mapping struct should not be nil")
	}

	if v, ok := attr["id"]; ok {
		i, err := convert.ToUint64(v)
		if err != nil {
			return err
		}
		resp.ID = i
	}
	if v, ok := attr["room_id"]; ok {
		resp.RoomID = fmt.Sprintf("%s", v)
	}
	if v, ok := attr["chat"]; ok {
		resp.Chat = fmt.Sprintf("%s", v)
	}
	if v, ok := attr["sender_id"]; ok {
		resp.SenderID = fmt.Sprintf("%s", v)
	}

	return nil
}

// StreamChatRoomRequestMessage define a stream chat room message request of api HandleWebSocketStreamConnect
type StreamWSMessageRequest struct {
	Type    WSMessageType   `json:"type"`
	Payload json.RawMessage `json:"payload"` // Allows for flexible payloads
}

type ChatRoomMessageRequestPayLoad struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Chat   string `json:"chat"`
}

type ChatRoomActionRequestPayload struct {
	Action ChatRoomAction `json:"action"`
	RoomID string         `json:"room_id"`
}

func (msg *ChatRoomMessageRequestPayLoad) convertToKeyValuePairs(userID string) map[string]interface{} {
	return map[string]interface{}{
		"id":        msg.ID,
		"room_id":   msg.RoomID,
		"sender_id": userID,
		"chat":      msg.Chat,
		"timestamp": time.Now().Unix(),
	}
}
