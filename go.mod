module github.com/everstake/cosmoscan-api

go 1.14

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

require (
	github.com/Masterminds/squirrel v1.4.0
	github.com/cosmos/cosmos-sdk v0.42.8
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang-migrate/migrate/v4 v4.11.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.1.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/mailru/go-clickhouse v1.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/persistenceOne/persistenceCore v0.1.3
	github.com/rs/cors v1.7.0
	github.com/rubenv/sql-migrate v0.0.0-20200429072036-ae26b214fa43
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	github.com/superoo7/go-gecko v1.0.0
	github.com/tendermint/tendermint v0.34.14
	github.com/urfave/negroni v1.0.0
	go.uber.org/zap v1.17.0
)

replace github.com/persistenceOne/persistenceCore => ../persistenceCore
