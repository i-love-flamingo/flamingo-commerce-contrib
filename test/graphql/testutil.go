package graphql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"unicode"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const (
	replacementScale = 2
)

func LoadGraphQL(t *testing.T, name string, replacements map[string]string) string {
	t.Helper()

	content, err := os.ReadFile(path.Join("testdata", name+".graphql"))
	if err != nil {
		t.Fatal(err)
	}

	r := make([]string, replacementScale*len(replacements))
	i := 0

	for key, val := range replacements {
		r[i] = fmt.Sprintf("###%s###", key)
		r[i+1] = val
		i += 2
	}

	replacer := strings.NewReplacer(r...)

	return replacer.Replace(string(content))
}

// PrepareCart adds a simple product via graphQl
func PrepareCart(t *testing.T, e *httpexpect.Expect) {
	t.Helper()
	helper.GraphQlRequest(t, e, LoadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": "fake_simple", "DELIVERY_CODE": "delivery"})).Expect().Status(http.StatusOK)
}

// PrepareCartWithPaymentSelection adds a simple product via graphQl
func PrepareCartWithPaymentSelection(t *testing.T, e *httpexpect.Expect, paymentMethod string, marketPlaceCode *string) {
	t.Helper()

	code := "fake_simple"
	if marketPlaceCode != nil {
		code = *marketPlaceCode
	}

	helper.GraphQlRequest(t, e, LoadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": code, "DELIVERY_CODE": "delivery"})).Expect().Status(http.StatusOK)
	helper.GraphQlRequest(t, e, LoadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": paymentMethod})).Expect().Status(http.StatusOK)
}

func GetValue(response *httpexpect.Response, queryName, key string) *httpexpect.Value {
	return response.JSON().Object().Value("data").Object().Value(queryName).Object().Value(key)
}

func GetArray(response *httpexpect.Response, queryName string) *httpexpect.Array {
	return response.JSON().Object().Value("data").Object().Value(queryName).Array()
}

func AssertResponseForExpectedState(t *testing.T, response *httpexpect.Response, expectedState map[string]interface{}) {
	t.Helper()

	data := make(map[string]interface{})
	require.NoError(t, json.Unmarshal([]byte(response.Body().Raw()), &data))

	var theData interface{}
	var ok bool

	if theData, ok = data["data"]; !ok || theData == nil {
		t.Fatalf("no data in response: %s", response.Body().Raw())
		return
	}

	data = theData.(map[string]interface{})

	if diff := cmp.Diff(data, expectedState); diff != "" {
		t.Errorf("diff mismatch (-want +got):\\n%s", diff)
	}
}

// SpaceMap strips all whitespace from given string
func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
