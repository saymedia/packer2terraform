
package packer2terraform

import (
    "bytes"
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "strconv"
    "strings"
    "text/template"
)


type LogLine struct {
    timestamp     string
    builderTarget string
    lineType      string
    messageType   string
    messageTypeI  int
    messageA      string
    messageB      string
}

type Artifact struct {
    BuilderTarget string
    BuilderId     string
    Id            string
    IdSplit       []string
    Message       string
    FilesCount    string
}

type TemplatePage struct {
    Artifacts []Artifact
}


var TemplateAmazonEBS = `variable "images" {
    default = {
{{range .Artifacts}}
        {{index .IdSplit 0}} = "{{index .IdSplit 1}}"{{end}}
    }
}`;


func ReadCSV(csvReader io.Reader) (ret [][]string, err error) {
    reader := csv.NewReader(csvReader)
    reader.FieldsPerRecord = -1
    reader.LazyQuotes = true
    return reader.ReadAll()
}

func Filter(parsed [][]string) (artifacts []Artifact, err error) {
    var errorCount int
    var errorMsg []string
    var artifactCount int

    for _, v := range parsed {
        // Build a LogLine
        line := LogLine{"", "", "", "", 0, "", ""}
        if len(v) > 0 {
            line.timestamp = v[0]
        }
        if len(v) > 1 {
            line.builderTarget = v[1]
        }
        if len(v) > 2 {
            line.lineType = v[2]
        }
        if len(v) > 3 {
            line.messageType = v[3]
        }
        if len(v) > 4 {
            line.messageA = v[4]
        }
        if len(v) > 5 {
            line.messageB = v[5]
        }
        if len(line.messageType) > 0 {
            line.messageTypeI, _ = strconv.Atoi(line.messageType)
        }

        // Artifacts:
        if line.lineType == "artifact-count" {
            artifactCount = line.messageTypeI
        }
        if line.lineType == "artifact" {

            if len(artifacts) < line.messageTypeI+1 {
                a := Artifact{}
                a.BuilderTarget = line.builderTarget
                artifacts = append(artifacts, a)
            }

            a := &artifacts[line.messageTypeI]
            if line.messageA == "id" {
                a.Id = line.messageB
                a.IdSplit = strings.Split(line.messageB, ":")
            }
            if line.messageA == "files-count" {
                a.FilesCount = line.messageB
            }
            if line.messageA == "builder-id" {
                a.BuilderId = line.messageB
            }
            if line.messageA == "string" {
                a.Message = line.messageB
            }
        }

        // Errors:
        if line.lineType == "error-count" && line.messageTypeI > 0 {
            errorCount = line.messageTypeI
        }
        if line.lineType == "error" {
            errorMsg = append(errorMsg, line.messageType)
        }
    }

    if artifactCount < len(artifacts) {
        artifactsMissing := artifactCount - len(artifacts)
        return nil, errors.New(fmt.Sprintf("Missing %s artifacts.", artifactsMissing))
    }

    if errorCount > 0 && len(errorMsg) > 0 {
        return nil, errors.New(strings.Join(errorMsg, "\n"))
    }

    // Clean up empty artifacts
    for i, artifact := range artifacts {
        if artifact.Id == "" {
            artifacts = append(artifacts[:i], artifacts[i+1:]...)
        }
    }
    if len(artifacts) == 0 {
        return nil, errors.New("No Artifacts found.")
    }

    return artifacts, nil
}

func ToTemplate(artifacts []Artifact, tmpl string) (ret string, err error) {
    // Setup the page vars
    var thePage = TemplatePage{}
    thePage.Artifacts = artifacts

    t := template.Must(template.New("tmpl").Parse(tmpl))

    var doc bytes.Buffer
    t.Execute(&doc, thePage)
    ret = doc.String()

    return ret, nil
}
