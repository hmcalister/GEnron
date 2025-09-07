package ticker

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type BaseTicker struct {
	name    string
	value   float64
	randGen *rand.Rand
	mu      sync.Mutex
}

// Initialize only the base ticker attributes using the given viper config.
// This initialization is shared between all tickers, so it is extracted here.
// Note an error is returned if initialization is malformed, be sure to check and return!
func (t *BaseTicker) initializeBase(tickerConfig *viper.Viper) error {
	if !tickerConfig.IsSet("Name") {
		return errors.New("error initializing ticker, name field not specified")
	}
	t.name = tickerConfig.GetString("Name")

	t.value = tickerConfig.GetFloat64("Value")
	if t.value < 0.0 {
		return errors.New("error initializing ticker, specified initial value is negative")
	}

	var randomSeed int64
	if tickerConfig.IsSet("RandomSeed") {
		randomSeed = tickerConfig.GetInt64("RandomSeed")
	} else {
		randomSeed = time.Now().UnixNano()
	}
	t.randGen = rand.New(rand.NewSource(randomSeed))

	return nil
}

func (t *BaseTicker) String() string {
	return t.name
}

func (t *BaseTicker) GetValue() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.value
}
