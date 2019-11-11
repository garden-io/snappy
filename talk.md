# The Snappiest Go & Kubernetes Development Workflow

## Introduction

There are tools like Draft, Garden, Skaffold, and Tilt.

When working with Go, their defaults are not ideal, and the time between iterations is longer than it could be.

This is focused strictly on creating a fast development workflow—production concerns are deliberately ignored.

Let's fix that.

## The Code

In every iteration of this project, we'll have the same Go code.

It'll have unused imports just so we have some dependencies:

```go
	_ "github.com/dgraph-io/badger"
	_ "github.com/influxdata/influxdb"
```

And a Unix timestamp. 

```go
	called := time.Unix(0, 1572547765952071764)
```

This timestamp will be updated every time we make a change to our file with the following bash script:

```sh
sed -re 's/time.Unix\(0, ([0-9]*)\)/time.Unix\(0, '"$(($(date +%s%N)))"'\)/' -i $1/main.go
```

When the code runs, it will calculate the time difference between the moment I made a change to my code, and the moment the binary ran in the container.

In another file, `a_main-packr.go`, a kubectl binary has been embedded into our source code. This is to simulate having a very large binary.


# The Set-Up

We'll use Garden, against a minikube cluster. Garden will re-build and re-deploy my container on every code change, and we'll do our best to make the time it takes to do that as small as possible.

Every step will be benchmarked three times, and averaged.

# Benchmarks

## Basic (A)

Here we're using a `golang:alpine` image, and copying our source code into it, and building.

This happens on every code change.

Note that there is no dependency cache (vendoring) or compiler cache.

It will create a new image, download dependencies, and compile from scratch. Every single time.

On my machine, times were:

- 16.902s
- 16.685s
- 16.845s

Average: 16.810s.

## Vendor (B)

On the above, downloading dependencies is probably the most time consuming part.

For part B, let's `go mod vendor` on our local codebase, then include that in our container and see how much it speeds things up.

We'll need to add `-mod=vendor` to our build command. E.g.

```
go build -o binary -mod=vendor
```

Every time: It will create a new image, and compile from scratch.

Times:

- 10.442s
- 10.327s
- 10.326s

Average: 10.365s.

Better.

## Debugging (C)

When you search for how to optimize Go binaries, the first suggestion that tends to come up is removing debugging information (DWARF) from your executables.

To test this, add `-ldflags '-w'` to your build command, like so:

```go
go build -o binary -mod=vendor -ldflags '-w'
```

Let's see if that makes any difference. Times:

- 10.869s
- 10.339s
- 10.442s

Average: 10.550s.

No difference. (The binary does turn up a bit smaller though.)

## Cache (D)

Something else that can significantly speed things up is making use of Go's compiler cache.

Testing it locally, there was a *10x* improvement in compilation speed.

So let's copy that cache into our container and see if it helps.

To do that, we can simply use the `GOCACHE` variable, e.g. prepending it to our like command:

```go
GOCACHE=/app/cache go build -o binary -mod=vendor -ldflags '-w'
```

Times:

- 14.911s
- 16.738s
- 15.344s

Average: 15.664s.

It got a lot worse!

Why? The compiler cache consists of a ton of tiny files, and the i/o kills it.

## Hot Reload (E)

So what if instead of having to make a new image with a billion files every time we change a file, we just keep a same container alive indefinitely, and simply sync our code changes into it?

This way, we won't have a billion cache files being copied around every time.

To implement this, we'll use Garden's hot reload feature. To activate it, we need this in our `garden.yml`:

```yaml
hotReload:
  sync:
    - target: /app/
	  source: src/
```

This will have two separate compile times. One for the first time, before there's a compiler cache inside the container, and then subsequent shorter times when there's already a cache there.

Times:

- 19.225s (first)
- 8.612s
- 7.511s
- 6.638s

Average (excluding first run): 7.587s.

Not bad.

But this is highly dependent on the amount of resources the container has available for compilation.

If instead of 1 cpu, the container has only 0.5 cpu available, here's what happens:

- 1m13s (first)
- 17.119s
- 17.593s
- 18.439s

Average (excluding first run): 17.717s.

So smaller container equals lots of waiting.

## Local Compilation, Hot-Reloaded Binary (F)

Let's try something else instead: We compile the binary locally, using our laptop's CPU. We then hot reload the binary into a container.

This way we have fast compilation, compiler cache, vendoring, we can use a very tiny base image for the container, and this container doesn't have to be restarted on updates.

The only potential bottleneck is the size of the binary.

[Implementation.]

- 6.921s
- 7.008s
- 7.601s

Average: 7.176s.

Not bad. 

## UPX

But the binary we're using is *huge.* Let's see if compressing will make it better, or whether the time it takes to compress offsets any gain in smaller file size:

- 8.702s
- 8.587s
- 8.458s

Average: 8.582s.

So nope. It gets worse.

# Conclusion

The best method is to compile locally and hot reload the results. 

The next best one, if resources allow, is to hot reload source files and compile in-container.

In local testing, the difference was small (if we discount the latter method's slower, first run): 7.176s vs 7.587s.