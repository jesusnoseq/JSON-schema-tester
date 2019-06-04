package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jesusnoseq/jsonschematester/clog"
	"github.com/kelseyhightower/envconfig"
	"github.com/xeipuuv/gojsonschema"
)

// PathConfig config
type PathConfig struct {
	SchemasDir  string `envconfig:"SCHEMA_DIR" default:"testdata/schemas"`
	SchemasURL  string `envconfig:"SCHEMA_URL" default:"/"`
	ExamplesDir string `envconfig:"EXAMPLE_DIR" default:"testdata/examples"`
	ExamplesURL string `envconfig:"EXAMPLE_URL" default:"/examples/"`
	ServerAddr  string `envconfig:"SERVER_ADDR" default:":8080"`
}

func main() {
	config := parseConfig()
	clog.Info("Initiating tester with configuration: %+v", config)
	go initWebServer(config)

	schemas := scanFolder(config.SchemasDir)
	filterPaths(schemas, config.SchemasDir)
	sort.Strings(schemas)
	examples := scanFolder(config.ExamplesDir)
	filterPaths(examples, config.ExamplesDir)
	sort.Strings(examples)
	var toValidate []string

	clog.Info("There are %d schemas and %d examples", len(schemas), len(examples))

	for i := 0; i < len(schemas); i++ {
		if contains(examples, schemas[i]) {
			toValidate = append(toValidate, schemas[i])
		} else {
			clog.Error("Schema %s does not have an example", schemas[i])
		}
	}
	for i := 0; i < len(examples); i++ {
		if contains(schemas, examples[i]) {
			// everything ok
		} else {
			clog.Error("Example %s does not have a schema", examples[i])
		}
	}
	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true
	for i := 0; i < len(schemas); i++ {
		schemaURL := "http://127.0.0.1:8080/" + config.SchemasURL + schemas[i]
		documentURL := "http://127.0.0.1:8080/" + config.ExamplesURL + schemas[i]
		loader := gojsonschema.NewReferenceLoader(schemaURL)
		errSchema := sl.AddSchemas(loader)
		if errSchema != nil {
			clog.Error("Schema %s is not valid. see error : \n %s", schemas[i], errSchema.Error())
			continue
		} else {
			clog.Success("Valid %s", schemas[i])
		}
		if !contains(toValidate, schemas[i]) {
			continue
		}
		documentLoader := gojsonschema.NewReferenceLoader(documentURL)
		result, err := gojsonschema.Validate(loader, documentLoader)
		if err != nil {
			clog.Error("Error testing %s schema; %s", schemas[i], err.Error())
		} else {
			if result.Valid() {
				clog.Success("Valid and tested %s ", schemas[i])
			} else {
				clog.Error("The document %s is not valid. see errors :\n", schemas[i])
				for _, desc := range result.Errors() {
					clog.Error("- %s", desc)
				}
			}
		}
	}

}

//loader := gojsonschema.NewReferenceLoader("http://127.0.0.1:8080/")
////loader.LoadJSON()
/*sl := gojsonschema.NewSchemaLoader()
sl.Validate = true
errSchema := sl.AddSchemas(gojsonschema.NewStringLoader(`{
     $id" : "http://some_host.com/invalid.json",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "multipleOf" : true
}`))
if errSchema != nil {
	clog.Error("The schema is not valid. see error : \n %s", errSchema.Error())
}*/
//sl.Compile()

/*schemaLoader := gojsonschema.NewReferenceLoader("file:///home/me/schema.json")
documentLoader := gojsonschema.NewReferenceLoader("file:///home/me/document.json")

result, err := gojsonschema.Validate(schemaLoader, documentLoader)
if err != nil {
	panic(err.Error())
}

if result.Valid() {
	fmt.Printf("The document is valid\n")
} else {
	fmt.Printf("The document is not valid. see errors :\n")
	for _, desc := range result.Errors() {
		fmt.Printf("- %s\n", desc)
	}
}
*/
/*
	fsEx := http.FileServer(http.Dir(config.ExamplesDir))
	http.Handle(config.ExamplesURL, http.StripPrefix(config.ExamplesURL, fsEx))

	fsSc := http.FileServer(http.Dir(config.SchemasDir))
	http.Handle(config.SchemasURL, http.StripPrefix(config.SchemasURL, fsSc))

	http.ListenAndServe(config.ServerAddr, nil)*/

/*
gojsonschema.NewStringLoader(`{
	     $id" : "http://some_host.com/invalid.json",
	    "$schema": "http://json-schema.org/draft-07/schema#",
	    "multipleOf" : true
	}`)
*/
func parseConfig() PathConfig {
	var r PathConfig
	err := envconfig.Process("", &r)
	if err != nil {
		log.Fatal(err.Error())
	}

	return r
}

func initWebServer(conf PathConfig) {
	srv := &http.Server{Addr: conf.ServerAddr, Handler: nil}

	fsEx := http.FileServer(http.Dir(conf.ExamplesDir))
	http.Handle(conf.ExamplesURL, http.StripPrefix(conf.ExamplesURL, fsEx))

	fsSc := http.FileServer(http.Dir(conf.SchemasDir))
	http.Handle(conf.SchemasURL, http.StripPrefix(conf.SchemasURL, fsSc))

	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		clog.Error("server on %s is closed", conf.ServerAddr)
	} else if err == http.ErrServerClosed {
		clog.Error("server error %s on %s", conf.ServerAddr, err.Error())
	} else {
		clog.Info("server listening on %s...", conf.ServerAddr)
	}
}

func scanFolder(rootPath string) []string {
	var files []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func contains(a []string, x string) bool {
	return sort.SearchStrings(a, x) != len(a)
}

func pathToURL(path string) string {
	return ""
}

func filterPaths(paths []string, filter string) {
	for i, path := range paths {
		// for windows
		path = strings.Replace(path, "\\", "/", -1)
		paths[i] = strings.TrimPrefix(path, filter)
	}
}
