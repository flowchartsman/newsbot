package main

import (
	log "github.com/kdar/factorlog"
	T "html/template"
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

func webserverInit() {
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
