package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
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
    "version": "0.0.0",
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
0.0.0
`}

var pagePHP = File{Name: "page.php",
	Contents: `<?php
/*
* File: mytheme/page.php
*/

if (!defined('ABSPATH')) {
	exit;
}

// Page content hook
add_action('wp', 'mytheme_page_content');
?>`}

var singlePHP = File{Name: "single.php",
	Contents: `<?php
/*
 * File: mytheme/single.php
 */
 
if (!defined('ABSPATH')) {
    exit;
}

// Single post content hook
add_action('wp', 'mytheme_single_content');
?>`}

var customFunc = File{Name: "custom-functions.php",
	Contents: `// mytheme/inc/custom-functions.php
function mytheme_custom_function() {
    // Function implementation goes here
}
`}

var customTemp = File{Name: "custom-templates.php",
	Contents: `
<?php
/**
 * File: mytheme/inc/custom-templates.php
 */
 
// Some custom template content
"Some Text"
`}

var indexPHP = File{Name: "index.php",
	Contents: `
<?php
/**
 * File: mytheme/index.php
 */
 
if (!defined('ABSPATH')) {
    exit;
}

// Index content hook
add_action('wp', 'mytheme_index_content');
?>
`}

var homePHP = File{Name: "home.php",
	Contents: `
<?php
/**
 * File: mytheme/home.php
 */
 
if (!defined('ABSPATH')) {
    exit;
}

// Home page content hook
add_action('wp', 'mytheme_home_content');
?>
`}

func main() {
	themeName := flag.String("name", "theme", "Theme for WordPress.")

	git := flag.String("git", "false", "Start git repository.")

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	author := flag.String("author", user.Username, "Author name")

	flag.Parse()

	rootFiles := []File{functionsPHP, headerPHP, footerPHP, styles, themeJSON, version, indexPHP, homePHP}
	templateFiles := []File{pagePHP, singlePHP}
	incFiles := []File{customFunc, customTemp}

	project := Theme{Name: strings.Replace(*themeName, " ", "-", -1), RootFiles: rootFiles, TemplateFiles: templateFiles, IncFiles: incFiles}

	os.Mkdir(project.Name, 0755)
	os.Mkdir(project.Name+"/template-parts", 0755)
	os.Mkdir(project.Name+"/inc", 0755)

	writeTemplate(project.Name, project.RootFiles[:], *author)
	writeTemplate(project.Name+"/template-parts", project.TemplateFiles[:], *author)
	writeTemplate(project.Name+"/inc", project.IncFiles[:], *author)

	if *git == "true" {
		gitSetup(project.Name)
	}

	fmt.Println("Theme started successfully!")
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
		_, err = file.WriteString(formatContent(authorName, themeName, files[i].Contents))
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

	fmt.Println("Git repo initialized successfully.")
}

func formatContent(authorName string, themeName string, content string) string {
	myTheme := "[Mm]y[Tt]heme"
	author := "[Yy]our [Nn]ame"

	re := regexp.MustCompile(myTheme)

	content = re.ReplaceAllString(content, themeName)

	re = regexp.MustCompile(author)

	content = re.ReplaceAllString(content, authorName)

	return content
}
