# Testing

```shell
make test
```

## Coverage

```shell
make coverage
```


## Debugging Mocker Calls
```shell
AWSMOCKER_DEBUG=1 go test -v ./cmd/... -run TestCleanCmd
```