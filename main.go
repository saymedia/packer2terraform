package main

import (
	"bufio"
	"flag"
	"github.com/saymedia/packer2terraform/packer2terraform"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func help() {
	log.Println(`Usage packer2terraform [options...]
packer2terraform turns Packer's machine-readable output into a Terraform-readable format.

Options:
    -f Filename of the input CSV. Alternatively use STDIN.
    -h This help information.
    -template Filename of the template to use in the output.

Example:
    packer -machine-readable build app.json | \
        packer2terraform -template templates/amazon-ebs.hcl > app.tfvars
`)
}

func main() {

	tmpl := flag.String("template", "", "a template file")
	csv := flag.String("f", "", "a csv file")
	helpMe := flag.Bool("h", false, "help")

	flag.Parse()

	if *helpMe {
		help()
		os.Exit(0)
	}

	// Read a file or use STDIN
	var csvFile io.Reader
	if len(*csv) > 0 {
		f, err := os.Open(*csv)
		if err != nil {
			log.Fatalf("CSV file read failed %s", err)
		}
		csvFile = bufio.NewReader(f)
	} else if stat, _ := os.Stdin.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
		// has STDIN data
		csvFile = bufio.NewReader(os.Stdin)
	} else {
		// No input data
		help()
		os.Exit(0)
	}

	// Get the CSV as a string array
	parsed, err := packer2terraform.ReadCSV(csvFile)
	if err != nil {
		log.Fatalf("CSV read failed %s", err)
	}

	// Extract the artifacts
	artifacts, err := packer2terraform.Filter(parsed)
	if err != nil {
		// fmt.Errorf("Packer build failed: %s", err)
		log.Fatalf("Packer build failed: %s", err)
	}

	// Print artifacts using a template
	var templateString string
	if len(*tmpl) == 0 {
		templateString = packer2terraform.TemplateAmazonEBS
	} else {
		buf, err := ioutil.ReadFile(*tmpl)
		if err != nil {
			log.Fatalf("Template file read failed: %s", err)
		}
		templateString = string(buf)
	}
	doc, err := packer2terraform.ToTemplate(artifacts, templateString)
	if err != nil {
		log.Fatalf("Template render failed: %s", err)
	}
	log.Println(doc)

	// Done
	os.Exit(0)

}
