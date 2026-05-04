---
name: corezoid-dashboard-manager
description: >
  Creates and manages Corezoid dashboards — adds charts (column, pie, funnel, table), binds metrics
  to process nodes, configures real-time mode, and sets up drill-down linking between dashboards.
  Activate whenever a user asks to create a dashboard, add a chart, visualize process metrics,
  set up reporting for a Corezoid process, configure real-time monitoring, or asks what
  dashboards or charts exist. Also activate when the user wants to show task counts, completion
  rates, error rates, or any other process statistics visually in Corezoid.
---

# Corezoid Dashboard Manager

## What dashboards are

A Corezoid dashboard visualizes **task counters in process nodes** — it shows how many tasks
are accumulated in (or have passed through) specific nodes. The data comes directly from
running processes, not from external databases.

Key implication: processes must be deployed and have tasks flowing through them for data to appear.

## Required credentials (ask if missing)

| Parameter | Description |
|-----------|-------------|
| `API_URL` | Corezoid admin URL (e.g. `https://admin.corezoid.com`) |
| `API_TOKEN` | Personal API token |
| `WORKSPACE_ID` | Workspace ID |

All dashboard API calls go to `<API_URL>/api/2/json`.

## Authentication

Every API request requires HMAC-SHA1 authentication:
- `X-API-Login` — API login
- `X-API-Timestamp` — current GMT Unix timestamp  
- `X-API-Sign` — HMAC-SHA1 of `timestamp + request_body` using the API secret key

---

## Workflow: Create a dashboard with charts

### Step 1 — Identify the processes and nodes to visualize

Ask the user:
- Which process(es) do they want to monitor?
- Which nodes in those processes represent the metrics they care about? (End nodes are best — they accumulate completed tasks)
- What kind of visualization: comparison (column), proportions (pie), sequential drop-off (funnel), or tabular (table)?

Export the target process to find node IDs:
```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh fetch-process <PROC_ID> .processes/target_process.json
```

Find node IDs in `scheme.nodes[].id` — note the IDs of the nodes you want to measure.

### Step 2 — Create the dashboard

```json
{
  "ops": [{
    "obj": "dashboard",
    "obj_type": "create",
    "title": "<Dashboard Name>",
    "description": "<Optional description>"
  }]
}
```

Note the returned `dashboard_id` — needed for adding charts.

### Step 3 — Add charts

One chart per visualization. Chart types:

| Type | When to use |
|------|------------|
| `column` | Comparing values across nodes or time periods (vertical bars) |
| `pie` | Showing how tasks split across outcomes |
| `funnel` | Visualizing drop-off through a sequential process flow |
| `table` | Tabular display of metric values |

> **Critical:** Use `column` (NOT `bar`) — `bar` is not a valid Corezoid chart type.

```json
{
  "ops": [{
    "obj": "chart",
    "obj_type": "create",
    "dashboard_id": <dashboard_id>,
    "title": "<Chart Title>",
    "chart_type": "column",
    "metrics": [
      {
        "conv_id": <process_id>,
        "node_id": "<node_id>",
        "title": "<Metric Label>"
      }
    ]
  }]
}
```

Multiple metrics on one chart = multiple bars/slices for comparison.

### Step 4 — Verify series after creation

After creating a chart, **always call `get`** to verify that `series` is populated. If `series` is empty, the chart will not render correctly.

```json
{
  "ops": [{
    "obj": "chart",
    "obj_type": "show",
    "id": <chart_id>
  }]
}
```

If `series` is empty or missing, re-create the chart with the metrics.

### Step 5 — Advise on real-time mode

Real-time mode works ONLY for these node types:
- **End nodes** (`obj_type: 2`) — tasks that finished the process
- **Waiting for Callback** — tasks waiting for external HTTP callback
- **Delay** — tasks paused for a time period
- **Set State** — tasks in a named state

If the user wants real-time monitoring, ensure the metrics point to these node types.
Intermediate nodes (Code, API Call, Condition) pass tasks through instantly — real-time shows 0.

---

## Modifying charts — full payload required

When modifying a chart, **always include `obj_type` and the full `series` array**. Partial updates are NOT supported — omitting either field returns a validation error.

```json
{
  "ops": [{
    "obj": "chart",
    "obj_type": "modify",
    "id": <chart_id>,
    "title": "<Updated Title>",
    "chart_type": "column",
    "series": [
      {
        "conv_id": <process_id>,
        "node_id": "<node_id>",
        "title": "<Metric Label>"
      }
    ]
  }]
}
```

---

## Dashboard grid layout

Charts are positioned on a grid. Use `width` and `height` fields (NOT `w`/`h`) for chart dimensions. Using `w`/`h` causes a validation error.

```json
{
  "width": 6,
  "height": 4
}
```

---

## Dashboard operations reference

| Goal | `obj_type` value | Required fields |
|------|-----------------|-----------------|
| Create dashboard | `create` | `title` |
| Get dashboard + charts | `show` | `id` (dashboard ID) |
| Update dashboard | `modify` | `id`, fields to change |
| Soft delete (trash) | `delete` | `id` |
| Restore from trash | `restore` | `id` |
| Permanent delete | `destroy` | `id` |
| Share with users/groups | `link` | `id`, `users`/`groups`/`api_keys` |

---

## Drill-down (linked dashboards)

For hierarchical data — when clicking a chart segment should open a more detailed view:
1. Create the high-level dashboard first
2. Create the detail dashboard
3. In the chart Metrics settings, use "Add chart link" to point a metric to the detail dashboard
4. Users can click chart segments to drill down

---

## Funnel chart — node ordering matters

For funnel charts, add metrics in the **same order as the process flow**:

```json
"metrics": [
  {"conv_id": 123, "node_id": "start-node-id", "title": "Started"},
  {"conv_id": 123, "node_id": "mid-node-id", "title": "Processed"},
  {"conv_id": 123, "node_id": "final-node-id", "title": "Completed"}
]
```

The funnel visualizes drop-off from each step to the next.

---

## Best practices

- Use **End nodes** as primary metric sources — they reliably hold completed tasks
- For error rate charts: add both Final (success) and Error end nodes as metrics on one column chart
- Name metrics descriptively — labels appear directly on charts
- Create separate dashboards: one for real-time ops monitoring, one for historical reporting
- Group all metrics from one business flow (e.g., payment processing) on a single dashboard
- Always verify `series` after chart creation — empty series means the chart won't render

---

## Reference

Full documentation: `docs/dashboards/README.md`
