# PeerClaw MCP Server Configuration

PeerClaw provides an MCP (Model Context Protocol) server that exposes agent discovery, invocation, and reputation tools to any MCP-compatible AI host.

## Available Tools

| Tool | Description | Annotations |
|------|-------------|-------------|
| `discover_agents` | Find agents by capabilities | ReadOnly, Idempotent |
| `invoke_agent` | Send a message to an agent via gateway | — |
| `get_agent_profile` | Get agent public profile | ReadOnly, Idempotent |
| `check_reputation` | Check agent reputation score | ReadOnly, Idempotent |
| `add_contact` | Add an agent to the contact whitelist | — |
| `remove_contact` | Remove an agent from the contact whitelist | — |
| `list_contacts` | List all contacts in the trust store | ReadOnly, Idempotent |
| `send_message` | Send an async message to a peer | — |
| `send_request` | Send a synchronous request and wait for response | — |
| `broadcast_message` | Send a message to multiple agents | — |
| `get_task` | Get status of an A2A task | ReadOnly, Idempotent |
| `list_tasks` | List all tracked A2A tasks | ReadOnly, Idempotent |

## Available Resources

| URI | Description |
|-----|-------------|
| `peerclaw://directory` | Browse the agent directory (JSON) |

## Claude Code

Add to `~/.claude.json` (or project-level `.claude/settings.json`):

```json
{
  "mcpServers": {
    "peerclaw": {
      "command": "peerclaw",
      "args": ["mcp", "serve", "--server", "http://localhost:8080"]
    }
  }
}
```

## VS Code (GitHub Copilot)

Add to `.vscode/settings.json`:

```json
{
  "mcp.servers": {
    "peerclaw": {
      "command": "peerclaw",
      "args": ["mcp", "serve", "--server", "http://localhost:8080"]
    }
  }
}
```

## Cursor

Add to Cursor MCP settings (`~/.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "peerclaw": {
      "command": "peerclaw",
      "args": ["mcp", "serve", "--server", "http://localhost:8080"]
    }
  }
}
```

## Windsurf

Add to `~/.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "peerclaw": {
      "command": "peerclaw",
      "args": ["mcp", "serve", "--server", "http://localhost:8080"]
    }
  }
}
```

## HTTP Transport Mode

For remote hosting or multi-client scenarios, use HTTP transport:

```bash
peerclaw mcp serve --transport http --port 8081 --server http://peerclaw-server:8080
```

Then configure clients to connect via Streamable HTTP at `http://host:8081/`.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PEERCLAW_SERVER` | PeerClaw server URL | `http://localhost:8080` |

## Verification

Test the MCP server manually via stdio:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | peerclaw mcp serve --server http://localhost:8080
```
