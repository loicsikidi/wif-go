package extract

import (
	"fmt"
	"testing"
)

func TestExtractIdentifier(t *testing.T) {
	tests := []struct {
		val      string
		expected string
	}{
		{
			val:      "projects/{project-id}/",
			expected: "{project-id}",
		},
		{
			val:      "projects/{project_id}/",
			expected: "{project_id}",
		},
		{
			val:      "projects/{project-id}",
			expected: "{project-id}",
		},
		{
			val:      "projects/{project.id}",
			expected: "",
		},
		{
			val:      "projects/{project/id}",
			expected: "",
		},
		{
			val:      "projects/0123456789",
			expected: "",
		},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			if out := Identifier(tc.val); out != tc.expected {
				t.Fatalf("Identifier(%s) = %s, expected %s", tc.val, out, tc.expected)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	source := "projects/_/buckets/acme-orders-aaa/objects/data_lake/orders/order_date=2019-11-03/aef87g87ae0876"

	tests := []struct {
		val      string
		expected string
	}{
		{
			val:      "/order_date={date}/",
			expected: "2019-11-03",
		},
		{
			val:      "buckets/{name}/",
			expected: "acme-orders-aaa",
		},
		{
			val:      "/orders/{empty}order_date",
			expected: "",
		},
		{
			val:      "{start}/objects/data_lake",
			expected: "projects/_/buckets/acme-orders-aaa",
		},
		{
			val:      "orders/{end}",
			expected: "order_date=2019-11-03/aef87g87ae0876",
		},
		{
			val:      "{all}",
			expected: source,
		},
		{
			val:      "/orders/{none}/order_date=",
			expected: "",
		},
		{
			val:      "/orders/order_date=2019-11-03/{id}/data_lake",
			expected: "",
		},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			if out, _ := extract(source, tc.val); out != tc.expected {
				t.Fatalf("extract(%s) = %s, expected %s", tc.val, out, tc.expected)
			}
		})
	}
}
