package realtime

import (
	"fmt"

	"gsm/pkg/util/convert"
)

// StreamChatRoomQueryParams defines query params of api HandleWebSocketStreamConnect
type StreamChatRoomQueryParams struct {
	RoomID string `form:"room_id" json:"room_id"`
}

// StreamChatRoomResponse defines resp body of api HandleWebSocketStreamConnect
type StreamChatRoomResponse struct {
	ID       uint64 `json:"id"`
	Chat     string `json:"chat"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	RoomID   string `json:"room_id"`
}

func (resp *StreamChatRoomResponse) convertRedisDataTo(attr map[string]interface{}) error {
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
	if v, ok := attr["user_id"]; ok {
		resp.UserID = fmt.Sprintf("%s", v)
	}
	if v, ok := attr["user_name"]; ok {
		resp.UserName = fmt.Sprintf("%s", v)
	}
	if v, ok := attr["room_id"]; ok {
		resp.RoomID = fmt.Sprintf("%s", v)
	}

	return nil
}
