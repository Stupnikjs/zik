package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"

	repo "github.com/Stupnikjs/zik/pkg/db"
	"github.com/Stupnikjs/zik/pkg/gstore"
	"github.com/go-chi/chi/v5"
)

var pathToTemplates = "./static/templates/"

type TemplateData struct {
	Data map[string]any
}

func render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {

	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "/base.layout.gohtml"))
	if err != nil {
		return err
	}
	err = parsedTemplate.Execute(w, td)
	if err != nil {
		return err
	}
	return nil

}

// template rendering

func (app *Application) RenderAccueil(w http.ResponseWriter, r *http.Request) {

	td := TemplateData{}
	tracks := app.DB.GetAllTracks()
	fmt.Println(tracks)
	td.Data = make(map[string]any)
	td.Data["Tracks"] = tracks
	_ = render(w, r, "/acceuil.gohtml", &td)
}

func (app *Application) RenderLoader(w http.ResponseWriter, r *http.Request) {

	_ = render(w, r, "/loader.gohtml", &TemplateData{})
}

func (app *Application) RenderSingleTrack(w http.ResponseWriter, r *http.Request) {
	trackid := chi.URLParam(r, "id")
	td := TemplateData{}
	td.Data["Track"] = app.DB.GetTrackFromId(trackid)

	_ = render(w, r, "/singletrack.gohtml", &td)
}

func (app *Application) UploadFile(w http.ResponseWriter, r *http.Request) {
	// load file to gcp bucket

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	err = app.LoadMultipartReqToBucket(r, BucketName)
	if err != nil {
		fmt.Println(err)
	}
}

func ListObjectHandler(w http.ResponseWriter, r *http.Request) {
	names, err := gstore.ListObjectsBucket(BucketName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	byteNames, err := json.Marshal(names)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Write(byteNames)

}

func (app *Application) LoadMultipartReqToBucket(r *http.Request, bucketName string) error {
	objNames, err := gstore.ListObjectsBucket(BucketName)

	if err != nil {
		return err
	}

	for _, headers := range r.MultipartForm.File {

		for _, h := range headers {

			if IsInSlice[string](h.Filename, objNames) {
				// format a messgage with already present files
				break
			}

			file, err := h.Open()

			if err != nil {
				return err
			}

			defer file.Close()

			finalByteArr, err := ByteFromMegaFile(file)

			if err != nil {
				return err
			}

			err = gstore.LoadToBucket(bucketName, h.Filename, finalByteArr)

			if err != nil {
				return err
			}

			track := repo.Track{}
			url, err := gstore.GetObjectURL(bucketName, h.Filename)
			track.StoreURL = url
			track.Name = h.Filename
			err = app.DB.PushTrackToSQL(track)
			if err != nil {
				return err
			}
		}

	}
	return nil

}

func (app *Application) UploadTrackFromGCPHandler(w http.ResponseWriter, r *http.Request) {
	trackid := chi.URLParam(r, "id")
	track := app.DB.GetTrackFromId(trackid)
	resp, err := http.Get(track.StoreURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusOK)

	_, _ = io.Copy(w, resp.Body)

}
