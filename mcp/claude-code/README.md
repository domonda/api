# domonda MCP — Claude Code

This directory contains setup instructions and configuration for connecting
[**Claude Code**](https://claude.ai/claude-code) — Anthropic's coding agent in
the terminal, IDE, Desktop, and web — to the [domonda](https://domonda.app) financial
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

- [Claude Code](https://claude.ai/claude-code) installed
- Either a domonda API key (JWT) from the domonda admin panel, or a
  domonda OAuth connector already signed in on claude.ai with an admin or
  super-admin account; see [OAuth setup](#oauth-through-a-claudeai-connector)

## Directory Structure

```
claude-code/
├── SKILL.md              # Skill definition (frontmatter + tool reference)
├── mcp.json.example      # Example project-scope .mcp.json
├── scripts/
│   └── health_check.sh   # Connectivity check script
└── README.md             # This file
```

The optional agent skill is separate from the MCP connection. To install it,
copy `SKILL.md` to `~/.claude/skills/domonda/SKILL.md` or to a project's
`.claude/skills/domonda/SKILL.md`. The folder must match `name: domonda`.

## Setup — choose an API-key path or the OAuth connector

### 1. CLI (user scope)

Adds the server to your **user** config (`~/.claude.json`), available in every
project you open with Claude Code:

```bash
claude mcp add --scope user --transport http domonda \
  https://domonda.app/api/mcp/ \
  --header 'Authorization: Bearer ${DOMONDA_API_KEY}'
```

The single quotes keep the key out of the command line and shell history.
Claude Code resolves the literal `${DOMONDA_API_KEY}` from its environment
when it loads the server entry. Set that variable in the process launching
Claude Code; `claude mcp list` can reveal a missing-variable warning.

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
keeps the file safe to commit. Each teammate supplies their own API key to
Claude Code's environment. In an interactive shell, this avoids putting the
value in shell history:

```bash
read -rsp 'Domonda API key: ' DOMONDA_API_KEY
export DOMONDA_API_KEY
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
        "Authorization": "Bearer ${DOMONDA_API_KEY}"
      }
    }
  }
}
```

Use the literal `${DOMONDA_API_KEY}` reference shown above. Do not paste the
key itself into `~/.claude.json`; a literal key there would be stored in
plaintext. The project `.mcp.json` path is convenient for teams.

### OAuth through a claude.ai connector

Add domonda at **claude.ai → Customize → Connectors** using the
[Claude connector setup](../claude-desktop/#option-1--custom-connector-ui-recommended).
Choose **Use Claude's published identity** for OAuth and sign in as a Domonda
admin or super-admin. When Claude Code uses the same claude.ai subscription
login, that connector appears automatically in `/mcp`. This route gives
`execute_query` access. If it does not appear, check `/status` for the active
login method and restart Claude Code. A directly configured server at the
same URL takes precedence, so remove it if you want the connector instead.

Direct OAuth from a server added with `claude mcp add` or `.mcp.json` still
has a callback limitation for domonda. Omitting the `Authorization` header
starts discovery and opens the browser, but the Auth0 authorization request
is rejected before a token is issued:

Claude Code identifies itself with a
[Client ID Metadata Document](../README.md#client-id-metadata-document-cimd)
that declares `http://localhost/callback`, and then listens on a port it picks
per session. domonda's OAuth proxy accepts the port-less form (RFC 8252
section 7.3 requires that of a native client), but the callback it can
register with Auth0 carries no port either, and Auth0 matches callbacks
exactly — so `http://localhost:<port>/callback` is refused at Auth0's
`/authorize`. Use the claude.ai connector above or an API key for the direct
server entry.

Other Claude surfaces are unaffected: **Claude Desktop, Cowork and claude.ai**
take a fixed HTTPS callback and sign in normally — see
[`../claude-desktop/`](../claude-desktop/).

For the API-key path, `execute_query` stays unavailable — it requires an
admin or super-admin **OAuth** user — and the key is bound to one company, so
`X-Selected-Client-Company-ID` has nothing to select. Use a separate server
entry per company instead.

### Base URL: production vs local

The production MCP endpoint is `https://domonda.app/api/mcp/` — note the
`/api` prefix. A non-production deployment may be mounted without it, in which
case the endpoint is `<base-url>/mcp/`; point the server entry at that URL
instead.

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
   tools don't cover. Requires an admin or super-admin **OAuth** user;
   API-key callers are rejected.

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
   an admin or super-admin at the target company. Supports the
   `X-Selected-Client-Company-ID` header for multi-company users. Claude Code
   can receive this through a signed-in claude.ai connector; direct Claude
   Code OAuth is subject to the callback limitation above.

```
Authorization: Bearer <DOMONDA_API_KEY_OR_OAUTH_TOKEN>
```

The server enforces:

- **Read-only transactions for query tools** — no INSERT, UPDATE, DELETE, or DDL
- **Schema restriction** — only the `api` schema is accessible
- **Row limit** — max 1,000 rows per query
- **Query timeout** — 30-second default
- **Query length** — max 10,000 characters
- **Admin requirement** — an OAuth user must be an admin or super-admin at
  the selected client company to use the MCP server at all; accountants are
  refused
- **`execute_query` is further restricted** — it is company-scoped but does
  not apply per-user ACL, so API-key callers are rejected outright and OAuth
  callers need one of the roles above. Its SQL is validated before it reaches
  the database; see [what it rejects](../README.md#what-execute_query-rejects)
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
| Direct OAuth opens the browser, then fails    | Claude Code's loopback callback cannot be registered        | Use a signed-in claude.ai connector or an API key; see [OAuth setup](#oauth-through-a-claudeai-connector) |
| HTTP 403                                      | Forbidden company access or insufficient OAuth scopes       | Check the selected company or required scopes                                         |
| HTTP 404                                      | Wrong endpoint URL                                          | Ensure the URL ends with `/mcp/` or `/mcp`                                            |
| "MCP access is restricted to admin users"         | The signed-in user is not an admin at the selected company  | Sign in as an admin or super-admin user                                           |
| Query returns no rows                         | Company scoping — data belongs to another company           | Confirm `get_my_company` returns the expected company                                 |
| `execute_query` not available to API-key callers | The tool requires an admin or super-admin OAuth user     | Use the signed-in claude.ai connector or specialized tools with your API key           |
| `execute_query` rejected                      | SQL validation failed                                       | See [what `execute_query` rejects](../README.md#what-execute_query-rejects)            |
