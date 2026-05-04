# Corezoid Dashboards

Dashboards are Corezoid's built-in reporting and visualization layer. They display statistics
from running processes by reading **task counters** accumulated in process nodes.

## How dashboards work

A dashboard is a collection of **charts**. Each chart pulls data from one or more process nodes —
specifically, it counts how many tasks are currently sitting in (or have passed through) a given node.

This means:
- Processes must be running and have tasks flowing through them for data to appear
- The data source is always a **process node**, not an external database
- Real-time mode only works for nodes that accumulate tasks: End nodes, Waiting for Callback nodes, Delay nodes, Set State nodes

## Dashboard object

```json
{
  "obj_type": "dashboard",
  "title": "Dashboard Name",
  "description": "Optional description",
  "status": "active"
}
```

## Managing dashboards via API

All dashboard management uses the Corezoid API at `/api/2/json` with HMAC-SHA1 authentication.

**Authentication headers required:**
- `X-API-Login` — your API login
- `X-API-Timestamp` — current GMT Unix timestamp
- `X-API-Sign` — HMAC-SHA1 signature of `timestamp + request_body` using your secret key

### Create a dashboard

```json
{
  "ops": [{
    "obj": "dashboard",
    "obj_type": "create",
    "title": "My Dashboard",
    "description": "Sales metrics"
  }]
}
```

### Add a chart to a dashboard

```json
{
  "ops": [{
    "obj": "chart",
    "obj_type": "create",
    "dashboard_id": 12345,
    "title": "Task Completion Rate",
    "chart_type": "column",
    "metrics": [
      {
        "conv_id": 67890,
        "node_id": "final-node-id",
        "title": "Completed"
      }
    ]
  }]
}
```

**After creating a chart, always verify** that `series` is populated by running a `get` operation on the chart. If `series` is empty, the chart may not render correctly.

### Modify a chart — full payload required

When modifying a chart, you **must include `obj_type` and the full `series` array** — partial updates are not supported. Omitting either field returns `Key 'obj_type' is required` or `Key 'series' is required`.

```json
{
  "ops": [{
    "obj": "chart",
    "obj_type": "modify",
    "id": 99999,
    "title": "Updated Title",
    "chart_type": "column",
    "series": [
      {
        "conv_id": 67890,
        "node_id": "final-node-id",
        "title": "Completed"
      }
    ]
  }]
}
```

### Dashboard grid layout

Charts are positioned on a grid. Use `width` and `height` (NOT `w`/`h`) for chart dimensions:

```json
{
  "width": 6,
  "height": 4
}
```

### Other dashboard operations

| Operation | `obj_type` value |
|-----------|-----------------|
| Get dashboard with charts | `show` |
| Update dashboard | `modify` |
| Move to trash | `delete` |
| Restore from trash | `restore` |
| Permanently delete | `destroy` |
| Share (link to users/groups/API keys) | `link` |
| Add to favorites | `favorite` |

## Chart types

| Type | Use for |
|------|---------|
| `column` | Comparing values across categories or over time (vertical bars) |
| `pie` | Showing proportions of a whole |
| `funnel` | Visualizing drop-off through a sequential process |
| `table` | Tabular view of metric values |

> **Note:** The chart type is `column`, not `bar`. Using `bar` will result in an unknown chart type error.

## Metrics (data binding)

Each chart metric points to a specific process node:

```json
{
  "conv_id": 67890,
  "node_id": "abc123def456",
  "title": "Label shown on chart"
}
```

Multiple metrics can be added to one chart for comparison (e.g., success vs. error node counts).

## Real-time mode

Enable in the Corezoid UI via the "Real-Time" toggle. Works only for nodes that hold tasks:
- **End nodes** (`obj_type: 2`) — tasks accumulate after process completion
- **Waiting for Callback** — tasks wait for external HTTP callback
- **Delay** — tasks pause for a set time
- **Set State** — tasks reside in a state

**Does NOT work for** intermediate logic nodes (Code, API Call, Condition, etc.) — tasks pass through these too quickly.

## Drill-down (linked dashboards)

Charts can link to other dashboards for hierarchical drill-down. When a user clicks a chart
segment, it opens the linked dashboard filtered to that data. Configure via "Add chart link"
in the chart Metrics settings.

## Time range

Set a time range on the dashboard that applies to all charts simultaneously. Options typically
include last hour, last 24 hours, last 7 days, last 30 days, and custom ranges.

## UI workflow

1. Workspace tab → **Create → Dashboard**
2. Enter name and optional description
3. Click **+** to add a chart
4. Set chart type (column / pie / funnel / table)
5. Click **+ Add Metrics** — select process, then select specific nodes
6. Save and enable Real-Time mode if needed

## Best practices

- Use **End nodes** as primary metric sources — they reliably accumulate completed tasks
- Name metrics descriptively (not just node IDs) — labels appear directly on charts
- Group related processes in one dashboard (e.g., all payment flow metrics together)
- For funnel charts, add metrics in the order of the actual process flow
- Create separate dashboards for operations (real-time) vs. reporting (historical ranges)
- After chart creation, always verify `series` via `get` before considering the chart complete
