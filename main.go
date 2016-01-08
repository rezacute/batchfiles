package main

import (
  "os"
  "github.com/codegangsta/cli"
  "path/filepath"
  "fmt"
  "strings"
  "text/template"
  "bytes"
  "bufio"
  "container/list"
)
type Snippet struct{
  NONBLOCKING string
  BLOCKING string
}
var (
  ext string
  src_dir string
  in_prefix string
  out_prefix string
  dst_dir string
  success_count int
  skip_count int
)
func main() {
  success_count = 0
  skip_count = 0
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

        if src_dir == "" {

          return
        }
        ext = c.String("src-extension")
        in_prefix = c.String("src-prefix")
        out_prefix = c.String("add-prefix")
        filepath.Walk(src_dir, addPrefix)
      },
    },
  },
},
{
Name:        "blend",
Usage:       "use it to batch copy with pattern",
Description: "use '_' as separator",
Subcommands: []cli.Command{
  {
  Name:        "merge",
  Usage:       "copy ",
  Description: "rename only files, fill skip directories",
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
      Name:  "base_destination",
      Value: "",
      Usage: "The base destination directory ",
    },
    cli.StringFlag{
      Name:  "src-extension",
      Value: "",
      Usage: "The filter extension source filename ",
    },
  },
  Action: func(c *cli.Context) {
    in_prefix = c.String("src-prefix")
    src_dir = c.String("source")
    dst_dir = c.String("base_destination")
    if src_dir == ""{
      src_dir,_ = os.Getwd()
    }
    if dst_dir == ""{
      dst_dir,_ = os.Getwd()
    }
    ext = ".swift"
    filepath.Walk(src_dir, blendMerge)
    fmt.Println("Success :",success_count)
    fmt.Println("Skip :",skip_count)
  },
},
},
},
}

app.Run(os.Args)
}

func addPrefix(path string, f os.FileInfo, err error) (e error) {
  if filepath.Ext(path) != ext || !strings.HasPrefix(f.Name(), in_prefix){
    return
  }
  dir := filepath.Dir(path)
  base := filepath.Base(path)
  newname := filepath.Join(dir, out_prefix + base)
  os.Rename(path, newname)
  return
}
func blendMerge(path string, f os.FileInfo, err error) (e error) {

  if filepath.Ext(path) != ext || !strings.HasPrefix(f.Name(), in_prefix){
    return
  }
  dir := filepath.Dir(path)
  var base string
  if strings.HasPrefix(f.Name(),"guides_ab-"){
    base = strings.Replace(f.Name(),"_","/",2)
    }else{
      base = strings.Replace(f.Name(),"_","/",-1)
    }
    r_name := filepath.Join(dir, f.Name())
    filename := strings.Replace(filepath.Join(dst_dir, base),".swift",".mkd",1)
    //log.Printf("merge \"%s\"\n", filename)
    f1, err := os.OpenFile(r_name, os.O_RDONLY, 0666)
    if err != nil {
      panic(err)
    }

    defer f1.Close()
    scanner := bufio.NewScanner(f1)
    var x list.List
    isParsingSingle := false
    str := ""
    blocking := ""
    non_blocking := ""
    var snip Snippet
    for scanner.Scan() {

      if strings.HasPrefix(scanner.Text(),"private func snippet") {
        str = ""
        //log.Println("start")
        if strings.HasSuffix(scanner.Text(),"blocking(){") {
          if strings.HasSuffix(scanner.Text(),"non_blocking(){") {
            non_blocking = " "
            }else{
              snip = Snippet{}
              blocking = " "
            }
            } else{
              isParsingSingle = true
              snip = Snippet{}
              blocking = " "
            }
            continue
          }
          if strings.HasPrefix(scanner.Text(),"}"){
            if blocking != "" {
              snip.BLOCKING = str
              blocking = ""
            }
            if non_blocking != "" {
              snip.NONBLOCKING = str
              non_blocking = ""
              x.PushBack(snip)
            }

            //log.Println("stop")
            //log.Println("\n\n")
            if isParsingSingle {
              isParsingSingle = false
              x.PushBack(snip)
            }
            continue
          }
          str = str +"\n"+ strings.Replace(scanner.Text(),"    "," ",3)

          //log.Println(scanner.Text())
        }
        f2, err := os.OpenFile(filename, os.O_RDONLY, 0666)
        if err != nil {
          panic(err)
        }

        defer f2.Close()
        scanner = bufio.NewScanner(f2)
        shouldSkip := false
        var element = x.Front()
        str = ""
        for scanner.Scan() {
          if strings.HasPrefix(scanner.Text(),"**Swift:**") {
            fmt.Println("Skipped ",filename)
            skip_count++
            return
          }
          if scanner.Text() == "#### Android" ||
          scanner.Text() == "#### Javascript" ||
          strings.HasPrefix(scanner.Text(),"```csharp") ||
          strings.HasPrefix(scanner.Text(),"```java") {
          shouldSkip = true
          }else if strings.HasPrefix(scanner.Text(),"#### iOS") ||
          strings.HasPrefix(scanner.Text(),"```objc") {
          shouldSkip = false
        }
        if scanner.Text() == "```objc" || strings.HasPrefix(scanner.Text(),"{% tabcontrol %}"){
          if element != nil && !shouldSkip{
            str = str + "\n**Objective-C:**\n"
          }

        }
        if str=="" {
          str = str + scanner.Text()
          }else{
            str = str +"\n"+ scanner.Text()
          }

          if element == nil || shouldSkip{
            continue;
          }
          if scanner.Text() == "```" || strings.HasPrefix(scanner.Text(),"{% endtabcontrol %}"){

            snip = element.Value.(Snippet)
            text := writeToString(snip)
            str = str +"\n"+ text
            element = element.Next()

          }

        }

        f3, err := os.OpenFile(filename, os.O_WRONLY, 0666)
        if err != nil {
          panic(err)
        }

        defer f3.Close()
        _, err = f3.WriteString(str)
        fmt.Println("Completed",filename)
        success_count++
        return
      }


const tmpl_tab = `
**Swift:**

{% tabcontrol %}

{% tabpage Blocking API %}
{% highlight swift %}
{{.BLOCKING}}
{% endhighlight %}
{% endtabpage %}

{% tabpage Non-Blocking API %}
{% highlight swift %}
{{.NONBLOCKING}}
{% endhighlight %}
{% endtabpage %}

{% endtabcontrol %}
`
const tmpl_single = `
**Swift:**

` + "```swift" +`
{{.BLOCKING}}
` + "```" + `
`

  func writeToString(snippet Snippet) (result string){
    var tmpl string
    if snippet.NONBLOCKING != "" {
      tmpl = tmpl_tab
      }else{
        tmpl = tmpl_single
      }
      t, err := template.New("person").Parse(tmpl)

      if err == nil {
        buff := bytes.NewBufferString("")
        t.Execute(buff, snippet)
        result = buff.String()
      }
      return
    }
