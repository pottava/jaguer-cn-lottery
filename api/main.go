package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pottava/jaguer-cn-lottery/api/internal/lib"
	"github.com/pottava/jaguer-cn-lottery/api/internal/logic"
	"github.com/pottava/jaguer-cn-lottery/api/internal/logs"
)

func main() {
	if len(lib.Config.ProjectID) == 0 {
		logs.Fatal("Missing required environment variable: PROJECT_ID", nil, nil)
	}
	if len(lib.Config.SheetID) == 0 {
		logs.Fatal("Missing required environment variable: SPREAD_SHEET_ID", nil, nil)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, lib.Config.Version)
	})
	http.HandleFunc("/swags", logic.GetSwags)
	http.HandleFunc("/requests", logic.PostRequests)

	logs.Info("Server started", nil, &logs.Map{"Port": lib.Config.Port})
	log.Fatal(http.ListenAndServe(":"+lib.Config.Port, nil))
}
