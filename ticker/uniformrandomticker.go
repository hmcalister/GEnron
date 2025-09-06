package ticker

import (
	"math/rand"
	"sync"
)

type UniformRandomTicker struct {
	name        string
	Value       float64 `mapstructure:"Value"`
	RandomRange float64 `mapstructure:"RandomRange"`
	randGen     *rand.Rand
	mu          sync.Mutex
}

func (t *UniformRandomTicker) String() string {
	return t.name
}

func (t *UniformRandomTicker) Validate() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Value < 0.0 {
		return ErrorNegativeValueTicker
	}

	if t.RandomRange < 0.0 {
		return ErrorGenericInvalidTicker
	}

	return nil
}

func (t *UniformRandomTicker) GetValue() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.Value
}

func (t *UniformRandomTicker) Update() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Value += -t.RandomRange + 2*t.RandomRange*t.randGen.Float64()
	if t.Value < 0 {
		t.Value = 0
	}
}
