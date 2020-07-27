package api

import "net/http"

func handleStatus(w http.ResponseWriter, request *http.Request) {
	RenderJson(w, map[string]bool{"ok": true})
}
