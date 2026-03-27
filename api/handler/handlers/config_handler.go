package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kickplate/api/lib"
)

type ConfigHandler struct {
	env lib.Env
}

func NewConfigHandler(env lib.Env) ConfigHandler {
	return ConfigHandler{env: env}
}

func (h ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.env.Customization)
}
