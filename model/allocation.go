package model

// Get agent online on queue to insert
// Table temporary, cronjob 3 hours to delete
type AgentAllocation struct {
	UserIdByApp   string `json:"user_id_by_app"`
	AgentId       string `json:"agent_id"`
	QueueId       string `json:"queue_id"`
	AllocatedTime int64  `json:"allocated_time"`
}
