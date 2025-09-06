# Synthetic Stock Data in Go
### Go-Enron (This project is not affiliated with Enron)

An implementation of real-time stock data served over [Connect-RPC](https://connectrpc.com/). 

### Config Specification

| Key | Datatype | Default | Meaning |
| --- | -------- | ------- |
| LogLevel | String Enum ("none", "error", "warn", "info", "debug") | "Info" | The level at which logs are recorded. None disables logging. |
| LogFile | String | "" | The filepath to write logs to. If left unset or empty, logs are sent to `stdout`. The file is truncated before logging begins. If the file cannot be opened for writing, the program panics. |

See `config/LoadConfig` for more information.

### Plans

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