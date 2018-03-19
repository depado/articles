title: Some Go snippets I find useful
description: Here are some tricks and best practices for you
slug: some-go-snippet-i-find-useful
author: Depado
date: 2015-07-23 09:11:00
tags:
    - go
    - dev
    - snippet

# Reading a whole line on stdin

```go
var stdioScanner = bufio.NewScanner(os.Stdin)

func simpleReadLine() (string, error) {
	if !stdioScanner.Scan() {
		return "", stdioScanner.Err()
	}
	return stdioScanner.Text(), nil
}
```
I thought it was a lot more simple to read a whole line on `stdin` actually, because `fmt.Scanln`'s name is quite... Self explanatory I guess. But, according to the [fmt package's GoDoc](https://golang.org/pkg/fmt/#Scanln) :

> Scanln is similar to Scan, but stops scanning at a newline and after the final item there must be a newline or EOF.

So let's have a look at the `fmt.Scan` GoDoc :

> Scan scans text read from standard input, storing successive space-separated values into successive arguments. Newlines count as space. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why.  

So `fmt.Scanln` would be more useful to do something like a IRC client. For example this line in weechat `/server add stuff irc.stuff.com -ssl` would be easily parsed by `fmt.Scanln`.

**Note about the changes made to this function**

I had a really interesting conversation with [Axel Wagner](https://plus.google.com/u/0/+AxelWagner_Merovius/posts) about why this was a terrible idea to define the reader inside the function. The fact is that the `bufio` package, as its name suggests it, uses buffers. The scanner would read more than only one line. The fact that the scanner was declared inside the function body would lost the said scanner's scope at the end of the function, so there would be data loss. Axel wrote [a small example on the Go Playground](http://play.golang.org/p/vcbczoIuSO) that shows this particular behaviour. You can find the whole conversation [here](https://plus.google.com/114932755645700075856/posts/SuAKruB9F95).
He also explained to me why he thought naked returns are a bad idea, and his arguments are actually quite convincing.

---------------------
# Function to unmarshall a Json URL to a struct

```go
func fetchURL(url string, out interface{}) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(out)
	return
}
```

Now your `out` argument must obviously be a pointer because otherwise there would be no point in doing that. Let's have a look at a small example !
*Note : I already have a struct named GithubUserAPIType which is generated using [json-to-go](http://mholt.github.io/json-to-go/) which is an amazing tool created by [mholt](https://github.com/mholt).*

```go
func main() {
    var err error
    var me GithubUserAPIType

    err = fetchURL("https://api.github.com/users/Depado", &me)
    if err != nil {
        log.Fatal(err)
    }
}
```

Of course you need to make sure that your Json URL returns the right data and can be stored in your struct.

---------------------------
# Load a yaml configuration file to a struct

One thing I commonly do with all my projects at one point is to provide an easy way to configure the program without opening the code itself. Instead of configuring the program using command line arguments, you can use a yaml file. Why yaml ? I don't know. I guess I kind of like the way yaml is structured and find it easy to create configuration files using this language.

```go
import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func Load(f string, i interface{}) error {
	conf, err := ioutil.ReadFile(f)
	if err != nil {
        return err
	}
	err = yaml.Unmarshal(conf, &i)
	if err != nil {
		return err
	}
}
```

For example, let's say the above code is stored in the `configuration` package. Here is an example on how to use it :

```
# conf.yml
host: irc.freenode.net
port: 6697
name: b0t
```

```go
// main.go
package main

import (
    "log"
    "place/where/you/stored/configuration"
)

type Configuration struct {
    Host string
    Port string
    Name string
}

var Conf = new(Configuration)

func main() {
    var err error
    err = configuration.Load("conf.yml", &Conf)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(Conf.Host, Conf.Port, Conf.Name)
}
```

---------------------------
# Calculate the Md5Sum for a file

```go
// Generate the md5sum of a file.
func GenerateMd5Sum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return string(hash.Sum(nil)), nil
}
```

---------------
# Create custom loggers to write to different files

```go
logfile, err := os.OpenFile("custom.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
if err != nil {
	log.Fatal(err)
}
defer logfile.Close()

custom = log.New(logfile, "", log.Ldate|log.Ltime)
custom.Println("Hello World !")
```
