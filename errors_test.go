package httperrors

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertRespCode is a test helper to assert HTTP response codes.
func assertRespCode(t *testing.T, actual, expected int) {
	if actual != expected {
		t.Fatalf(
			"HTTP response code was %d and expected %d\n",
			actual,
			expected,
		)
	}
}

// assertJSONBytes is a test helper to assert two JSON byte slices are equal.
func assertJSONBytes(t *testing.T, actual, expected []byte) {
	var j, j2 interface{}
	if err := json.Unmarshal(actual, &j); err != nil {
		t.Fatal("Could not unmarshal actual response")
	}
	if err := json.Unmarshal(expected, &j2); err != nil {
		t.Fatal("Could not unmarshal actual response")
	}
	equal := reflect.DeepEqual(j2, j)
	if !equal {
		t.Fatal("The expected and actual JSON responses are not equal")
	}
}

func TestWriteErrorInvalidtokenPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	w := httptest.NewRecorder()
	WriteError(w, -1, "blah")
}

func TestWriteErrorMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, http.StatusBadRequest, "Parameters were missing")
	expectedResponse := []byte(`{
		"code":   400,
		"message": "Bad Request",
		"reason":  "Parameters were missing"
	}`)
	assertRespCode(t, w.Code, http.StatusBadRequest)
	assertJSONBytes(t, w.Body.Bytes(), expectedResponse)
}

func TestUnableToMarshalJSON(t *testing.T) {
	originalMarshalJSON := marshalJSON
	defer func() { marshalJSON = originalMarshalJSON }()
	marshalJSON = func(interface{}) ([]byte, error) { return []byte{}, fmt.Errorf("TEST_ERROR") }
	w := httptest.NewRecorder()
	WriteError(w, http.StatusNotFound, "Entity Not Found")
	assertRespCode(t, w.Code, 404)
}

func TestErrorList(t *testing.T) {
	errs := []error{
		goerrors.New("first"),
		goerrors.New("second"),
		goerrors.New("third"),
	}
	errlist := New(errs)
	assert.Equal(t, "errors: [first second third]", errlist.Error())
}
