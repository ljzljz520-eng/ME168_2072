# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	gymrecommend/cmd/server	[no test files]
?   	gymrecommend/internal/analytics	[no test files]
?   	gymrecommend/internal/catalog	[no test files]
ok  	gymrecommend/internal/config	0.002s
ok  	gymrecommend/internal/httpapi	0.027s
ok  	gymrecommend/internal/importcsv	0.005s
ok  	gymrecommend/internal/model	0.002s
?   	gymrecommend/internal/planning	[no test files]
?   	gymrecommend/internal/profile	[no test files]
ok  	gymrecommend/internal/recommend	0.003s
ok  	gymrecommend/internal/storage	0.036s
--- FAIL: TestWorkflow12 (0.00s)
    workflow12_test.go:21: summary left workflow in summarized
FAIL
FAIL	gymrecommend/internal/workflow	0.049s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
- Frontend build (web): exit `0`
