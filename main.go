package main

import (
	"fmt"
	"github.com/MaxKQC/Gator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	
	err = cfg.SetUser("max")
	if err != nil {
		log.Fatal(err)
	}

	new_cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(new_cfg)
}