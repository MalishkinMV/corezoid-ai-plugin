package main

// Two read-only cz-structure actions for change tracking:
//
//   list-changed-processes  walk a stage (or any folder) and return the processes
//                           whose change_time is at or after `since`
//   list-process-history    the change log of one process (obj=history):
//                           who created / modified / deleted which node, commits
//
// Both send only `list` ops. Output is one JSON document so a script can read it
// as is; nothing is written to disk.

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	changedMaxFolders = 3000
	historyPageSize   = 200
	historyMaxPages   = 25
)

// optInt returns an optional integer argument or def when absent.
func optInt(args map[string]interface{}, key string, def int) (int, error) {
	if _, ok := args[key]; !ok {
		return def, nil
	}
	return intArg(args, key)
}

// sinceArg reads `since` (unix seconds) or `since_hours` (default 24) into unix seconds.
func sinceArg(args map[string]interface{}) (int64, error) {
	if _, ok := args["since"]; ok {
		s, err := intArg(args, "since")
		return int64(s), err
	}
	h, err := optInt(args, "since_hours", 24)
	if err != nil {
		return 0, err
	}
	if h <= 0 || h > 24*90 {
		return 0, fmt.Errorf("since_hours must be 1..2160, got %d", h)
	}
	return time.Now().Unix() - int64(h)*3600, nil
}

func listRaw(v *Executor, folderID int) ([]map[string]any, error) {
	ops := []map[string]any{{"type": "list", "obj": "folder", "obj_id": folderID, "company_id": v.WorkspaceID}}
	resp, err := v.req("list_folder", ops)
	if err != nil {
		return nil, err
	}
	first, err := firstOp(resp)
	if err != nil {
		return nil, err
	}
	list, _ := first["list"].([]any)
	out := make([]map[string]any, 0, len(list))
	for _, it := range list {
		if m, ok := it.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// kindOf: servers differ — some tag entries with obj ("folder"/"conv"),
// others with a string obj_type.
func kindOf(m map[string]any) string {
	if s, ok := m["obj"].(string); ok && s != "" {
		return s
	}
	if s, ok := m["obj_type"].(string); ok {
		return s
	}
	return ""
}

func num(m map[string]any, k string) int64 {
	if f, ok := m[k].(float64); ok {
		return int64(f)
	}
	return 0
}

func str(m map[string]any, k string) string {
	if s, ok := m[k].(string); ok {
		return s
	}
	return ""
}

func handleListChangedProcesses(ctx context.Context, args map[string]interface{}) (string, bool) {
	folderID, err := intArg(args, "folder_id")
	if err != nil {
		return "Error: " + err.Error(), true
	}
	since, err := sinceArg(args)
	if err != nil {
		return "Error: " + err.Error(), true
	}
	v := NewValidator(ctx, 0)
	type item struct {
		ProcessID  int64  `json:"process_id"`
		Title      string `json:"title"`
		Path       string `json:"path"`
		ChangeTime int64  `json:"change_time"`
		ChangedAt  string `json:"changed_at"`
		Version    int64  `json:"version,omitempty"`
		Status     string `json:"status,omitempty"`
		ConvType   string `json:"conv_type,omitempty"`
	}
	type qe struct {
		id   int
		path string
	}
	queue := []qe{{folderID, ""}}
	changed := []item{}
	folders, convs := 0, 0
	errs := []string{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		folders++
		if folders > changedMaxFolders {
			errs = append(errs, fmt.Sprintf("stopped after %d folders", changedMaxFolders))
			break
		}
		list, err := listRaw(v, cur.id)
		if err != nil {
			errs = append(errs, fmt.Sprintf("folder %d: %v", cur.id, err))
			continue
		}
		for _, m := range list {
			title := str(m, "title")
			switch kindOf(m) {
			case "folder":
				queue = append(queue, qe{int(num(m, "obj_id")), cur.path + "/" + title})
			case "conv":
				convs++
				if ct := num(m, "change_time"); ct >= since {
					changed = append(changed, item{
						ProcessID: num(m, "obj_id"), Title: title, Path: cur.path, ChangeTime: ct,
						ChangedAt: time.Unix(ct, 0).UTC().Format(time.RFC3339), Version: num(m, "version"),
						Status: str(m, "status"), ConvType: str(m, "conv_type"),
					})
				}
			}
		}
	}
	out := map[string]any{
		"folder_id": folderID, "since": since, "since_utc": time.Unix(since, 0).UTC().Format(time.RFC3339),
		"folders_scanned": folders, "processes_scanned": convs, "changed": changed, "errors": errs,
	}
	b, _ := json.MarshalIndent(out, "", " ")
	return string(b), false
}

var historyChangeActions = map[string]bool{
	"create node": true, "modify node": true, "delete node": true, "confirm commit": true,
	"modify conv": true, "create conv": true, "deploy conv": true,
}

func handleListProcessHistory(ctx context.Context, args map[string]interface{}) (string, bool) {
	processID, err := intArg(args, "process_id")
	if err != nil {
		return "Error: " + err.Error(), true
	}
	since, err := sinceArg(args)
	if err != nil {
		return "Error: " + err.Error(), true
	}
	all := boolishArg(args, "include_reads")
	v := NewValidator(ctx, 0)
	type event struct {
		Time      string `json:"time"`
		TS        int64  `json:"ts"`
		UserID    int64  `json:"user_id,omitempty"`
		User      string `json:"user"`
		Action    string `json:"action"`
		Version   int64  `json:"version,omitempty"`
		NodeID    string `json:"node_id,omitempty"`
		NodeTitle string `json:"node_title,omitempty"`
	}
	events := []event{}
	pages, done := 0, false
	for offset := 0; !done && pages < historyMaxPages; offset += historyPageSize {
		pages++
		ops := []map[string]any{{"type": "list", "obj": "history", "obj_id": processID, "obj_type": "conv",
			"limit": historyPageSize, "offset": offset, "company_id": v.WorkspaceID}}
		resp, err := v.req("list_history", ops)
		if err != nil {
			return fmt.Sprintf("Error: %v", err), true
		}
		first, err := firstOp(resp)
		if err != nil {
			return fmt.Sprintf("Error: %v", err), true
		}
		list, _ := first["list"].([]any)
		for _, it := range list {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			ts := num(m, "create_time")
			if ts < since {
				done = true
				break
			}
			action := str(m, "type") + " " + str(m, "obj")
			if !all && !historyChangeActions[action] {
				continue
			}
			e := event{Time: time.Unix(ts, 0).UTC().Format(time.RFC3339), TS: ts, UserID: num(m, "user_id"),
				User: str(m, "user_name"), Action: action}
			if r, ok := m["request"].(map[string]any); ok {
				e.Version = num(r, "version")
				if str(m, "obj") == "node" {
					e.NodeTitle = str(r, "title")
					if id, ok := r["obj_id"]; ok {
						e.NodeID = fmt.Sprint(id)
					}
				}
			}
			events = append(events, e)
		}
		if len(list) < historyPageSize {
			done = true
		}
	}
	out := map[string]any{"process_id": processID, "since": since,
		"since_utc": time.Unix(since, 0).UTC().Format(time.RFC3339), "events": events, "complete": done}
	b, _ := json.MarshalIndent(out, "", " ")
	return string(b), false
}
