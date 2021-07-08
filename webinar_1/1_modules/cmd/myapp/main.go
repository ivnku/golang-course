package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rvasily/examplerepo"

	// "myapp/pkg/user"
	// "gitlab.com/rvasily/go-stepik-2021q1/pkg/user"
	"gitlab.com/rvasily/go-sm-stepik/webinar_1/1_modules/pkg/user"
)

func main() {
	u := user.NewUser(42, "rvasily")
	fmt.Println("my user:", u)

	fmt.Println("const:", examplerepo.FirstName)

}
