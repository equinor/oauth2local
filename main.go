package main

import (
	"github.com/equinor/oauth2local/cmd"

	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)
	cmd.Execute()
}
