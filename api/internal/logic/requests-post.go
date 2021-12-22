package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	googlecloud "github.com/pottava/jaguer-cn-lottery/api/internal/google-cloud"
	"github.com/pottava/jaguer-cn-lottery/api/internal/lib"
	"github.com/pottava/jaguer-cn-lottery/api/internal/logs"
)

type SwagRequest struct {
	SwagID int64  `json:"swag"`
	UserID string `json:"user"`
}

// PostRequests SWAG の申請を受け付けます
func PostRequests(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, http.MethodPost) {
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
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}
	var swagReq SwagRequest
	if err := json.NewDecoder(r.Body).Decode(&swagReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	swagName, found := checkStock(r.Context(), swagReq)
	if !found {
		status := http.StatusNotFound
		http.Error(w, http.StatusText(status), status)
		return
	}
	if err := googlecloud.UpdateSheetCell(
		r.Context(),
		lib.Config.SheetID,
		lib.Config.SheetTabName,
		[]interface{}{
			swagReq.UserID,
			swagName,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error("access", err, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	logs.Info("access", nil, &logs.Map{"method": r.Method, "path": r.URL.Path})
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, *swagName)
}

// 本アプリケーションはサンプルであり "在庫を減らす処理がない" ため、ロックは実装していません
func checkStock(ctx context.Context, swagReq SwagRequest) (*string, bool) {
	swags, err := googlecloud.ListSwags(
		ctx,
		lib.Config.ProjectID,
		lib.Config.SpannerInstance,
		lib.Config.SpannerDatabase)
	if err != nil {
		return nil, false
	}
	for _, swag := range swags {
		if swag.ID == swagReq.SwagID {
			return &swag.Name, swag.Stock > 0
		}
	}
	return nil, false
}
