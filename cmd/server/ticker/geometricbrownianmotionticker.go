package ticker

import (
	"errors"
	"math"

	"github.com/spf13/viper"
)

type GeometricBrownianMotionTicker struct {
	BaseTicker
	drift      float64
	volatility float64
}

func (t *GeometricBrownianMotionTicker) Initialize(tickerConfig *viper.Viper) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.initializeBase(tickerConfig); err != nil {
		return err
	}

	// Drift is allowed to be negative.
	t.drift = tickerConfig.GetFloat64("drift")

	t.volatility = tickerConfig.GetFloat64("volatility")
	if t.volatility < 0.0 {
		return errors.New("error initializing geometric brownian motion ticker, volatility term is negative")
	}

	return nil
}

func (t *GeometricBrownianMotionTicker) Update() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Geometric Brownian Motion update looks like
	// X_{t+dt} = X_t * exp((drift - 0.5 * volatility**2)dt + volatility*sqrt(dt)*Z)
	// For a random gaussian number Z (simulating the random walk)
	//
	// We will discretize this to assume dt=1, and have the user set drift and volatility accordingly

	dt := 1.0
	exponent := (t.drift-0.5*math.Pow(t.volatility, 2))*dt + t.volatility*math.Sqrt(dt)*t.randGen.NormFloat64()
	t.value *= math.Exp(exponent)
	if t.value < 0 {
		t.value = 0
	}
}
