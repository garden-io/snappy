This is the example code for Ellen Körbes' talk, [The Quest For The Snappiest Go & Kubernetes Workflow](https://garden.slides.com/ellenkorbes/snappy?token=hxfCQmd8). See [talk.md](talk.md) for a script of what's explained in the talk.

Usage:

For reasonableness sake, `vendor` and `cache` directories in this repo are empty. Before running benchmarks, 1) go through services `b` through `f` and run `go mod vendor` on each, to download dependencies (on service `e`, do that on `e-hot/src/`); and 2) on service `d`, run `GOCACHE=$PWD/cache` to populate the compiler cache. 

Then, start with `garden dev --hot e,f`.

Then use:
- `./change.sh a-basic`
- `./change.sh b-vendor`
- `./change.sh c-debugging`
- `./change.sh d-cache`
- `./change.sh e-hot/src`
- `./change.sh f-localhot`

...to make changes to the code, and keep an eye on the logs for each service to see how long it took between code change and the new process being up and running ([stern](https://github.com/wercker/stern) is recommended here).

If you're a Garden user just looking for the best solution, mimic the two-module setup found in `f-localhot` and `f-localhot-binary`.