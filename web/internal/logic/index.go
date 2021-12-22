package logic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/pottava/jaguer-cn-lottery/web/internal/lib"
	"github.com/pottava/jaguer-cn-lottery/web/internal/logs"
)

func Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		status := http.StatusNotFound
		http.Error(w, http.StatusText(status), status)
		return
	}
	if strings.EqualFold(r.Method, http.MethodGet) {
		getIndex(w, r)
		return
	}
	if strings.EqualFold(r.Method, http.MethodPost) {
		postIndex(w, r)
		return
	}
	status := http.StatusMethodNotAllowed
	http.Error(w, http.StatusText(status), status)
}

type Swag struct {
	ID    int64
	Name  string
	Stock int64
}

func getIndex(w http.ResponseWriter, r *http.Request) {

	userID := "anonymous"
	if lib.Config.LogLevel != "debug" {
		if candidate, err := metadata.Email(""); err == nil {
			userID = candidate
		}
	}
	headers := &http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(userID, "")))

	status, bytes, err := lib.HTTPGet(
		r.Context(), http.DefaultClient,
		lib.Config.APIEndpoint+"/swags",
		headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error("access", err, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	if status != http.StatusOK {
		http.Error(w, http.StatusText(status), status)
		logs.Error("access", nil, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	swags := []Swag{}
	if err = json.Unmarshal(bytes, &swags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error("access", err, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	fmt.Fprintln(w, toHTML(swags))
	logs.Info("access", nil, &logs.Map{"method": r.Method, "path": r.URL.Path})
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func toHTML(swags []Swag) string {
	html := "<!DOCTYPE html><html><body style=\"margin: 25px 30px;\">"
	html += "<h3>Jagu'e'r クラウド ネイティブ分科会 Meetup #3</h3>"
	html += "<div>ノベルティを選択し、応募してください</div><br>"
	html += "<form action=\"/\" method=\"post\">"
	for idx, swag := range swags {
		checked := ""
		if idx == 0 {
			checked = "checked"
		}
		html += fmt.Sprintf("<input id=\"swag-%d\" type=\"radio\" name=\"swag\""+
			"value=\"%d\" %s><label for=\"swag-%d\">%s</label><br>",
			swag.ID, swag.ID, checked, swag.ID, swag.Name)
	}
	html += "<br><label>参加登録時メールアドレス<span style=\"color: red;font-weight: bold;\">（* 必須）</span></label>"
	html += "<br><input type=\"text\" name=\"user\" style=\"margin: 3px 0;width: 200px;\"><br>"
	html += "<br><button type=\"submit\">応募</button>"
	return html + "</form></body></html>"
}

func postIndex(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	userID := r.Form.Get("user")
	if lib.Config.LogLevel != "debug" {
		if candidate, err := metadata.Email(""); err == nil {
			userID = candidate
		}
	}
	if userID == "" {
		w.Header().Set("location", "/")
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json; charset=UTF-8")
	headers.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(userID, "")))

	data := fmt.Sprintf(`{"swag": %s, "user": "%s"}`, r.Form.Get("swag"), userID)

	status, bytes, err := lib.HTTPPost(
		r.Context(), http.DefaultClient,
		lib.Config.APIEndpoint+"/requests",
		bytes.NewBuffer([]byte(data)),
		headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error("access", err, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	if status != http.StatusCreated {
		http.Error(w, http.StatusText(status), status)
		logs.Error("access", nil, &logs.Map{"method": r.Method, "path": r.URL.Path})
		return
	}
	html := "<!DOCTYPE html><html><body style=\"margin: 25px 30px;\">"
	html += "<h3>Jagu'e'r クラウド ネイティブ分科会 Meetup #3</h3>"
	html += "<div>以下の通り申請をお受けしました</div><br><ul>"
	html += "<li>" + userID + "</li>"
	html += "<li>" + string(bytes) + "</li>"
	html += "</ul></body></html>"
	fmt.Fprint(w, html)

	logs.Info("access", nil, &logs.Map{"method": r.Method, "path": r.URL.Path})
}
