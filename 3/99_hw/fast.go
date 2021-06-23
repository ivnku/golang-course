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

	seenBrowsers := make(map[string]bool, 100)

	reader := bufio.NewReader(file)

	fmt.Fprintf(out, "found users:\n")

	i := -1
	for {
		i++
		line, err := reader.ReadSlice('\n')
		if err != nil {
			break
		}
		user := &User{}

		err = user.UnmarshalJSON(line)

		if err != nil {
			panic(err)
		}

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

			_, isSeenBefore := seenBrowsers[browser]

			if !isSeenBefore {
				seenBrowsers[browser] = true
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, strings.Replace(user.Email, "@", " [at] ", 1))
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
