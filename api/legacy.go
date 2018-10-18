package api

//Legacy Helper Functions
import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var legacySupport bool

//Check the current module downloaded for refs,
//If present we need to download legacybridge.

func CheckforLegacySupport(args string) bool {

	path := os.Getenv("GOPATH")
	currDir, _ := os.Getwd()
	truePath := getTruePath(currDir, args)
	dirPath := Concat(path, "/pkg/mod/", truePath)
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {

		tag := strings.Split(strings.Split(truePath, "/")[2], "@")[1]

		for i := 3; os.IsNotExist(err); i++ {

			elements := strings.Split(args, "/")

			elements[i] = Concat(elements[i], "@", tag)

			tempPath := strings.Join(elements, "/")
			dirPath = Concat(path, "/pkg/mod/", tempPath)

			_, err = os.Stat(dirPath)

		}
	}

	files, err := ioutil.ReadDir(dirPath)
	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {

			listsOfRefs := GetRefsFromFile(Concat(dirPath, "/", f.Name()))
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

	os.Chdir(Concat(path, "/src"))
	cliCmd, err := exec.Command("go", "mod", "tidy").Output()
	os.Chdir(path)

	fmt.Println("Opening go.mod")
	file, err := os.Open(Concat(path, "/src/go.mod"))

	if err != nil {
		fmt.Println(string(cliCmd))
		fmt.Println(err)
	}
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
