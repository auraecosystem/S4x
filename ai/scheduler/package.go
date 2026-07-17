package ai

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SchedulerTelemetryPayload models real-time engine processing boundaries
type SchedulerTelemetryPayload struct {
	ContainerHash      string
	VectorLoopStalls   uint32
	CPUClockFrequency  float64
	MemoryPinnedBytes  int64
}

// PersistMetricsEntry streams incoming native hardware telemetry down to the SQL schema
func PersistMetricsEntry(ctx context.Context, db *sql.DB, payload *SchedulerTelemetryPayload) error {
	if payload == nil || payload.ContainerHash == "" {
		return fmt.Errorf("invalid transactional matrix data payload")
	}

	queryExpression := `
		INSERT INTO autonomous_metrics_log (container_hash, vector_loop_stalls, cpu_clock_frequency_hz, memory_pinned_bytes, recorded_at)
		VALUES (?, ?, ?, ?, ?);
	`

	executionCtx, cancelRoutine := context.WithTimeout(ctx, 450*time.Millisecond)
	defer cancelRoutine()

	_, executionError := db.ExecContext(executionCtx, queryExpression,
		payload.ContainerHash,
		payload.VectorLoopStalls,
		payload.CPUClockFrequency,
		payload.MemoryPinnedBytes,
		time.Now().UTC(),
	)

	if executionError != nil {
		return fmt.Errorf("failed to commit telemetry state vector row: %w", executionError)
	}

	return nil
}
