package destroy

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/flect"
)

func confirm(msg string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	text, _ := reader.ReadString('\n')

	return (text == "y\n" || text == "Y\n")
}

func removeTemplates(confirmed bool, fileName string) {
	if !confirmed && !confirm("Want to remove templates? (y/N)") {
		return
	}

	templatesFolder := filepath.Join("templates", fileName)
	fmt.Printf("- Deleted %v folder\n", templatesFolder)
	os.RemoveAll(templatesFolder)
}

func removeActions(confirmed bool, fileName string) error {
	if !confirmed && !confirm("Want to remove actions? (y/N)") {
		return nil
	}

	fmt.Printf("- Deleted %v\n", fmt.Sprintf("actions/%v.go", fileName))
	os.Remove(filepath.Join("actions", fmt.Sprintf("%v.go", fileName)))

	fmt.Printf("- Deleted %v\n", fmt.Sprintf("actions/%v_test.go", fileName))
	os.Remove(filepath.Join("actions", fmt.Sprintf("%v_test.go", fileName)))

	content, err := os.ReadFile(filepath.Join("actions", "app.go"))
	if err != nil {
		fmt.Println("error reading app.go content")
		return err
	}

	resourceExpression := fmt.Sprintf("app.Resource(\"/%v\", %vResource{})", fileName, flect.Pascalize(fileName))
	newContents := strings.Replace(string(content), resourceExpression, "", -1)

	err = os.WriteFile(filepath.Join("actions", "app.go"), []byte(newContents), 0)
	if err != nil {
		fmt.Println("error writing new app.go content")
		return err
	}

	fmt.Printf("- Deleted References for %v in actions/app.go\n", fileName)

	return nil
}

func removeLocales(confirmed bool, fileName string) {
	if confirmed || confirm("Want to remove locales? (y/N)") {
		removeMatch("locales", fmt.Sprintf("%v.*.yaml", fileName))
	}
}

func removeModel(confirmed bool, name string) {
	if !confirmed && !confirm("Want to remove model? (y/N)") {
		return
	}

	modelFileName := flect.Singularize(flect.Underscore(name))

	os.Remove(filepath.Join("models", fmt.Sprintf("%v.go", modelFileName)))
	os.Remove(filepath.Join("models", fmt.Sprintf("%v_test.go", modelFileName)))

	fmt.Printf("- Deleted %v\n", fmt.Sprintf("models/%v.go", modelFileName))
	fmt.Printf("- Deleted %v\n", fmt.Sprintf("models/%v_test.go", modelFileName))
}

func removeMigrations(confirmed bool, fileName string) {
	if !confirmed && !confirm("Want to remove migrations? (y/N)") {
		return
	}

	removeMatch("migrations", fmt.Sprintf("*_create_%v.up.*", fileName))
	removeMatch("migrations", fmt.Sprintf("*_create_%v.down.*", fileName))
}

func removeMatch(folder, pattern string) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return
	}

	for _, f := range files {
		matches, _ := filepath.Match(pattern, f.Name())
		if f.IsDir() || !matches {
			continue
		}

		path := filepath.Join(folder, f.Name())
		os.Remove(path)
		fmt.Printf("- Deleted %v\n", path)
	}
}
