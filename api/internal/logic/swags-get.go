package logic

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	googlecloud "github.com/pottava/jaguer-cn-lottery/api/internal/google-cloud"
	"github.com/pottava/jaguer-cn-lottery/api/internal/lib"
	"github.com/pottava/jaguer-cn-lottery/api/internal/logs"
)

// GetSwags SWAG 一覧を返します
func GetSwags(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, http.MethodGet) {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}
	userID := retrieveUserID(r.Header.Get("Authorization"))
	if userID == "" {
		status := http.StatusUnauthorized
		http.Error(w, http.StatusText(status), status)
		return
	}
	swags, err := googlecloud.ListSwags(
		r.Context(),
		lib.Config.ProjectID,
		lib.Config.SpannerInstance,
		lib.Config.SpannerDatabase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error("access", err, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	bytes, err := json.Marshal(swags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error("access", err, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	fmt.Fprintln(w, string(bytes))
	logs.Info("access", nil, &logs.Map{"method": r.Method, "path": r.URL.Path})
}

func retrieveUserID(header string) string {
	auth := strings.SplitN(header, " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return ""
	}
	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return ""
	}
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) < 1 {
		return ""
	}
	return pair[0]
}
