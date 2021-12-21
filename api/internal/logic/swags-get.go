package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	gc "github.com/pottava/jaguer-cn-lottery/api/internal/google-cloud"
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
	swags, err := gc.ListSwags(
		context.Background(),
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
