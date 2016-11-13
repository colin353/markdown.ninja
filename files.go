/*
  files.go

  Defines operations related to files, such as uploading,
  deleting, etc.
*/

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/colin353/portfolio/models"
	"github.com/colin353/portfolio/requesthandler"
)

// NewFileHandler returns an instance of the edit handler, with
// the routes populated.
func NewFileHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"file":   file,
		"files":  files,
		"upload": upload,
		"rename": renameFile,
		"delete": deleteFile,
	}
	return &a
}

func file(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	return nil
}

func files(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	f := models.File{}
	f.Domain = u.Domain
	iterator, err := models.GetList(&f)
	if err != nil {
		log.Printf("Tried to load files under `%s`, but it failed.", f.RegistrationKey())
		return requesthandler.ResponseError
	}

	fileList := make([]map[string]interface{}, 0, iterator.Count())
	for iterator.Next() {
		fileList = append(fileList, iterator.Value().Export())
	}

	return fileList
}

func upload(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	// We will allow for 32 MiB of memory allocated
	// to file reading, if the file is much bigger than that,
	// it will be put into a temp directory.
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")

	// Figure out the file size. Since the file may be in memory
	// or may be an actual file, we can use Seek to find the size.
	// The second argument of Seek() tells it where to reference the first
	// argument as an offset from. 2 --> end of the file, 1 --> start.
	size, err := file.Seek(0, 2)
	if err != nil {
		log.Printf("Had trouble seeking to end in uploaded file (!)")
		return requesthandler.ResponseError
	}
	// Important to seek back to the start, so we don't mess up the future
	// read operations.
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("Had trouble seeking back to start in uploaded file (!)")
		return requesthandler.ResponseError
	}
	defer file.Close()

	if size > (10 << 20) {
		log.Printf("You can't upload such a big file (%d bytes) to the server. It's not allowed.", size)
		return requesthandler.ResponseError
	}

	if err != nil {
		log.Printf("Trying to load uploaded file, but it failed.")
		return requesthandler.ResponseError
	}

	f := models.File{}

	// Need to ensure that the filename is acceptable. In order to
	// do this, I'll set the filename "safely", which is guaranteed
	// to result in a valid filename by stripping illegal characters.
	f.SetNameSafely(handler.Filename)
	f.Domain = u.Domain

	// We need to check if that file already exists. If it does,
	// we'll delete it so we can replace it.
	err = models.Load(&f)
	if err == nil {
		// The file DOES exist. So we'll delete the redis record
		// and the actual file.
		os.Remove(f.GetPath())
		models.Delete(&f)
	}

	f.Size = int(size)

	log.Printf("Uploaded file: size = %v bytes", f.Size)

	// Now we need to compute the file hash using MD5.
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Unable to read uploaded file.")
		return requesthandler.ResponseError
	}
	f.Hash = fmt.Sprintf("%x", md5.Sum(data))

	if !f.Validate() {
		log.Printf("Validation failed on attempted upload: `%s`", f.Key())
		return requesthandler.ResponseError
	}

	// Seek back to the beginning of the file so we can copy it.
	_, err = file.Seek(0, 0)

	// Finally, we'll move the file to the appropriate location
	// on the server.
	targetFile, err := os.OpenFile(f.GetPath(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Unable to create a new file and copy the upload.")
		return requesthandler.ResponseError
	}
	defer targetFile.Close()

	// Copy the file over.
	_, err = io.Copy(targetFile, file)
	if err != nil {
		log.Printf("Unable to copy the upload.")
		return requesthandler.ResponseError
	}

	// Okay, everything was a success so far. So our final step will
	// be to create the new record of the file in redis.
	err = models.Insert(&f)
	if err != nil {
		log.Printf("Unable to create a record of the upload in the database.")
		return requesthandler.ResponseError
	}

	return requesthandler.ResponseOK
}

func renameFile(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	return nil
}

func deleteFile(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type deleteArgs struct {
		Name string `json:"name"`
	}
	args := deleteArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	// Try to load the file record.
	f := models.File{}
	f.Domain = u.Domain
	f.Name = args.Name
	err = models.Load(&f)
	if err != nil {
		log.Printf("Tried to delete file `%s`, but failed", f.Key())
		return requesthandler.ResponseInvalidArgs
	}

	// Now try to delete the file from the hard drive.
	err = os.Remove(f.GetPath())
	if err != nil {
		log.Printf("Failed to delete file `%s` from hard drive.", f.Key())
		http.Error(w, "", http.StatusInternalServerError)
		return requesthandler.ResponseError
	}

	// And now delete it from the database.
	err = models.Delete(&f)
	if err != nil {
		log.Printf("Failed to delete file `%s` from database.", f.Key())
		http.Error(w, "", http.StatusInternalServerError)
		return requesthandler.ResponseError
	}

	return requesthandler.ResponseOK
}
