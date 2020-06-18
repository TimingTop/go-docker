package main

import (
	"github.com/urfave/cli/v2"
	"github.com/go-docker/container"
	"fmt"
)

var runCommand = cli.Command {
	Name: "run",
	Usage: `go-docker run xxx`,
	Flags: []cli.Flag {
		&cli.BoolFlag {
			Name: "it",
			Usage: "enable tty",
		},
		&cli.StringFlag {
			Name: "m",
			Usage: "memory limit",
		},
		&cli.StringFlag {
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		&cli.StringFlag {
			Name: "cpuset",
			Usage: "cpuset limit",
		},
		&cli.StringFlag {
			Name: "v",
			Usage: "volume",
		},
	},
	Action: func(c *cli.Context) error {

		fmt.Println("run command.....")

		
		// if len(c.Args().Slice()) < 1 {
		// 	return fmt.Errorf("Missing command")
		// }
		if c.Args().Len() < 1 {
			return fmt.Errorf("Missing command")
		}

		
		// 把 标签以外的所有参数读出来，第一个参数
		// 是 image 的名字
		var cmdArray []string
		for _, arg := range c.Args().Slice() {
			fmt.Println(arg)
			cmdArray = append(cmdArray, arg)
		}

		// 获取镜像名字
		// imageName := cmdArray[0]
		// cmdArray = cmdArray[1:]
		
		createTty := c.Bool("it")
		// detach := c.Bool("d")
		volume := c.String("v")
		// cmd := c.Args().Get(0)
			
		Run(createTty, cmdArray, volume)
		return nil
	},

}

var initCommand = cli.Command {
	Name: "init",
	Usage: "init container process ",
	Action: func(c *cli.Context) error {
		fmt.Println("init command....")
		cmd := c.Args().Get(0)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}