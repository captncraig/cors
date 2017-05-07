package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func getReq(method string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, "http://foo.com", nil)
	r.Header.Set(requestMethodKey, "POST")
	return w, r
}

func TestDefault_NoOrigin(t *testing.T) {
	w, r := getReq("OPTIONS")
	Default().HandleRequest(w, r)
	if w.Header().Get(allowOriginKey) != "" {
		t.Fatal("Should not have origin header when no origin passed in")
	}
}

func TestDefault_Origin_Supplied(t *testing.T) {
	w, r := getReq("OPTIONS")
	r.Header.Set(originKey, "http://bar.com")
	Default().HandleRequest(w, r)
	if w.Header().Get(allowOriginKey) != "*" {
		t.Fatalf("Expect origin of \"*\". Got \"%s\".", w.Header().Get(allowOriginKey))
	}
	if w.Header().Get(varyKey) != "Origin" {
		t.Fatal("Must include Vary:Origin if allow origin header set to specific domain.")
	}
}

func TestDefault_Dissallowed_Origin(t *testing.T) {
	w, r := getReq("OPTIONS")
	r.Header.Set(originKey, "http://bar.com")
	c := Default()
	c.AllowedOrigins = []string{"http://blog.bar.com"}
	c.HandleRequest(w, r)
	if w.Header().Get(allowOriginKey) != "" {
		t.Fatal("Should not have origin header when no origin not explicitely allowed")
	}
}

func TestDefault_Methods(t *testing.T) {
	w, r := getReq("OPTIONS")
	r.Header.Set(originKey, "http://bar.com")
	c := Default()
	c.HandleRequest(w, r)
	if w.Header().Get(allowMethodsKey) != "POST, GET, OPTIONS, PUT, DELETE" {
		t.Fatal("Allow methods should be set")
	}
}
