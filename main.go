package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	portEnv := os.Getenv("PORT")
	var (
		port int
		err  error
	)
	if port, err = strconv.Atoi(portEnv); err != nil {
		port = 8080
	}
	sv := &ServeTracer{}
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), sv))
}

type ServeTracer struct {
	logs []*RequestLog
}

func (sv *ServeTracer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == "GET" && r.URL.Path == "/log" {
		q := r.URL.Query()
		if q.Has("id") {
			id, err := strconv.Atoi(q.Get("id"))
			if err != nil {
				log.Println(err)
				http.Error(w, "ID inválido", http.StatusNotFound)
				return
			}
			data, err := json.MarshalIndent((sv.logs[id].JSON()), "", "    ")
			if err != nil {
				log.Println(err)
				http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Log Detail</title>
    <style>
        body { font-family: Arial, sans-serif; }
        pre { background-color: #f0f0f0; padding: 10px; border-radius: 5px; }
    </style>
</head>
<body>
    <header>
        <h1>Log Detail</h1>
    </header>
    <main>
        <pre>` + strings.ReplaceAll(string(data), "\n", "<br/>") + `</pre>
    </main>
    <footer>
        <p>&copy; 2024 Your Company</p>
    </footer>
</body>
</html>`))
		} else {
			logs := ""
			for i, lg := range sv.logs {
				logs += `<a href="/log?id=` + fmt.Sprintf("%d", i) + `">` + lg.Title() + `</a><br/>`
			}
			w.WriteHeader(200)
			w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Logs</title>
    <style>
        body { font-family: Arial, sans-serif; }
        a { display: block; margin-bottom: 10px; }
    </style>
</head>
<body>
    <header>
        <h1>Logs</h1>
    </header>
    <main>
        ` + logs + `
    </main>
    <footer>
        <p>&copy; 2024 Your Company</p>
    </footer>
</body>
</html>`))
		}
	} else {
		sv.logs = append(sv.logs, &RequestLog{
			Request: r,
			Time:    time.Now(),
		})
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Log Added</title>
</head>
<body>
    <header>
        <h1>Log Added</h1>
    </header>
    <main>
        <p>Log has been successfully added.</p>
    </main>
    <footer>
        <p>&copy; 2024 Your Company</p>
    </footer>
</body>
</html>`))
	}
}

type RequestLog struct {
	Request *http.Request
	Time    time.Time
}

func (rl *RequestLog) Title() string {
	return fmt.Sprintf("%v %v %v %v", rl.Time, rl.Request.Method, rl.Request.Proto, rl.Request.RemoteAddr)
}

func (rl *RequestLog) JSON() RequestJSON {
	r := rl.Request
	return RequestJSON{
		Method:        r.Method,
		URL:           r.URL.String(),
		Proto:         r.Proto,
		ProtoMajor:    r.ProtoMajor,
		ProtoMinor:    r.ProtoMinor,
		Header:        r.Header,
		ContentLength: int(r.ContentLength),
		Host:          r.Host,
		Form:          r.Form,
		PostForm:      r.PostForm,
		Trailer:       r.Trailer,
		RemoteAddr:    r.RemoteAddr,
		RequestURI:    r.RequestURI,
	}
}

type RequestJSON struct {
	Method        string      `json:"Method"`
	URL           string      `json:"URL"`
	Proto         string      `json:"Proto"`
	ProtoMajor    int         `json:"ProtoMajor"`
	ProtoMinor    int         `json:"ProtoMinor"`
	Header        http.Header `json:"Header"`
	ContentLength int         `json:"ContentLength"`
	Host          string      `json:"Host"`
	Form          url.Values  `json:"Form"`
	PostForm      url.Values  `json:"PostForm"`
	Trailer       http.Header `json:"Trailer"`
	RemoteAddr    string      `json:"RemoteAddr"`
	RequestURI    string      `json:"RequestURI"`
}
