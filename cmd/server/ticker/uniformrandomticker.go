package ticker

import (
	"errors"

	"github.com/spf13/viper"
)

type UniformRandomTicker struct {
	BaseTicker
	randomRange float64
}

func (t *UniformRandomTicker) Initialize(tickerConfig *viper.Viper) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.initializeBase(tickerConfig); err != nil {
		return err
	}

	t.randomRange = tickerConfig.GetFloat64("randomrange")
	if t.randomRange < 0.0 {
		return errors.New("error initializing uniform random ticker, random range is negative")
	}

	return nil
}

func (t *UniformRandomTicker) Update() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.value += -t.randomRange + 2*t.randomRange*t.randGen.Float64()
	if t.value < 0 {
		t.value = 0
	}
}
