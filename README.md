# packer2terraform

packer2terraform turns Packer's [machine-readable output](https://packer.io/docs/command-line/machine-readable.html) into [Terraform-readable tfvars](https://terraform.io/docs/configuration/variables.html). For example, you have Packer build an AMI that Terraform deploys to AWS.

## Usage

packer2terraform reads from STDIN and writes to STDOUT.

    packer -machine-readable build app.json | packer2terraform > app.tfvars

## Install

    go get github.com/saymedia/packer2terraform

## Test

    go test ./..

## License

Copyright © 2015 Say Media Ltd. All Rights Reserved. See the LICENSE file for distribution terms.