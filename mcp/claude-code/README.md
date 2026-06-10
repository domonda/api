# domonda MCP — Claude Code

This directory contains setup instructions and configuration for connecting
[**Claude Code**](https://claude.ai/claude-code) — Anthropic's CLI, Desktop,
Web, and IDE-extension agent — to the [domonda](https://domonda.app) financial
data platform via its
[Model Context Protocol (MCP)](https://modelcontextprotocol.io) server.

Claude Code speaks Streamable HTTP MCP natively, so no `mcp-remote` bridge and
no Node.js runtime are required. If you are configuring **Claude Desktop**,
**Cowork**, or the **claude.ai web app**, use
[`../claude-desktop/`](../claude-desktop/) instead — those surfaces share the
Custom Connector UI (and fall back to `mcp-remote` for localhost).

All query operations are **read-only**. The server exposes 28 tools covering
documents, invoices, payments, money transactions, companies, GL accounts,
partner companies, document workflows, web app deep links, ad-hoc SQL
queries, and a `list_capabilities` discovery tool. The single write tool is
`add_document`, which uploads a document via the same pipeline as the public
upload endpoint.

> **Live tool descriptions are in German.** The tool and parameter
> descriptions surfaced by the running MCP server (and therefore what your
> agent actually matches against user prompts) are written in German with
> domain vocabulary pulled from the domonda-app UI (Sachkonto, Rechnung,
> Gutschrift, Geschäftspartner, Fälligkeitsdatum, Mandant, Freigabe, …) and
> a "Nutzen wenn …" trigger-phrase block for each tool. End users chat in
> German, so the descriptions are optimised for the words they actually
> use. Agents that are uncertain which tool to call should call
> `list_capabilities` first for a topic-grouped catalog.

## Prerequisites

- [Claude Code](https://claude.ai/claude-code) installed (CLI, Desktop app,
  Web app, or IDE extension — all share the same MCP config)
- A domonda API key (JWT) — obtain from the domonda admin panel — or an
  Auth0 OAuth access token for the MCP resource

## Directory Structure

```
claude-code/
├── SKILL.md              # Skill definition (frontmatter + tool reference)
├── mcp.json.example      # Example project-scope .mcp.json
├── scripts/
│   └── health_check.sh   # Connectivity check script
└── README.md             # This file
```

## Setup — pick one of three paths

### 1. CLI (recommended, user scope, one-liner)

Adds the server to your **user** config (`~/.claude.json`), available in every
project you open with Claude Code:

```bash
claude mcp add --scope user --transport http domonda \
  https://domonda.app/api/mcp/ \
  --header "Authorization: Bearer $DOMONDA_API_KEY"
```

Verify it landed:

```bash
claude mcp list           # should show: domonda (http)
```

To remove it later: `claude mcp remove domonda`.

### 2. Project scope (checked in, shared with teammates)

Copy the example into the root of a repo and rename it to `.mcp.json`:

```bash
cp mcp/claude-code/mcp.json.example /path/to/your/project/.mcp.json
```

The example uses `${DOMONDA_API_KEY}` environment-variable expansion, which
keeps the file safe to commit. Each teammate exports their own API key:

```bash
export DOMONDA_API_KEY=eyJhbGciOiJIUzI1NiIs...
```

Claude Code picks up `.mcp.json` automatically when opened at the repo root
and prompts for trust on first load.

### 3. Manual global config

If you prefer editing config files directly, paste the `domonda` block into
the top-level `mcpServers` key of `~/.claude.json`:

```json
{
  "mcpServers": {
    "domonda": {
      "type": "http",
      "url": "https://domonda.app/api/mcp/",
      "headers": {
        "Authorization": "Bearer YOUR_API_KEY_HERE"
      }
    }
  }
}
```

This is exactly what `claude mcp add --scope user` writes for you. The bearer
token is stored in **plaintext** in `~/.claude.json` — prefer path 1 (CLI) or
path 2 (project `.mcp.json` with env expansion) when possible.

### OAuth instead of an API key

All three paths above assume a static JWT API key. If you'd rather sign in
with your Auth0 account, **omit the `Authorization` header** — the domonda
MCP server exposes RFC 9728 OAuth discovery metadata at
`.well-known/oauth-protected-resource`, and Claude Code follows the
`WWW-Authenticate` challenge to complete the flow in your browser on first
use:

```bash
claude mcp add --scope user --transport http domonda \
  https://domonda.app/api/mcp/
```

This avoids persisting a bearer token on disk and is the right choice for
multi-company users who rely on the `X-Selected-Client-Company-ID` header.

> Some Claude Code versions in the v2.1.85+ range have regressions in the
> OAuth discovery flow
> ([anthropics/claude-code#44830](https://github.com/anthropics/claude-code/issues/44830));
> if OAuth gets stuck, fall back to the Bearer-token form above.

### Base URL: production vs local

In production the domonda-web-server is mounted behind an `/api` path prefix,
so the MCP endpoint is `https://domonda.app/api/mcp/`. When running the web
server locally there is **no** `/api` prefix — use `http://localhost:5001/mcp/`
(5001 is the default port; override with `PORT=…`).

## Verify Connectivity

### Health check script

```bash
 export DOMONDA_API_KEY="your-api-key-here"   # leading space keeps this out of shell history
./scripts/health_check.sh
```

Expected output:

```
Checking https://domonda.app/api/mcp/ ...
OK (HTTP 200)
```

### Inside Claude Code

Start a new Claude Code session and ask:

> "Use the domonda tools to get company info."

Claude should call `get_my_company` and return your company details.

## Usage

Once configured, Claude can call any of the 28 domonda tools.
See `SKILL.md` for the complete tool reference with parameters.

### Recommended first steps

1. **`list_capabilities`** — When unsure which tool fits, get a grouped
   catalog of all available tools first.
2. **`get_my_company`** — Confirm authentication and see which company the
   API key is scoped to.
3. **`list_tables`** / **`describe_table`** — Explore the available database
   schema before writing queries.
4. **`list_documents`** or **`list_invoices`** — Browse recent data.
5. **`execute_query`** — Run ad-hoc read-only SQL for anything the specialized
   tools don't cover.

### Tool categories

| Category  | Tools                                                                                                                                                                      |
|-----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Discovery | `list_capabilities`                                                                                                                                                        |
| Schema    | `list_tables`, `describe_table`                                                                                                                                            |
| Query     | `execute_query`                                                                                                                                                            |
| Documents | `list_documents`, `get_document`, `get_document_fulltext`, `search_documents`, `download_document_pdf`, `get_document_thumbnail`, `get_document_preview`, `add_document` (write) |
| Invoices  | `list_invoices`, `get_invoice`                                                                                                                                             |
| Payments  | `list_money_accounts`, `list_money_transactions`, `get_invoice_payments`                                                                                                   |
| Company   | `get_my_company`, `list_document_categories`, `list_gl_accounts`, `list_partner_companies`                                                                                 |
| User      | `get_user`                                                                                                                                                                 |
| Workflows | `list_document_workflows`, `list_document_workflow_steps`, `list_document_workflow_states`                                                                                 |
| App URLs  | `get_app_urls`, `get_document_urls`, `get_document_selection_url`                                                                                                          |

### Example prompts

- "Show me all invoices from January 2026 over 1,000"
- "Download the PDF for document 8a2f..."
- "Which GL accounts are in use?"
- "Search documents for 'electricity bill'"
- "Run a query to find the top 10 partners by invoice count"
- "Upload this invoice PDF into the OTHER_DOCUMENTS category"

## Authentication

The domonda MCP server uses **Streamable HTTP** transport with **Bearer token**
authentication. Two token types are supported:

1. **API key (HS256 JWT)** — the standard domonda API key. The JWT `sub` claim
   identifies the client company directly. This is the simplest option for
   programmatic access.
2. **OAuth (Auth0 RS256)** — an Auth0 access token issued for the MCP resource.
   The token's Auth0 `sub` claim is mapped to a Domonda user, who must be
   enabled and have access to the target company. Supports the
   `X-Selected-Client-Company-ID` header for multi-company users. Clients that
   support OAuth discovery will be guided automatically via the
   `WWW-Authenticate` header and the `.well-known/oauth-protected-resource`
   metadata endpoint (RFC 9728).

```
Authorization: Bearer <DOMONDA_API_KEY_OR_OAUTH_TOKEN>
```

The server enforces:

- **Read-only transactions for query tools** — no INSERT, UPDATE, DELETE, or DDL
- **Schema restriction** — only the `api` schema is accessible
- **Row limit** — max 1,000 rows per query
- **Query timeout** — 30-second default
- **Query length** — max 10,000 characters
- **Writes are limited to `add_document`** — processed through the same OCR,
  duplicate-detection, and category-ownership checks as the public API

## Troubleshooting

| Symptom                                       | Cause                                                       | Fix                                                                                   |
|-----------------------------------------------|-------------------------------------------------------------|---------------------------------------------------------------------------------------|
| `claude mcp list` doesn't show `domonda`      | Config not written, or written to a different scope         | Re-run the `claude mcp add` command and check `--scope`; inspect `~/.claude.json`     |
| Tools don't appear in a Claude Code session   | Session started before the config was written               | Restart the Claude Code session / reload the window                                   |
| `.mcp.json` not picked up                     | Session opened at the wrong working directory, or trust denied | Open Claude Code at the repo root; approve the trust prompt                        |
| `${DOMONDA_API_KEY}` literally in requests    | Env var not exported in the shell that launched Claude Code | `export DOMONDA_API_KEY=…` before starting Claude Code                                |
| HTTP 401                                      | Invalid, expired, or blocked token                          | Verify your API key or OAuth token; contact admin if blocked                          |
| HTTP 403                                      | Forbidden company access or insufficient OAuth scopes       | Check the selected company or required scopes                                         |
| HTTP 404                                      | Wrong endpoint URL                                          | Ensure the URL ends with `/mcp/` or `/mcp`                                            |
| Query returns no rows                         | Company scoping — data belongs to another company           | Confirm `get_my_company` returns the expected company                                 |
| `execute_query` rejected                      | SQL validation failed (DML, wrong schema)                   | Use only SELECT/WITH against the `api` schema                                         |
