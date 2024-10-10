module github.com/tel4vn/fins-microservices

go 1.22.0

toolchain go1.22.5

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.35.1-20240920164238-5a7b106cbb87.1
	github.com/bufbuild/protovalidate-go v0.7.2
	github.com/coder/websocket v1.8.12
	github.com/gin-gonic/gin v1.10.0
	github.com/go-co-op/gocron v1.37.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0
	github.com/hibiken/asynq v0.24.1
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.6.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	github.com/uptrace/bun v1.2.3
	github.com/uptrace/bun/dialect/mysqldialect v1.2.3
	github.com/uptrace/bun/dialect/pgdialect v1.2.3
	github.com/uptrace/bun/driver/pgdriver v1.2.3
	github.com/xuri/excelize/v2 v2.8.1
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
	google.golang.org/genproto/googleapis/api v0.0.0-20241007155032-5fefd90f89a9
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

require (
	cloud.google.com/go/compute v1.25.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/adjust/uniuri v0.0.0-20130923163420-498743145e60 // indirect
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/alicebob/miniredis/v2 v2.33.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.6.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/validator/v10 v10.22.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cel-go v0.21.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gookit/goutil v0.6.17 // indirect
	github.com/gookit/gsr v0.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.10 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.4.0 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/xuri/efp v0.0.0-20240408161823-9ad904a10d6d // indirect
	github.com/xuri/nfp v0.0.0-20240318013403-ab9948c2c4a7 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.22.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20170922094635-f56db5e73a5e // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/redis.v3 v3.6.4 // indirect
)

require (
	github.com/adjust/rmq v1.0.0
	github.com/adjust/rmq/v5 v5.2.0
	github.com/bytedance/sonic v1.12.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/go-elasticsearch v0.0.0
	github.com/elastic/go-elasticsearch/v8 v8.15.0
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-resty/resty/v2 v2.15.3
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/google/uuid v1.6.0
	github.com/gookit/slog v0.5.6
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/jellydator/ttlcache/v3 v3.3.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/minio-go/v7 v7.0.77
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/rabbitmq/rabbitmq-stream-go-client v1.4.10
	github.com/shaj13/go-guardian/v2 v2.11.6
	github.com/shaj13/libcache v1.2.1
	github.com/spf13/cobra v1.8.1
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/arch v0.11.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.30.0
	golang.org/x/sys v0.26.0
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/time v0.7.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools v2.2.0+incompatible
	mellium.im/sasl v0.3.2 // indirect
)
