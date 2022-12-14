module TJChatRoom

go 1.19

require github.com/gorilla/websocket v1.5.0

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.77 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v2.2.6+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.11.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.3.0 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
)

require (
	github.com/db v0.0.0
	github.com/dto v0.0.0
	github.com/entity v0.0.0
	github.com/method v0.0.0
	github.com/socket v0.0.0
	github.com/tool v0.0.0
)

replace (
	github.com/db => ./db
	github.com/dto => ./dto
	github.com/entity => ./entity
	github.com/method => ./method
	github.com/socket => ./socket
	github.com/tool => ./tool
)
