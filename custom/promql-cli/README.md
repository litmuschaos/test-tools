## Promql Cli

- It contains binaries for the [promql cli](https://github.com/nalbury/promql-cli), which can be used to run prometheus queries on the provided host to extract out the prometheus metrices.
- These binaries can be pulled into [litmus-go](https://github.com/litmuschaos/litmus-go) repository while building the go-runner image.
- Adding these binaries here as it's official repository doesn't contain amd64 binaries yet. It will be removed after addition of amd64 binaries in the official repository.