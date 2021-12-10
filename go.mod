module github.com/mises-id/sns-apigateway

go 1.16

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/alexflint/go-filemutex v1.1.0 // indirect
	github.com/bluele/factory-go v0.0.1
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fasthttp-contrib/websocket v0.0.0-20160511215533-1f3b11f56072 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/gavv/httpexpect v2.0.0+incompatible
	github.com/go-kit/kit v0.12.0
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/mock v1.4.4
	github.com/google/go-github/v33 v33.0.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/joho/godotenv v1.4.0
	github.com/khaiql/dbcleaner v2.3.0+incompatible
	github.com/labstack/echo v3.3.10+incompatible
	github.com/lib/pq v1.10.4 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/mises-id/sns-socialsvc v0.0.0-20211208063925-3b844bc52949
	github.com/mises-id/sns-storagesvc v0.0.0-00010101000000-000000000000
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.5
	github.com/valyala/fasthttp v1.31.0 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	go.mongodb.org/mongo-driver v1.8.0
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c
	google.golang.org/grpc v1.42.0
)

replace github.com/go-kit/kit => github.com/mises-id/kit v0.12.1-0.20211203081751-bc5397e8a165

replace github.com/mises-id/sns-socialsvc => ../socialsvc

replace github.com/mises-id/sns-storagesvc => ../storagesvc
