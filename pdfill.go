package pdfill

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const pdftool = "pdftk"

// Form holds field names and field values
type Form map[string]interface{}

// FdfContent returns fdf content of PDF form
func (f Form) FdfContent() string {
	var content string
	for k, v := range f {
		content += fmt.Sprintf("<</T(%s)/V(%v)>>", k, v)
	}
	return content
}

// Fill pdf file with Form and write filled file to new destination
func Fill(input Form, src string, dest string) error {
	//check pdftk command
	_, err := exec.LookPath(pdftool)
	if err != nil {
		return fmt.Errorf("pdftk is not installed")
	}

	//check src file

	//path for fdf file
	tmpDir, err := ioutil.TempDir("", "pdfill-")
	if err != nil {
		return err
	}

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	fdfFile := filepath.Clean(tmpDir + "/data.fdf")

	err = createFdfFile(fdfFile, input)
	if err != nil {
		return err
	}

	return runCommand(pdftool, src, "fill_form", fdfFile, "output", dest, "flatten")
}

func createFdfFile(fdfFile string, input Form) (err error) {
	file, err := os.Create(fdfFile)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.WriteString(file, fdfHeader+input.FdfContent()+fdfFooter)
	return
}

func runCommand(bn string, args ...string) error {
	cmd := exec.Command(bn, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(stderr.String())
	}

	return nil
}

const fdfHeader = `%FDF-1.2
1 0 obj<</FDF<< /Fields[`

const fdfFooter = `] >> >>
endobj
trailer
<</Root 1 0 R>>
%%EOF`
