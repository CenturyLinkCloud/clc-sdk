package main

import (
	"encoding/json"
	"fmt"

	"github.com/mikebeyer/clc-sdk/clc"
)

func main() {
	config, err := clc.EnvConfig()
	if err != nil {
		panic(err)
	}
	client := clc.New(config)
	resp, err := client.GetDatacenter("VA1")
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", b[:])

}
