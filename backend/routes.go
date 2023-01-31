package main

import "net/http"


type information struct {
	test string
}

func testHandler(responseWriter http.ResponseWriter, request *http.Request) {
	info := &information{}
	info.test = "Dies ist eine Test Information."
	serveTemplate(responseWriter, "views/site/home.html", nil, info)
}
