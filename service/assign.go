package service

type (
	IAssignConversation interface {
		// GetUserInQueue(ctx context.Context, data model.GetUserInQueueRequest) (int, any)
		// AssignConversation(ctx context.Context, req *model.AssignConversationRequest) (*model.AssignConversationResponse, error)
	}
	AssignConversation struct{}
)

func NewAssignConversation() IAssignConversation {
	return &AssignConversation{}
}

// func (s *AssignConversation) GetUserInQueue(ctx context.Context, data model.GetUserInQueueRequest) (int, any) {
// 	return nil, nil
// }
