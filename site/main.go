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
	folderName := fmt.Sprintf("%x", guid)
	fmt.Println("GUID>")

	// Create a folder with the GUID as its name
	fmt.Println(">folder-GUID")
	err = os.Mkdir(folderName, 0700)
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

	filePath := fmt.Sprintf("%s/%s", folderName, header.Filename)
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
		Value: folderName,
	})
	fmt.Println("GUID-cookie>")
	fmt.Fprintln(w, "File uploaded successfully")
	fmt.Println(">unzip")
	err = os.Mkdir(folderName + "/unzipped", 0700)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	error := Unzip(guid.String()+"/site.zip", guid.String()+"/unzipped")
	if error != nil {
		log.Fatal(err)
	}
	fmt.Println("unzip>")
}

func showUploadForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "upload.html")
}

func Unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0755)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        path := filepath.Join(dest, f.Name)

        // Check for ZipSlip (Directory traversal)
        if !strings.HasPrefix(path, filepath.Clean(dest) + string(os.PathSeparator)) {
            return fmt.Errorf("illegal file path: %s", path)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}

func main() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/", showUploadForm)
	http.ListenAndServe(":8080", nil)
}
