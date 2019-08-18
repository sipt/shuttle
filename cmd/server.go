package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sipt/shuttle/conf"
)

func main() {
	ctx := context.Background()
	params := map[string]string{"path": "config1.toml"}
	config, err := conf.LoadConfig(ctx, "file", "toml", params, func() {
		fmt.Println("config file change")
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*config)
	err = http.ListenAndServe(":8080", &handler{})
	if err != nil {
		panic(err)
	}
}

type handler struct {
	http.Handler
}
