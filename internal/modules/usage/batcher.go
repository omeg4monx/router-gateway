package usage

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
)

const flushTimeout = 10 * time.Second

type batcher struct {
	repo          *repository
	pricing       *pricingStore
	bus           *eventbus.Bus
	logger        *slog.Logger
	now           func() time.Time
	batchSize     int
	flushInterval time.Duration
	bufferSize    int

	mu      sync.RWMutex
	started bool
	stopped bool
	ch      chan contracts.UsageRecord
	done    chan struct{}

	dropped atomic.Int64
}

func newBatcher(
	repo *repository,
	pricing *pricingStore,
	bus *eventbus.Bus,
	logger *slog.Logger,
	now func() time.Time,
	batchSize int,
	flushInterval time.Duration,
	bufferSize int,
) *batcher {
	return &batcher{
		repo:          repo,
		pricing:       pricing,
		bus:           bus,
		logger:        logger,
		now:           now,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		bufferSize:    bufferSize,
	}
}

func (b *batcher) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.started && !b.stopped {
		return
	}
	b.ch = make(chan contracts.UsageRecord, b.bufferSize)
	b.done = make(chan struct{})
	b.started = true
	b.stopped = false

	go b.run(b.ch, b.done)
}

func (b *batcher) Enqueue(rec contracts.UsageRecord) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.started || b.stopped {
		return
	}

	select {
	case b.ch <- rec:
	default:
		count := b.dropped.Add(1)
		if b.logger != nil && (count == 1 || count%100 == 0) {
			b.logger.Warn("usage buffer full; dropping usage telemetry record", "dropped_total", count, "model", rec.Model)
		}
	}
}

func (b *batcher) Stop(ctx context.Context) error {
	b.mu.Lock()
	if !b.started || b.stopped {
		b.mu.Unlock()
		return nil
	}
	b.stopped = true
	ch := b.ch
	done := b.done
	close(ch)
	b.mu.Unlock()

	if ctx == nil {
		<-done
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *batcher) DroppedCount() int64 {
	return b.dropped.Load()
}

func (b *batcher) run(ch <-chan contracts.UsageRecord, done chan struct{}) {
	defer close(done)

	ticker := time.NewTicker(b.flushInterval)
	defer ticker.Stop()

	batch := make([]contracts.UsageRecord, 0, b.batchSize)
	for {
		select {
		case rec, ok := <-ch:
			if !ok {
				if len(batch) > 0 {
					b.flush(batch)
				}
				return
			}
			batch = append(batch, rec)
			if len(batch) >= b.batchSize {
				b.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				b.flush(batch)
				batch = batch[:0]
			}
		}
	}
}

func (b *batcher) flush(records []contracts.UsageRecord) {
	if len(records) == 0 {
		return
	}

	entries := make([]UsageLog, len(records))
	for i, rec := range records {
		ts := rec.Timestamp
		if ts.IsZero() {
			ts = b.now()
		}
		inTok := max(rec.InputTokens, 0)
		outTok := max(rec.OutputTokens, 0)
		cachedTok := max(rec.CachedTokens, 0)

		cost := rec.CostUSD
		if cost < 0 {
			cost = 0
		} else if cost == 0 {
			cost = b.pricing.CalculateCost(rec.Model, inTok, outTok, cachedTok)
		}

		latencyMS := max(rec.Latency.Milliseconds(), 0)
		ttftMS := max(rec.TTFT.Milliseconds(), 0)

		entries[i] = UsageLog{
			RequestID:    rec.RequestID,
			KeyID:        rec.KeyID,
			Model:        rec.Model,
			ProviderID:   rec.ProviderID,
			ConnectionID: rec.ConnectionID,
			Protocol:     rec.Protocol,
			InputTokens:  inTok,
			OutputTokens: outTok,
			CachedTokens: cachedTok,
			CostUSD:      cost,
			LatencyMS:    latencyMS,
			TTFTMS:       ttftMS,
			StatusCode:   rec.StatusCode,
			ErrorSummary: rec.ErrorSummary,
			Stream:       rec.Stream,
			CreatedAt:    ts.UTC().Format(rfc3339NanoFixed),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), flushTimeout)
	defer cancel()

	inserted, err := b.repo.InsertBatch(ctx, entries)
	if err != nil {
		if b.logger != nil {
			b.logger.Error("usage batch insert failed", "count", len(entries), "err", err)
		}
		return
	}

	if b.bus != nil {
		for _, item := range inserted {
			b.bus.Publish(TopicUsageRecorded, item)
		}
	}
}
