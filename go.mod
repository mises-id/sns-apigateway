module github.com/mises-id/sns-apigateway

go 1.16

require (
	github.com/alexflint/go-filemutex v1.1.0 // indirect
	github.com/bluele/factory-go v0.0.1
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gavv/httpexpect v2.0.0+incompatible
	github.com/go-kit/kit v0.12.0
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v33 v33.0.0
	github.com/google/uuid v1.2.0
	github.com/joho/godotenv v1.4.0
	github.com/khaiql/dbcleaner v2.3.0+incompatible
	github.com/labstack/echo v3.3.10+incompatible
	github.com/lib/pq v1.10.4 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/mises-id/sns-socialsvc v0.0.0-20211230063117-75473fcba06e
	github.com/mises-id/sns-storagesvc v0.0.0-20211229064402-41052e86ccfb
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.5
	go.mongodb.org/mongo-driver v1.8.1
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	google.golang.org/grpc v1.43.0
)

replace github.com/go-kit/kit => github.com/mises-id/kit v0.12.1-0.20211203081751-bc5397e8a165

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tendermint => github.com/mises-id/tendermint v0.34.15-0.20211207033151-1f29b59c0edf

replace github.com/cosmos/cosmos-sdk => github.com/mises-id/cosmos-sdk v0.44.6-0.20211209094558-a7c9c77cfc17
