package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	mpReader, err := r.MultipartReader()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("not a multipart request"))
		return
	}

	form, err := mpReader.ReadForm(int64(readerMaxMemory)) // what gover over this limit will be stored in temporary files on disk
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("error reading from form"))
		return
	}
	totalSize := int64(0)
	nbFiles := 0
	for key, fileHeaders := range form.File {
		log.Println("form key:", key)
		for _, fileHeader := range fileHeaders {
			log.Println("  fileName: ", fileHeader.Filename)
			log.Printf("  size: %d KB\n", fileHeader.Size / 1000)
			log.Println("  ---------------")
			totalSize += fileHeader.Size
			nbFiles++
		}
		log.Println("---------------")
	}
	log.Printf("upload request for took %s\n", time.Since(start).String())
	fmt.Fprintf(w, "uploaded %d files, with total size %d bytes\n", nbFiles, totalSize)
}

func Download(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sizeParameter := r.URL.Query()["size"]
	var contentSize int64
	if len(sizeParameter) > 0 {
		if v, err := strconv.Atoi(sizeParameter[0]); err != nil {
			w.Write([]byte("size query parameters must be a number (bytes)"))
			return
		} else {
			contentSize = int64(v)
		}
	} else {
		w.Write([]byte("must include size query parameter"))
		return
	}
	w.Header().Set("Content-Type", "application/octet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"file.randombytes\""))
	written := int64(0)
	for  written < contentSize {
		remaining := contentSize - written
		chunkLimit := int64(chunkSize)
		if remaining < chunkLimit {
			// write remaining bytes
			buff := make([]byte, remaining)
			rand.Read(buff)
			w.Write(buff)
			break
		}
		buff := make([]byte, chunkLimit)
		rand.Read(buff)
		sent, _ := w.Write(buff)
		written += int64(sent)
	}
	log.Printf("download request for %d bytes took %s", contentSize, time.Since(start).String())
}

func Headers(w http.ResponseWriter, r *http.Request) {
	var out io.Writer = w
	if printParam := r.URL.Query().Get("print"); len(printParam) > 0 {
		if v, err := strconv.ParseBool(printParam); err == nil && v == true {
			out =  io.MultiWriter(w, os.Stdout)
		}
	}
	fmt.Fprintf(out, "%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	headerKeys := make([]string, len(r.Header))
	i := 0
	for key, _ := range r.Header {
		headerKeys[i] = key
		i++
	}
	sort.Strings(headerKeys)
	for _, headerKey := range headerKeys {
		headerVal := r.Header[headerKey]
		fmt.Fprintf(out, "%s: %s\n", headerKey, strings.Join(headerVal, ";"))
	}
	fmt.Fprintln(out)
}
