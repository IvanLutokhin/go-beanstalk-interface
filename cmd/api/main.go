package main

import (
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)

			os.Exit(1)
		}
	}()

	api.New().Run()
}
