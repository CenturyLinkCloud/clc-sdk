package main

import (
	"fmt"

	"github.com/mikebeyer/clc-sdk/clc"
)

func main() {
	config, err := clc.EnvConfig()
	if err != nil {
		panic(err)
	}
	client := clc.New(config)
	resp, err := client.Auth()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
