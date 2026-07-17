-- Create an optimized table configuration for localized engine tracking
CREATE EXTERNAL TABLE IF NOT EXISTS aura_moby_telemetry.container_metrics (
    container_id      STRING,
    cpu_utilization   DOUBLE,
    memory_pinned_kb  BIGINT,
    vector_loop_stalls INT,
    tensor_load_ratio FLOAT
)
PARTITIONED BY (execution_date STRING)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY '\t'
STORED AS ORC; -- Optimized Row Columnar for high performance read patterns

-- High-throughput query to identify optimization drops
FROM aura_moby_telemetry.container_metrics
SELECT container_id, AVG(cpu_utilization) AS avg_cpu, MAX(tensor_load_ratio) AS max_tensor
WHERE execution_date = '2026-07-17'
GROUP BY container_id
HAVING avg_cpu > 85.0
ORDER BY max_tensor DESC
LIMIT 50;
