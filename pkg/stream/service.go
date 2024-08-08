package stream

import "fmt"

// ConstructChatRoomSubscriptionID construct subscription group id for customer
func ConstructChatRoomSubscriptionID(chatroomID, userID string) string {
	return fmt.Sprintf("streamSubscription:%s:%s", chatroomID, userID)
}

// ConstructChatRoomTopicID construct topic stream for chatroom
func ConstructChatRoomTopicID(chatroomID string) string {
	return fmt.Sprintf("streamTopic:%s", chatroomID)
}
