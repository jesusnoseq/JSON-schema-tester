package checker

import (
	"github.com/jesusnoseq/JSON-schema-tester/clog"
	"github.com/jesusnoseq/JSON-schema-tester/config"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Check(config config.PathConfig) int {
	clog.Info("date-time format checker has been removed due to inconsistencies between json schema validators")
	// Temporally removed because it does not validate "2017-07-04T13:23:55" value as date-time format
	gojsonschema.FormatCheckers.Remove("date-time")

	schemas := scanFolder(config.SchemasDir)
	filterPaths(schemas, config.SchemasDir)
	sort.Strings(schemas)
	examples := scanFolder(config.ExamplesDir)
	filterPaths(examples, config.ExamplesDir)
	sort.Strings(examples)

	clog.Info("There are %d schemas and %d examples", len(schemas), len(examples))

	for i := 0; i < len(schemas); i++ {
		if contains(examples, schemas[i]) {
			clog.Success(" Schema have its example %s", schemas[i])
		} else {
			clog.Error("Schema %s does not have an example", schemas[i])
		}
	}
	for i := 0; i < len(examples); i++ {
		if contains(schemas, examples[i]) {
			clog.Success("Example have its schema %s", examples[i])
		} else {
			clog.Error("Example %s does not have a schema", examples[i])
		}
	}
	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true
	for i := 0; i < len(schemas); i++ {
		schemaURL := "http://127.0.0.1:8080" + config.SchemasURL + schemas[i]
		documentURL := "http://127.0.0.1:8080" + config.ExamplesURL + schemas[i]
		loader := gojsonschema.NewReferenceLoader(schemaURL)
		errSchema := sl.AddSchemas(loader)
		if errSchema != nil {
			clog.Error("Schema %s is not valid. see error: \n %s", schemas[i], errSchema.Error())
			continue
		}
		checkAdditionalPropertiesFalse(schemas[i], loader)

		clog.Success("Valid %s", schemaURL)

		if !contains(examples, schemas[i]) {
			continue
		}
		documentLoader := gojsonschema.NewReferenceLoader(documentURL)
		result, err := gojsonschema.Validate(loader, documentLoader)
		if err != nil {
			clog.Error("Error testing example %s against schema; %s", schemas[i], err.Error())
		} else {
			if result.Valid() {
				clog.Success("Valid and tested %s", documentURL)
			} else {
				clog.Error("The document %s is not valid. see errors:", schemas[i])
				for _, desc := range result.Errors() {
					clog.Error("- %s", desc)
				}
			}
		}
	}
	return clog.GetErrorsPrinted()
}

func checkAdditionalPropertiesFalse(schema string, l gojsonschema.JSONLoader) {
	i, err := l.LoadJSON()
	if err != nil {
		clog.Error("Error loading schema " + schema)
	}
	m, ok := i.(map[string]interface{})
	if !ok {
		clog.Error("Error parsing schema " + schema)
		return
	}
	ap, ok := m["additionalProperties"]
	if !ok {
		clog.Warn("Schema " + schema + " allow unknown additional properties, add property \"additionalProperties\": false")
		return
	}

	if app, ok := ap.(bool); !ok || app {
		clog.Warn("Schema " + schema + " allow unknown additional properties, change property additionalProperties to false")
		return
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
	index := sort.SearchStrings(a, x)
	return index < len(a) && a[index] == x
}

func filterPaths(paths []string, filter string) {
	for i, path := range paths {
		// for windows
		path = strings.Replace(path, "\\", "/", -1)
		paths[i] = strings.TrimPrefix(path, filter)
	}
}
