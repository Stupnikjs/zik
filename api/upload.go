package api

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Stupnikjs/zik/gstore"
	"github.com/Stupnikjs/zik/repo"
	"github.com/Stupnikjs/zik/util"
)

// Compare speed with io.ReadAll()
func ByteFromMegaFile(file io.Reader) ([]byte, error) {

	reader := bufio.NewReader(file)
	finalByteArr := make([]byte, 0, 2048*1000)
	for {
		soloByte, err := reader.ReadByte()
		if err != nil {
			log.Println(err)
			break
		}
		finalByteArr = append(finalByteArr, soloByte)
	}
	return finalByteArr, nil
}

func (app *Application) LoadMultipartReqToBucket(r *http.Request, bucketName string) (string, error) {
	objNames, err := gstore.ListObjectsBucket(app.BucketName)

	// filename with already present storage object names
	already := []string{}
	if err != nil {
		return "", err
	}

	for _, headers := range r.MultipartForm.File {

		for _, h := range headers {

			if util.IsInSlice(h.Filename, objNames) {
				already = append(already, h.Filename)
				break
			}

			file, err := h.Open()
			if err != nil {
				return "", err
			}
			defer file.Close()

			finalByteArr, err := ByteFromMegaFile(file)
			if err != nil {
				return "", err
			}

			err = gstore.LoadToBucket(bucketName, h.Filename, finalByteArr)
			if err != nil {
				return "", err
			}

			url, err := gstore.GetObjectURL(bucketName, h.Filename)
			if err != nil {
				return "", err
			}

			track := repo.Track{}
			track.StoreURL = url
			track.Name = h.Filename
			err = app.DB.PushTrackToSQL(track)

			if err != nil {
				err = fmt.Errorf(" error pushing track to sql and %w", err)
				_ = gstore.DeleteObjectInBucket(bucketName, h.Filename)
				return "", err
			}
		}

	}
	msg := ""
	for _, s := range already {
		msg += fmt.Sprintf(" %s were alreaddy in bucket ", s)
	}
	return msg, nil

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
	msg, err := app.LoadMultipartReqToBucket(r, app.BucketName)
	if err != nil {
		WriteErrorToResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(msg))
}
