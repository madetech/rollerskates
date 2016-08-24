package rollerskates

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
)

func main() {
	app := cli.NewApp()
	app.Name = "rollerskates"
	app.Usage = "rollerskates load-balancer-name"

	app.Action = func(c *cli.Context) error {
		fmt.Println("Bear with me whilst I put my skates on...")
		//RestartLoadBalancerInstances(c.Args().Get(0))
		return nil
	}

	app.Run(os.Args)
}