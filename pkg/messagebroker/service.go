package messagebroker

import "fmt"

// ConstructSubscriptionGroupID construct subscription group id for single customer
func ConstructSubscriptionGroupID(topicID, userID string) string {
	return fmt.Sprintf("%s-%s", topicID, userID)
}
