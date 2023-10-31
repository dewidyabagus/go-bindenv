# Go Bind Env
[![Continuous Integration](https://github.com/dewidyabagus/go-bindenv/actions/workflows/code-analysis.yml/badge.svg)](https://github.com/dewidyabagus/go-bindenv/actions/workflows/code-analysis.yml)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=dewidyabagus_go-bindenv&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=dewidyabagus_go-bindenv)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=dewidyabagus_go-bindenv&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=dewidyabagus_go-bindenv)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=dewidyabagus_go-bindenv&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=dewidyabagus_go-bindenv)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=dewidyabagus_go-bindenv&metric=coverage)](https://sonarcloud.io/summary/new_code?id=dewidyabagus_go-bindenv)

Bind environment variables into a struct.

## Installation
Add packages to module that have been created with a minimum requirement of Go version 1.18
```shell
go get -u github.com/dewidyabagus/go-bindenv
```

## Usage
To use the `go-bindenv` package, make sure you have the module ready to use. In this example we will use a new `main.go` file as follows :
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dewidyabagus/go-bindenv"
)

type GCS struct {
	CloudStorageDomain string `env:"CLOUD_STORAGE_DOMAIN"`
	CloudStorageBucket string `env:"CLOUD_STORAGE_BUCKET"`
}

func init() {
	os.Setenv("CLOUD_STORAGE_DOMAIN", "storage.googleapis.com")
	os.Setenv("CLOUD_STORAGE_BUCKET", "my-example-bucket")
}

func main() {
	env := bindenv.New()

	gcs := GCS{}
	if err := env.Bind(&gcs); err != nil {
		log.Fatalln("Error binding env:", err.Error())
	}
	fmt.Println("Storage Domain:", gcs.CloudStorageDomain)
	// [Output] Storage Domain: storage.googleapis.com

	fmt.Println("Storage Bucket:", gcs.CloudStorageBucket)
	// [Output] Storage Bucket: my-example-bucket
}
```
You can also add a default value if the `env` variable is not found.
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dewidyabagus/go-bindenv"
)

type Database struct {
	Host string `env:"DATABASE_HOST"`
	Port int    `env:"DATABASE_PORT"`
}

func init() {
	os.Setenv("DATABASE_HOST", "10.17.19.1")
}

func main() {
	env := bindenv.New()
	env.SetDefault("DATABASE_HOST", "127.0.0.1")
	env.SetDefault("DATABASE_PORT", "5432")

	dbConfig := Database{}
	if err := env.Bind(&dbConfig); err != nil {
		log.Fatalln("Error binding env:", err.Error())
	}
	fmt.Println("Database Host:", dbConfig.Host)
	// [Output] Database Host: 10.17.19.1

	fmt.Println("Database Port:", dbConfig.Port)
	// [Output] Database Port: 5432
}
```

## Contributing
I am very happy and open to your contributions. Immediately open a new issue if you find a bug or feature addition and let's discuss it there.
