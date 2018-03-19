title: Writing a markov-chain IRC bot in Go
description: A way to make your IRC bot awesome.
slug: writing-a-markov-chain-irc-bot-in-go
author: Depado
date: 2015-10-10 14:42:00
tags:
    - dev
    - irc
    - go

# Introduction

*The original idea was not mine. While looking for tutorials on markov chains, I stumbled across [this article](http://charlesleifer.com/blog/building-markov-chain-irc-bot-python-and-redis/) which is great and helped me a lot understanding what was actually going on.*

For the past few years, I've been aggregating a quite a lot of IRC logs and I was wondering what to do with those. Markov chains looked like a pretty cool principle so I decided to work on that for a few days and see if I could use those logs to create random sentences from them.

# A basic Go bot
*Configuration file handling and writing the base of the bot*

First of all, let's start by creating the `conf` package (`conf` being short for configuration obviously). It's quite a handy little snippet of code that I usually include in all my projects that needs configuration files. The configuration file format is `yaml`.

Now create a new folder called `conf`, and edit `conf.go`.

```go
// project/conf/conf.go
package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Configuration is the main struct that represents a configuration.
type Configuration struct {
	Server      string
	Channel     string
	BotName     string
	TLS         bool
	InsecureTLS bool
}

// C is the Configuration instance that will be exposed to the other packages.
var C = new(Configuration)

// Load parses the yml file passed as argument and fills the Config.
func Load(cp string) error {
	conf, err := ioutil.ReadFile(cp)
	if err != nil {
		return fmt.Errorf("Conf : Could not read configuration : %v", err)
	}
	if err = yaml.Unmarshal(conf, &C); err != nil {
		return fmt.Errorf("Conf : Error while parsing yaml : %v", err)
	}
	return nil
}
```

This package is intented to work that way : Load a configuration file, and fill a `C` instance (which is exposed outside the package) of the `Configuration` struct. Now create a new file at the root of your project and name it `conf.yml`. The configuration file must contain the following information :

```yaml
botname: my-bot-nickname
server: irc.freenode.net:6667
channel: #my-awesome-channel
tls: false
insecuretls: false
```

I'd advise you not to version that file as it may contain information about the channels you go to and/or the servers you're connected on.  

Now, back to the `project/main.go` file, we will load the configuration file named `conf.yml` placed at the root of your project. We will also initialise the IRC connection to suit your configuration file. Note that in the conf.yml file, you'll also define whether or not to use TLS. If you decide to use TLS, don't forget to change the port of the server in the server variable. Also if the server you're connecting to doesn't have a valid certificate or whatever, set the insecuretls variable to true.
```go
// project/main.go
package main

import (
	"crypto/tls"
	"log"

	"github.com/yourusername/project/conf"
	"github.com/thoj/go-ircevent"
)

func main() {
	var err error

	// Load the configuration.
	if err = conf.Load("conf.yml"); err != nil {
		log.Fatal(err)
	}

	// Initialize the bot and setup the TLS parameters if needed
	ib := irc.IRC(conf.C.BotName, conf.C.BotName)
	if conf.C.TLS {
		ib.UseTLS = true
		if conf.C.InsecureTLS {
			ib.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}

	// Connect to the server
	if err = ib.Connect(conf.C.Server); err != nil {
		log.Fatal(err)
	}

	// On connection to the server, automatically join the configured channel
	ib.AddCallback("001", func(e *irc.Event) {
		ib.Join(conf.C.Channel)
	})

	// Callback to execute when a message is received either on the channel or
	// directly to the bot (query/msg for example)
	ib.AddCallback("PRIVMSG", func(e *irc.Event) {
		m := e.Message()
		log.Printf("I just received : '%v'", m)
	})
	ib.Loop()
}
```

The code is commented and should allow you to understand what's going on. First we load the `conf.yaml` file and exit the program if there is an error. Then we initialize the irc bot by calling successive functions. You can already compile your bot and start it. Once it receives a message it will just print it out in the console it was launched in. Nothing special.

Your project should now look like this :
```
project
├── conf
│   └── conf.go
├── conf.yml
└── main.go
```

# Markov chains !
Finally ! The base of our bot is ready so now we can switch to the more interesting part which is the ellaboration of the markov chain. Once again, I did not write all the code by myself and took a large part of the code from [here](https://golang.org/doc/codewalk/markov/) which is an example of how to implement markov chains in Go. So let's get started. Let's create a new package named `markov` and edit the file `markov.go` :

```go
// project/markov/markov.go
package markov

import (
	"math/rand"
	"strings"
	"time"
)

// PrefixLen is the number of words per Prefix defined as the key for the map.
const PrefixLen = 2

// MainChain is the chain that will be available outside the package.
var MainChain *Chain

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	Chain map[string][]string
}

// Build builds the chain using the given string parameter
func (c *Chain) Build(s string) {
	p := make(Prefix, PrefixLen)
	for _, v := range strings.Split(s, " ") {
		key := p.String()
		c.Chain[key] = append(c.Chain[key], v)
		p.Shift(v)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate() string {
	p := make(Prefix, PrefixLen)
	var words []string
	for {
		choices := c.Chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain() *Chain {
	return &Chain{make(map[string][]string)}
}

// Init initializes the markov chain
func Init() {
	rand.Seed(time.Now().UnixNano())
	MainChain = NewChain()
}
```

Here things are getting a bit complicated. As usual the code is already commented but let's explain a bit more the concept here. What will actually happen when we build the markov chain ?

```
> "Hello."
{" ": ["Hello."]}
> "Hello !"
{
    " ": ["Hello.", "Hello"],
    "Hello": ["!"]
}
> "Hello World !"
{
    " ": ["Hello.", "Hello", "Hello"],  
    "Hello": ["!", "World"],
    "Hello World": ["!"]
}
```

As you can see, each time a word or a sentence is typed in, it will modify the chain to include all the possible changes and paths. Then, when we will as the chain to generate a new sentence from the previous ones, it will walk through the chain, performing a random operation on each node. Of course things are starting to get interesting once the chain gets a bit more complicated than these three "Hello"s. As you can see, each node is naturally weighted by repetition. There is a higher chance that the first word will be "Hello" than "Hello.". Also there is a 50% chance that the said "Hello" will turn into a "Hello World !" or "Hello !".

# Integrating the markov chain

Let's edit the `project/main.go` file, and more specifically the `PRIVMSG` callback so that it doesn't just log the message it has received, but also update the markov chain.

```go
// project/main.go

    // Place this right after loading the configuration :
    markov.Init()

    // Callback to execute when a message is received either on the channel or
	// directly to the bot (query/msg for example)
	ib.AddCallback("PRIVMSG", func(e *irc.Event) {
		m := e.Message()
		if strings.HasPrefix(m, "!") {
			if strings.HasPrefix(m, "!mk") {
				ib.Privmsg(conf.C.Channel, markov.MainChain.Generate())
			}
		} else if strings.HasPrefix(m, conf.C.BotName) {
			ib.Privmsg(conf.C.Channel, markov.MainChain.Generate())
		} else {
			markov.MainChain.Build(m)
		}
	})
```

When the bot receives a message, it will check if there is a specific command or its own name as the first word. If so, then it will generate a random sentence using the markov chain. If neither of these conditions are met, then it will append the received message to the chain.
