# Synthetic Stock Data in Go
### Go-Enron (This project is not affiliated with Enron)

An implementation of real-time stock data served over [Connect-RPC](https://connectrpc.com/). 

## Getting Started

To start the server (which creates and updates tickers), first create a `config.yaml` file and set it according to the [section below](#config-specification), or use the example config file given. Alternative config files may be loaded using the command line flag `-configFilePath=...`. 

Compile the server using `go build cmd/server/main.go -o server`, then run it using `./server`. Depending on the configuration, the program will either start printing logs to `stdout`, or create a log file, but as long as the process does not immediately terminate the server should be running.

Make requests to the server on `localhost:8080` (the port of which may be changed using the config file, see the `port` option). A Go client is given, but other clients may be created using the prototbuf definitions given under the `api` directory.

To generate new ConnectRPC bindings for your client of choice, alter the `buf.gen.yaml` file to target your client of choice. See [buf](https://buf.build/) and [ConnectRPC](https://connectrpc.com/) for details on what clients are available.

## Config Specification

See `config/LoadConfig` for more information. See `config.yaml` for an example configuration.

| Key | Datatype | Default | Meaning |
| --- | -------- | ------- | ------- |
| loglevel | String Enum ("none", "error", "warn", "info", "debug") | "Info" | The level at which logs are recorded. None disables logging. |
| logfile | String | "" | The filepath to write logs to. If left unset or empty, logs are sent to `stdout`. The file is truncated before logging begins. If the file cannot be opened for writing, the program panics. |
| updateperiod | int64 | 1_000_000 | The amount of time (in nanoseconds) to between updates of the tickers. Must be greater than 0. If the update period is too small, the program may not be able to achieve the required period. |
| tickers | Dictionary[String, Ticker] | Empty | The tickers to create and manage. Ticker names are used to request data from the server, and tickers have unique specifications based on the ticker type. See below for a list of ticker types and their specifications. The key string is the ticker `name`, which must be unique for each ticker. All tickers have the fields `type`, `value`, and `randomseed`. The `type` field that identifies the ticker type. The `value` field specifies the initial value, and must be non-negative. The valid ticker types are listed below. In general, all ticker fields are required. The except is `randomseed` which may be left unset to specify a random seed based on the current timestamp. |

Ticker Types:
- "UniformRandom"

### Uniform Random Ticker

`type: "UniformRandom"`

Update the ticker value with a uniformly chosen random number at every step. This is a very simple, but very unrealistic model.

| Key | Datatype | Meaning |
| --- | -------- | ------- | ------- |
| type | String | The ticker type. Must be explicitly the above type to be processed at this ticker variety. |
| value | float64 | The initial value for the ticker. Must be non-negative. |
| randomseed | int64 | None | The random seed to use for the generator. If left unset, the current unix timestamp is used instead. |
| randomrange | float64 | The upper and lower bound on the random number. Must be non-negative. |


## Plans

- Several tickers, each modelled with a different synthetic approach.
    - Naive Monte Carlo
    - Geometric Brownian Motion
    - Jump-Diffusion
    - Stochastic Volatility
    - ARIMA, GARCH Models
- Sampling of time series data to a persistent database.
    - SQLite
    - Postgresql with extensions
- Benchmarking of update methods for good real-time applications.

## Technologies

- [Slog](https://pkg.go.dev/log/slog): For structured logging.
- [Viper](https://github.com/spf13/viper): For configuration setting.