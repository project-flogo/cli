package cmd

//Legacy Helper Functions
import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//Check the current module downloaded for refs,
//If present we need to download legacybridge.
func NeedsLegacySupport(args string) bool {

	path := os.Getenv("GOPATH")
	currDir, _ := os.Getwd()
	truePath := getTruePath(currDir, args)

	files, err := ioutil.ReadDir(Concat(path, "/pkg/mod/", truePath))

	if err != nil {
		log.Fatal(err)
		return false
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {

			listsOfRefs := GetRefsFromFile(Concat(path, "/pkg/mod/", truePath, "/", f.Name()))
			if len(listsOfRefs) == 1 {
				return true
			}
			return false
		}
	}
	return false
}

//True Path becase the packages are stored in the version format and we need to get the version
//in order to navigate to that path.
func getTruePath(path string, pkg string) string {
	fmt.Println("Opening go.mod")
	file, err := os.Open(Concat(path, "/src/go.mod"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		modName := strings.Split(strings.TrimSpace(line), " ")

		if strings.Contains(pkg, modName[0]) && len(modName) != 1 {
			tempPath := strings.Split(pkg, "/")
			tempPath = makeItLowerCase(tempPath)
			tempPath[2] = Concat(tempPath[2], "@", modName[1])
			return strings.Join(tempPath, "/")
		}
	}
	return ""
}

//This function converts capotal letters in package name
// to !(smallercase). Eg C => !c . As this is the way
// go.mod saves every repository in the $GOPATH/pkg/mod.
func makeItLowerCase(s []string) []string {
	result := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		var b bytes.Buffer
		for _, c := range s[i] {
			if c >= 65 && c <= 90 {
				b.WriteRune(33)
				b.WriteRune(c + 32)
			} else {
				b.WriteRune(c)
			}
		}
		result[i] = b.String()
	}
	return result
}
