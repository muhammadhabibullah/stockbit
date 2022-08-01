# stockbit

## Setup

1. Setup config file on the [config](config) directory.

```shell
$ make config
```

2. Run docker compose to create/run the broker

```shell
$ make dependency
```

3. Run processor. On the first time, it must be run first so HTTP can serve view.

```shell
$ go run main.go -command processor -processor balance
```

On separate terminal run another processor

```shell
$ go run main.go -command processor -processor above_threshold
```

4. Run HTTP service

```shell
$ go run main.go -command http
```

## API

Access OpenAPI after running the HTTP service at [localhost:8000/swagger/](http://localhost:8000/swagger/)
