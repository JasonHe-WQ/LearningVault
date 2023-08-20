from prometheus_client import CollectorRegistry, Gauge, push_to_gateway

registry = CollectorRegistry()
g = Gauge(
    'job_last_success_unixtime',
    'Last time a batch job successfully finished',
    ['job'],
    registry=registry)
g.labels(job='my_batch_job').set_to_current_time()
push_to_gateway('localhost:9091', job='my_batch_job', registry=registry)
