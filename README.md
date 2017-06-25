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
  -b, --batch int           Number of threads per batch. (default 100)
  -c, --collection string   Collection to run queries against (default "journal")
  -d, --database string     Database to run queries against (default "journaldb")
  -h, --help                help for mongobench
  -H, --host string         IP addresses or Hostnames and ports of the mongo hosts to connect to separated by commas, example: mongo1:27017, mongo2:27017 (default "localhost:27017")
  -q, --queryFile string    Path to the query file, one query per line. Only the query string, example: {"branchCode":230}" (default "/tmp/query")
  -t, --threads int         Total number of threads to use. Equal to number of queries against mongodb (default 100)
  -T, --timeout int         db query timeout in seconds (default 15)
  -v, --version             Prints version
```
