package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/scanner"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	// Embed is used to import the template files
	_ "embed"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/pkg/mathutil"
	"github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/utils"
	"golang.org/x/tools/imports"
)

//go:embed migration.go.tpl
var migrationFile string

//go:embed migration_impl.go.tpl
var migrationImplFile string

//go:embed migration_test.go.tpl
var migrationTestFile string

//go:embed seq_num.go.tpl
var seqNumFile string

var (
	migrationTemplate     = newTemplate(migrationFile)
	migrationImplTemplate = newTemplate(migrationImplFile)
	migrationTestTemplate = newTemplate(migrationTestFile)
	seqNumTemplate        = newTemplate(seqNumFile)
)

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing file %q\n", file.Name())
	}
}

func main() {
	c := &cobra.Command{
		Use: "bootstrap datastore migration",
	}
	var migrationName string
	c.Flags().StringVar(&migrationName, "description", "", "a sentence describing the migration")
	utils.Must(c.MarkFlagRequired("description"))
	var rootDirectory string
	c.Flags().StringVar(&rootDirectory, "root", "", "the root directory of the source tree")
	utils.Must(c.MarkFlagRequired("root"))

	c.RunE = func(*cobra.Command, []string) error {
		startVersion := migrations.CurrentDBVersionSeqNum()
		migrationDirName := getMigrationDirName(startVersion, migrationName)
		_ = migrationDirName

		templateMap := map[string]interface{}{
			"nextSeqNum":          startVersion + 1,
			"packageName":         getPackageName(startVersion),
			"startSequenceNumber": startVersion,
		}

		var err error
		// Create migration directory
		fullMigrationDirPath := path.Join(rootDirectory, "migrator", "migrations", migrationDirName)
		err = os.MkdirAll(fullMigrationDirPath, 0755)
		if err != nil {
			return err
		}
		// Write migration file
		migrationFilePath := path.Join(fullMigrationDirPath, "migration.go")
		err = renderFile(templateMap, migrationTemplate, migrationFilePath)
		if err != nil {
			return err
		}
		// Write migration impl file
		migrationImplFilePath := path.Join(fullMigrationDirPath, "migration_impl.go")
		err = renderFile(templateMap, migrationImplTemplate, migrationImplFilePath)
		if err != nil {
			return err
		}
		// Write migration test file
		migrationTestFilePath := path.Join(fullMigrationDirPath, "migration_test.go")
		err = renderFile(templateMap, migrationTestTemplate, migrationTestFilePath)
		if err != nil {
			return err
		}
		// Overwrite seqence number file
		seqNumFilePath := path.Join(rootDirectory, "pkg", "migrations", "internal", "seq_num.go")
		err = renderFile(templateMap, seqNumTemplate, seqNumFilePath)
		if err != nil {
			return err
		}
		// Register migration
		registrationFilePath := path.Join(rootDirectory, "migrator", "runner", "all.go")
		err = registerMigration(registrationFilePath, migrationDirName)
		if err != nil {
			return err
		}

		return nil
	}
	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addWord(elems []string, builder *strings.Builder) []string {
	word := strings.ToLower(builder.String())
	if len(word) > 0 {
		elems = append(elems, word)
		builder.Reset()
	}
	return elems
}

func convertDescriptionToSuffix(description string) string {
	elems := make([]string, 0)
	var builder strings.Builder
	for _, r := range description {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			builder.WriteRune(r)
		} else {
			elems = addWord(elems, &builder)
		}
	}
	elems = addWord(elems, &builder)
	return strings.Join(elems, "_")
}

func getPackageName(startVersion int) string {
	return fmt.Sprintf("m%dtom%d", startVersion, startVersion+1)
}

func getMigrationDirName(startVersion int, description string) string {
	suffix := convertDescriptionToSuffix(description)
	return fmt.Sprintf("m_%d_to_m_%d_%s", startVersion, startVersion+1, suffix)
}

func readRegistrationFileAndRegisterMigration(registrationFilePath string, migrationDirName string) ([]string, error) {
	newFileLines := make([]string, 0)
	readFile, err := os.Open(registrationFilePath)
	defer closeFile(readFile)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	isInImports := false
	registered := false
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.HasPrefix(line, "import (") {
			isInImports = true
		}
		if isInImports && strings.HasPrefix(line, ")") {
			if !registered {
				registrationPrefix := "github.com/stackrox/rox/migrator/migrations"
				registeredPath := path.Join(registrationPrefix, migrationDirName)
				newFileLines = append(newFileLines, fmt.Sprintf("\t_ %q", registeredPath))
				registered = true
			}
			isInImports = false
		}
		newFileLines = append(newFileLines, line)
	}
	if isInImports {
		return nil, errors.New("Bad registration file format: import section is not properly closed")
	}
	if !registered {
		return nil, errors.New("Failed to register migration")
	}
	return newFileLines, nil
}

func registerMigration(registrationFilePath string, migrationDirName string) error {
	newFileLines, err := readRegistrationFileAndRegisterMigration(registrationFilePath, migrationDirName)
	if err != nil {
		return err
	}
	// Write back
	writeFile, err := os.OpenFile(registrationFilePath, os.O_RDWR, 0644)
	defer closeFile(writeFile)
	if err != nil {
		return err
	}
	for _, line := range newFileLines {
		fmt.Fprintln(writeFile, line)
	}
	return nil
}

func renderFile(templateMap map[string]interface{}, temp func(s string) *template.Template, templateFileName string) error {
	buf := bytes.NewBuffer(nil)
	if err := temp(templateFileName).Execute(buf, templateMap); err != nil {
		return err
	}
	file := buf.Bytes()

	importProcessingStart := time.Now()
	formatted, err := imports.Process(templateFileName, file, nil)
	importProcessingDuration := time.Since(importProcessingStart)

	if err != nil {
		target := scanner.ErrorList{}
		if !errors.As(err, &target) {
			fmt.Println(string(file))
			return err
		}
		e := target[0]
		fileLines := strings.Split(string(file), "\n")
		fmt.Printf("There is an error in following snippet: %s\n", e.Msg)
		fmt.Println(strings.Join(fileLines[mathutil.MaxInt(0, e.Pos.Line-2):mathutil.MinInt(len(fileLines), e.Pos.Line+1)], "\n"))
		return err
	}
	if err := os.WriteFile(templateFileName, formatted, 0644); err != nil {
		return err
	}
	if importProcessingDuration > time.Second {
		absTemplatePath, err := filepath.Abs(templateFileName)
		if err != nil {
			absTemplatePath = templateFileName
		}
		log.Panicf("Import processing for file %q took more than 1 second (%s). This typically indicates that an import was "+
			"not added to the Go template, which forced import processing to search through all types and magically "+
			"add the import. Please add the import to the template; you can compare the imports in the generated file "+
			"with the ones in the template, and add the missing one(s)", absTemplatePath, importProcessingDuration)
	}
	return nil
}

func newTemplate(tpl string) func(name string) *template.Template {
	return func(name string) *template.Template {
		return template.Must(template.New(name).Option("missingkey=error").Parse(tpl))
	}
}
