package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	themeName := flag.String("name", "Rosepine", "Rosepine theme for wordpress")
	git := flag.String("git", "true", "Start git repository. must be boolean 'true' or 'false'")

	flag.Parse()

	os.Mkdir(*themeName, 0755)
	os.Mkdir(*themeName+"/template-parts", 0755)
	os.Mkdir(*themeName+"/inc", 0755)

	root_files := [6]string{"functions.php", "header.php", "footer.php", "style.css", "theme.json", "__version__"}
	template_files := [2]string{"page.php", "single.php"}
	inc_files := [2]string{"custom-functions.php", "custom-templates.php"}

	dirFiller(*themeName, root_files[:])
	dirFiller(*themeName+"/template-parts", template_files[:])
	dirFiller(*themeName+"/inc", inc_files[:])

	if *git == "true" {
		gitSetup(*themeName)
	}
}

func dirFiller(themeName string, files []string) {
	for i := range files {
		// Create a new file
		file, err := os.Create(themeName + "/" + files[i])
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Write data to the file
		_, err = file.WriteString("These are the contents of " + files[i])
		if err != nil {
			panic(err)
		}
	}

}

func gitSetup(themeName string) {

	gitInit := exec.Command("git", "init")
	gitAdd := exec.Command("git", "add", ".")
	gitCommit := exec.Command("git", "commit", "-m", "Initial Commit for "+themeName+" theme")

	gitInit.Dir = themeName
	gitAdd.Dir = themeName
	gitCommit.Dir = themeName

	output, err := gitInit.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	output, err = gitAdd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	output, err = gitCommit.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
