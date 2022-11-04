//go:build generate
// +build generate

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"
)

const (
	filename = "sizes_gen.go"
)

type fgResource struct {
	Cpu    int
	Memory int
}

type TemplateData struct {
	Sizes       []fgResource
	SmallestCPU int
}

func main() {
	fmt.Printf("Generating internal/fargate/%s\n", filename)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file (%s): %s", filename, err)
	}

	tplate, err := template.New("fargetsizes").Parse(tmpl)
	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	sizeList := generateSizes()

	td := TemplateData{
		Sizes:       sizeList,
		SmallestCPU: sizeList[0].Cpu,
	}

	var buffer bytes.Buffer
	err = tplate.Execute(&buffer, td)
	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	contents, err := format.Source(buffer.Bytes())
	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	if _, err := f.Write(contents); err != nil {
		f.Close()
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("error closing file (%s): %s", filename, err)
	}
}

func generateSizes() []fgResource {

	type fargateSizeDef struct {
		cpu      int
		memStart int
		memEnd   int
		memInc   int
	}

	sizelist := []fgResource{}

	// https://docs.aws.amazon.com/AmazonECS/latest/userguide/task_definition_parameters.html
	sizelist = append(
		sizelist,
		fgResource{Cpu: 256, Memory: 512},
		fgResource{Cpu: 256, Memory: 1024},
		fgResource{Cpu: 256, Memory: 2048},
	)

	steps := []fargateSizeDef{
		{cpu: 512, memStart: 1, memEnd: 4, memInc: 1},
		{cpu: 1024, memStart: 2, memEnd: 8, memInc: 1},
		{cpu: 2048, memStart: 4, memEnd: 16, memInc: 1},
		{cpu: 4096, memStart: 8, memEnd: 30, memInc: 1},
		{cpu: 8192, memStart: 16, memEnd: 60, memInc: 4},
		{cpu: 16384, memStart: 32, memEnd: 120, memInc: 8},
	}

	for _, step := range steps {
		for i := step.memStart; i <= step.memEnd; i += step.memInc {
			sizelist = append(sizelist, fgResource{
				Cpu:    step.cpu,
				Memory: i * 1024,
			})
		}
	}
	return sizelist
}

//go:embed file.tmpl
var tmpl string
