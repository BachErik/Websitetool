package main

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Generate a new GUID
	fmt.Println(">GUID")
	guid := uuid.New()
	_, err := io.ReadFull(rand.Reader, []byte(guid.String()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GUID>")

	// Create a folder with the GUID as its name
	fmt.Println(">folder-GUID")
	err = os.Mkdir(guid.String(), 0700)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("folder-GUID>")

	// Upload the file to the folder
	fmt.Println(">file to folder")
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filePath := fmt.Sprintf("%s/%s", guid.String(), header.Filename)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("file to folder>")

	// Set the GUID as a cookie
	fmt.Println(">GUID-cookie")
	http.SetCookie(w, &http.Cookie{
		Name:  "GUID",
		Value: guid.String(),
	})
	fmt.Println("GUID-cookie>")
	// fmt.Fprintln(w, "File uploaded successfully")

	// Unzip zip
	fmt.Println(">unzip")
	err = os.Mkdir(guid.String()+"/unzipped", 0700)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(guid.String() + "/site.zip")
	fmt.Println(guid.String() + "/unzipped")
	error := unzip(guid.String()+"/site.zip", guid.String()+"/unzipped")
	if error != nil {
		log.Fatal(err)
	}
	fmt.Println("unzip>")

	// Show Index
	// fmt.Fprintln(w, "File unzipped successfully")

	showIndexForm(w, r)
}

func showUploadForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "upload.html")
}

func showIndexForm(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "index.html")
	cookie, err := r.Cookie("GUID")
	if err != nil {
		http.ServeFile(w, r, "upload.html")
		fmt.Println(err)
	}
	fmt.Println(cookie)
	fmt.Println(cookie)
	fmt.Println(cookie)
	fmt.Println(cookie)
	fmt.Println(cookie)
	fmt.Println(cookie)

	http.ServeFile(w, r, "index.html")
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello :D")
	http.HandleFunc("/index", showIndexForm)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/", showUploadForm)
	http.ListenAndServe(":8989", nil)
}
