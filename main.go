package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-openapi/spec"
	"github.com/hahnlee/stamped/stamped-cli"
	"github.com/urfave/cli"
)

type SwaggerConfig struct {
	Version string `json:"version"`
}

type PostManConfig struct {
	Version string `json:"version"`
	Host    string `json:"host"`
}

type Config struct {
	SwaggerConfig SwaggerConfig `json:"swagger"`
	PostManConfig PostManConfig `json:"postman"`
}

func main() {
	var url string
	var configFile string
	var swaggerFilePath string
	var swagger spec.Swagger
	var saveFilePath string

	var config Config

	app := cli.NewApp()

	app.Name = "s2p"
	app.Usage = "Swagger to postman"
	app.Version = "0.1.0"
	app.Author = "hahnlee <hanlee.dev@gmail.com>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Usage:       "Load configuration from `file`.",
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "url",
			Usage:       "Load swagger json from `url`",
			Destination: &url,
		},
		cli.StringFlag{
			Name:        "file",
			Usage:       "Load swagger json from `file`",
			Destination: &swaggerFilePath,
		},
		cli.StringFlag{
			Name:        "host",
			Value:       "API-HOST",
			Usage:       "Your api endpoint",
			Destination: &config.PostManConfig.Host,
		},
		cli.StringFlag{
			Name:        "out",
			Value:       "./postman.json",
			Usage:       "Postman file output",
			Destination: &saveFilePath,
		},
	}

	app.Action = func(c *cli.Context) error {
		if configFile != "" {
			jsonFile, err := os.Open(configFile)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &config)
		}

		if url != "" {
			swagger = stamped.DownloadSwaggerFile(url)
		}

		if swaggerFilePath != "" {
			swaggerFile, _ := os.Open(swaggerFilePath)
			defer swaggerFile.Close()

			byteValue, _ := ioutil.ReadAll(swaggerFile)
			json.Unmarshal(byteValue, &swagger)
		}

		postman := stamped.NewPostMan(config.PostManConfig.Host)
		postman.SwaggerToPostman(&swagger)
		postman.Save(saveFilePath)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
