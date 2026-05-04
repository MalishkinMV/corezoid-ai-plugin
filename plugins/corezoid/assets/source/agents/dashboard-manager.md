# Agent: dashboard-manager

## Purpose

Create and manage Corezoid dashboards that visualize task counters from running processes. Dashboards show how many tasks are in specific nodes — useful for monitoring process health and business metrics.

## When to invoke

- User asks to "create a dashboard" or "visualize" process data
- User wants to monitor how many tasks succeed/fail in a process
- User wants to track business metrics from Corezoid processes
- After deploying a new process, to add observability

## Required inputs

| Parameter | Description |
|-----------|-------------|
| `API_LOGIN`, `SECRET`, `BASE_URL`, `COMPANY_ID` | Auth credentials |
| Process ID(s) | Which process(es) to measure |
| Metric goals | What the user wants to track (completion rate, error count, etc.) |

## MCP tools used

```
get_process_scheme    — find node IDs to measure
create_dashboard      — create dashboard container
add_chart             — add visualization charts
get_dashboard         — verify chart series after creation
```

## Workflow

### Step 1 — Identify nodes to measure

```python
result = client.get_process_scheme(process_id)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']

# End nodes (obj_type: 2) accumulate tasks — best for measurement
end_nodes = [n for n in nodes if n['obj_type'] == 2]
# Note: intermediate nodes pass tasks through instantly — measuring them shows 0 in real-time
```

### Step 2 — Create dashboard

```python
result = client.make_request({'ops': [{
    'obj': 'dashboard',
    'obj_type': 'create',
    'title': 'Weather Connector Monitoring',
    'description': 'Task success and error rates'
}]})
dashboard_id = result['ops'][0]['obj_id']
```

### Step 3 — Add charts

```python
result = client.make_request({'ops': [{
    'obj': 'chart',
    'obj_type': 'create',
    'dashboard_id': dashboard_id,
    'title': 'Success vs Error Tasks',
    'chart_type': 'column',   # ALWAYS 'column', NEVER 'bar'
    'metrics': [
        {'conv_id': process_id, 'node_id': success_node_id, 'title': 'Completed'},
        {'conv_id': process_id, 'node_id': error_node_id, 'title': 'Errors'}
    ]
}]})
chart_id = result['ops'][0]['obj_id']
```

### Step 4 — Verify series

After creating, ALWAYS verify that `series` is populated:

```python
verify = client.make_request({'ops': [{'obj': 'chart', 'obj_type': 'show', 'id': chart_id}]})
series = verify['ops'][0].get('series', [])
if not series:
    # Re-create chart — metrics didn't attach correctly
    pass
```

## Chart types

| Type | When to use |
|------|------------|
| `column` | Compare values — success vs error, node A vs node B |
| `pie` | Show proportions — what % of tasks succeed |
| `funnel` | Show drop-off through sequential steps |
| `table` | Tabular view of multiple metrics |

**`bar` does NOT exist** — will error silently.

## Modifying charts — full payload required

When modifying a chart, ALWAYS include `obj_type` and the complete `series` array. Partial modify returns validation error.

```python
client.make_request({'ops': [{
    'obj': 'chart',
    'obj_type': 'modify',
    'id': chart_id,
    'title': 'Updated Title',
    'chart_type': 'column',
    'series': [
        {'conv_id': process_id, 'node_id': node_id, 'title': 'Metric Label'}
    ]
}]})
```

## Grid layout

Use `width`/`height` fields (NOT `w`/`h`) — wrong field names cause validation errors.

```python
{'width': 6, 'height': 4}
```

## Real-time monitoring

Real-time mode works **only** for these node types:
- End nodes (`obj_type: 2`) — tasks finished the process
- Waiting for Callback — tasks waiting for external callback
- Delay nodes — tasks paused for a time period

Intermediate nodes (Code, API Call, Condition) pass tasks through instantly — they always show 0 in real-time.

## Knowledge references

- `knowledge/gotchas.md` — #15: dashboard chart type is `column` not `bar`
- `knowledge/auth.md` — credentials and batching
- `docs/dashboards/README.md` — full dashboard documentation
