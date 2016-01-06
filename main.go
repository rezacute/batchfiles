package main

import (
  "os"
  "github.com/codegangsta/cli"
  "fmt"
  "path/filepath"
  "log"
  "strings"
)

var (
  ext string
  src_dir string
  in_prefix string
  out_prefix string
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
                            Usage: "The directory source path ",
                        },
                        cli.StringFlag{
                            Name:  "src-prefix",
                            Value: "",
                            Usage: "The prefix source filename ",
                        },
                        cli.StringFlag{
                            Name:  "add-prefix",
                            Value: "",
                            Usage: "The prefix dest filename ",
                        },
                        cli.StringFlag{
                            Name:  "src-extension",
                            Value: "",
                            Usage: "The filter extension source filename ",
                        },
                    },
                    Action: func(c *cli.Context) {
                        src_dir = c.String("source")

                        if src_dir == ""{
                          return
                        }
                        ext = c.String("src-extension")
                        in_prefix = c.String("src-prefix")
                        out_prefix = c.String("add-prefix")
                        fmt.Println("extension", ext)

                        filepath.Walk(src_dir, visit)

                    },
                },
            },
        },
    }

  app.Run(os.Args)
}

func visit(path string, f os.FileInfo, err error) (e error) {
    if filepath.Ext(path) != ext || !strings.HasPrefix(f.Name(), in_prefix){

        return
    }
    dir := filepath.Dir(path)
    base := filepath.Base(path)
    newname := filepath.Join(dir, out_prefix + base)
    log.Printf("mv \"%s\" \"%s\"\n", path, newname)
    os.Rename(path, newname)
    return
}
