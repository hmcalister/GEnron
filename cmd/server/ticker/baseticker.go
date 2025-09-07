package ticker

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type BaseTicker struct {
	name                string
	value               float64
	lastUpdateTimestamp time.Time
	randGen             *rand.Rand
	mu                  sync.RWMutex
}

// Initialize only the base ticker attributes using the given viper config.
// This initialization is shared between all tickers, so it is extracted here.
// Note an error is returned if initialization is malformed, be sure to check and return!
//
// Does not lock the mutex, since this method will be called from the parent Initalize method, which already locks.
func (t *BaseTicker) initializeBase(tickerConfig *viper.Viper) error {
	if !tickerConfig.IsSet("name") {
		return errors.New("error initializing ticker, name field not specified")
	}
	t.name = tickerConfig.GetString("name")

	t.value = tickerConfig.GetFloat64("value")
	if t.value < 0.0 {
		return errors.New("error initializing ticker, specified initial value is negative")
	}

	var randomSeed int64
	if tickerConfig.IsSet("randomseed") {
		randomSeed = tickerConfig.GetInt64("randomseed")
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
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.value
}

func (t *BaseTicker) GetInfo() (string, float64, time.Time) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.name, t.value, t.lastUpdateTimestamp
}

func (t *BaseTicker) SetLastUpdatedTimestamp(timestamp time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastUpdateTimestamp = timestamp
}
