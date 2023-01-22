package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"log"
	"net/http"
	"net/http/httptest"
)

func ExampleServerHandler_FullURLHashBy() {
	cnf := config.NewConfig()
	controller := repository.NewStorage(cnf)
	r := NewRouter(controller, cnf)
	hash, err := controller.InsertURL(context.Background(), "http://test.test/test", "sadASdQeAWDwdAs")
	if err != nil {
		log.Fatal(err)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+hash, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func ExampleServerHandler_ShortURLTextBy() {
	cnf := config.NewConfig()
	controller := repository.NewStorage(cnf)
	r := NewRouter(controller, cnf)
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewBuffer([]byte("http://www.test.test/test")))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func ExampleServerHandler_ShortURLJSONBy() {
	cnf := config.NewConfig()
	controller := repository.NewStorage(cnf)
	r := NewRouter(controller, cnf)
	ts := httptest.NewServer(r)
	defer ts.Close()
	b, err := json.Marshal(repository.FullURL{
		Full: ""})
	if err != nil {
		log.Fatal(err)
	}
	_, err = controller.InsertURL(context.Background(), "http://www.test.net/test", "sadASdQeAWDwdAs")
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPatch, ts.URL+"/api/shorten", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func ExampleServerHandler_Ping() {
	cnf := config.NewConfig()
	controller := repository.NewStorage(cnf)
	r := NewRouter(controller, cnf)
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/ping", bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func ExampleServerHandler_GetAllUserURLs() {
	cnf := config.NewConfig()
	controller := repository.NewStorage(cnf)
	r := NewRouter(controller, cnf)
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewBuffer([]byte("http://www.test.test/test")))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	cookies := resp.Cookies()
	var c *http.Cookie
	for _, v := range cookies {
		if v.Name == "shortener" {
			c = v
		}
	}
	req, err = http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewBuffer([]byte("http://www.test.test/test2")))
	if err != nil {
		log.Fatal(err)
	}
	req.AddCookie(c)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	req, err = http.NewRequest(http.MethodGet, ts.URL+"/api/user/urls", bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Fatal(err)
	}
	req.AddCookie(c)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp2.Body.Close()
}
