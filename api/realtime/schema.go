package realtime

import (
	"fmt"

	"gsm/pkg/util/convert"
)

type ChatRoomAction string

const (
	ActionJoinChatRoom    ChatRoomAction = "JOIN_CHAT_ROOM"
	ActionLeaveChatRoom   ChatRoomAction = "LEAVE_CHAT_ROOM"
	ActionChatRoomMessage ChatRoomAction = "CHAT_ROOM_MESSAGE"
)

// // StreamChatRoomQueryParams defines query params of api HandleWebSocketStreamConnect
// type StreamChatRoomQueryParams struct {
// 	Action string
// 	RoomID string `form:"room_id" json:"room_id"`
// }

// // BindToStreamChatRoomQueryParams bind queryParams to StreamChatRoomQueryParams
// func BindToStreamChatRoomQueryParams(queryParams url.Values) (params *StreamChatRoomQueryParams) {
// 	return &StreamChatRoomQueryParams{
// 		RoomID: queryParams.Get("room_id"),
// 	}
// }

// StreamChatRoomMessageResponse defines a stream chat room message response of api HandleWebSocketStreamConnect
type StreamChatRoomMessageResponse struct {
	ID         uint64 `json:"id"`
	Chat       string `json:"chat"`
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	RoomID     string `json:"room_id"`
	Status     string
	ErrCode    string `json:"err_code"`
	ErrMsg     string `json:"err_msg"`
}

func (resp *StreamChatRoomMessageResponse) convertFromKeyValuePairs(attr map[string]interface{}) error {
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
	if v, ok := attr["sender_name"]; ok {
		resp.SenderName = fmt.Sprintf("%s", v)
	}

	return nil
}

// StreamChatRoomRequestMessage define a stream chat room message request of api HandleWebSocketStreamConnect
type StreamChatRoomMessageRequest struct {
	Action   ChatRoomAction `json:"action"`
	ID       string         `json:"id"`
	RoomID   string         `json:"room_id"`
	SenderID string         `json:"sender_id"`
	Chat     string         `json:"chat"`
}

type StreamChatRoomMessageData struct {
}

func (msg *StreamChatRoomMessageRequest) convertToKeyValuePairs() map[string]interface{} {
	return map[string]interface{}{
		"id":        msg.ID,
		"room_id":   msg.RoomID,
		"sender_id": msg.SenderID,
		"chat":      msg.Chat,
	}
}

// // Chat Define a chat object
// type Chat struct {
// 	ID        string `json:"id"`
// 	From      string `json:"from"`
// 	To        string `json:"to"`
// 	Msg       string `json:"message"`
// 	MsgType   string `json:"msg_type"`
// 	Timestamp int64  `json:"timestamp"`
// }
