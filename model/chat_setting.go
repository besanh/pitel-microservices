package model

// Get all setting with a conversation
type ChatSetting struct {
	RoutingAlias    string
	QueueUser       []ChatQueueUser
	Conversation    Conversation
	Message         Message
	ConnectionApp   ChatConnectionApp
	ConnectionQueue ConnectionQueue
	PreviousAssign  UserAllocate
}
