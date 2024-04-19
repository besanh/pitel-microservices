package model

type User struct {
	AuthUser       *AuthUser     `json:"auth_user"`
	ConnectionId   string        `json:"connection_id"`
	QueueId        string        `json:"queue_id"`
	IsOk           bool          `json:"is_ok"`
	PreviousAssign *UserAllocate `json:"previous_assign"`
}
