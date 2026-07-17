-- ==============================================================================
-- Aura Moby Embedded Telemetry Schema Map
-- ==============================================================================

-- 1. Tracks localized node state metadata profiles
CREATE TABLE IF NOT EXISTS node_runtime_state (
    node_uuid TEXT PRIMARY KEY,
    vector_registers_active INTEGER DEFAULT 0,
    tensor_cores_saturated REAL DEFAULT 0.0,
    last_heartbeat_timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Maintains hot metrics parsed by the 18% Assembly layer
CREATE TABLE IF NOT EXISTS autonomous_metrics_log (
    log_id INTEGER PRIMARY KEY AUTOINCREMENT,
    container_hash TEXT NOT NULL,
    vector_loop_stalls INTEGER UNSIGNED,
    cpu_clock_frequency_hz REAL,
    memory_pinned_bytes BIGINT,
    recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Provides declarative threshold paths for the AI optimizer sub-module
CREATE TABLE IF NOT EXISTS optimizer_policies (
    policy_id TEXT PRIMARY KEY,
    target_namespace TEXT UNIQUE,
    max_hardware_allocation_ratio REAL CHECK(max_hardware_allocation_ratio <= 1.0),
    enforce_sandbox_isolation BOOLEAN DEFAULT 1
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed mock data for development testing
INSERT INTO users (name, email) VALUES 
('Alice Vance', 'alice@example.com'),
('Bob Smith', 'bob@example.com')
ON CONFLICT (email) DO NOTHING;
