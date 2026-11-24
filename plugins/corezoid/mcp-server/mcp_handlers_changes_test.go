package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func okList(items ...map[string]interface{}) interface{} {
	list := make([]interface{}, 0, len(items))
	for _, it := range items {
		list = append(list, it)
	}
	return map[string]interface{}{"request_proc": "ok",
		"ops": []interface{}{map[string]interface{}{"proc": "ok", "list": list}}}
}

// ---- list-changed-processes --------------------------------------------------

func TestHandleToolCall_ListChangedProcesses_WalksAndFilters(t *testing.T) {
	resetGlobals(t)
	now := float64(time.Now().Unix())
	old := now - 10*24*3600
	var seen []map[string]interface{}
	srv, _ := mockAPIServer(t, func(ops []map[string]interface{}) interface{} {
		seen = append(seen, ops...)
		switch ops[0]["obj_id"] {
		case float64(10): // the stage: a subfolder, one fresh and one old process
			return okList(
				map[string]interface{}{"obj": "folder", "obj_id": float64(11), "title": "sub"},
				map[string]interface{}{"obj": "conv", "obj_id": float64(101), "title": "fresh", "change_time": now - 60},
				map[string]interface{}{"obj": "conv", "obj_id": float64(102), "title": "old", "change_time": old},
			)
		case float64(11): // servers that tag entries by a string obj_type instead of obj
			return okList(
				map[string]interface{}{"obj_type": "conv", "obj_id": float64(103), "title": "nested", "change_time": now - 120},
			)
		}
		return okList()
	})
	setProjectAuth(t, srv.URL)
	setVarTestAuth(t)

	result, isErr := handleToolCall(context.Background(), "list-changed-processes", map[string]interface{}{
		"folder_id": float64(10), "since_hours": float64(24),
	})
	if isErr {
		t.Fatalf("unexpected error: %s", result)
	}
	var out struct {
		FoldersScanned   int `json:"folders_scanned"`
		ProcessesScanned int `json:"processes_scanned"`
		Changed          []struct {
			ProcessID int64  `json:"process_id"`
			Path      string `json:"path"`
		} `json:"changed"`
	}
	if err := json.Unmarshal([]byte(result), &out); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, result)
	}
	if out.FoldersScanned != 2 || out.ProcessesScanned != 3 {
		t.Errorf("scanned folders=%d processes=%d, want 2 and 3", out.FoldersScanned, out.ProcessesScanned)
	}
	got := map[int64]string{}
	for _, c := range out.Changed {
		got[c.ProcessID] = c.Path
	}
	if _, ok := got[101]; !ok {
		t.Errorf("fresh process 101 missing: %s", result)
	}
	if p, ok := got[103]; !ok || p != "/sub" {
		t.Errorf("nested process 103 missing or wrong path %q: %s", p, result)
	}
	if _, ok := got[102]; ok {
		t.Errorf("old process 102 must be filtered out: %s", result)
	}
	for _, op := range seen {
		if op["type"] != "list" {
			t.Errorf("only list ops may be sent, got %v", op["type"])
		}
	}
}

func TestHandleToolCall_ListChangedProcesses_BadArgs(t *testing.T) {
	resetGlobals(t)
	srv, _ := mockAPIServer(t, func(ops []map[string]interface{}) interface{} { return okList() })
	setProjectAuth(t, srv.URL)
	setVarTestAuth(t)
	if res, isErr := handleToolCall(context.Background(), "list-changed-processes", map[string]interface{}{}); !isErr {
		t.Errorf("missing folder_id must fail, got %s", res)
	}
	if res, isErr := handleToolCall(context.Background(), "list-changed-processes", map[string]interface{}{
		"folder_id": float64(1), "since_hours": float64(0)}); !isErr || !strings.Contains(res, "since_hours") {
		t.Errorf("since_hours=0 must fail, got %s", res)
	}
}

// ---- list-process-history ----------------------------------------------------

func TestHandleToolCall_ListProcessHistory_ChangesOnlyAndStopsAtSince(t *testing.T) {
	resetGlobals(t)
	now := float64(time.Now().Unix())
	var captured []map[string]interface{}
	srv, _ := mockAPIServer(t, func(ops []map[string]interface{}) interface{} {
		captured = append(captured, ops...)
		return okList(
			map[string]interface{}{"type": "confirm", "obj": "commit", "create_time": now - 10, "user_name": "a",
				"request": map[string]interface{}{"version": float64(7)}},
			map[string]interface{}{"type": "modify", "obj": "node", "create_time": now - 20, "user_name": "a",
				"request": map[string]interface{}{"version": float64(7), "obj_id": "n1", "title": "Step"}},
			map[string]interface{}{"type": "get", "obj": "obj_scheme", "create_time": now - 30, "user_name": "b"},
			map[string]interface{}{"type": "modify", "obj": "node", "create_time": now - 3*24*3600, "user_name": "c"},
		)
	})
	setProjectAuth(t, srv.URL)
	setVarTestAuth(t)

	result, isErr := handleToolCall(context.Background(), "list-process-history", map[string]interface{}{
		"process_id": float64(55), "since_hours": float64(24),
	})
	if isErr {
		t.Fatalf("unexpected error: %s", result)
	}
	var out struct {
		Events []struct {
			Action    string `json:"action"`
			Version   int64  `json:"version"`
			NodeTitle string `json:"node_title"`
			User      string `json:"user"`
		} `json:"events"`
		Complete bool `json:"complete"`
	}
	if err := json.Unmarshal([]byte(result), &out); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, result)
	}
	if len(out.Events) != 2 || !out.Complete {
		t.Fatalf("want 2 change events (read and older entry dropped), complete; got %s", result)
	}
	if out.Events[0].Action != "confirm commit" || out.Events[0].Version != 7 {
		t.Errorf("first event = %+v, want confirm commit v7", out.Events[0])
	}
	if out.Events[1].NodeTitle != "Step" {
		t.Errorf("node title not read: %+v", out.Events[1])
	}
	op := captured[0]
	for k, want := range map[string]interface{}{"type": "list", "obj": "history", "obj_type": "conv",
		"obj_id": float64(55)} {
		if op[k] != want {
			t.Errorf("op[%s] = %v, want %v", k, op[k], want)
		}
	}

	result, _ = handleToolCall(context.Background(), "list-process-history", map[string]interface{}{
		"process_id": float64(55), "since_hours": float64(24), "include_reads": "true",
	})
	if !strings.Contains(result, "get obj_scheme") {
		t.Errorf("include_reads must keep read entries: %s", result)
	}
}
