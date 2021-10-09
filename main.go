package main

import (
	"fmt"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigtoml"

	"github.com/simba-fs/go-oauth/server"
	"github.com/simba-fs/go-oauth/types"
)

func main() {
	config := &types.Config{}
	loader := aconfig.LoaderFor(config, aconfig.Config{
		SkipEnv:    true,
		FlagPrefix: "",
		Files:      []string{"config.toml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".toml": aconfigtoml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		panic(err)
	}

	fmt.Printf("config: %#v\n", *config)

	server.New(config)
}
