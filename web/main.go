package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pottava/jaguer-cn-lottery/web/internal/lib"
	"github.com/pottava/jaguer-cn-lottery/web/internal/logic"
	"github.com/pottava/jaguer-cn-lottery/web/internal/logs"
)

func main() {
	if len(lib.Config.APIEndpoint) == 0 {
		logs.Fatal("Missing required environment variable: API_ENDPOINT", nil, nil)
	}
	http.HandleFunc("/", logic.Index)
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, lib.Config.Version)
	})
	logs.Info("Server started", nil, &logs.Map{"Port": lib.Config.Port})
	log.Fatal(http.ListenAndServe(":"+lib.Config.Port, nil))
}
