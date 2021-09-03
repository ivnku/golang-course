package helpers

import (
	"encoding/json"
	"net/http"
)

func SerializeAndReturn(w http.ResponseWriter, data interface{}) {
	serializedEntity, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(serializedEntity)
}
