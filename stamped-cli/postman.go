package stamped

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/go-openapi/spec"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type Url struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type Header struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
}

type Request struct {
	Method string   `json:"method"`
	Header []Header `json:"header"`
	Body   Body     `json:"body"`
	URL    Url      `json:"url"`
}

type RequestInfo struct {
	Name    string  `json:"name"`
	Request Request `json:"request"`
}

type Item struct {
	Name string        `json:"name"`
	Item []RequestInfo `json:"item"`
}

type Info struct {
	Schema      string `json:"schema"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PostmanFile struct {
	Item []Item `json:"item"`
	Info Info   `json:"info"`
}

type IndexMap map[string]int

type Postman struct {
	host     string
	indexMap IndexMap
	file     PostmanFile
}

func NewPostMan(host string) *Postman {
	postman := new(Postman)
	postman.host = "{{" + host + "}}"
	postman.file = PostmanFile{}
	postman.file.Info = Info{Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"}
	postman.file.Item = []Item{}
	postman.indexMap = IndexMap{}
	return postman
}

func (this *Postman) setRequest(operation *spec.Operation, path string, method string) {
	selectedItem := &this.file.Item[this.indexMap[operation.Tags[0]]]
	request := Request{}
	request.Method = method

	if operation.Produces != nil {
		header := Header{Key: "Accept", Value: operation.Produces[0], Disabled: false}
		request.Header = []Header{header}
	}

	for _, parameter := range operation.Parameters {
		if parameter.In == "header" {
			header := Header{Key: parameter.Name, Value: "{{" + parameter.Name + "}}", Disabled: parameter.Required}
			request.Header = append(request.Header, header)
		}
	}

	paths := []string{}
	urlPaths := strings.Split(path, "/")

	replacer := strings.NewReplacer("{", ":", "}", "")
	for _, urlPath := range urlPaths {
		transPath := urlPath
		if strings.HasPrefix(urlPath, "{") {
			transPath = replacer.Replace(urlPath)
		}
		paths = append(paths, transPath)
	}

	request.URL.Raw = strings.Join(paths, "/")
	request.URL.Host = []string{this.host}
	request.URL.Path = paths

	requestInfo := RequestInfo{Name: operation.Summary, Request: request}

	selectedItem.Item = append(selectedItem.Item, requestInfo)
}

func (this *Postman) SwaggerToPostman(swagger *spec.Swagger) {
	this.file.Info.Name = swagger.Info.Title
	this.file.Info.Description = swagger.Info.Description

	for index, tag := range swagger.Tags {
		item := Item{}
		item.Name = tag.Name
		this.file.Item = append(this.file.Item, item)
		this.indexMap[tag.Name] = index
	}

	for path, requests := range swagger.Paths.Paths {
		if requests.Get != nil {
			this.setRequest(requests.Get, path, "get")
		}

		if requests.Post != nil {
			this.setRequest(requests.Post, path, "post")
		}

		if requests.Put != nil {
			this.setRequest(requests.Put, path, "put")
		}

		if requests.Patch != nil {
			this.setRequest(requests.Patch, path, "patch")
		}

		if requests.Options != nil {
			this.setRequest(requests.Options, path, "options")
		}
	}
}

func (this *Postman) Save(path string) {
	postmanJSON, _ := json.MarshalIndent(this.file, "", "  ")
	_ = ioutil.WriteFile(path, postmanJSON, 0644)
}
