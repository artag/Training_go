package main

import (
	"bufio" // Read data from STDIN input stream (os.Stdin)
	"bytes"
	"flag"
	"fmt"
	"html/template" // To replace the definition of the header and footer with a template.
	"io"            // To use io.Writer interface
	"io/ioutil"
	"os"
	"os/exec" // To execute a separate process
	"runtime" // Uses constant GOOS to determine OS
	"time"    // To add a small delay

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	EnvironmentVariable = "MDP_TEMPLATE"
	defaultTemplate     = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
  <p><i>{{ .PreviewMessage }}</i></p>
{{ .Body }}
  </body>
</html>
`
	// File permissions.
	// Readable and writable by the owner and only readable by anyone else.
	filemode = 0644
)

// Content type represents the HTML content to add into the template
type content struct {
	Title          string
	PreviewMessage string
	Body           template.HTML
}

func main() {
	// Parse flags
	stdin := flag.Bool("in", false, "Provide the input Markdown via STDIN")
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	templateFilename := flag.String("t", "", "Alternate template name")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" && !*stdin {
		flag.Usage()
		os.Exit(1)
	}

	if *stdin {
		if err := run_stdin(os.Stdin, *templateFilename, os.Stdout, *skipPreview); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if *filename != "" {
		if err := run_file(*filename, *templateFilename, os.Stdout, *skipPreview); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func run_stdin(r io.Reader, templateFilename string, out io.Writer, skipPreview bool) error {
	htmlData, err := parseContentFromStdin(r, templateFilename)
	if err != nil {
		return err
	}

	return run_internal(htmlData, templateFilename, out, skipPreview)
}

func run_file(filename string, templateFilename string, out io.Writer, skipPreview bool) error {
	htmlData, err := parseContentFromFile(filename, templateFilename)
	if err != nil {
		return err
	}

	return run_internal(htmlData, templateFilename, out, skipPreview)
}

func run_internal(htmlData []byte, templateFilename string, out io.Writer, skipPreview bool) error {
	// Create temporary file and check for errors
	// "" - system-defined temporary directory
	// "mdp*.html" - filename pattern
	temp, err := ioutil.TempFile("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outFilename := temp.Name()
	fmt.Fprintln(out, outFilename)

	if err := saveHtml(outFilename, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outFilename)

	return preview(outFilename)
}

func parseContentFromStdin(r io.Reader, templateFilename string) ([]byte, error) {
	var bytes bytes.Buffer

	s := bufio.NewScanner(r)
	for {
		s.Scan()
		line := s.Text()
		if len(line) == 0 {
			break
		}

		bytes.WriteString(line)
		bytes.WriteString("\n")
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	previewMessage := "Preview from stdin"
	return parseBytes_internal(bytes.Bytes(), templateFilename, previewMessage)
}

func parseContentFromFile(filename string, templateFilename string) ([]byte, error) {
	// Read all data from the input and check for errors
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	previewMessage := fmt.Sprintf("Preview filename: %s", filename)
	return parseBytes_internal(input, templateFilename, previewMessage)
}

func parseBytes_internal(input []byte, templateFilename string, previewMessage string) ([]byte, error) {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Parse the contents of the defaultTemplate const into a new Template
	tmpl := getDefaultTemplate()
	t, err := template.New("mdp").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	// If user provided alternate template file, replace template
	if templateFilename != "" {
		t, err = template.ParseFiles(templateFilename)
		if err != nil {
			return nil, err
		}
	}

	// Instantiate the content type, adding the title and body
	c := content{
		Title:          "Markdown Preview Tool",
		PreviewMessage: previewMessage,
		Body:           template.HTML(body),
	}

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer
	// Execute the template with the content type
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func getDefaultTemplate() string {
	env := os.Getenv(EnvironmentVariable)
	if env != "" {
		return env
	}

	return defaultTemplate
}

func saveHtml(outFilename string, data []byte) error {
	// Write the bytes to the file
	return ioutil.WriteFile(outFilename, data, filemode)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)
	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using default program
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)
	return err
}
