package main

import (
	"testing"
)

func TestParseSSCmd(t *testing.T) {
	stdout := []byte(`Netid  State   Recv-Q  Send-Q   Local Address:Port    Peer Address:Port  Process
tcp    LISTEN  0       128      127.0.0.1:3000         0.0.0.0:*          users:(("node",pid=1234,fd=20))
tcp    LISTEN  0       128      0.0.0.0:5432           0.0.0.0:*          users:(("postgres",pid=5678,fd=10))
tcp    LISTEN  0       128      [::1]:8080             [::]:*             users:(("go",pid=9012,fd=6))
`)

	processes, err := ParseSSCmd(stdout)
	if err != nil {
		t.Fatalf("ParseSSCmd returned error: %v", err)
	}

	if len(processes) != 3 {
		t.Fatalf("Expected 3 processes, got %d", len(processes))
	}

	// Test first process (node)
	if processes[0].Name != "node" {
		t.Errorf("Expected name 'node', got '%s'", processes[0].Name)
	}
	if processes[0].PID != 1234 {
		t.Errorf("Expected PID 1234, got %d", processes[0].PID)
	}
	if processes[0].Port != 3000 {
		t.Errorf("Expected port 3000, got %d", processes[0].Port)
	}
	if processes[0].Address != "127.0.0.1" {
		t.Errorf("Expected address '127.0.0.1', got '%s'", processes[0].Address)
	}

	// Test third process (go with IPv6)
	if processes[2].Name != "go" {
		t.Errorf("Expected name 'go', got '%s'", processes[2].Name)
	}
	if processes[2].Port != 8080 {
		t.Errorf("Expected port 8080, got %d", processes[2].Port)
	}
	if processes[2].Address != "::1" {
		t.Errorf("Expected address '::1', got '%s'", processes[2].Address)
	}
}

func TestParseSSCmd_EmptyOutput(t *testing.T) {
	stdout := []byte(`Netid  State   Recv-Q  Send-Q   Local Address:Port    Peer Address:Port  Process
`)

	processes, err := ParseSSCmd(stdout)
	if err != nil {
		t.Fatalf("ParseSSCmd returned error: %v", err)
	}

	if len(processes) != 0 {
		t.Errorf("Expected 0 processes, got %d", len(processes))
	}
}

func TestParseLsofCmd(t *testing.T) {
	// Sample output from `lsof -i -P -n -sTCP:LISTEN`
	stdout := []byte(`COMMAND   PID   USER   FD   TYPE   DEVICE  SIZE/OFF  NODE  NAME
node     1234   user   20u  IPv4   12345   0t0       TCP   127.0.0.1:3000 (LISTEN)
postgres 5678   user   10u  IPv4   23456   0t0       TCP   *:5432 (LISTEN)
`)

	processes, err := ParseLsofCmd(stdout)
	if err != nil {
		t.Fatalf("ParseLsofCmd returned error: %v", err)
	}

	if len(processes) != 2 {
		t.Fatalf("Expected 2 processes, got %d", len(processes))
	}

	// Test first process
	if processes[0].Name != "node" {
		t.Errorf("Expected name 'node', got '%s'", processes[0].Name)
	}
	if processes[0].PID != 1234 {
		t.Errorf("Expected PID 1234, got %d", processes[0].PID)
	}
	if processes[0].Port != 3000 {
		t.Errorf("Expected port 3000, got %d", processes[0].Port)
	}

	// Test second process
	if processes[1].Name != "postgres" {
		t.Errorf("Expected name 'postgres', got '%s'", processes[1].Name)
	}
	if processes[1].Port != 5432 {
		t.Errorf("Expected port 5432, got %d", processes[1].Port)
	}
}

func TestParseNetstatCmd(t *testing.T) {
	// Sample output from `netstat -ano` on Windows
	stdout := []byte(`Active Connections

  Proto  Local Address          Foreign Address        State           PID
  TCP    0.0.0.0:135            0.0.0.0:0              LISTENING       1234
  TCP    0.0.0.0:445            0.0.0.0:0              LISTENING       5678
  TCP    192.168.1.1:50000      192.168.1.2:443        ESTABLISHED     9999
`)

	processes, err := ParseNetstatCmd(stdout)
	if err != nil {
		t.Fatalf("ParseNetstatCmd returned error: %v", err)
	}

	// Should only get LISTENING processes
	if len(processes) != 2 {
		t.Fatalf("Expected 2 processes (LISTENING only), got %d", len(processes))
	}

	if processes[0].PID != 1234 {
		t.Errorf("Expected PID 1234, got %d", processes[0].PID)
	}
	if processes[0].Port != 135 {
		t.Errorf("Expected port 135, got %d", processes[0].Port)
	}

	if processes[1].PID != 5678 {
		t.Errorf("Expected PID 5678, got %d", processes[1].PID)
	}
}
