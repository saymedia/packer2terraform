
package packer2terraform

import (
    "fmt"
    "strings"
    "testing"
)

func csvToStrings(data string) (out [][]string) {
    splitData := strings.Split(data, "\n")
    for _, v := range splitData {
        out = append(out, strings.Split(v, ","))
    }
    return out
}

func TestFilterFail(t *testing.T) {
    data := csvToStrings(`1432149127,,ui,message,    amazon-ebs: [2015-05-20T19:12:07+00:00] FATAL: Chef::Exceptions::ChildConvergeError: Chef run process exited unsuccessfully (exit code 1)
1432149127,,ui,say,==> amazon-ebs: Terminating the source AWS instance...
1432149151,,ui,say,==> amazon-ebs: Deleting temporary keypair...
1432149151,,ui,error,Build 'amazon-ebs' errored: Error executing Chef: Non-zero exit status: 1
1432149151,,error-count,1
1432149151,,ui,error,\n==> Some builds didn't complete successfully and had errors:
1432149151,amazon-ebs,error,Error executing Chef: Non-zero exit status: 1
1432149151,,ui,error,--> amazon-ebs: Error executing Chef: Non-zero exit status: 1
1432149151,amazon-ebs,error,Error executing Chef: Non-zero exit status: 2
1432149151,,ui,error,--> amazon-ebs: Error executing Chef: Non-zero exit status: 2
1432149151,,ui,say,\n==> Builds finished but no artifacts were created.`)

    artifacts, err := Filter(data)
    if err == nil {
        t.Log("Error data didn't produce a filter error")
        t.Log("Error:", err)
        t.Fail()
    }
    if fmt.Sprintf("%s", err) != "Error executing Chef: Non-zero exit status: 1\nError executing Chef: Non-zero exit status: 2" {
        t.Log("Error data produced wrong flter error message")
        t.Log(fmt.Sprintf("%s", err))
        t.Fail()
    }
    if len(artifacts) > 0 {
        t.Log("Error data produced a filter artifact")
        t.Fail()
    }
}

func TestFilterSuccess(t *testing.T) {
    data := csvToStrings(`1432168580,,ui,say,==> amazon-ebs: Modifying attributes on AMI (ami-df79909b)...
1432168580,,ui,message,    amazon-ebs: Modifying: description
1432168581,,ui,say,==> amazon-ebs: Terminating the source AWS instance...
1432168589,,ui,say,==> amazon-ebs: Deleting temporary keypair...
1432168589,,ui,say,Build 'amazon-ebs' finished.
1432168589,,ui,say,\n==> Builds finished. The artifacts of successful builds are:
1432168589,amazon-ebs,artifact-count,1
1432168589,amazon-ebs,artifact,0,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,0,id,us-west-1:ami-df79909b
1432168589,amazon-ebs,artifact,0,string,AMIs were created:\n\nus-west-1: ami-df79909b
1432168589,amazon-ebs,artifact,0,files-count,0
1432168589,amazon-ebs,artifact,0,end
1432168589,,ui,say,--> amazon-ebs: AMIs were created:\n\nus-west-1: ami-df79909b`)

    artifacts, err := Filter(data)
    if err != nil {
        t.Log("Success data produced a filter error")
        t.Log("Error:", err)
        t.Fail()
    }
    if len(artifacts) == 0 {
        t.Log("Success data didn't produce a filter artifact")
        t.Fail()
    }

    artifact := artifacts[0]
    // t.Log(artifact)
    if artifact.BuilderType != "amazon-ebs" {
        t.Log("Success data didn't produce the correct builderType")
        t.Fail()
    }
    if artifact.Id != "us-west-1:ami-df79909b" {
        t.Log("Success data didn't produce the correct id")
        t.Fail()
    }
}

func TestFilterMultiSuccess(t *testing.T) {
    data := csvToStrings(`1432168589,amazon-ebs,artifact-count,2
1432168589,amazon-ebs,artifact,0,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,0,id,us-west-1:ami-df79909b
1432168589,amazon-ebs,artifact,0,string,AMIs were created:\n\nus-west-1: ami-df79909b
1432168589,amazon-ebs,artifact,0,files-count,0
1432168589,amazon-ebs,artifact,0,end
1432168589,amazon-ebs,artifact,1,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,1,id,us-west-2:ami-df79909c
1432168589,amazon-ebs,artifact,1,string,AMIs were created:\n\nus-west-2: ami-df79909c
1432168589,amazon-ebs,artifact,1,files-count,0
1432168589,amazon-ebs,artifact,1,end`)

    artifacts, err := Filter(data)
    if err != nil {
        t.Log("Success data produced a filter error")
        t.Log("Error:", err)
        t.Fail()
    }
    if len(artifacts) == 0 {
        t.Log("Success data didn't produce a filter artifact")
        t.Fail()
    }
    if len(artifacts) < 2 {
        t.Log("Success data didn't produce (2) filter artifacts")
        t.Fail()
    }
    if len(artifacts) > 2 {
        t.Log("Success data produced too many artifacts")
        t.Fail()
    }

    artifact := artifacts[1]
    // t.Log(artifact)
    if artifact.BuilderType != "amazon-ebs" {
        t.Log("Success data didn't produce the correct builderType")
        t.Fail()
    }
    if artifact.Id != "us-west-2:ami-df79909c" {
        t.Log("Success data didn't produce the correct id")
        t.Fail()
    }
}

func TestToTerraformVars(t *testing.T) {
    data := csvToStrings(`1432168589,amazon-ebs,artifact-count,2
1432168589,amazon-ebs,artifact,0,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,0,id,us-west-1:ami-df79909b
1432168589,amazon-ebs,artifact,0,string,AMIs were created:\n\nus-west-1: ami-df79909b
1432168589,amazon-ebs,artifact,0,files-count,0
1432168589,amazon-ebs,artifact,0,end
1432168589,amazon-ebs,artifact,1,builder-id,mitchellh.amazonebs
1432168589,amazon-ebs,artifact,1,id,us-west-2:ami-df79909c
1432168589,amazon-ebs,artifact,1,string,AMIs were created:\n\nus-west-2: ami-df79909c
1432168589,amazon-ebs,artifact,1,files-count,0
1432168589,amazon-ebs,artifact,1,end`)
    out := `variable "images" {
    default = {
            us-west-1 = "ami-df79909b"
            us-west-2 = "ami-df79909c"
    }
}`

    artifacts, err := Filter(data)
    vars, err := ToTerraformVars(artifacts)
    if err != nil {
        t.Log("Vars transform produced an error")
        t.Log("Error:", err)
        t.Fail()
    }
    if vars != out {
        t.Log("Vars transform didn't produce correct output")
        t.Log("Output:", vars)
        t.Fail()
    }
}
