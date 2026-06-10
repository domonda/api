# MCP Server (AI assistant access)

The Domonda **Model Context Protocol (MCP)** server lets AI assistants —
Claude, ChatGPT, and any MCP-capable agent — read your company's financial
data and (optionally) upload documents, using natural language. It is a
fourth interface alongside the [GraphQL](../README.md#graphql-api) and
[REST](../README.md#rest-api) APIs and the [Go SDK](../README.md#go-sdk).

The server speaks the [Model Context Protocol](https://modelcontextprotocol.io)
over **Streamable HTTP**. All query tools are **read-only**; the single write
tool is `add_document`, which uploads a file through the same pipeline as the
public REST upload endpoint.

| Environment | Base URL                       |
|-------------|--------------------------------|
| Production  | `https://domonda.app/api/mcp/` |

The trailing slash is optional — both `/api/mcp/` and `/api/mcp` work.

## Table of Contents

- [Quick start](#quick-start)
- [Authentication](#authentication)
  - [OAuth flow](#oauth-flow)
  - [Client ID Metadata Document (CIMD)](#client-id-metadata-document-cimd)
  - [API key flow](#api-key-flow)
- [Client setup guides](#client-setup-guides)
- [Tools](#tools)
- [Resources](#resources)
- [curl examples](#curl-examples)
- [Why access is limited](#why-access-is-limited)

## Quick start

Point any MCP client at `https://domonda.app/api/mcp/` and authenticate with a
Domonda API key (the same Bearer token used by the REST and GraphQL APIs) or by
signing in via OAuth. Then jump to the guide for your client:

- **[Claude Code](./claude-code/)** — Anthropic's CLI / IDE agent (native HTTP MCP)
- **[Claude Desktop, Cowork & claude.ai](./claude-desktop/)** — Custom Connector UI
- **[ChatGPT & OpenAI API](./chatgpt/)** — developer-mode connector + Responses API
- **[OpenClaw](./openclaw/)** — OpenClaw skill / plugin

## Authentication

The MCP server uses **Bearer token** authentication. Two token types are
supported:

1. **API key (HS256 JWT)** — the **same** Domonda API key used by the REST and
   GraphQL APIs (see the top-level
   [Authentication section](../README.md#authentication)). The token's `sub`
   claim identifies the client company directly. This is the simplest option
   for programmatic access.
2. **OAuth** — an access token obtained by signing in with your Domonda
   account. The token is mapped to a Domonda user, who must be enabled and have
   access to the target company. This is what interactive AI clients (Claude
   Desktop, ChatGPT, …) obtain automatically.

```
Authorization: Bearer <DOMONDA_API_KEY_OR_OAUTH_TOKEN>
```

### OAuth flow

Interactive MCP clients discover how to sign in automatically, following the
OAuth 2.0 Protected Resource Metadata flow ([RFC 9728](https://www.rfc-editor.org/rfc/rfc9728)):

1. The client sends a request without an `Authorization` header (or with an
   invalid token).
2. The server responds `401 Unauthorized` with a `WWW-Authenticate` header
   pointing at the protected-resource metadata:
   ```
   WWW-Authenticate: Bearer resource_metadata="https://domonda.app/api/mcp/.well-known/oauth-protected-resource"
   ```
3. The client fetches that metadata, discovers the authorization server, and
   completes a standard OAuth 2.1 authorization-code flow with PKCE.
4. The client retries with `Authorization: Bearer <access token>`.

The `.well-known/oauth-protected-resource` and
`.well-known/oauth-authorization-server` metadata endpoints are public; every
other path requires authentication.

To act on a client company other than your default, set the
`X-Selected-Client-Company-ID` header to the target company UUID. The server
verifies you have access before proceeding.

### Client ID Metadata Document (CIMD)

The server supports the
[Client ID Metadata Document](https://modelcontextprotocol.io/specification/2025-11-25/basic/authorization)
mechanism from the 2025-11-25 MCP authorization spec, used by clients that
identify themselves with an `https://` URL as their `client_id` (e.g. Claude
and ChatGPT). No client registration step is required on your side — the server
advertises CIMD support and registers the client transparently. Clients that do
not use CIMD fall back to OAuth 2.0 Dynamic Client Registration
([RFC 7591](https://www.rfc-editor.org/rfc/rfc7591)) automatically.

### API key flow

For programmatic access, send your Domonda API key as a Bearer token:

```
Authorization: Bearer <HS256-signed-JWT>
```

The JWT `sub` claim identifies the client company directly; no interactive
sign-in is involved. This is the same key used by the REST and GraphQL APIs.

## Client setup guides

Step-by-step instructions, a tool reference (`SKILL.md`), and a
`scripts/health_check.sh` connectivity check live in each client folder:

- [`claude-code/`](./claude-code/) — Claude Code CLI / IDE (native HTTP MCP, no bridge)
- [`claude-desktop/`](./claude-desktop/) — Claude Desktop, Cowork, and claude.ai (Custom Connector UI, with an `mcp-remote` fallback)
- [`chatgpt/`](./chatgpt/) — ChatGPT developer-mode connector and the OpenAI Responses API / Agents SDK
- [`openclaw/`](./openclaw/) — OpenClaw skill and plugin

> **Tool descriptions are in German.** The tool and parameter descriptions the
> server surfaces at runtime are written in German with Domonda domain
> vocabulary (Sachkonto, Rechnung, Gutschrift, Geschäftspartner, …), because
> end users chat in German. The English text in these guides is reference only.
> When unsure which tool fits, call `list_capabilities` for a topic-grouped
> catalog.

## Tools

### Discovery
- **`list_capabilities`** — Topic-grouped catalog of all available tools

### Schema introspection
- **`list_tables`** — List tables/views in the `api` schema
- **`describe_table`** — Describe columns of an `api` schema table/view

### SQL query
- **`execute_query`** — Execute a read-only SQL query (api schema only)

### Documents
- **`list_documents`** — List documents with filtering (archived, category, fulltext, dates, amounts, status, …)
- **`get_document`** — Get a single document by ID (excludes `fulltext`)
- **`get_document_fulltext`** — Get the extracted fulltext (OCR) of a document. **Can be very large.**
- **`search_documents`** — Fulltext search across documents
- **`download_document_pdf`** — Download a document PDF (optional audit trail)
- **`get_document_thumbnail`** — Default JPEG thumbnail (256-wide)
- **`get_document_preview`** — Preview image (964×1364 JPEG)
- **`add_document`** — *(write)* Upload a document into the authenticated company

### Invoices
- **`list_invoices`** — List invoices with filtering (date range, partner, amounts)
- **`get_invoice`** — Invoice details with line items and matched payments

### Payments
- **`list_money_accounts`** — List money accounts (bank, credit card, PayPal, …)
- **`list_money_transactions`** — List money transactions with filtering
- **`get_invoice_payments`** — Payment matches for a specific invoice

### Company
- **`get_my_company`** — Authenticated company joined with master data and headquarters
- **`list_document_categories`** — List document categories
- **`list_gl_accounts`** — List general ledger accounts
- **`list_partner_companies`** — List/search partner companies

### User
- **`get_user`** — The authenticated user (OAuth only; API-key auth returns an error)

### Workflows
- **`list_document_workflows`** — Document workflows defined for the company
- **`list_document_workflow_steps`** — Workflow steps (optionally filtered by `workflow_id`)
- **`list_document_workflow_states`** — Current workflow step and state per document

### App URLs
- **`get_app_urls`** — Web-app deep links scoped to the authenticated company
- **`get_document_urls`** — Web-app deep links for a single document
- **`get_document_selection_url`** — Web-app deep links opening a selection of documents

See each client folder's `SKILL.md` for the full parameter reference.

## Resources

The server advertises the `resources` capability via RFC 6570 URI templates.
Each resource read runs through the same authentication and company-scoped
ownership check as the tools.

- **`document://{document_id}/thumb_256x0.jpeg`** — Default JPEG thumbnail (256-wide)
- **`document://{document_id}/preview.jpeg`** — Preview image (964×1364 JPEG)
- **`document://{document_id}/pdf{?audit_trail,audit_trail_lang}`** — Document PDF.
  Optional query parameters: `audit_trail` (`append`, `prepend`, `only`) and
  `audit_trail_lang` (ISO 639-1, default `de`).

## curl examples

The MCP server uses Streamable HTTP: all calls are `POST /api/mcp/` with a
JSON-RPC body. Set your API key first:

```bash
 export DOMONDA_API_KEY="eyJhbGciOiJIUzI1NiIs..."   # leading space keeps this out of shell history
```

Discover the OAuth protected-resource metadata (public, no auth):

```bash
curl -s "https://domonda.app/api/mcp/.well-known/oauth-protected-resource" | jq .
```

Initialize a session / verify a token:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"curl","version":"1.0"}}}' | jq .
```

List the available tools:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | jq .
```

List tables in the `api` schema:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_tables","arguments":{}}}' | jq .
```

Run a read-only SQL query:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"execute_query","arguments":{"sql":"SELECT id, title, document_date FROM api.document ORDER BY document_date DESC LIMIT 5"}}}' | jq .
```

List recent invoices:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_invoices","arguments":{"from_date":"2025-01-01","until_date":"2025-12-31","limit":10}}}' | jq .
```

Read a document PDF as an MCP resource:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"resources/read","params":{"uri":"document://DOCUMENT_UUID/pdf"}}' | jq .
```

## Why access is limited

The MCP server is a deliberately narrow, safe window onto your data. The limits
below exist to protect your data and keep responses fast and predictable — not
to make the API harder to use.

- **Read-only, so an assistant can explore freely.** Every tool except
  `add_document` can only read. You can point an AI assistant at your data and
  ask anything without worrying it will change, move, or delete a record. The
  one write tool, `add_document`, goes through the same validation and
  duplicate checks as a normal upload.
- **Only ever your company's data.** Every request sees just the data of the
  company its token belongs to. There is no way to phrase a question that
  reaches another customer's data.
- **A curated view, not the raw database.** Queries run against the documented
  `api` schema — the same one behind the GraphQL and REST APIs — rather than
  internal tables. Your queries keep working across product updates, and you
  only ever see fields meant for you.
- **Sized to stay fast.** A single request returns at most **1,000 rows**, runs
  for at most **30 seconds**, accepts up to **10,000 characters** of SQL, and
  returns up to **10 MiB**. These bounds stop one broad question from slowing
  things down for you or anyone else — narrow with filters (date range, partner,
  status) and paginate when you need more.
