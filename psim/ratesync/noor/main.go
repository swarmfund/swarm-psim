package noor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/ratesync/providers"
)

var isProviderRunning = false

type Provider struct {
	host   string
	port   int
	pairs  []Pair
	tick   providers.Tick
	ticks  chan providers.Tick
	errors chan error

	logger *logan.Entry
}

func NewProvider(logger *logan.Entry, host string, port int, pairs []Pair) *Provider {
	provider := &Provider{
		host:   host,
		port:   port,
		pairs:  pairs,
		ticks:  make(chan providers.Tick),
		errors: make(chan error),
		logger: logger,
	}
	return provider
}

func (c *Provider) Errors() chan error {
	return c.errors
}

func (c *Provider) Ticks() chan providers.Tick {
	if isProviderRunning {
		c.logger.Panic("Noor index provider already running")
	}

	isProviderRunning = true
	go c.run()
	return c.ticks
}

func (c *Provider) run() {
	var tick Tick
	for {
		err := c.runOnce(&tick)
		if err != nil {
			c.errors <- err
		}
	}
}

func (c *Provider) runOnce(tick *Tick) error {
	defer func() {
		if rec := recover(); rec != nil {
			c.logger.WithField("recover", rec).Error("RateSync runOnce panic, but recovered")
		}
	}()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return errors.Wrap(err, "failed to resolve addr")
	}

	c.logger.Info("Trying to connect")
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		c.logger.Warn("Failed to dial")
		return errors.Wrap(err, "failed to dial")
	}

	defer conn.Close()
	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)
	for {
		conn.SetDeadline(time.Now().UTC().Add(time.Duration(30) * time.Second))

		var symbol Symbol
		c.logger.Debug("Waiting for tick")
		err := decoder.Decode(&symbol)
		if err != nil {
			if err == io.EOF {
				c.logger.Debug("Received EOF")
				return nil
			}
			c.logger.WithField("current_time", time.Now().UTC().String()).WithError(err).Error("Failed to decode the tick")
			// do not return the error to speed up reconnection
			return nil
		}

		if !c.tryUpdateTick(&symbol, tick) {
			c.logger.WithField("symbol", symbol).Debug("Skipping")
			continue
		}

		c.logger.WithField("tick", tick).Info("Ticked")
		c.ticks <- tick
	}
}

func (c *Provider) tryUpdateTick(symbol *Symbol, tick *Tick) bool {
	currentTime := time.Now().UTC().Add(time.Duration(-5) * time.Second)
	if !symbol.UpdateTime.After(currentTime) {
		return false
	}

	pair := c.findMatchingPair(symbol.Symbol)
	if pair == nil {
		return false
	}

	c.logger.WithField("symbol", symbol).Debug("Received valid symbol")

	price := pair.PhysicalPrice((symbol.Bid + symbol.Ask) / 2)

	op := horizon.SetRateOp{
		BaseAsset:     pair.Code,
		QuoteAsset:    pair.Quote,
		PhysicalPrice: price,
	}

	for i, storedOp := range tick.ops {
		if storedOp.BaseAsset != pair.Code || storedOp.QuoteAsset != pair.Quote {
			continue
		}

		tick.ops[i] = op
		return true
	}

	tick.ops = append(tick.ops, op)
	return true
}

func (c *Provider) findMatchingPair(symbol string) *Pair {
	for i := range c.pairs {
		if c.pairs[i].Symbol != symbol {
			continue
		}

		return &c.pairs[i]
	}

	return nil
}
