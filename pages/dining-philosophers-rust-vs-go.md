title: "Dining Philosophers : Rust vs Go"
description: A silly comparison between Rust and Go where nothing interesting happen.
slug: dining-philosophers-rust-vs-go
date: 2015-08-19 18:54:00
tags:
    - go
    - rust

This article will be about Rust and Go. Today I decided to have a look at Rust, and while I was reading the [official tutorial](https://doc.rust-lang.org/book/), I saw [this example](https://doc.rust-lang.org/book/dining-philosophers.html) about the dining philosophers problem. This article isn't about performances because it's already well known that [Rust is faster than Go](http://benchmarksgame.alioth.debian.org/u64q/compare.php?lang=rust&lang2=go), but it will be about code readability, length and complexity.

**I'm a total Rust beginner and this is what I can tell from the language after a day of reading docs, practicing, etc... I know there is a lot of things I don't know and maybe I don't even have the right approach to Rust programming.**

# The problem itself
In ancient times, a wealthy philanthropist endowed a College to accommodate five eminent philosophers. Each philosopher had a room in which they could engage in their professional activity of thinking; there was also a common dining room, furnished with a circular table, surrounded by five chairs, each labelled by the name of the philosopher who was to sit in it. They sat anticlockwise around the table. To the left of each philosopher there was laid a golden fork, and in the centre stood a large bowl of spaghetti, which was constantly replenished. A philosopher was expected to spend most of their time thinking; but when they felt hungry, they went to the dining room, sat down in their own chair, picked up their own fork on their left, and plunged it into the spaghetti. But such is the tangled nature of spaghetti that a second fork is required to carry it to the mouth. The philosopher therefore had also to pick up the fork on their right. When they were finished they would put down both their forks, get up from their chair, and continue thinking. Of course, a fork can be used by only one philosopher at a time. If the other philosopher wants it, they just have to wait until the fork is available again.

# Code

## Rust Code

*Edit : Thanks to Luigi M. for pointing the errors in the Rust code. You can find his blog [here](https://grigio.org/).*

```rust
use std::thread;
use std::sync::{Mutex, Arc};

struct Table {
    forks: Vec<Mutex<()>>,
}

struct Philosopher {
    name: String,
    left: usize,
    right: usize,
}

impl Philosopher {
    fn new(name: &str, left: usize, right: usize) -> Philosopher {
        Philosopher {
            name: name.to_string(),
            left: left,
            right: right,
        }
    }

    fn eat(&self, table: &Table) {
        let _left = table.forks[self.left].lock().unwrap();
        let _right = table.forks[self.right].lock().unwrap();

        println!("{} is eating.", self.name);

        thread::sleep_ms(1000);

        println!("{} is done eating.", self.name);
    }
}

fn main() {
    let table = Arc::new(Table { forks: vec![
        Mutex::new(()),
        Mutex::new(()),
        Mutex::new(()),
        Mutex::new(()),
        Mutex::new(()),
    ]});

    let philosophers = vec![
            Philosopher::new("Judith Butler", 0, 1),
            Philosopher::new("Gilles Deleuze", 1, 2),
            Philosopher::new("Karl Marx", 2, 3),
            Philosopher::new("Emma Goldman", 3, 4),
            Philosopher::new("Michel Foucault", 0, 4),
        ];

        let handles: Vec<_> = philosophers.into_iter().map(|p| {
            let table = table.clone();

            thread::spawn(move || {
                p.eat(&table);
            })
        }).collect();

        for h in handles {
            h.join().unwrap();
        }
}
```

## Go Code

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

type table struct {
	forks []sync.Mutex
}

type philosopher struct {
	name  string
	left  int
	right int
}

func (p philosopher) eat(t *table) {
	defer wg.Done()
	t.forks[p.left].Lock()
	defer t.forks[p.left].Unlock()
	t.forks[p.right].Lock()
	defer t.forks[p.right].Unlock()

	fmt.Println(p.name, "is eating.")

	time.Sleep(1 * time.Second)

	fmt.Println(p.name, "finished eating.")
}

func main() {
	philosophers := [...]philosopher{
		philosopher{"One", 0, 1},
		philosopher{"Two", 1, 2},
		philosopher{"Three", 2, 3},
		philosopher{"Four", 3, 4},
		philosopher{"Five", 0, 4},
	}
	t := table{forks: make([]sync.Mutex, len(philosophers))}

	for _, p := range philosophers {
		wg.Add(1)
		go p.eat(&t)
	}
	wg.Wait()
}
```

# Comparison

## Length

Alright. A source code number of line doesn't represent anything except maybe the keystrokes the developper had to type to write it down. It has nothing to do with how easy it is to code in a said language, how optimized it is or whatever. It also depends on how you code, whether you want to write clean code or compact code...

| Rust | Go |
|-------|------|
| 63 Lines of code | 50 Lines of code |
As you can see, the Rust code is slightly longer. The gap isn't that big as you can see. I was told that Rust was really explicit and could produce a really large gap between programming languages. Not on that example anyway. Note that I didn't try to reduce the Go code's size, or increase the Rust one, I'm trying to stay objective even though that's not a very important metric.

## Readability and Complexity

There's quite a lot to say on those points. First of all, I didn't know quite what I was doing when following the tutorial on Rust because everything looks like some kind of magic trick. For example, in the `eat` method, you can see that the locks are acquired and assigned to a variable that starts with a `_`. What does this mean ? That means that those variables are going to be destroyed when they are no longer in the scope and so the locks will be released. It's actually a trick to tell the Rust compiler not to complain about those variables not being used. The whole `handles` part is getting a bit too much complicated too. A lambda programmer, at first glance, couldn't tell what's going on. Of course you can get the principle quite easily but that's not a way of programming I would memorize easily I think.

For the concurrency part I'll just say that Go is really more simple. I would say that Rust sticks to the old way of doing concurrency, even if that's not entirely true. You still need to create manually threads, wait for them, tell them to start, all of that explicitely. With the Go version there is actually a drawback because it's a little more tricky to wait for all the goroutines to finish, using a WaitGroup (I guess, from a Go beginner point of view, that would seem like a magic trick too).

I'd say the Go version, despite being shorter, actually shows what's going on at first glance.

# Conclusion
Rust looks like a pretty nice language, with a LOT of features. Despite the fact that it's a modern language I really wouldn't recommend it as a first programming language. You actually need to have strong knowledge of what's going on exactly, you need to know how references work, how generic works and such things. It looks promising though I'll keep digging that Rust language and see if I can accomodate with its principles and understand the right way to code.
