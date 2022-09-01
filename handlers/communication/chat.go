package communication

type ChatGroup struct {
	GroupID string `json:"group_id"`
}

type ChatMessage struct {
	TopicID string `json:"topic_id"` // The topicID will be mapped to the group id
}
