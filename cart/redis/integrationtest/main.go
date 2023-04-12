package main

import (
	"os"

	"flamingo.me/flamingo-commerce-contrib/cart/redis/integrationtest/helper"
)

func main() {
	if os.Getenv("RUN") == "1" {
		info := helper.BootupDemoProject("../../config/")
		<-info.Running
	} else {
		helper.GenerateGraphQL()
	}
}
