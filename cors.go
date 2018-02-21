package cors

import (
	"fmt"
	"net/http"
)

const (
	allowOriginKey      string = "Access-Control-Allow-Origin"
	allowCredentialsKey        = "Access-Control-Allow-Credentials"
	allowHeadersKey            = "Access-Control-Allow-Headers"
	allowMethodsKey            = "Access-Control-Allow-Methods"
	maxAgeKey                  = "Access-Control-Max-Age"

	originKey         = "Origin"
	varyKey           = "Vary"
	requestMethodKey  = "Access-Control-Request-Method"
	requestHeadersKey = "Access-Control-Request-Headers"
	exposeHeadersKey  = "Access-Control-Expose-Headers"
	options           = "OPTIONS"
)

type Config struct {
	AllowedOrigins   []string
	AllowedMethods   string
	AllowedHeaders   string
	ExposedHeaders   string
	AllowCredentials *bool
	MaxAge           int
}

func Default() *Config {
	return &Config{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   "POST, GET, OPTIONS, PUT, DELETE",
		AllowedHeaders:   "",
		ExposedHeaders:   "",
		MaxAge:           0,
		AllowCredentials: nil,
	}
}

// Read the request, setting response headers as appropriate.
// Will NOT write anything to response in any circumstances.
func (c *Config) HandleRequest(w http.ResponseWriter, r *http.Request) {
	requestOrigin := r.Header.Get(originKey)
	if requestOrigin == "" {
		return
	}
	//check origin against allowed origins
	for _, ao := range c.AllowedOrigins {
		if ao == "*" || ao == requestOrigin {
			responseOrigin := "*"
			if ao != "*" {
				responseOrigin = requestOrigin
			}
			w.Header().Set(allowOriginKey, responseOrigin)
			w.Header().Add(varyKey, originKey)
			break
		}
	}

	if w.Header().Get(allowOriginKey) == "" {
		return //if we didn't set a valid allow-origin, none of the other headers matter
	}

	if IsPreflight(r) {
		w.Header().Set(allowMethodsKey, c.AllowedMethods)
		if c.AllowedHeaders != "" {
			if c.AllowedHeaders != "*" {
				w.Header().Set(allowHeadersKey, c.AllowedHeaders)
			} else {
				w.Header().Set(allowHeadersKey, r.Header.Get(requestHeadersKey))
			}

		}
		if c.MaxAge > 0 {
			w.Header().Set(maxAgeKey, fmt.Sprint(c.MaxAge))
		}
	} else {
		//regular request
		if c.ExposedHeaders != "" {
			w.Header().Set(exposeHeadersKey, c.ExposedHeaders)
		}
	}

	if c.AllowCredentials != nil {
		w.Header().Set(allowCredentialsKey, fmt.Sprint(*c.AllowCredentials))
	}

}

func IsPreflight(r *http.Request) bool {
	return r.Method == options && r.Header.Get(requestMethodKey) != ""
}
