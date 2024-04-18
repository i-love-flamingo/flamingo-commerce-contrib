//go:build integration
// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce-contrib/test/graphql"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CartSummary(t *testing.T) {
	t.Parallel()

	baseURL := "http://" + FlamingoURL
	expect := integrationtest.NewHTTPExpect(t, baseURL)

	marketPlaceCode := "fake_simple_with_fixed_price"
	graphql.PrepareCartWithPaymentSelection(t, expect, "creditcard", &marketPlaceCode)

	response := helper.GraphQlRequest(t, expect, graphql.LoadGraphQL(t, "cart_summary", map[string]string{"METHOD": "creditcard"})).Expect().Status(http.StatusOK)
	response.Status(http.StatusOK)

	graphql.AssertResponseForExpectedState(t, response, map[string]interface{}{
		"Commerce_Cart_DecoratedCart": map[string]interface{}{
			"cartSummary": map[string]interface{}{
				"total": map[string]interface{}{
					"amount":   10.49,
					"currency": "EUR",
				},
			},
		},
	})
}
