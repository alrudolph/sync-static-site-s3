package main

import (
	"github.com/alrudolph/snyc-static-site-s3/cmd"
	_ "github.com/alrudolph/snyc-static-site-s3/cmd/config"
	_ "github.com/alrudolph/snyc-static-site-s3/cmd/setup"
)

func main() {
	cmd.Execute()
}
