package templator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

func GetGopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func GetPackagePath() string {
	currentPath, _ := filepath.Abs("./")
	pathAfterGopath := strings.Split(currentPath, GetGopath()+"/src/")
	return pathAfterGopath[1]
}

func WriteFileIfNotExist(templateFile, outputFile string, data interface{}) error {

	if !IsExist(outputFile) {
		return WriteFile(templateFile, outputFile, data)
	}

	return nil
}

// function that used in templates
var FuncMap = template.FuncMap{
	"CamelCase":  CamelCase,
	"PascalCase": PascalCase,
	"SnakeCase":  SnakeCase,
	"UpperCase":  UpperCase,
	"LowerCase":  LowerCase,
}

func WriteFile(templateFile, outputFile string, data interface{}) error {

	var buffer bytes.Buffer

	// this process can be refactor later
	{
		// open template file
		file, err := os.Open(fmt.Sprintf("%s/src/github.com/mirzaakhena/gogen/templates/%s", GetGopath(), templateFile))
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			row := scanner.Text()
			buffer.WriteString(row)
			buffer.WriteString("\n")
		}
	}

	tpl := template.Must(template.New("something").Funcs(FuncMap).Parse(buffer.String()))

	fileOut, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	{
		err := tpl.Execute(fileOut, data)
		if err != nil {
			return err
		}
	}

	return nil

}

// CamelCase is
func CamelCase(name string) string {

	// force it!
	// this is bad. But we can figure out later
	{
		if name == "IPAddress" {
			return "ipAddress"
		}

		if name == "ID" {
			return "id"
		}
	}

	out := []rune(name)
	out[0] = unicode.ToLower([]rune(name)[0])
	return string(out)
}

// UpperCase is
func UpperCase(name string) string {
	return strings.ToUpper(name)
}

// LowerCase is
func LowerCase(name string) string {
	return strings.ToLower(name)
}

// PascalCase is
func PascalCase(name string) string {
	return name
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase is
func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func CreateFolder(format string, a ...interface{}) {
	folderName := fmt.Sprintf(format, a...)
	if err := os.MkdirAll(folderName, 0755); err != nil {
		panic(err)
	}
}

func GenerateMock(packagePath, usecaseName string) {
	fmt.Printf("mockery %s\n", usecaseName)

	lowercaseUsecaseName := strings.ToLower(usecaseName)

	cmd := exec.Command(
		"mockery",

		// read file under path usecase/USECASE_NAME/port/
		"--dir", fmt.Sprintf("usecase/%s/port/", lowercaseUsecaseName),

		// specifically interface with name that has suffix 'Outport'
		"--name", fmt.Sprintf("%sOutport", usecaseName),

		// put the mock under usecase/%s/mocks/
		"-output", fmt.Sprintf("usecase/%s/mocks/", lowercaseUsecaseName),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("mockery failed with %s\n", err)
	}
}

func IsExist(fileOrDir string) bool {
	_, err := os.Stat(fileOrDir)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false

}

func GoFormat(path string) {
	fmt.Println("go fmt")
	cmd := exec.Command("go", "fmt", fmt.Sprintf("%s/...", path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

// func ReadYAML(usecaseName string) (*Usecase, error) {

// 	content, err := ioutil.ReadFile(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName))
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, fmt.Errorf("cannot read %s.yml", usecaseName)
// 	}

// 	tp := Usecase{}

// 	if err = yaml.Unmarshal(content, &tp); err != nil {
// 		log.Fatalf("error: %+v", err)
// 		return nil, fmt.Errorf("%s.yml is unrecognized usecase file", usecaseName)
// 	}

// 	tp.Name = usecaseName
// 	tp.PackagePath = GetPackagePath()

// 	return &tp, nil

// }

// func ReadLineByLine(filepath string) []string {
// 	var lineOfCodes []string
// 	{
// 		file, err := os.Open(filepath)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer file.Close()

// 		scanner := bufio.NewScanner(file)

// 		scanner.Split(bufio.ScanLines)

// 		for scanner.Scan() {
// 			lineOfCodes = append(lineOfCodes, scanner.Text())
// 		}
// 	}

// 	return lineOfCodes
// }
