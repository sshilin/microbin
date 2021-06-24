package routes

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbout(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	About("v.dev")(rr, req)
	data, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err)

	assert.Equal(t, "microbin v.dev", string(data))
}
