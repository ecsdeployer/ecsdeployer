# Testing

```shell
make test
```

## Coverage

```
make coverage
```


## Debugging Mocker Calls
```
AWSMOCKER_DEBUG=1 go test -v ./cmd/... -run TestCleanCmd
```