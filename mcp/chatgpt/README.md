# domonda MCP — ChatGPT & OpenAI API

This directory contains setup instructions and configuration for connecting
**ChatGPT** (developer-mode custom MCP connectors, sometimes surfaced as
"apps") and the **OpenAI Responses API / Agents SDK** to the
[domonda](https://domonda.app) financial data platform via its
[Model Context Protocol (MCP)](https://modelcontextprotocol.io) server.

domonda's MCP server already implements everything ChatGPT requires of a remote
MCP server — OAuth 2.1 with PKCE, RFC 9728 protected-resource metadata, RFC 8414
authorization-server metadata, Client ID Metadata Documents (CIMD), and Dynamic
Client Registration (DCR) — so **no server-side changes are needed** to connect
ChatGPT. See the [main README](../README.md#authentication) for the protocol
details.

All query operations are **read-only**. The server exposes 28 tools covering
documents, invoices, payments, money transactions, companies, GL accounts,
partner companies, document workflows, web app deep links, ad-hoc SQL
queries, and a `list_capabilities` discovery tool. The single write tool is
`add_document`, which uploads a document via the same pipeline as the public
upload endpoint.

> **Live tool descriptions are in German.** The tool and parameter
> descriptions surfaced by the running MCP server (and therefore what the
> model actually matches against user prompts) are written in German with
> domain vocabulary pulled from the domonda-app UI (Sachkonto, Rechnung,
> Gutschrift, Geschäftspartner, Fälligkeitsdatum, Mandant, Freigabe, …) and
> a "Nutzen wenn …" trigger-phrase block for each tool. End users chat in
> German, so the descriptions are optimised for the words they actually
> use. When uncertain which tool to call, call `list_capabilities` first for
> a topic-grouped catalog.

## Prerequisites

- **For the ChatGPT app connector (Option 1):** a paid ChatGPT plan (Plus,
  Pro, Business, Enterprise, or Edu) with **developer mode / custom MCP
  connectors enabled**. On Business/Enterprise/Edu workspaces a workspace
  owner or admin must enable it first. Authentication is **OAuth only** — the
  ChatGPT UI has no field for pasting a static API key. See the
  [OpenAI help article](https://help.openai.com/en/articles/12584461-developer-mode-apps-and-full-mcp-connectors-in-chatgpt-beta)
  for the current plan-by-plan availability, as the UI changes frequently.
- **For the OpenAI Responses API / Agents SDK (Options 2 & 3):** an OpenAI API
  key, plus either a domonda API key (JWT — obtain from the domonda admin
  panel) or an Auth0 OAuth access token for the MCP resource.
- `curl` (for the health check script).

## Directory Structure

```
chatgpt/
├── SKILL.md                          # Skill / tool reference (frontmatter + tool docs)
├── openai_responses_mcp.json.example # Example `mcp` tool object for the Responses API
├── scripts/
│   └── health_check.sh               # Connectivity check script
└── README.md                         # This file
```

> The ChatGPT app connector (Option 1) is configured entirely in the ChatGPT
> UI, so it has no config file. The `.example` file is for the API path only.

## Setup — pick one of three paths

### Option 1 — ChatGPT developer-mode connector (recommended, OAuth)

The connector originates from OpenAI's cloud and dials the MCP endpoint on your
behalf — so it needs a **publicly reachable** URL and drives an OAuth sign-in.

1. **Enable developer mode.** On Plus/Pro: **Settings → Connectors → Advanced →
   Developer mode**. On Business/Enterprise/Edu: a workspace owner/admin enables
   it under workspace settings (Permissions & Roles → Connected data → custom
   MCP connectors), then individual members enable it under their own
   **Settings → Connectors → Advanced**.
2. Open **Settings → Connectors → Add custom connector** (labelled **Create** in
   some builds).
3. Enter the MCP URL: `https://domonda.app/api/mcp/`
4. Complete the authentication prompt. domonda serves the Auth0 OAuth discovery
   metadata (`.well-known/oauth-protected-resource` per RFC 9728), so ChatGPT
   redirects you to Auth0 to sign in and registers itself via CIMD (preferred)
   or DCR automatically — no client ID or secret to enter.
5. Start a new chat and confirm the domonda tools appear in the connector / tool
   picker.

> **Doesn't work for localhost.** Because the connector dials from OpenAI's
> cloud, `http://localhost:5001/mcp/` is unreachable. Use Option 2 or 3 against
> your local server, or expose it with a tunnel (e.g. ngrok) over HTTPS.
>
> **OAuth only.** The UI has no raw-Bearer-JWT field. To use a static domonda
> API key, use the API paths below.

### Option 2 — OpenAI Responses API (programmatic, API key or OAuth)

The Responses API connects to remote MCP servers over Streamable HTTP. Add an
`mcp` tool pointing at the domonda endpoint and pass the bearer token in
`headers`. See `openai_responses_mcp.json.example` for the ready-to-paste tool
object.

```python
from openai import OpenAI
import os

client = OpenAI()  # reads OPENAI_API_KEY

response = client.responses.create(
    model="gpt-5",
    tools=[
        {
            "type": "mcp",
            "server_label": "domonda",
            "server_url": "https://domonda.app/api/mcp/",
            "headers": {"Authorization": f"Bearer {os.environ['DOMONDA_API_KEY']}"},
            "require_approval": "never",
        },
    ],
    input="Show me the 5 most recent invoices over €1,000.",
)
print(response.output_text)
```

> OpenAI discards the header values after each request, so the
> `Authorization` header must be included with **every** API call. A domonda
> API key (HS256 JWT) works directly here; an Auth0 OAuth access token works
> too.

### Option 3 — OpenAI Agents SDK (programmatic)

```python
import os
from agents import Agent, Runner
from agents.mcp import MCPServerStreamableHttp

async def main():
    async with MCPServerStreamableHttp(
        name="domonda",
        params={
            "url": "https://domonda.app/api/mcp/",
            "headers": {"Authorization": f"Bearer {os.environ['DOMONDA_API_KEY']}"},
        },
        cache_tools_list=True,
    ) as server:
        agent = Agent(
            name="Assistant",
            instructions="Use the domonda tools to answer.",
            mcp_servers=[server],
        )
        result = await Runner.run(agent, "Which GL accounts are in use?")
        print(result.final_output)
```

## How ChatGPT registers as an OAuth client

For the Option 1 connector, ChatGPT needs to register itself as an OAuth client
with domonda's authorization server. domonda advertises everything ChatGPT
checks for, so registration is automatic:

- **CIMD (preferred).** ChatGPT uses an `https://` URL as its `client_id`;
  domonda fetches that document and registers the client with Auth0 via DCR on
  ChatGPT's behalf. Selected automatically because domonda's
  authorization-server metadata advertises `"client_id_metadata_document_supported": true`
  and `"none"` in `token_endpoint_auth_methods_supported`.
- **DCR (fallback, RFC 7591).** ChatGPT calls domonda's `/register` endpoint
  (proxied to Auth0) once per connector instance.
- **PKCE (S256)** is mandatory and advertised via
  `"code_challenge_methods_supported": ["S256"]`.

No client ID or secret is entered in the ChatGPT UI. See the
[main README CIMD / DCR sections](../README.md#client-id-metadata-document-cimd)
for the full flow and the Auth0 setup it depends on.

### Base URL: production vs local

In production the domonda-web-server is mounted behind an `/api` path prefix,
so the MCP endpoint is `https://domonda.app/api/mcp/`. When running the web
server locally there is **no** `/api` prefix — use `http://localhost:5001/mcp/`
(5001 is the default port; override with `PORT=…`). The ChatGPT app connector
cannot reach localhost (see Option 1); the API paths can.

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

### Inside ChatGPT

Start a new conversation with the connector enabled and ask:

> "Use the domonda tools to get company info."

ChatGPT should call `get_my_company` and return your company details.

## Usage

Once connected, the model can call any of the 28 domonda tools. See `SKILL.md`
for the complete tool reference with parameters.

### Recommended first steps

1. **`list_capabilities`** — When unsure which tool fits, get a grouped
   catalog of all available tools first.
2. **`get_my_company`** — Confirm authentication and see which company the
   token is scoped to.
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

- "Show me all invoices from January 2026 over €1,000"
- "Download the PDF for document 8a2f..."
- "Which GL accounts are in use?"
- "Search documents for 'electricity bill'"
- "Run a query to find the top 10 partners by invoice count"
- "Upload this invoice PDF into the OTHER_DOCUMENTS category"

## Authentication

The domonda MCP server uses **Streamable HTTP** transport with **Bearer token**
authentication. Two token types are supported:

1. **API key (HS256 JWT)** — the standard domonda API key. The JWT `sub` claim
   identifies the client company directly. Use it in the Responses API / Agents
   SDK `headers` (Options 2 & 3). The ChatGPT app connector UI does **not**
   accept it.
2. **OAuth (Auth0 RS256)** — an Auth0 access token issued for the MCP resource.
   The token's Auth0 `sub` claim is mapped to a Domonda user, who must be
   enabled and have access to the target company. This is what the ChatGPT app
   connector (Option 1) obtains automatically. Supports the
   `X-Selected-Client-Company-ID` header for multi-company users.

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

| Symptom                                          | Cause                                                        | Fix                                                                                  |
|--------------------------------------------------|--------------------------------------------------------------|--------------------------------------------------------------------------------------|
| No "Add custom connector" / "Create" option      | Developer mode not enabled, or plan/workspace doesn't allow it | Enable developer mode (Settings → Connectors → Advanced); on workspaces ask an admin |
| Connector can't reach `localhost:…`              | Connector dials from OpenAI's cloud, not your machine        | Use Option 2/3, or expose the local server over HTTPS via a tunnel                    |
| No place to paste an API key in ChatGPT          | The app connector UI is OAuth-only                           | Use the Responses API / Agents SDK (Options 2 & 3) for a static Bearer JWT            |
| OAuth sign-in loops or fails                     | Auth0 DCR/CIMD prerequisites missing on the tenant           | See the main README "DCR and CIMD" Auth0 setup steps                                  |
| HTTP 401                                         | Invalid, expired, or blocked token                           | Verify your API key or OAuth token; contact admin if blocked                         |
| HTTP 403                                         | Forbidden company access or insufficient OAuth scopes        | Check the selected company or required scopes                                        |
| HTTP 404                                         | Wrong endpoint URL                                           | Ensure the URL ends with `/mcp/` or `/mcp` (production: `/api/mcp/`)                  |
| Query returns no rows                            | Company scoping — data belongs to another company           | Confirm `get_my_company` returns the expected company                                |
| `execute_query` rejected                         | SQL validation failed (DML, wrong schema)                   | Use only SELECT/WITH against the `api` schema                                        |
