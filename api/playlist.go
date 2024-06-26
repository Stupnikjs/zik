package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Stupnikjs/zik/util"
)

func (app *Application) CreatePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	reqJson := JsonReq{}
	bytes, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(bytes, &reqJson)

	if err != nil {
		util.WriteErrorToResponse(w, err, http.StatusBadRequest)
	}

	body, ok := reqJson.Object.Body.(string)
	if reqJson.Action == "create" && reqJson.Object.Type == "playlist" && ok {
	err =	app.DB.InsertPlaylist(body)

}

func (app *Application) AppendToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	// call to app

}

func (app *Application) RemoveToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	// call to app

}
