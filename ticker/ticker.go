package ticker

// The default ticker interface.
// All tickers must implement the below methods, but may use vastly different implementations.
// This allows for, say, different simulation methods.
type Ticker interface {
	// Return the name of the ticker.
	String() string

	// Validate the ticker. Ensures the ticker is a meaningful representation in its current state.
	// Returns nil if the ticker is valid, and an ErrorInvalidTicker otherwise.
	Validate() error

	// Get the current value of a ticker.
	//
	// Note that no ticker value may be below zero, as a rule of business logic.
	GetValue() float64

	// Update the value of a ticker.
	//
	// All updates should finish by ensuring the ticker value is non-negative,
	// that is, clamp the updated value to be zero or larger.
	Update()
}

// Factory pattern to initialize the new ticker.
// Returns:
// - (Ticker, nil) if all checks pass and the ticker is initialized
// - (nil, ErrorUnknownTickerType) if the ticker type in tickerConfig is unknown
// - (nil, ErrorInvalidTicker) if the ticker is initialized with invalid parameters (e.g. a negative initial value)
func NewTicker(name string, tickerConfig *viper.Viper) (Ticker, error) {
	tickerType := tickerConfig.GetString("Type")

	switch tickerType {
	default:
		return nil, ErrorUnknownTickerType
	}
}

