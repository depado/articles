title: Closeable time ticker and multiple receivers in Go
description: Tutorial on how to have a multiple receivers listening on a single time.Ticker
slug: closable-time-ticker-and-multiple-receivers-in-go
author: Depado
date: 2015-07-15 14:47:00
tags:
    - go
    - dev

# Problem
Once, I wanted to have multiple goroutines to listen on the same time.Ticker channel and have different behaviours. This can't be achieved that easily because if you pass the same channel to two different goroutines, they will alternate the receive action. The first one will catch the first signal, the second one the second, then the first one will receive the third signal, etc...  
The second problem I found is that I wanted to stop the said ticker. If the ticker stops, there is no reason for all the goroutines started using this ticker to keep running.

Let's say that we have two functions that will be run as goroutines. These two functions will be the base for this article :

```go
func doThing(c <-chan time.Time) {
	for range c {
		log.Println("Hello from doThing")
	}
}

func doOtherThing(c <-chan time.Time) {
	for range c {
		log.Println("Hello from doOtherThing")
	}
}
```
Note the use of the `range` keyword. This is useful to detect if a channel has been closed (among other things).

# The multiTicker base function
As a basic solution, someone (sorry I can't remember his name) on the #go-nuts irc channel told me that it was pretty easy to create a function that takes a time.Ticker as an argument and returns two channels. Let's see how to do that :

```go
func multiTicker(c <-chan time.Time) (chan time.Time, chan time.Time) {
	a := make(chan time.Time)
	b := make(chan time.Time)
	go func() {
		for t := range c {
			a <- t
			b <- t
		}
	}()
	return a, b
}

func main() {
    tc := time.Tick(1 * time.Second)
    a, b := multiTicker(tc)
    go doThing(a)
    go doOtherThing(b)
    fmt.Scanln()
    log.Println("Exiting")
}
```
That solution works just fine. Once the `multiTicker` function is called, a goroutine is started and listens on the `tc` channel. Once it receives data from that channel it immediatly forwards it to the two channels created earlier `a` and `b`.

There is one problem though, what if my program grows bigger and at one point I want to stop that ticker ? First of all we shouldn't use the `time.Tick` function because it doesn't allow us to stop the said ticker. But even if we used a real ticker (using `tc := time.NewTicker` and then `tc.C.Stop()`, the two goroutines we started will keep running and wait for data, as well as the internal goroutine of `multiTicker`. This isn't efficient at all and it limits the control we have over all these elements.

# Another channel to rule them all
The solution is actually pretty simple. First of all let's declare a new structure and the associated `stop()` method :

```go
type closableTicker struct {
	ticker *time.Ticker
	halt   chan bool
}

func (ct *closableTicker) stop() {
	ct.ticker.Stop()
	close(ct.halt)
}
```
As you can see we the `stop` method will do two things. It will stop the ticker and close the `halt` channel. Doesn't make any sense for now but let's modify our `multiTicker` function as follow :

```go
func multiTicker(ct closableTicker) (chan time.Time, chan time.Time) {
	a := make(chan time.Time)
	b := make(chan time.Time)
	go func() {
		for {
			select {
			case t := <-ct.ticker.C:
				a <- t
				b <- t
			case <-ct.halt:
				close(a)
				close(b)
				return
			}
		}
	}()
	return a, b
}
```
Here we're using a subtle mecanism of Go about channels. When you close a channel it actually sends something to that channel. So our `halt` channel will never be used directly but once it's closed, we can close the two channels we declared (`a` and `b`) and terminate the goroutine. Using the `range` technique on `doThing` and `doOtherThing`, once the `a` and `b` channels will be closed, these two goroutines will return. Let's have a look at the main function :

```go
func main() {
	ct := closableTicker{
		ticker: time.NewTicker(1 * time.Second),
		halt:   make(chan bool, 1),
	}
	a, b := multiTicker(ct)
	go doThing(a)
	go doOtherThing(b)
	time.Sleep(3 * time.Second)
	ct.stop()
	fmt.Scanln()
	log.Println("Exit")
}
```
Now, the program will work as expected for 3 seconds and then all the goroutines will return, the channels will be closed and back to normal !
