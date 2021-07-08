https://github.com/golang-standards/project-layout


go mod init myapp
# go mod init github.com/rvasily/myapp
go build
go mod download
go mod verify
go mod tidy

go build  -o ./bin/myapp ./cmd/myapp
go test -v -coverpkg=./... ./...

# если хотите чтобы у вас все зависимости жили рядом с вашим проектом и не подкачивались из-вне
go mod vendor
# с версии 1.16 не надо указывтаь мод-вендор

# обновление модулей
go get github.com/rvasily/examplerepo@v0.1.1

# goproxy
GOPROXY=https://proxy.golang.org

GOPROXY=https://proxy.corp.example.com
GONOSUMDB=*.corp.example.com,*.gitlab.example.com

# приватные репозитории
git config \
    --global \
    url."ssh://username@gitlab.corp.example.com".insteadOf "https://gitlab.corp.example.com"

git config \
    --global \
    url."https://username:<access-token>@gitlab.corp.example.com".insteadOf "https://gitlab.corp.example.com"

$HOME/.netrc
machine gitlab.corp.example.com login USERNAME password APIKEY

https://golang.org/doc/faq#git_https

# директивы gomod
replace

# почитать 
https://blog.golang.org/using-go-modules
https://golang.org/cmd/go/
https://golang.org/ref/mod
https://golang.org/doc/modules/developing
https://github.com/golang/go/wiki/Modules
