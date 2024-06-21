package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) Routes() http.Handler {

	mux := chi.NewRouter()

	// register routes
	mux.Get("/", app.RenderAccueil)
	mux.Get("/loader", app.RenderDragDrop)

	mux.Get("/urls", app.ListObjectHandler)

	mux.Post("/upload", app.UploadFile)
	mux.Get("/stream/sound/{id}", app.UploadTrackFromGCPHandler)
	mux.Get("/youtube/mp3/{id}", app.YoutubeToGCPHandler)

	mux.Post("/api/track/tag", app.UpdateTrackTagHandler)
	mux.Post("/api/track/remove", app.DeleteTrackHandler)

	mux.Post("/api/playlist/create", app.CreatePlaylistHandler)
	mux.Post("/api/playlist/append", app.AppendToPlaylistHandler)
	mux.Post("/api/playlist/remove", app.RemoveToPlaylistHandler)

	// static assets

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}