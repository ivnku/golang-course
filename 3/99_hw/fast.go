package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	var seenBrowsers []string
	uniqueBrowsers := 0
	foundUsers := ""

	reader := bufio.NewReader(file)

	users := make([]User, 0, 1000)
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			break
		}
		user := &User{}

		err = user.UnmarshalJSON(line)

		if err != nil {
			panic(err)
		}
		users = append(users, *user)
	}

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
			} else if ok = strings.Contains(browser, "MSIE"); ok {
				isMSIE = true
			} else {
				continue
			}

			notSeenBefore := true
			for _, item := range seenBrowsers {
				if item == browser {
					notSeenBefore = false
				}
			}
			if notSeenBefore {
				seenBrowsers = append(seenBrowsers, browser)
				uniqueBrowsers++
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, strings.Replace(user.Email, "@", " [at] ", 1))
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
