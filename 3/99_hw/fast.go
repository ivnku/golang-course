package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	var seenBrowsers []string
	uniqueBrowsers := 0
	foundUsers := ""

	scanner := bufio.NewScanner(file)

	users := make([]map[string]interface{}, 0)
	for scanner.Scan() {
		line := scanner.Text()
		user := make(map[string]interface{}, 10)
		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				continue
			}
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
		email := strings.Replace(user["email"].(string), "@", " [at] ", 1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
