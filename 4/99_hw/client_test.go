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

type TestCase struct {
	Description string
	SearchRequest *SearchRequest
	Result *SearchResponse
	IsError bool
	ErrorMessage string
}

func initFixtures(token string, handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *SearchClient, SearchRequest) {
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

	return searchServer, &SearchClient{token, searchServer.URL}, searchRequest
}

// test limit and offset
func TestParameters(t *testing.T) {
	cases := []TestCase{
		{
			Description: "Test negative limit",
			SearchRequest: &SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    0,
			},
			Result: nil,
			IsError: true,
			ErrorMessage: "limit must be > 0",
		},
		{
			Description: "Test negative offset",
			SearchRequest: &SearchRequest{
				Limit:      2,
				Offset:     -1,
				Query:      "",
				OrderField: "",
				OrderBy:    0,
			},
			Result: nil,
			IsError: true,
			ErrorMessage: "offset must be > 0",
		},
		{
			Description: "Test limit > 25",
			SearchRequest: &SearchRequest{
				Limit:      30,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    0,
			},
			Result: &SearchResponse{
				Users:    make([]User, 25),
			},
			IsError: false,
		},
		{
			Description: "Test bad order field",
			SearchRequest: &SearchRequest{
				Limit:      5,
				Offset:     0,
				Query:      "",
				OrderField: "InvalidOrderField",
				OrderBy:    0,
			},
			Result: nil,
			IsError: true,
			ErrorMessage: "OrderFeld InvalidOrderField invalid",
		},
	}
	searchServer, searchClient, _ := initFixtures(token, nil)
	defer searchServer.Close()
	for caseNum, item := range cases {
		result, err := searchClient.FindUsers(*item.SearchRequest)

		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if  err != nil && item.ErrorMessage != "" && !strings.Contains(err.Error(), item.ErrorMessage) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if  item.Result != nil && len(result.Users) != len(item.Result.Users) {
			t.Errorf("[%d] the wrong amount of users!", caseNum)
		}
	}
}

func TestFindUsersByLastName(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, nil)
	defer searchServer.Close()

	_, err := searchClient.FindUsers(searchRequest)

	if err != nil {
		t.Error(err)
	}
}

func TestSearchServerFatalError(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	})
	defer searchServer.Close()

	_, err := searchClient.FindUsers(searchRequest)

	if err == nil || err.Error() != "SearchServer fatal error" {
		t.Error(err)
	}
}

func TestCantUnpackErrorJson(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("some wrong structured error message json"))
		return
	})
	defer searchServer.Close()

	_, err := searchClient.FindUsers(searchRequest)

	if err == nil || err.Error() != "cant unpack error json: invalid character 's' looking for beginning of value" {
		t.Error(err)
	}
}

func TestUnknownBadRequestError(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		errorJson, _ := json.Marshal(&SearchErrorResponse{"SomeUnknownErrorMessage"})
		_, _ = w.Write(errorJson)
		return
	})
	defer searchServer.Close()

	_, err := searchClient.FindUsers(searchRequest)

	if err == nil || err.Error() != "unknown bad request error: SomeUnknownErrorMessage" {
		t.Error(err)
	}
}

func TestCantUnpackResultJson(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, func(w http.ResponseWriter, r *http.Request) {
		resultJson, _ := json.Marshal("resultInvalidJson")
		_, _ = w.Write(resultJson)
		return
	})
	defer searchServer.Close()

	_, err := searchClient.FindUsers(searchRequest)

	if err == nil || !strings.Contains(err.Error(), "cant unpack result json") {
		t.Error(err)
	}
}

func TestDataLengthNotEqualLimit(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, func(w http.ResponseWriter, r *http.Request) {
		user := []User{{14, "UserName", 23, "about info", "male"}}
		resultJson, _ := json.Marshal(user)
		_, _ = w.Write(resultJson)
		return
	})
	defer searchServer.Close()

	res, err := searchClient.FindUsers(searchRequest)

	if err != nil || len(res.Users) != 1 {
		t.Error(err)
	}
}

func TestTimeoutError(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, nil)
	searchServer.Close()

	_, err := searchClient.FindUsers(searchRequest)

	if err == nil || !strings.Contains(err.Error(), "timeout for") {
		t.Error(err)
	}
}

func TestUnknownHttpRequestError(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures(token, nil)
	searchServer.Close()

	searchClient.URL = ""

	_, err := searchClient.FindUsers(searchRequest)

	if err == nil || !strings.Contains(err.Error(), "unknown error") {
		t.Error(err)
	}
}

func TestInvalidToken(t *testing.T) {
	searchServer, searchClient, searchRequest := initFixtures("InvalidToken", nil)
	defer searchServer.Close()

	_, err := searchClient.FindUsers(searchRequest)

	if err == nil || err.Error() != "Bad AccessToken" {
		t.Error(err)
	}
}
