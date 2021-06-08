module ps-chartdata

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/shopspring/decimal v1.2.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.7.1
	github.com/valyala/fasttemplate v1.2.0 // indirect
	go.mongodb.org/mongo-driver v1.4.5
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	google.golang.org/api v0.37.0 // indirect
)

replace google.golang.org/grpc v1.34.0 => google.golang.org/grpc v1.29.0
