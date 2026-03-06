**English** | [中文](README_zh.md)

# peerclaw-cli

The PeerClaw command-line tool. Interact with PeerClaw Server via REST API to manage agents, send messages, and check service status.

## Installation

```bash
cd cli
go build -o peerclaw ./cmd/peerclaw
```

## Usage

### Configuration

Connects to `http://localhost:8080` by default. You can change this with an environment variable or the config file:

```bash
# Environment variable
export PEERCLAW_SERVER=http://my-server:8080

# Or config file
peerclaw config set server http://my-server:8080
peerclaw config show
```

### Agent Management

```bash
# List all agents
peerclaw agent list

# Filter by protocol
peerclaw agent list -protocol a2a

# View agent details
peerclaw agent get <agent-id>

# Register an agent
peerclaw agent register -name "MyAgent" -url http://localhost:3000 -protocols a2a,mcp

# Delete an agent
peerclaw agent delete <agent-id>
```

### Sending Messages

```bash
peerclaw send -from agent-a -to agent-b -protocol a2a -payload '{"message": "hello"}'
```

### Health Check

```bash
peerclaw health

# JSON output
peerclaw health -output json
```

### Output Formats

All list commands support the `-output` flag:

- `table` (default): table format
- `json`: JSON format

```bash
peerclaw agent list -output json
```
