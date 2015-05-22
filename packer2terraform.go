
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "strconv"
    "strings"
)

type LogLine struct {
    time         string
    builderType  string
    lineType     string
    messageType  string
    messageTypeI int
    messageA     string
    messageB     string
}

type Artifact struct {
    BuilderType string
    BuilderId   string
    Id          string
    Message     string
    FilesCount  string
}

func Parse(csvReader io.Reader) (ret [][]string, err error) {
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
        // line := new(LogLine)
        line := LogLine{v[0], v[1], v[2], v[3], 0, "", ""}
        if len(v) > 4 {
            line.messageA = v[4]
        }
        if len(v) > 5 {
            line.messageB = v[5]
        }
        if len(line.messageType) > 0 {
            line.messageTypeI, _ = strconv.Atoi(line.messageType)
        }
        // fmt.Printf("Line: %s", line)

        // Artifacts:
        if line.lineType == "artifact-count" {
            artifactCount = line.messageTypeI
            // fmt.Printf("Artifact Count: %d\n", artifactCount)
        }
        if line.lineType == "artifact" {

            if len(artifacts) < line.messageTypeI+1 {
                a := Artifact{}
                a.BuilderType = line.builderType
                artifacts = append(artifacts, a)
            }

            a := &artifacts[line.messageTypeI]
            if line.messageA == "id" {
                a.Id = line.messageB
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

    if len(artifacts) == 0 {
        return nil, errors.New("No Artifact Id found.")
    }

    return artifacts, nil
}

func ToTerraformVars(artifacts []Artifact) (ret string, err error) {
    // fmt.Println("ttv:", artifacts)
    retSlice := []string{}

    tfTemplate := `variable "images" {
    default = {
%s
    }
}`

    for _, artifact := range artifacts {
        // ret = fmt.Sprintf("id=%s", artifact.Id)
        // fmt.Println("id=", artifact)
        if artifact.BuilderType == "amazon-ebs" {
            stringSplit := strings.Split(artifact.Id, ":")
            // fmt.Sprintf("ss=%s", stringSplit)
            retSlice = append(retSlice, fmt.Sprintf("            %s = \"%s\"", stringSplit[0], stringSplit[1]))
        }
    }

    ret = strings.Join(retSlice, "\n")
    ret = fmt.Sprintf(tfTemplate, ret)
    // fmt.Println("ttv:", ret)
    return ret, nil
}
