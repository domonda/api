# domonda MCP — OpenClaw Skill & MCP Connection

This directory contains an [OpenClaw](https://openclaw.ai) skill that connects
AI agents to the [domonda](https://domonda.app) financial data platform via its
[Model Context Protocol (MCP)](https://modelcontextprotocol.io) server.

Query tools are **read-only**. The skill exposes 28 tools covering documents,
invoices, payments, money transactions, companies, GL accounts, partner
companies, document workflows, web app deep links, ad-hoc SQL queries, and a
`list_capabilities` discovery tool. The single write tool is `add_document`,
which uploads a file into the authenticated company's document pipeline.

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

- OpenClaw installed with an agent configured
- A domonda API key (JWT), or a Domonda admin/super-admin account for OAuth
- `curl` (for the health check script)

## Directory Structure

```
openclaw/
├── SKILL.md                  # Skill definition (frontmatter + tool reference)
├── .env.example              # Environment variable template
├── openclaw.json.example     # OpenClaw MCP server configuration template
├── scripts/
│   └── health_check.sh       # Connectivity check script
└── README.md                 # This file
```

## Installation

### 1. Register the skill in OpenClaw

Copy `SKILL.md` into `<workspace>/skills/domonda/SKILL.md` or
`~/.openclaw/skills/domonda/SKILL.md`. OpenClaw discovers skills from those
directories; the `domonda` folder must match the skill name in its frontmatter.
The skill supplies instructions, while the MCP server entry below supplies
the live tools.

### 2. Set up environment variables

For an API-key connection, add the key to `~/.openclaw/.env` (or the Gateway
process environment). Use `.env.example` as a template, but do not put the
credential in the skill folder:

```bash
mkdir -p ~/.openclaw
chmod 700 ~/.openclaw
```

Edit `~/.openclaw/.env`:

```
DOMONDA_API_KEY=eyJhbGciOiJIUzI1NiIs...    # your JWT API key
```

| Variable          | Required     | Description                                           |
|-------------------|--------------|-------------------------------------------------------|
| `DOMONDA_API_KEY` | API key path | JWT API key for Bearer token authentication           |
| `DOMONDA_MCP_URL` | no           | Full URL for the health check (default: production)   |

> **Security:** Keep `~/.openclaw/.env` private (`chmod 600 ~/.openclaw/.env`).
> The API key is a secret and must not be committed, logged, displayed, or
> passed as a CLI argument. Skip the key entirely if you use OAuth below.

### 3. Configure the MCP server connection

Merge the `domonda` entry from `openclaw.json.example` into
`~/.openclaw/openclaw.json`. OpenClaw loads that config path by default; a
copy in this skill directory is not loaded. The production entry is:

```json
{
  "mcp": {
    "servers": {
      "domonda": {
        "url": "https://domonda.app/api/mcp/",
        "transport": "streamable-http",
        "headers": {
          "Authorization": "Bearer ${DOMONDA_API_KEY}"
        }
      }
    }
  }
}
```

OpenClaw substitutes `${DOMONDA_API_KEY}` from its process environment or
global `.env` at runtime. `transport` is required here: omitting it selects
SSE. Do not use the old top-level `mcpServers` key.

For local development, replace the `url` value with the full local MCP URL,
for example `http://localhost:<port>/mcp/`. Set `DOMONDA_MCP_URL` to that full
URL in `~/.openclaw/.env` if you also want to use the health check script;
OpenClaw's config uses the literal production URL above unless you edit it.

For OAuth, use this server entry instead of the API-key `headers` entry:

```json
{
  "mcp": {
    "servers": {
      "domonda": {
        "url": "https://domonda.app/api/mcp/",
        "transport": "streamable-http",
        "auth": "oauth"
      }
    }
  }
}
```

Then run `openclaw mcp login domonda` and sign in as a Domonda admin or
super-admin. OAuth also enables `execute_query`; API-key access does not.

### 4. Verify connectivity

Run `openclaw mcp doctor domonda --probe` to test the saved MCP connection.
For an API-key connection, you can also run the included health check:

```bash
./scripts/health_check.sh
```

Expected output:

```
Checking https://domonda.app/api/mcp/ ...
OK (HTTP 200)
```

The script reads `~/.openclaw/.env` if the key is not already exported, sends
an MCP `initialize` request, and reports the HTTP status code. It checks an
API key independently of OpenClaw; use `openclaw mcp doctor domonda --probe`
for OAuth.

## Usage

Once installed, your OpenClaw agent can call any of the 28 domonda tools.
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

| Category  | Tools                                                                                                                                                    |
|-----------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| Discovery | `list_capabilities`                                                                                                                                      |
| Schema    | `list_tables`, `describe_table`                                                                                                                          |
| Query     | `execute_query`                                                                                                                                          |
| Documents | `list_documents`, `get_document`, `get_document_fulltext`, `search_documents`, `download_document_pdf`, `get_document_thumbnail`, `get_document_preview`, `add_document` (write) |
| Invoices  | `list_invoices`, `get_invoice`                                                                                                                           |
| Payments  | `list_money_accounts`, `list_money_transactions`, `get_invoice_payments`                                                                                 |
| Company   | `get_my_company`, `list_document_categories`, `list_gl_accounts`, `list_partner_companies`                                                               |
| User      | `get_user`                                                                                                                                               |
| Workflows | `list_document_workflows`, `list_document_workflow_steps`, `list_document_workflow_states`                                                               |
| App URLs  | `get_app_urls`, `get_document_urls`, `get_document_selection_url`                                                                                        |

### Example prompts for your agent

- "Show me all invoices from January 2026 over €1,000"
- "Download the PDF for document 8a2f..."
- "Which GL accounts are in use?"
- "Search documents for 'electricity bill'"
- "Run a query to find the top 10 partners by invoice count"

## Authentication

The domonda MCP server uses **Streamable HTTP** transport with **Bearer token**
authentication. Two token types are supported:

1. **API key (HS256 JWT)** — the standard domonda API key. The JWT `sub` claim
   identifies the client company directly. This is the simplest option for
   programmatic access.
2. **OAuth (Auth0 RS256)** — an Auth0 access token issued for the MCP resource.
   The token's Auth0 `sub` claim is mapped to a Domonda user, who must be
   an admin or super-admin at the target company. Supports the
   `X-Selected-Client-Company-ID` header for multi-company users. Clients that
   support OAuth discovery will be guided automatically via the
   `WWW-Authenticate` header and the `.well-known/oauth-protected-resource`
   metadata endpoint (RFC 9728).

```
Authorization: Bearer <DOMONDA_API_KEY_OR_OAUTH_TOKEN>
```

The server enforces:

- **Read-only transactions for query tools** — no INSERT, UPDATE, DELETE, or DDL through `execute_query` or any specialized query tool
- **Schema restriction** — only the `api` schema is accessible
- **Row limit** — max 1,000 rows per query
- **Query timeout** — 30-second default
- **Query length** — max 10,000 characters
- **Admin requirement** — an OAuth user must be an admin or super-admin at the selected client company to use the MCP server at all; accountants are refused
- **`execute_query` is further restricted** — it is company-scoped but does not apply per-user ACL, so API-key callers are rejected outright and OAuth callers need one of the roles above. Its SQL is validated before it reaches the database; see [what it rejects](../README.md#what-execute_query-rejects)
- **Writes are limited to `add_document`** — processed through the same OCR, duplicate-detection, and category-ownership checks as the public REST upload endpoint

## Troubleshooting

| Symptom                         | Cause                                      | Fix                                                  |
|---------------------------------|--------------------------------------------|------------------------------------------------------|
| `Error: DOMONDA_API_KEY is not set` | Missing API key for the health check   | Set `DOMONDA_API_KEY` in `~/.openclaw/.env` or your environment |
| HTTP 401                        | Invalid, expired, or blocked token         | Verify your API key or OAuth token; contact admin if blocked |
| HTTP 403                        | Forbidden company access or insufficient OAuth scopes | Check the selected company or required scopes        |
| HTTP 404                        | Wrong endpoint URL                         | Ensure the URL ends with `/api/mcp/` or `/api/mcp`             |
| "MCP access is restricted to admin users" | The signed-in user is not an admin at the selected company | Sign in as an admin or super-admin user |
| Query returns no rows           | Company scoping — data belongs to another company | Confirm `get_my_company` returns the expected company |
| `execute_query` not available to API-key callers | The tool requires an admin or super-admin OAuth user | Use an OAuth token, or use the specialized tools with your API key |
| `execute_query` rejected        | SQL validation failed                      | See [what `execute_query` rejects](../README.md#what-execute_query-rejects) |
