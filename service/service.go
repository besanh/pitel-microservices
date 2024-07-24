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
}

// MAP TENANT_ID SQL_CONN
var (
	SECRET_KEY_SUPERADMIN string = ""
	MapDBConn             map[string]sqlclient.ISqlClientConn
	ERR_EMPTY_CONN        = errors.New("empty_conn")
	ERR_DB_CONN_FAIL      = errors.New("db_conn_fail")

	// ES
	ES_INDEX              = "" //             = "pitel_bss_chat"
	ES_INDEX_CONVERSATION = "" // = "pitel_bss_conversation"

	// Redis
	CONVERSATION               = "conversation"
	CONVERSATION_EXPIRE        = 30 * time.Minute
	CHAT_QUEUE                 = "chat_queue"
	CHAT_QUEUE_EXPIRE          = 30 * time.Minute
	CHAT_ROUTING               = "chat_routing"
	CHAT_ROUTING_EXPIRE        = 1 * time.Hour
	CHAT_QUEUE_USER            = "chat_queue_user"
	CHAT_QUEUE_USER_EXPIRE     = 10 * time.Minute
	CHAT_APP                   = "chat_app"
	CHAT_APP_EXPIRE            = 6 * time.Hour
	USER_ALLOCATE              = "user_allocate"
	USER_ALLOCATE_EXPIRE       = 1 * time.Hour
	MANAGE_QUEUE_USER          = "manage_queue_user"
	MANAGE_QUEUE_USER_EXPIRE   = 1 * time.Hour
	CHAT_CONNECTION            = "chat_connection"
	CHAT_CONNECTION_EXPIRE     = 5 * time.Minute
	CHAT_AUTO_SCRIPT           = "chat_auto_script"
	CHAT_AUTO_SCRIPT_EXPIRE    = 30 * time.Minute
	CHAT_POLICY_SETTING        = "chat_policy_setting"
	CHAT_POLICY_SETTING_EXPIRE = 24 * time.Hour // a day since policies are rarely changed

	ORIGIN_LIST = []string{"localhost:*", "*.tel4vn.com"}

	OTT_URL                          string        = ""
	OTT_VERSION                      string        = ""
	API_SHARE_INFO_HOST              string        = ""
	API_DOC                          string        = ""
	ENABLE_PUBLISH_ADMIN             bool          = false
	ENABLE_CHAT_AUTO_SCRIPT_REPLY    bool          = false
	ENABLE_CHAT_POLICY_SETTINGS      bool          = false
	AAA_HOST                         string        = ""
	CHAT_INTEGRATE_SYSTEM            string        = "chat_integrate_system"
	CHAT_INTEGRATE_SYSTEM_EXPIRE     time.Duration = 1 * time.Hour
	CHAT_APP_INTEGRATE_SYSTEM        string        = "chat_app_integrate_system"
	CHAT_APP_INTEGRATE_SYSTEM_EXPIRE time.Duration = 1 * time.Hour

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
	BSS_CHAT_QUEUE_NAME string = "bss_chat_request_queue"
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
		err = ERR_DB_CONN_FAIL
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

func HandleGetDBConSource(authUser *model.AuthUser) (sqlclient.ISqlClientConn, error) {
	var dbCon sqlclient.ISqlClientConn
	if authUser == nil {
		return nil, errors.New("authUser is nil")
	}
	dbCon = repository.DBConn
	// dbConTmp, err := GetDBConnOfUser(*authUser)
	// if err != nil {
	// 	return dbCon, err
	// }
	// dbCon = dbConTmp
	return dbCon, nil
}
