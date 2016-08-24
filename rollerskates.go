package main

import (
	"os"
	"github.com/urfave/cli"
	"madetech.com/rollerskates/aws"
	"fmt"
	"github.com/joho/godotenv"
)


func main() {
	app := cli.NewApp()
	app.Name = "rollerskates"
	app.Usage = "rollerskates load-balancer-name"

	app.Action = func(c *cli.Context) error {
		godotenv.Load()
		arg1 := c.Args().Get(0)
		arg2 := c.Args().Get(1)
		fmt.Println( arg1 )
		fmt.Println( arg2 )
		status := rollerskates.DeregisterInstancesFromLoadBalancer(arg1, arg2)
		fmt.Println(status)
		return nil
	}

	app.Run(os.Args)
}