package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/time/rate"
)

func InitServices() {
	ChatAuthService = NewChatAuth()
	ExampleService = NewExample()
	ChatRoleService = NewChatRole()
	ChatUserService = NewChatUser()
	ChatVendorService = NewChatVendor()
	ChatTenantService = NewChatTenant()
	ChatIntegrateSystemService = NewChatIntegrateSystem()
	ChatAppService = NewChatApp()
	AssignConversationService = NewAssignConversation()
	ChatMessageSampleService = NewChatMsgSample()
	ChatScriptService = NewChatScript()
	ChatAutoScriptService = NewChatAutoScript()
	ChatEmailService = NewChatEmail()
	ChatConnectionPipelineService = NewChatConnectionPipeline()
	ChatConnectionAppService = NewChatConnectionApp()
	ChatQueueService = NewChatQueue()
	ConversationService = NewConversation()
	MessageService = NewMessage()
	ChatPolicySettingService = NewChatPolicySetting()
	ChatLabelService = NewChatLabel()
	ChatRoutingService = NewChatRouting()
	ProfileService = NewProfile()
	ChatConnectionQueueService = NewChatConnectionQueue()
	ManageQueueService = NewManageQueue()
	ChatQueueUserService = NewChatQueueUser()
	ShareInfoService = NewShareInfo()
	NotesListService = NewNotesList()
	ChatReportService = NewChatReport()
}

// MAP TENANT_ID SQL_CONN
var (
	SECRET_KEY_SUPERADMIN string = ""
	MapDBConn             map[string]sqlclient.ISqlClientConn
	ErrEmptyConn          error = errors.New("empty_conn")
	ErrDBConnFail         error = errors.New("db_conn_fail")

	ORIGIN_LIST []string = []string{"localhost:*", "*.tel4vn.com"}

	// ES
	ES_INDEX_MESSAGE      string = "pitel_bss_message"
	ES_INDEX_CONVERSATION string = "pitel_bss_conversation"

	// Redis
	CONVERSATION                     string        = "conversation"
	CONVERSATION_EXPIRE              time.Duration = 30 * time.Minute
	CHAT_QUEUE                       string        = "chat_queue"
	CHAT_QUEUE_EXPIRE                time.Duration = 30 * time.Minute
	CHAT_ROUTING                     string        = "chat_routing"
	CHAT_ROUTING_EXPIRE              time.Duration = 1 * time.Hour
	CHAT_QUEUE_USER                  string        = "chat_queue_user"
	CHAT_QUEUE_USER_EXPIRE           time.Duration = 10 * time.Minute
	CHAT_APP                         string        = "chat_app"
	CHAT_APP_EXPIRE                  time.Duration = 6 * time.Hour
	USER_ALLOCATE                    string        = "user_allocate"
	USER_ALLOCATE_EXPIRE             time.Duration = 1 * time.Hour
	MANAGE_QUEUE_USER                string        = "manage_queue_user"
	MANAGE_QUEUE_USER_EXPIRE         time.Duration = 1 * time.Hour
	CHAT_CONNECTION                  string        = "chat_connection"
	CHAT_CONNECTION_EXPIRE           time.Duration = 5 * time.Minute
	CHAT_AUTO_SCRIPT                 string        = "chat_auto_script"
	CHAT_AUTO_SCRIPT_EXPIRE          time.Duration = 30 * time.Minute
	CHAT_POLICY_SETTING              string        = "chat_policy_setting"
	CHAT_POLICY_SETTING_EXPIRE       time.Duration = 24 * time.Hour // a day since policies are rarely changed
	CHAT_INTEGRATE_SYSTEM            string        = "chat_integrate_system"
	CHAT_INTEGRATE_SYSTEM_EXPIRE     time.Duration = 1 * time.Hour
	CHAT_APP_INTEGRATE_SYSTEM        string        = "chat_app_integrate_system"
	CHAT_APP_INTEGRATE_SYSTEM_EXPIRE time.Duration = 1 * time.Hour

	OTT_URL                       string = ""
	OTT_VERSION                   string = ""
	API_SHARE_INFO_HOST           string = ""
	API_DOC                       string = ""
	ENABLE_PUBLISH_ADMIN          bool   = false
	ENABLE_CHAT_AUTO_SCRIPT_REPLY bool   = false
	ENABLE_CHAT_POLICY_SETTINGS   bool   = false
	CONVERSATION_NOTES_LIST_LIMIT int    = 1000

	// Storage
	S3_ENDPOINT    string = ""
	S3_BUCKET_NAME string = ""
	S3_ACCESS_KEY  string = ""
	S3_SECRET_KEY  string = ""

	// Zalo
	ZALO_SHARE_INFO_SUBTITLE string = ""
	ZALO_POLICY_CHAT_WINDOW  int    = 0

	// Facebook
	FACEBOOK_GRAPH_API_VERSION  string = ""
	FACEBOOK_POLICY_CHAT_WINDOW int    = 0

	// DB PG
	DB_HOST     string = ""
	DB_DATABASE string = ""
	DB_USERNAME string = ""
	DB_PASSWORD string = ""
	DB_PORT     int    = 0

	// Email default
	SMTP_SERVER         string = ""
	SMTP_USERNAME       string = ""
	SMTP_PASSWORD       string = ""
	SMTP_MAILPORT       int    = 0
	SMTP_INFORM         bool   = false
	ENABLE_NOTIFY_EMAIL bool   = false

	// Queue
	BSS_CHAT_CONVERSATION_QUEUE_NAME string = "bss_chat_conversation_request_queue"
	BSS_CHAT_MESSAGE_QUEUE_NAME      string = "bss_chat_message_request_queue"

	// RabbitMQ stream
	RABBITMQ_STREAM_NAME string = ""
)

type (
	DBConfig struct {
		Host     string
		Port     int
		Database string
		Username string
		Password string
	}

	Subscriber struct {
		Id                 string      `json:"id"`
		TenantId           string      `json:"tenant_id"`
		BusinessUnitId     string      `json:"business_unit_id"`
		UserId             string      `json:"user_id"`
		Username           string      `json:"username"`
		Services           []string    `json:"services"`
		Level              string      `json:"level"`
		Message            chan []byte `json:"-"`
		CloseSlow          func()      `json:"-"`
		SubscribeAt        time.Time   `json:"subscribe_at"`
		IsAssignRoundRobin bool        `json:"is_assign_round_robin"`
		Source             string      `json:"source"`
		QueueId            string      `json:"queue_id"`      //use for allocate manager
		ConnectionId       string      `json:"connection_id"` // use for allocate
		RoleId             string      `json:"role_id"`
		ApiUrl             string      `json:"api_url"`
		Status             string      `json:"status"`
	}

	Subscribers struct {
		SubscriberMessageBuffer int
		PublishLimiter          *rate.Limiter
		SubscribersMu           sync.Mutex
		Subscribers             SubscriberItem
	}

	SubscriberItem map[*Subscriber]struct{}
)

func NewDBConn(tenantId string, config DBConfig) (dbConn sqlclient.ISqlClientConn, err error) {
	sqlClientConfig := sqlclient.SqlConfig{
		Host:         config.Host,
		Port:         config.Port,
		Database:     config.Database,
		Username:     config.Username,
		Password:     config.Password,
		DialTimeout:  20,
		ReadTimeout:  30,
		WriteTimeout: 30,
		Timeout:      30,
		PoolSize:     10,
		MaxIdleConns: 10,
		MaxOpenConns: 10,
		Driver:       sqlclient.POSTGRESQL,
	}
	dbConn = sqlclient.NewSqlClient(sqlClientConfig)
	if err = dbConn.Connect(); err != nil {
		log.Error(err)
		err = ErrDBConnFail
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		repository.InitTables(ctx, dbConn)
	}()
	wg.Wait()
	MapDBConn[tenantId] = dbConn
	return
}

func GetDBConnOfUser(user model.AuthUser) (dbConn sqlclient.ISqlClientConn, err error) {
	// if len(user.DatabaseHost) < 1 {
	// 	err = ERR_EMPTY_CONN
	// 	return
	// }
	// dbConn, ok := MapDBConn[user.TenantId]
	// if !ok {
	// 	dbConn, err = NewDBConn(user.TenantId, DBConfig{
	// 		Host:     user.DatabaseHost,
	// 		Port:     user.DatabasePort,
	// 		Database: user.DatabaseName,
	// 		Username: user.DatabaseUser,
	// 		Password: user.DatabasePassword,
	// 	})
	// 	return
	// }
	return
}

func HandleGetDBConSource(authUser *model.AuthUser) (dbCon sqlclient.ISqlClientConn, err error) {
	if authUser == nil {
		err = errors.New("authUser is nil")
		return
	}
	dbCon = repository.DBConn
	return
}
