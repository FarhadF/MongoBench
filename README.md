# MongoBench
Benchmark your MongoDB Deployment.

Build steps:
1. Clone the repository in your `$GOPATH/src/`
2. `go build mongobench.go`
3. `./mongobench`

```
Usage:
  mongobench [flags]

Flags:
  -b, --batch int          Number of threads per batch. (default 100)
  -h, --help               help for mongobench
  -q, --queryFile string   Path to the query file, one query per line. Only the query string, example: {"branchCode": 230}" (default "/tmp/query")
  -t, --threads int        Total number of threads to use. Equal to number of queries against mongodb. (default 100)
  -v, --version            Prints version
```


