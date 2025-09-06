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
	Value() float64

	// Update the value of a ticker.
	//
	// All updates should finish by ensuring the ticker value is non-negative,
	// that is, clamp the updated value to be zero or larger.
	Update()
}
