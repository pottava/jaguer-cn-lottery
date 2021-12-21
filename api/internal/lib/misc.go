package lib

import (
	"net/http"
	"strings"
)

func Header(r *http.Request, key string) (values []string, found bool) {
	if r.Header == nil {
		return
	}
	for k, v := range r.Header {
		if strings.EqualFold(k, key) && len(v) > 0 {
			return v, true
		}
	}
	return
}
