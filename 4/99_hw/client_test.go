package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
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
		errorJson, _ := json.Marshal(&SearchErrorResponse{err.Error()})
		_, _ = w.Write(errorJson)
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
			return searchRequest, &MyError{"ErrorBadOrderField"}
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
				return users[i].Id < users[j].Id
			} else {
				return users[i].Id > users[j].Id
			}
		case "Age":
			if direction == -1 {
				return users[i].Age < users[j].Age
			} else {
				return users[i].Age > users[j].Age
			}
		case "Name":
			if direction == -1 {
				return (users[i].FirstName + " " + users[i].LastName) < (users[j].FirstName + " " + users[j].LastName)
			} else {
				return users[i].FirstName+" "+users[i].LastName > users[j].FirstName+" "+users[j].LastName
			}
		default:
			return users[i].Id > users[j].Id
		}
	})
}

// TESTS

type TestCase struct {
	Description   string
	SearchRequest *SearchRequest
	Result        *SearchResponse
	IsError       bool
	ErrorMessage  string
	Handler       func(w http.ResponseWriter, r *http.Request)
	Token         string
	Url           *URL
}

type URL struct {
	url string
}

func initFixtures(token string, handler func(w http.ResponseWriter, r *http.Request), url *URL) (*httptest.Server, *SearchClient, SearchRequest) {
	if handler == nil {
		handler = SearchServer
	}
	searchServer := httptest.NewServer(http.HandlerFunc(handler))

	searchRequest := SearchRequest{
		Limit:      10,
		Offset:     0,
		Query:      "",
		OrderField: "",
		OrderBy:    0,
	}

	var urlAddress string
	if url == nil {
		urlAddress = searchServer.URL
	} else {
		urlAddress = url.url
	}

	return searchServer, &SearchClient{token, urlAddress}, searchRequest
}

// test different search parameters
func TestParameters(t *testing.T) {
	cases := []TestCase{
		{
			Description:   "Test negative limit",
			SearchRequest: &SearchRequest{Limit: -1},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "limit must be > 0",
		},
		{
			Description:   "Test negative offset",
			SearchRequest: &SearchRequest{Offset: -1},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "offset must be > 0",
		},
		{
			Description:   "Test limit > 25",
			SearchRequest: &SearchRequest{Limit: 30},
			Result: &SearchResponse{
				Users: make([]User, 25),
			},
			IsError: false,
		},
		{
			Description:   "Test bad order field",
			SearchRequest: &SearchRequest{OrderField: "InvalidOrderField"},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "OrderFeld InvalidOrderField invalid",
		},
		{
			Description:   "Test data length not equal limit",
			SearchRequest: &SearchRequest{Limit: 10},
			Result: &SearchResponse{
				Users: make([]User, 1),
			},
			IsError: false,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				user := []User{{14, "UserName", 23, "about info", "male"}}
				resultJson, _ := json.Marshal(user)
				_, _ = w.Write(resultJson)
				return
			},
		},
	}

	for caseNum, item := range cases {
		searchServer, searchClient, _ := initFixtures(token, item.Handler, item.Url)
		result, err := searchClient.FindUsers(*item.SearchRequest)

		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if err != nil && item.ErrorMessage != "" && !strings.Contains(err.Error(), item.ErrorMessage) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if item.Result != nil && len(result.Users) != len(item.Result.Users) {
			t.Errorf("[%d] the wrong amount of users!", caseNum)
		}
		searchServer.Close()
	}
}

func TestErrors(t *testing.T) {
	cases := []TestCase{
		{
			Description:   "Test SearchServerFatalError",
			SearchRequest: &SearchRequest{},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "SearchServer fatal error",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			},
			Token: token,
		},
		{
			Description:   "Test cant unpack error's json",
			SearchRequest: &SearchRequest{},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "cant unpack error json: invalid character",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("some wrong structured error message json"))
				return
			},
			Token: token,
		},
		{
			Description:   "Test unknown bad request error",
			SearchRequest: &SearchRequest{},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "unknown bad request error: SomeUnknownErrorMessage",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				errorJson, _ := json.Marshal(&SearchErrorResponse{"SomeUnknownErrorMessage"})
				_, _ = w.Write(errorJson)
				return
			},
			Token: token,
		},
		{
			Description:   "Test cant unpack result json",
			SearchRequest: &SearchRequest{},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "cant unpack result json",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				resultJson, _ := json.Marshal("resultInvalidJson")
				_, _ = w.Write(resultJson)
				return
			},
			Token: token,
		},
		{
			Description:   "Test timeout error",
			SearchRequest: &SearchRequest{},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "timeout for",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
				return
			},
			Token: token,
		},
		{
			Description:   "Test invalid token",
			SearchRequest: &SearchRequest{},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "Bad AccessToken",
			Token:         "InvalidToken",
		},
		{
			Description:   "Test UnknownHttpRequestError",
			SearchRequest: &SearchRequest{},
			Result:        nil,
			IsError:       true,
			ErrorMessage:  "unknown error",
			Token:         token,
			Url:           &URL{url: ""},
		},
	}

	for caseNum, item := range cases {
		searchServer, searchClient, _ := initFixtures(item.Token, item.Handler, item.Url)
		_, err := searchClient.FindUsers(*item.SearchRequest)

		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if err != nil && item.ErrorMessage != "" && !strings.Contains(err.Error(), item.ErrorMessage) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}

		searchServer.Close()
	}
}

// TEST LOGIC

func TestFindUsersByQuery(t *testing.T) {
	cases := []TestCase{
		{
			Description:   "Test search user by 'Name' field",
			SearchRequest: &SearchRequest{Limit: 1, Query: "Hilda Mayer"},
			Result:        &SearchResponse{Users: []User{{Id: 1, Name: "Hilda Mayer"}}},
			IsError:       false,
			Token:         token,
		},
		{
			Description:   "Test search user by 'About' field",
			SearchRequest: &SearchRequest{Limit: 1, Query: "Velit ullamco est aliqua voluptate nisi do"},
			Result:        &SearchResponse{Users: []User{{Id: 1, About: "Velit ullamco est aliqua voluptate nisi do"}}},
			IsError:       false,
			Token:         token,
		},
	}

	for caseNum, item := range cases {
		searchServer, searchClient, _ := initFixtures(item.Token, item.Handler, item.Url)
		res, err := searchClient.FindUsers(*item.SearchRequest)

		if err != nil {
			t.Errorf("[%d] unexpected error!", caseNum)
		}
		if res == nil {
			t.Errorf("[%d] result is nil, value expected!", caseNum)
		}
		if len(res.Users) == 0 {
			t.Errorf("[%d] result is empty, value expected!", caseNum)
		}
		if item.Result.Users[0].Name != res.Users[0].Name && !strings.Contains(res.Users[0].About, item.Result.Users[0].About) {
			t.Errorf("[%d] unexpected value: %#v", caseNum, res.Users[0])
		}

		searchServer.Close()
	}
}

func TestOrder(t *testing.T) {
	cases := []TestCase{
		{
			Description:   "Test order by name DESC",
			SearchRequest: &SearchRequest{Limit: 3, Query: "err", OrderField: "Name", OrderBy: OrderByDesc},
			Result: &SearchResponse{Users: []User{
				{Id: 18, Name: "Terrel Hall"},
				{Id: 11, Name: "Gilmore Guerra"},
				{Id: 12, Name: "Cruz Guerrero"},
			}},
			Token:   token,
		},
		{
			Description:   "Test order by name ASC",
			SearchRequest: &SearchRequest{Limit: 3, Query: "err", OrderField: "Name", OrderBy: OrderByAsc},
			Result: &SearchResponse{Users: []User{
				{Id: 12, Name: "Cruz Guerrero"},
				{Id: 11, Name: "Gilmore Guerra"},
				{Id: 18, Name: "Terrel Hall"},
			}},
			Token:   token,
		},
		{
			Description:   "Test order by id DESC",
			SearchRequest: &SearchRequest{Limit: 3, Query: "err", OrderField: "Id", OrderBy: OrderByDesc},
			Result: &SearchResponse{Users: []User{
				{Id: 18, Name: "Terrel Hall"},
				{Id: 12, Name: "Cruz Guerrero"},
				{Id: 11, Name: "Gilmore Guerra"},
			}},
			Token:   token,
		},
		{
			Description:   "Test order by id ASC",
			SearchRequest: &SearchRequest{Limit: 3, Query: "err", OrderField: "Id", OrderBy: OrderByAsc},
			Result: &SearchResponse{Users: []User{
				{Id: 11, Name: "Gilmore Guerra"},
				{Id: 12, Name: "Cruz Guerrero"},
				{Id: 18, Name: "Terrel Hall"},
			}},
			Token:   token,
		},
		{
			Description:   "Test order by age DESC",
			SearchRequest: &SearchRequest{Limit: 3, Query: "err", OrderField: "Age", OrderBy: OrderByDesc},
			Result: &SearchResponse{Users: []User{
				{Id: 12, Name: "Cruz Guerrero", Age: 36},
				{Id: 11, Name: "Gilmore Guerra", Age: 32},
				{Id: 18, Name: "Terrel Hall", Age: 27},
			}},
			Token:   token,
		},
		{
			Description:   "Test order by age ASC",
			SearchRequest: &SearchRequest{Limit: 3, Query: "err", OrderField: "Age", OrderBy: OrderByAsc},
			Result: &SearchResponse{Users: []User{
				{Id: 18, Name: "Terrel Hall", Age: 27},
				{Id: 11, Name: "Gilmore Guerra", Age: 32},
				{Id: 12, Name: "Cruz Guerrero", Age: 36},
			}},
			Token:   token,
		},
	}

	for caseNum, item := range cases {
		searchServer, searchClient, _ := initFixtures(item.Token, item.Handler, item.Url)
		res, err := searchClient.FindUsers(*item.SearchRequest)

		if err != nil {
			t.Errorf("[%d] unexpected error!", caseNum)
		}
		if res == nil {
			t.Errorf("[%d] result is nil, value expected!", caseNum)
		}
		if len(res.Users) == 0 {
			t.Errorf("[%d] result is empty, value expected!", caseNum)
		}
		if len(item.Result.Users) == len(res.Users) {
			for i, user := range res.Users {
				if item.Result.Users[i].Id != user.Id {
					t.Errorf("[%d] The order is wrong! %v", caseNum, res.Users)
					break
				}
			}
		} else {
			t.Errorf("[%d] length of result is not equal length of testcase!", caseNum)
		}

		searchServer.Close()
	}
}

func TestOffset(t *testing.T) {
	expected := &SearchResponse{Users: []User{
		{Id: 2, Name: "Brooks Aguilar"},
		{Id: 3, Name: "Everett Dillard"},
	}}

	searchServer, searchClient, _ := initFixtures(token, nil, nil)
	defer searchServer.Close()

	res, err := searchClient.FindUsers(SearchRequest{Limit: 2, Offset: 2})

	if err != nil {
		t.Error("Unexpected error!")
	}

	if res.Users[0].Id != expected.Users[0].Id && res.Users[1].Id != expected.Users[1].Id {
		t.Error("Wrong values of offset!")
	}
}