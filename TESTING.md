# Testing

```shell
task test
```

## Coverage

```shell
task test:coverage
```


## Debugging Mocker Calls
```shell
AWSMOCKER_DEBUG=1 go test -v ./cmd/... -run TestCleanCmd
```