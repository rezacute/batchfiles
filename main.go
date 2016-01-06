package main

import (
  "os"
  "github.com/codegangsta/cli"
  "fmt"
)

func main() {
  app := cli.NewApp()
  app.Name = "batchfiles"
  app.Usage = "Batch File operation CLI"
  app.Commands = []cli.Command{
    {
            Name:        "rename",
            Usage:       "use it to batch rename",
            Description: "This will do batch operation for rename",
            Subcommands: []cli.Command{
                {
                    Name:        "files",
                    Usage:       "rename files in a folder",
                    Description: "rename onli files, fill skip directories",
                    Flags: []cli.Flag{
                        cli.StringFlag{
                            Name:  "source",
                            Value: "",
                            Usage: "The directory sorce path ",
                        },
                    },
                    Action: func(c *cli.Context) {
                        fmt.Println("renaming", c.String("source"))
                    },
                },
            },
        },
    }

  app.Run(os.Args)
}
