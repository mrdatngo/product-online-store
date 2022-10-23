package integrate

import (
	"encoding/json"
	"fmt"
	model "github.com/mrdatngo/gin-products/models"
	route "github.com/mrdatngo/gin-products/routes"
	assert "github.com/mrdatngo/gin-products/test"
	util "github.com/mrdatngo/gin-products/utils"
	"github.com/restuwahyu13/go-supertest/supertest"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router = route.SetupRouter()

func Test_ViewProducts(t *testing.T) {
	productID := int64(2)
	test := supertest.NewSuperTest(router, t)
	test.Get(fmt.Sprintf(`/api/v1/product/%v`, productID))
	test.Send(nil)
	test.Set("Content-Type", "application/json")
	test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {
		response := rr.Body.String()
		assert.Equals(t, http.StatusOK, rr.Code)
		assert.Equals(t, http.MethodGet, req.Method)
		type ProductData struct {
			Product model.EntityProduct
		}
		type ResponseData struct {
			Data ProductData `json:"data"`
			Meta util.Meta   `json:"meta"`
		}
		responseData := &ResponseData{}
		err := json.Unmarshal([]byte(response), responseData)
		assert.Ok(t, err)
		assert.Equals(t, productID, responseData.Data.Product.ID)
		assert.Equals(t, http.StatusOK, responseData.Meta.StatusCode)
		assert.Equals(t, "OK", responseData.Meta.Message)
	})
}
