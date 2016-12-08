package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/OpenBazaar/openbazaar-go/test"
	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	// Create, Read, Update
	runJSONAPIBlackboxTests(t, jsonAPIBlackboxTests{
		{"POST", "/ob/settings", settingsJSON, 200, "{}"},                      // Create
		{"GET", "/ob/settings", "", 200, settingsJSON},                         // Read
		{"POST", "/ob/settings", settingsJSON, 409, settingsAlreadyExistsJSON}, // Fail 2nd Create
		{"PUT", "/ob/settings", settingsUpdateJSON, 200, "{}"},                 // Update
		{"GET", "/ob/settings", "", 200, settingsUpdateJSON},                   // Read
		{"PUT", "/ob/settings", settingsUpdateJSON, 200, "{}"},                 // Update idempotency
		{"GET", "/ob/settings", "", 200, settingsUpdateJSON},                   // Read
	})

	// Invalid JSON
	runJSONAPIBlackboxTests(t, jsonAPIBlackboxTests{
		{"POST", "/ob/settings", settingsMalformedJSON, 400, settingsMalformedJSONErr},
	})

	// Invalid JSON
	runJSONAPIBlackboxTests(t, jsonAPIBlackboxTests{
		{"POST", "/ob/settings", settingsJSON, 200, "{}"},
		{"GET", "/ob/settings", "", 200, settingsJSON},
		{"PUT", "/ob/settings", settingsMalformedJSON, 400, settingsMalformedJSONErr},
	})
}

func TestProfile(t *testing.T) {
	// Create, Update
	runJSONAPIBlackboxTests(t, jsonAPIBlackboxTests{
		{"POST", "/ob/profile", profileJSON, 200, profileJSON},              // Create
		{"POST", "/ob/profile", profileJSON, 409, profileAlreadyExistsJSON}, // Fail recreating
		{"PUT", "/ob/profile", profileUpdateJSON, 200, profileUpdatedJSON},  // Update
		{"PUT", "/ob/profile", profileUpdatedJSON, 200, profileUpdatedJSON}, // Update idempotency
	})
}

func Test404(t *testing.T) {
	// Test undefined endpoints
	runJSONAPIBlackboxTests(t, jsonAPIBlackboxTests{
		{"GET", "/ob/a", "{}", 404, notFoundJSON},
		{"PUT", "/ob/a", "{}", 404, notFoundJSON},
		{"POST", "/ob/a", "{}", 404, notFoundJSON},
		{"PATCH", "/ob/a", "{}", 404, notFoundJSON},
		{"DELETE", "/ob/a", "{}", 404, notFoundJSON},
	})
}

//
// JSON API testing blackbox
//

// jsonAPIBlackboxTest is a test case to be run against the api blackbox
type jsonAPIBlackboxTest struct {
	method      string
	path        string
	requestBody string

	expectedResponseCode int
	expectedResponseBody string
}

// jsonAPIBlackboxTests is a slice of jsonAPIBlackboxTest
type jsonAPIBlackboxTests []jsonAPIBlackboxTest

// jsonAPIBlackbox is a testing oracle for the JSON API
type jsonAPIBlackbox struct {
	*testing.T
	handlerFunc http.HandlerFunc
}

// newJSONAPIBlackbox creates a new testing blackbox with a mock node
func newJSONAPIBlackbox(t *testing.T) jsonAPIBlackbox {
	// Create a test node, cookie, and config
	node := test.NewNode(t)
	var authCookie http.Cookie
	apiConfig, err := test.NewAPIConfig()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new jsonAPIHandler to test
	apiHandler, err := newJsonAPIHandler(node, authCookie, *apiConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Return a new blackbox
	return jsonAPIBlackbox{
		T:           t,
		handlerFunc: apiHandler.ServeHTTP,
	}
}

// request issues an http request directly to the blackbox handler
func (api *jsonAPIBlackbox) request(method string, path string, body string) *http.Response {
	// Create a JSON request to the given endpoint
	req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	if err != nil {
		api.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Execute the call and return the result and the response body
	api.handlerFunc(rr, req)

	return rr.Result()
}

//
// Test runners
//

func runJSONAPIBlackboxTest(t *testing.T, blackbox jsonAPIBlackbox, test jsonAPIBlackboxTest) {
	// Make the request and read the response
	resp := blackbox.request(test.method, test.path, test.requestBody)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		blackbox.Fatal(err)
	}

	// Perform assertions
	assert.Equal(t, test.expectedResponseCode, resp.StatusCode)
	assertJSONEqual(t, test.expectedResponseBody, string(respBody))
}

func runJSONAPIBlackboxTests(t *testing.T, tests jsonAPIBlackboxTests) {
	blackbox := newJSONAPIBlackbox(t)
	for _, test := range tests {
		runJSONAPIBlackboxTest(t, blackbox, test)
	}
}

//
// Assertions
//

// stripWhiteSpace removes all whitespace characters from a string
var stripWhiteSpace = func(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// assertJSONEqual asserts that two given strings represent equal JSON objects
func assertJSONEqual(t *testing.T, expected string, actual string) bool {
	var aval interface{}
	var eval interface{}

	err := json.Unmarshal([]byte(actual), &aval)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal([]byte(expected), &eval)
	if err != nil {
		t.Fatal(err)
	}

	equal := assert.True(t, reflect.DeepEqual(aval, eval))
	if !equal {
		fmt.Println("expected:", expected)
		fmt.Println("actual:", actual)
	}

	return equal
}
