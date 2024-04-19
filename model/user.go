package model

type User struct {
	AuthUser           *AuthUser     `json:"auth_user"`
	ConnectionId       string        `json:"connection_id"`
	QueueId            string        `json:"queue_id"`
	IsOk               bool          `json:"is_ok"`
	IsReassignNew      bool          `json:"is_reassign_new"`
	IsReassignSame     bool          `json:"is_reassign_same"`
	UserIdReassignNew  string        `json:"user_id_reassign_new"`
	UserIdReassignSame string        `json:"user_id_reassign_same"`
	UserIdRemove       string        `json:"user_id_remove"`
	PreviousAssign     *UserAllocate `json:"previous_assign"`
}
