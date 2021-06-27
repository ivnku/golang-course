package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type MyError struct {
	errorMessage string
}

func (me *MyError) Error() string {
	return me.errorMessage
}

type UserXML struct {
	Id        int    `xml:"id"`
	Name      string `xml:"-"`
	FirstName string `xml:"first_name" json:"-"`
	LastName  string `xml:"last_name" json:"-"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type Users struct {
	Users []UserXML `xml:"row"`
}

const token = "servertoken"

func SearchServer(w http.ResponseWriter, r *http.Request) {

	accessToken := r.Header.Get("AccessToken")

	if accessToken != token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	searchRequest, err := prepareParameters(*r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	foundUsers, err := getUsers(searchRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	serializedUsers, err := json.Marshal(foundUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(serializedUsers)
}

func getUsers(request SearchRequest) ([]UserXML, error) {
	bytesData, err := ioutil.ReadFile("dataset.xml")

	if err != nil {
		return nil, err
	}

	var users Users
	if err := xml.Unmarshal(bytesData, &users); err != nil {
		return nil, err
	}

	for i, user := range users.Users {
		users.Users[i].Name = user.FirstName + " " + user.LastName
	}

	resultUsersList := users.Users

	if request.Query != "" {
		resultUsersList = filterByQuery(request.Query, users.Users)
	}
	if request.OrderField != "" {
		order(request.OrderField, request.OrderBy, resultUsersList)
	}
	if request.Offset+request.Limit > len(resultUsersList) {
		if request.Offset >= len(resultUsersList) {
			return make([]UserXML, 0), nil
		}
		return resultUsersList[request.Offset:], nil
	}

	return resultUsersList[request.Offset:request.Limit], nil
}

func prepareParameters(r http.Request) (SearchRequest, error) {
	searchRequest := SearchRequest{}

	var err error
	if limit := r.URL.Query().Get("limit"); limit != "" {
		searchRequest.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return searchRequest, err
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		searchRequest.Offset, err = strconv.Atoi(offset)
		if err != nil {
			return searchRequest, err
		}
	}
	if query := r.URL.Query().Get("query"); query != "" {
		searchRequest.Query = query
	}
	if orderField := r.URL.Query().Get("order_field"); orderField != "" {
		if orderField != "Id" && orderField != "Age" && orderField != "Name" {
			return searchRequest, &MyError{errorMessage: "ErrorBadOrderField"}
		}
		searchRequest.OrderField = orderField
	}
	if orderBy := r.URL.Query().Get("order_by"); orderBy != "" {
		searchRequest.OrderBy, err = strconv.Atoi(orderBy)
		if err != nil || searchRequest.OrderBy < -1 || searchRequest.OrderBy > 1 {
			return searchRequest, err
		}
	}

	return searchRequest, nil
}

func filterByQuery(search string, users []UserXML) []UserXML {
	resultList := make([]UserXML, 0)
	for _, user := range users {
		if strings.Contains(user.Name, search) || strings.Contains(user.About, search) {
			resultList = append(resultList, user)
		}
	}
	return resultList
}

func order(field string, direction int, users []UserXML) {
	sort.Slice(users, func(i, j int) bool {
		switch field {
		case "Id":
			if direction == -1 {
				return users[i].Id > users[j].Id
			} else {
				return users[i].Id < users[j].Id
			}
		case "Age":
			if direction == -1 {
				return users[i].Age > users[j].Age
			} else {
				return users[i].Age < users[j].Age
			}
		case "Name":
			if direction == -1 {
				return (users[i].FirstName + " " + users[i].LastName) > (users[j].FirstName + " " + users[j].LastName)
			} else {
				return users[i].FirstName+" "+users[i].LastName < users[j].FirstName+" "+users[j].LastName
			}
		default:
			return users[i].Id > users[j].Id
		}
	})
}

// TESTS

func TestTmp(t *testing.T) {
	searchServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer searchServer.Close()

	searchClient := &SearchClient{"servertoken", searchServer.URL}
	searchRequest := SearchRequest{
		Limit:      3,
		Offset:     1,
		Query:      "Guerr",
		OrderField: "Age",
		OrderBy:    -1,
	}

	res, _ := searchClient.FindUsers(searchRequest)

	fmt.Printf("response is %v", *res)
}

func TestInvalidToken(t *testing.T) {
	//searchServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	//defer searchServer.Close()
	//
	//searchClient := &SearchClient{"InvalidToken", searchServer.URL}
	//searchRequest := SearchRequest{}
	//
	//_, err := searchClient.FindUsers(searchRequest)
	//
	//if err == nil {
	//	t.Error("No error for invalid token")
	//}
	//
	//if err.Error() != "Bad AccessToken" {
	//	t.Error("Invalid error message")
	//}
}
