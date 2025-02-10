package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
)

type File struct {
	Name     string
	Contents string
}

type Theme struct {
	Name          string
	RootFiles     []File
	TemplateFiles []File
	IncFiles      []File
}

var functionsPHP = File{
	Name: "functions.php",
	Contents: `<?php
/**
 * File: mytheme/functions.php
 *
 * @package    mytheme
 * @subpackage  mytheme
 */
 
if (!defined('ABSPATH')) {
    exit;
}
 
// Add custom functions here
function mytheme_header() {
    // Header HTML goes here
}

// Enqueue CSS
add_action('wp_enqueue_scripts', function() {
    wp_enqueue_style(
        'mytheme-style',
        get_theme_path() . '/style.css',
        array(),
        filemtime(get_theme_path() . '/style.css')
    );
});
?>
`}

var headerPHP = File{Name: "header.php",
	Contents: `<?php
/**
 * File: mytheme/header.php
 */
 
if (!defined('ABSPATH')) {
    exit;
}

// Header HTML
$HeaderContent = '<nav class="main-nav">
    <div class="nav-brand">MyTheme</div>
    <ul class="nav-menu">
        <li><a href="/">Home</a></li>
        <li><a href="/about">About Us</a></li>
        <li><a href="/services">Services</a></li>
    </ul>
</nav>';
?>
`}

var footerPHP = File{Name: "footer.php",
	Contents: `<?php
/**
 * File: mytheme/footer.php
 */
 
if (!defined('ABSPATH')) {
    exit;
}

// Footer HTML
$FooterContent = '<footer class="main-footer">
    <div class="footer-content">
        <p>&copy; 2023 MyTheme. All rights reserved.</p>
    </div>
</footer>';
?>`}

var styles = File{Name: "style.css",
	Contents: `/* File: mytheme/style.css */
 
body {
    font-family: 'Arial', sans-serif;
    margin: 0;
    padding: 0;
}

header {
    background-color: #f0f0f0;
    padding: 20px;
}

nav.main-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
`}

var themeJSON = File{Name: "theme.json",
	Contents: `{
    "name": "MyTheme",
    "version": "1.0.0",
    "description": "A Custom WordPress Theme.",
    "author": "Your Name",
    "theme-features": {
        "custom-html": true,
        "custom-css": true,
        "template-parts": true
    }
}
`}

var version = File{Name: "__version__",
	Contents: ` __VERSION__
1.0.0
`}

var pagePHP = File{Name: "page.php",
	Contents: "Some Text"}

var singlePHP = File{Name: "single.php",
	Contents: "Some Text"}

var customFunc = File{Name: "custom-functions.php",
	Contents: `// mytheme/inc/custom-functions.php
function mytheme_custom_function() {
    // Function implementation goes here
}
`}

var customTemp = File{Name: "custom-templates.php",
	Contents: "Some Text"}

func main() {
	themeName := flag.String("name", "theme", "Theme for WordPress.")

	git := flag.String("git", "false", "Start git repository.")

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	author := flag.String("author", user.Username, "Start git repository.")

	flag.Parse()

	rootFiles := []File{functionsPHP, headerPHP, footerPHP, styles, themeJSON, version}
	templateFiles := []File{pagePHP, singlePHP}
	incFiles := []File{customFunc, customTemp}

	project := Theme{Name: *themeName, RootFiles: rootFiles, TemplateFiles: templateFiles, IncFiles: incFiles}

	os.Mkdir(project.Name, 0755)
	os.Mkdir(project.Name+"/template-parts", 0755)
	os.Mkdir(project.Name+"/inc", 0755)

	writeTemplate(project.Name, project.RootFiles[:], *author)
	writeTemplate(project.Name+"/template-parts", project.TemplateFiles[:], *author)
	writeTemplate(project.Name+"/inc", project.IncFiles[:], *author)

	if *git == "true" {
		gitSetup(project.Name)
	}
}

func writeTemplate(themeName string, files []File, authorName string) {
	for i := range files {
		// Create a new file
		file, err := os.Create(themeName + "/" + files[i].Name)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Write data to the file
		_, err = file.WriteString(formatContent(files[i].Contents, themeName, authorName))
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

	_, err := gitInit.Output()
	if err != nil {
		panic(err)
	}

	_, err = gitAdd.Output()
	if err != nil {
		panic(err)
	}

	_, err = gitCommit.Output()
	if err != nil {
		panic(err)
	}

	fmt.Print("Git repo initialized successfully.")
}

func formatContent(content string, themeName string, authorName string) string {
	myTheme := "[Mm]y[Tt]heme"
	author := "[Yy]our [Nn]ame"

	re := regexp.MustCompile(myTheme)

	content = re.ReplaceAllString(content, themeName)

	re = regexp.MustCompile(author)

	return re.ReplaceAllString(content, authorName)
}
