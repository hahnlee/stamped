package stamped

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/spec"
)

func DownloadSwaggerFile(url string) spec.Swagger {
	var swagger spec.Swagger
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&swagger)
	if err != nil {
		panic(err)
	}
	return swagger
}
