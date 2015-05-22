
package main

import (
    "bufio"
    "fmt"
    "os"
    "github.com/saymedia/packer2terraform/packer2terraform"
)

func main() {

    stats, _ := os.Stdin.Stat()
    if stats.Size() == 0 {
        // No Stdin data
        fmt.Println(`packer2terraform`)
        fmt.Println(`example: packer -machine-readable build app.json | packer2terraform > app.tfvars`)
        os.Exit(0)
    }

    reader := bufio.NewReader(os.Stdin)

    parsed, err := packer2terraform.Parse(reader)
    if err != nil {
        fmt.Printf("File read failed %s", err)
        os.Exit(1)
    }

    artifacts, err := packer2terraform.Filter(parsed)
    if err != nil {
        // fmt.Errorf("Packer build failed: %s", err)
        fmt.Printf("Packer build failed: %s", err)
        os.Exit(2)
    }

    vars, err := packer2terraform.ToTerraformVars(artifacts)
    if err != nil {
        os.Exit(3)
    }

    fmt.Println(vars)
    os.Exit(0)

}
