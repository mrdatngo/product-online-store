package unit

import (
	product_controllers "github.com/mrdatngo/gin-products/controllers/product-controllers"
	testpackage "github.com/mrdatngo/gin-products/test"
	util "github.com/mrdatngo/gin-products/utils"
	"testing"
)

func Test_GetMessageFromInputGetProduct(t *testing.T) {
	type InputTest struct {
		userId   int64
		input    *product_controllers.InputGetProduct
		expected string
	}
	var inputTests = []InputTest{
		{
			userId: 0,
			input: &product_controllers.InputGetProduct{
				ProductID: 1,
			},
			expected: `{"user_id":0,"event_type":"VIEWING","detail_logs":[{"event_type":"NUMBER","param":"PRODUCT_ID","value":"1"}]}`,
		},
		{
			userId: 10,
			input: &product_controllers.InputGetProduct{
				ProductID: 100000,
			},
			expected: `{"user_id":10,"event_type":"VIEWING","detail_logs":[{"event_type":"NUMBER","param":"PRODUCT_ID","value":"100000"}]}`,
		},
	}

	for _, test := range inputTests {
		got := util.GetMessageFromInputGetProduct(test.userId, test.input)
		testpackage.Equals(t, test.expected, got)
	}
}
