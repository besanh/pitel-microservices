package repository

import (
	"context"
	"regexp"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

var DBConn sqlclient.ISqlClientConn
var ESClient elasticsearch.IElasticsearchClient

func CreateTable(ctx context.Context, db sqlclient.ISqlClientConn, entity any) (err error) {
	_, err = db.GetDB().NewCreateTable().Model(entity).
		IfNotExists().
		Exec(ctx)
	return
}

func InitRepositories() {
	ExampleRepo = NewExample()
	AuthSourceRepo = NewAuthSource()
	ChatAppRepo = NewChatApp()
	ChatConnectionAppRepo = NewConnectionApp()
	ChatQueueRepo = NewChatQueue()
	ChatQueueUserRepo = NewChatQueueUser()
	ChatRoutingRepo = NewChatRouting()
	UserAllocateRepo = NewUserAllocate()
	ConnectionQueueRepo = NewConnectionQueue()
	ShareInfoRepo = NewShareInfo()
	ManageQueueRepo = NewManageQueue()
	ChatMsgSampleRepo = NewChatMsgSample()
	ChatScriptRepo = NewChatScript()
	ChatLabelRepo = NewChatLabel()
	ChatAutoScriptRepo = NewChatAutoScript()
	ChatPolicySettingRepo = NewChatPolicySetting()
	VendorRepo = NewVendor()
	ChatIntegrateSystemRepo = NewChatIntegrateSystem()
	ChatAppIntegrateSystemRepo = NewChatAppIntegrateSystem()
	ChatRoleRepo = NewChatRole()
	ChatUserRepo = NewChatUser()
	ChatTenantRepo = NewChatTenant()
	ChatConnectionPipelineRepo = NewConnectionPipeline()
}

func InitRepositoriesES() {
	ESRepo = NewES()
	ConversationESRepo = NewConversationES()
	MessageESRepo = NewMessageES()
}

func InitTables(ctx context.Context, dbConn sqlclient.ISqlClientConn) {
	if err := CreateTable(ctx, dbConn, (*model.Example)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatApp)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatConnectionApp)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatQueue)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatRouting)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatQueueUser)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.AllocateUser)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ConnectionQueue)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ShareInfoForm)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.FacebookPage)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatManageQueueUser)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatEmail)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatMsgSample)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatScript)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatAutoScriptToChatScript)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatAutoScriptToChatLabel)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatAutoScript)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatLabel)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatPolicySetting)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatVendor)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatAppIntegrateSystem)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatIntegrateSystem)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatRole)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatUser)(nil)); err != nil {
		log.Error(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.ChatTenant)(nil)); err != nil {
		log.Error(err)
	}
	log.Println("TABLES WERE CREATED")
}

func InitColumn(ctx context.Context, db sqlclient.ISqlClientConn) {
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
