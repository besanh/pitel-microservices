package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	INotesList interface {
		InsertNoteInConversation(ctx context.Context, authUser *model.AuthUser, data *model.ConversationNoteRequest) (string, error)
		UpdateNoteInConversationById(ctx context.Context, authUser *model.AuthUser, id string, data model.ConversationNoteRequest) error
		DeleteNoteInConversationById(ctx context.Context, authUser *model.AuthUser, id string, data model.ConversationNoteRequest) error
		GetConversationNotesList(ctx context.Context, authUser *model.AuthUser, filter model.ConversationNotesListFilter, limit, offset int) (total int, result *[]model.NotesList, err error)
	}
	NotesList struct{}
)

var NotesListService INotesList

func NewNotesList() INotesList {
	return &NotesList{}
}

const eventNotesList = "notes_list"

func (s *NotesList) InsertNoteInConversation(ctx context.Context, authUser *model.AuthUser, data *model.ConversationNoteRequest) (id string, err error) {
	notesListFilter := model.ConversationNotesListFilter{
		TenantId:       authUser.TenantId,
		ConversationId: data.ConversationId,
		AppId:          data.AppId,
		OaId:           data.OaId,
	}
	total, _, err := repository.NotesListRepo.GetNotesLists(ctx, repository.DBConn, notesListFilter, 0, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > CONVERSATION_NOTES_LIST_LIMIT {
		err = errors.New("out of limit of notes list in conversation")
		log.Error(err)
		return
	}

	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, data.AppId, data.ConversationId)
	if err != nil {
		log.Error(err)
		return
	} else if len(conversationExist.ConversationId) < 1 {
		err = errors.New("conversation " + data.ConversationId + " not found")
		log.Error(err)
		return
	}
	newEntry := model.NotesList{
		Base:           model.InitBase(),
		TenantId:       authUser.TenantId,
		Content:        data.Content,
		ConversationId: data.ConversationId,
		AppId:          data.AppId,
		OaId:           data.OaId,
	}
	id = newEntry.Id
	if err = repository.NotesListRepo.Insert(ctx, repository.DBConn, newEntry); err != nil {
		log.Error(err)
		return
	}

	filter := model.AllocateUserFilter{
		AppId:          data.AppId,
		ConversationId: data.ConversationId,
		MainAllocate:   "active",
	}
	_, allocateUser, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*allocateUser) > 0 {
		allocateUserTmp := (*allocateUser)[0]
		// Event to manager
		manageQueueUser, errTmp := GetManageQueueUser(ctx, allocateUserTmp.QueueId)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if len(manageQueueUser.Id) < 1 {
			err = errors.New("queue " + allocateUserTmp.QueueId + " not found")
			log.Error(err)
			return
		}
		PublishEventToHighLevel(authUser, manageQueueUser, variables.EVENT_CHAT["conversation_note_created"], eventNotesList, &newEntry, variables.MANAGER_LEVEL, variables.ADMIN_LEVEL)
	} else {
		PublishEventToHighLevel(authUser, nil, variables.EVENT_CHAT["conversation_note_created"], eventNotesList, &newEntry, variables.ADMIN_LEVEL)
	}
	return
}

func (s *NotesList) UpdateNoteInConversationById(ctx context.Context, authUser *model.AuthUser, id string, data model.ConversationNoteRequest) (err error) {
	noteExist, err := repository.NotesListRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if noteExist == nil {
		err = errors.New("note " + id + " not found")
		log.Error(err)
		return
	}

	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, data.AppId, data.ConversationId)
	if err != nil {
		log.Error(err)
		return
	} else if len(conversationExist.ConversationId) < 1 {
		err = errors.New("conversation " + data.ConversationId + " not found")
		log.Error(err)
		return
	}
	noteExist.Content = data.Content
	if err = repository.NotesListRepo.Update(ctx, repository.DBConn, *noteExist); err != nil {
		log.Error(err)
		return
	}

	filter := model.AllocateUserFilter{
		AppId:          data.AppId,
		ConversationId: data.ConversationId,
		MainAllocate:   "active",
	}
	_, allocateUser, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*allocateUser) > 0 {
		allocateUserTmp := (*allocateUser)[0]
		// Event to manager
		manageQueueUser, errTmp := GetManageQueueUser(ctx, allocateUserTmp.QueueId)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if len(manageQueueUser.Id) < 1 {
			err = errors.New("queue " + allocateUserTmp.QueueId + " not found")
			log.Error(err)
			return
		}
		PublishEventToHighLevel(authUser, manageQueueUser, variables.EVENT_CHAT["conversation_note_updated"], eventNotesList, noteExist, variables.MANAGER_LEVEL, variables.ADMIN_LEVEL)
	} else {
		PublishEventToHighLevel(authUser, nil, variables.EVENT_CHAT["conversation_note_updated"], eventNotesList, noteExist, variables.ADMIN_LEVEL)
	}
	return
}

func (s *NotesList) DeleteNoteInConversationById(ctx context.Context, authUser *model.AuthUser, id string, data model.ConversationNoteRequest) (err error) {
	noteExist, err := repository.NotesListRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	}

	if err = repository.NotesListRepo.Delete(ctx, repository.DBConn, id); err != nil {
		log.Error(err)
		return
	}

	filter := model.AllocateUserFilter{
		AppId:          data.AppId,
		ConversationId: data.ConversationId,
		MainAllocate:   "active",
	}
	_, allocateUser, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*allocateUser) > 0 {
		allocateUserTmp := (*allocateUser)[0]
		// Event to manager
		manageQueueUser, errTmp := GetManageQueueUser(ctx, allocateUserTmp.QueueId)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if len(manageQueueUser.Id) < 1 {
			err = errors.New("queue " + allocateUserTmp.QueueId + " not found")
			log.Error(err)
			return
		}
		PublishEventToHighLevel(authUser, manageQueueUser, variables.EVENT_CHAT["conversation_note_removed"], eventNotesList, noteExist, variables.MANAGER_LEVEL, variables.ADMIN_LEVEL)
	} else {
		PublishEventToHighLevel(authUser, nil, variables.EVENT_CHAT["conversation_note_removed"], eventNotesList, noteExist, variables.ADMIN_LEVEL)
	}
	return
}

func (s *NotesList) GetConversationNotesList(ctx context.Context, authUser *model.AuthUser, filter model.ConversationNotesListFilter, limit, offset int) (total int, result *[]model.NotesList, err error) {
	filter.TenantId = authUser.TenantId
	total, result, err = repository.NotesListRepo.GetNotesLists(ctx, repository.DBConn, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}
