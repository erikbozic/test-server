package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func(){
		log.Printf("upload request for took %s\n", time.Since(start).String())
	}()
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
	fmt.Fprintf(w, "uploaded %d files, with total size %d bytes\n", nbFiles, totalSize)
}

func Download(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var contentSize int64
	defer func(){
		log.Printf("download request for %d bytes took %s", contentSize, time.Since(start).String())
	}()
	sizeParameter := r.URL.Query()["size"]
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

func Service(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func(){
		log.Printf("service request for took %s\n", time.Since(start).String())
	}()
	requestUrl, err := url.ParseRequestURI(serviceBaseUrl + serviceCallPath)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("error parsing url: %v", err)))
		return
	}
	req := http.Request{
		Method:           http.MethodGet,
		URL:              requestUrl,
		Header:           http.Header{},
	}
	if v := r.URL.Query()["xb3"]; len(v) > 0 && v[0] == "true" {
		copyXb3Headers(r, &req)
	}
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("error during http call: %v", err)))
		return
	}
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("error reading response body: %v", err)))
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBytes)
}

func Error(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func(){
		log.Printf("error request for took %s\n", time.Since(start).String())
	}()
	codeParameters := r.URL.Query()["code"]
	if len(codeParameters) == 0 {
		w.WriteHeader(500)
		w.Write([]byte("default status code 500"))
		return
	}
	randomIndex := 0
	if len(codeParameters)> 1 {
		randomIndex = mathrand.Intn(len(codeParameters))
	}
	code, err := strconv.Atoi(codeParameters[randomIndex])
	if err != nil || code < 100 || code > 599 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("invalid status code: %d. defaulting to 500", code)))
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("randomly returned %d code from provided codes in query string\n", code)))
}
