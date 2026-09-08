# domonda MCP — Claude Desktop & Cowork

This directory contains setup instructions and configuration for connecting
**Claude Desktop**, [**Claude Cowork**](https://claude.com/product/cowork),
and the **claude.ai web app** to the [domonda](https://domonda.app) financial
data platform via its
[Model Context Protocol (MCP)](https://modelcontextprotocol.io) server. All
three surfaces share the same Custom Connector UI.

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

- **Claude Desktop** app ([download](https://claude.ai/download)),
  [**Cowork**](https://claude.com/product/cowork), or a paid plan on
  [claude.ai](https://claude.ai/) (Free plan is limited to one custom
  connector)
- A domonda API key (JWT) — obtain from the domonda admin panel — or an
  Auth0 OAuth access token for the MCP resource
- **For the stdio-bridge fallback only:** Node.js and `npx` (see the
  [What is `npx`?](#what-is-npx) section below for how to install them)

## Directory Structure

```
claude-desktop/
├── SKILL.md                              # Skill definition (frontmatter + tool reference)
├── claude_desktop_config.json.example    # Example Claude Desktop configuration
├── scripts/
│   └── health_check.sh                  # Connectivity check script
└── README.md                            # This file
```

## Claude Desktop Setup

Claude Desktop supports two paths. Use the **Custom Connector UI** for
production; use the **`mcp-remote` stdio bridge** only when the UI path can't
reach your endpoint (localhost) or when you need to inject a raw Bearer JWT
without an OAuth flow.

### Option 1 — Custom Connector UI (recommended)

Available on Free, Pro, Max, Team, and Enterprise plans (Free is limited to
one custom connector). Anthropic's cloud dials the MCP endpoint on your
behalf — your machine does not need Node.js, and there is no config file to
edit.

1. Open **Settings → Connectors** (inside Claude Desktop or at
   [claude.ai](https://claude.ai/) → Settings → Connectors).
2. Scroll to **Add custom connector**.
3. Enter the MCP URL: `https://domonda.app/api/mcp/`
4. Complete the authentication prompt. domonda's MCP endpoint serves the
   OAuth discovery metadata (`.well-known/oauth-protected-resource` per
   RFC 9728) and proxies the `/authorize`, `/token` and `/register` endpoints
   to Auth0, so the UI will redirect you to Auth0 to sign in.
5. Confirm the tools appear in the tool picker (hammer icon) in a new chat.

> **Doesn't work for localhost.** The connector originates from Anthropic's
> cloud, not your machine, so `http://localhost:<port>/mcp/` is unreachable.
> Use Option 2 for local development.
>
> **Plain API-key JWTs aren't supported here.** The UI drives an OAuth flow;
> use Option 2 if you want to paste a raw `Authorization: Bearer <jwt>`
> header instead.

### Option 2 — `mcp-remote` stdio bridge (fallback)

Use this when the Custom Connector UI doesn't fit: **localhost dev endpoints**
or **raw Bearer JWT** authentication without OAuth. Claude Desktop talks
stdio to the local `mcp-remote` process, which in turn speaks Streamable
HTTP to the domonda server.

See [mcp-remote](https://www.npmjs.com/package/mcp-remote) for the upstream
package.

#### What is `npx`?

`npx` is a small command that ships with [Node.js](https://nodejs.org/) — the
JavaScript runtime that powers tools like VS Code, Slack's desktop app, and
countless developer utilities. When Claude Desktop's config says
`"command": "npx"`, it is asking your computer to run a JavaScript program
(in this case `mcp-remote`, the little bridge that lets Claude Desktop reach
a web-based MCP server). `npx` downloads that program on demand the first
time you use it and caches it for later runs — no separate installation
step, no app in your Applications folder.

To install Node.js (which also gives you `npx`):

- **macOS / Linux:** download the LTS installer from
  [nodejs.org](https://nodejs.org/) and run it, or use a package manager
  like [Homebrew](https://brew.sh/) (`brew install node`) or `apt install
  nodejs npm`.
- **Windows:** download the LTS installer from
  [nodejs.org](https://nodejs.org/) and run it. The installer adds `node`
  and `npx` to your `PATH` automatically, so Claude Desktop can find them.

Verify after install: open a new terminal and run `npx --version`. You
should see something like `10.x.x`. If you see "command not found," restart
your terminal (or sign out and back in on Windows) so the updated `PATH`
takes effect.

> If installing Node.js on the end-user's machine isn't realistic, prefer
> Option 1 (Custom Connector UI) — it requires no local runtime at all.

#### 1. Locate the configuration file

| Platform | Path                                                                 |
|----------|----------------------------------------------------------------------|
| macOS    | `~/Library/Application Support/Claude/claude_desktop_config.json`    |
| Windows  | `%APPDATA%\Claude\claude_desktop_config.json`                        |
| Linux    | `~/.config/Claude/claude_desktop_config.json`                        |

Create the file if it does not exist.

#### 2. Add the domonda MCP server

Edit `claude_desktop_config.json` and add (or merge) the `domonda` entry under
`mcpServers`. See `claude_desktop_config.json.example` for a ready-to-copy
template.

```json
{
  "mcpServers": {
    "domonda": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "https://domonda.app/api/mcp/",
        "--header",
        "Authorization: Bearer YOUR_API_KEY_HERE"
      ]
    }
  }
}
```

Replace `YOUR_API_KEY_HERE` with your actual domonda JWT API key.

> **Base URL.** The production MCP endpoint is
> `https://domonda.app/api/mcp/` — note the `/api` prefix. A non-production
> deployment may be mounted without it, in which case the endpoint is
> `<base-url>/mcp/`.

For local development, replace the URL with your local endpoint (note the
missing `/api` segment):

```json
{
  "mcpServers": {
    "domonda": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:<port>/mcp/",
        "--header",
        "Authorization: Bearer YOUR_API_KEY_HERE"
      ]
    }
  }
}
```

#### 3. Restart Claude Desktop

Quit and reopen Claude Desktop for the configuration to take effect.
The domonda tools should appear in the tool picker (hammer icon).

## Cowork and claude.ai

[Claude Cowork](https://claude.com/product/cowork) and the
[claude.ai](https://claude.ai/) web app use the same Custom Connector flow as
Option 1. Open Settings → Connectors (direct link:
[claude.ai/settings/connectors](https://claude.ai/settings/connectors)), click
**Add custom connector**, paste `https://domonda.app/api/mcp/`, and complete
the Auth0 OAuth prompt. Any connector you add at claude.ai is also available
inside Claude Desktop and Cowork on the same account.

> See the [Claude custom connectors help article](https://support.claude.com/en/articles/11175166-get-started-with-custom-connectors-using-remote-mcp)
> for the latest plan-by-plan details, as the UI may change.

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

### Inside Claude Desktop

Start a new conversation and ask:

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
   tools don't cover. Requires an admin, super-admin, or accountant **OAuth**
   user; API-key callers are rejected.

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

### OAuth client registration (CIMD, DCR)

When you add domonda through the Custom Connector UI, Claude has to register
itself as an OAuth client with domonda's authorization server before the
sign-in flow can run. This happens automatically — there is no client ID or
secret for you to enter, and nothing to configure on your side:

1. **Client ID Metadata Document (CIMD)** — *preferred, automatic.* Claude
   uses an `https://` URL as its `client_id`, and domonda's server fetches
   that document and registers the client on Claude's behalf. Claude selects
   CIMD automatically because domonda's authorization-server metadata
   advertises **both** `"client_id_metadata_document_supported": true` and
   `"none"` in `token_endpoint_auth_methods_supported` (required for a public
   client). See the [main README CIMD section](../README.md#client-id-metadata-document-cimd).
2. **Dynamic Client Registration (DCR, RFC 7591)** — *automatic fallback.*
   Claude calls domonda's `/register` endpoint on each new connection. Used
   automatically if the CIMD metadata above is ever missing.

The server enforces:

- **Read-only transactions for query tools** — no INSERT, UPDATE, DELETE, or DDL
- **Schema restriction** — only the `api` schema is accessible
- **Row limit** — max 1,000 rows per query
- **Query timeout** — 30-second default
- **Query length** — max 10,000 characters
- **Admin/accountant requirement** — an OAuth user must be an admin,
  super-admin, or accountant at the selected client company to use the MCP
  server at all
- **`execute_query` is further restricted** — it is company-scoped but does
  not apply per-user ACL, so API-key callers are rejected outright and OAuth
  callers need one of the roles above. Its SQL is validated before it reaches
  the database; see [what it rejects](../README.md#what-execute_query-rejects)
- **Writes are limited to `add_document`** — processed through the same OCR,
  duplicate-detection, and category-ownership checks as the public API

## Troubleshooting

| Symptom                                             | Cause                                                       | Fix                                                                                       |
|-----------------------------------------------------|-------------------------------------------------------------|-------------------------------------------------------------------------------------------|
| Custom Connector can't reach `localhost:…`          | Connector dials from Anthropic's cloud, not your machine    | Switch to Option 2 (`mcp-remote` stdio bridge)                                            |
| OAuth redirect fails / no OAuth prompt              | Using a local / non-public URL, or plan doesn't support it  | Use Option 2 with a raw JWT in the `Authorization` header                                 |
| Tools don't appear (Option 2) in Claude Desktop     | Config file not found or malformed JSON                     | Verify the path and validate JSON syntax                                                  |
| `npx` not found (Option 2)                          | Node.js not installed, or terminal opened before install    | See [What is `npx`?](#what-is-npx); install Node.js and restart the terminal              |
| `mcp-remote` connection error (Option 2)            | Network or firewall issue                                   | Check internet connectivity; run `./scripts/health_check.sh`                              |
| `curl https://domonda.app/api/mcp/` returns 401     | Not a fault — an unauthenticated GET is challenged          | That is the healthy response; add `-H 'Accept: text/html'` for the info page              |
| HTTP 401 from a client                              | Invalid, expired, or blocked token                          | Verify your API key or OAuth token; contact admin if blocked                              |
| HTTP 403                                            | Forbidden company access or insufficient OAuth scopes       | Check the selected company or required scopes                                             |
| HTTP 404                                            | Wrong endpoint URL                                          | Ensure the URL ends with `/mcp/` or `/mcp`                                                |
| "MCP access is restricted to admin or accountant"   | The signed-in user has neither role at the selected company | Sign in as an admin, super-admin, or accountant user                                      |
| Query returns no rows                               | Company scoping — data belongs to another company           | Confirm `get_my_company` returns the expected company                                     |
| `execute_query` not available to API-key callers    | The tool requires an admin/accountant OAuth user            | Use Option 1 (OAuth), or use the specialized tools with your API key                      |
| `execute_query` rejected                            | SQL validation failed                                       | See [what `execute_query` rejects](../README.md#what-execute_query-rejects)               |
