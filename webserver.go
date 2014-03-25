package main

import (
	T "html/template"
    log "github.com/kdar/factorlog"
	"net/http"
	"regexp"
)

var (
	home      *T.Template
	hasSuffix = regexp.MustCompile("/([^/]*\\.[^/]*)$")
	chttp     = http.NewServeMux()
)

func homeHandler(w http.ResponseWriter, req *http.Request) {
	if hasSuffix.MatchString(req.URL.Path) {
		chttp.ServeHTTP(w, req)
	} else {
		home.Execute(w, req.Host)
	}
}

func init() {
	homeTemplate, err := T.ParseFiles(BinPath + "/template/homepage.html")
	if err != nil {
		log.Fatalln(err)
	}
	home = homeTemplate
	chttp.Handle("/", http.FileServer(http.Dir(BinPath+"/static/")))
	http.HandleFunc("/", homeHandler)
	go func() {
		log.Fatalln(http.ListenAndServe(config.Port, nil))
	}()
}
