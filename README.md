# Synthetic Stock Data in Go
### Go-Enron (This project is not affiliated with Enron)

An implementation of real-time stock data served over [Connect-RPC](https://connectrpc.com/). 

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