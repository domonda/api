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
- [Protocol versions](#protocol-versions)
- [Authentication](#authentication)
  - [OAuth flow](#oauth-flow)
    - [Public endpoints](#public-endpoints)
    - [`GET` on the base path](#get-on-the-base-path)
    - [Token endpoint (`/token`)](#token-endpoint-token)
  - [Client ID Metadata Document (CIMD)](#client-id-metadata-document-cimd)
    - [Dynamic Client Registration (`/register`)](#dynamic-client-registration-register)
  - [API key flow](#api-key-flow)
- [Client setup guides](#client-setup-guides)
- [Tools](#tools)
  - [What `execute_query` rejects](#what-execute_query-rejects)
  - [Why `execute_query`: two measured examples](#why-execute_query-two-measured-examples)
- [Resources](#resources)
- [curl examples](#curl-examples)
- [Why access is limited](#why-access-is-limited)

## Quick start

Point any MCP client at `https://domonda.app/api/mcp/` and authenticate with a
Domonda API key (the same Bearer token used by the REST and GraphQL APIs) or by
signing in via OAuth. Then jump to the guide for your client:

- **[Claude Code](./claude-code/)** — Anthropic's CLI / IDE agent (native HTTP MCP)
- **[Claude Desktop, Cowork & claude.ai](./claude-desktop/)** — Custom Connector UI
- **[ChatGPT & OpenAI API](./chatgpt/)** — custom MCP plugin + Responses API
- **[OpenClaw](./openclaw/)** — OpenClaw skill and MCP connection

## Protocol versions

The same endpoint serves the current stateless MCP revision and the older
revisions with an `initialize` handshake. The server picks the revision per
request, so clients of either kind connect without any configuration:

| Revision | How a client talks to the server |
|----------|----------------------------------|
| **2026-07-28** (stateless) | No handshake and no session. The client calls `server/discover` to learn the supported versions and capabilities, then sends every request on its own, with the protocol version, client info and client capabilities in the request's `_meta`. |
| **2025-11-25, 2025-06-18, 2025-03-26, 2024-11-05** | The client sends `initialize`, receives an `Mcp-Session-Id` header and passes it on later requests. `GET` on the endpoint opens the server-to-client stream. |
| Plain JSON-RPC | A request with neither `_meta` nor a session ID (the [curl examples](#curl-examples)) is accepted as well, so scripts using an API key can call a tool in a single request. |

Over HTTP, a 2026-07-28 request repeats its protocol version in the
`Mcp-Protocol-Version` header and its JSON-RPC method in the `Mcp-Method`
header; a `tools/call` also names the tool, and a `resources/read` the
resource URI, in `Mcp-Name`. A request whose headers are missing or do not
match the body gets `400`. MCP client libraries set these headers themselves;
see the [curl example](#curl-examples) for a hand-written request. A plain
JSON-RPC call must not send `Mcp-Protocol-Version: 2026-07-28`: the header
alone marks a request as 2026-07-28, and without the matching `_meta` it gets
`400`.

For 2026-07-28 requests, errors also set the HTTP status: an unknown method
gets `404` and an invalid request or invalid parameters get `400`. Older
clients receive the same errors with status `200`, as JSON-RPC over HTTP
always did.

Authentication is the same for every revision.

**Cache hints.** For 2026-07-28 clients, the results of `server/discover`,
`tools/list`, `resources/list` and `resources/templates/list` carry
`"ttlMs": 300000` and `"cacheScope": "private"`: a client may reuse them for
5 minutes instead of asking again before every tool call. The lists only change
when Domonda deploys a new server version, which therefore reaches clients
within 5 minutes. `resources/read` returns document contents and carries
`"ttlMs": 0`, so clients fetch a document again on every read. The `private`
scope keeps shared caches from serving one caller's response to another.

## Authentication

The MCP server uses **Bearer token** authentication. Two token types are
supported:

1. **API key (HS256 JWT)** — the **same** Domonda API key used by the REST and
   GraphQL APIs (see the top-level
   [Authentication section](../README.md#authentication)). The token's `sub`
   claim identifies the client company directly. This is the simplest option
   for programmatic access.
2. **OAuth** — an access token obtained by signing in with your Domonda
   account. The token is mapped to an enabled Domonda user who must be an admin
   or super-admin at the target company. This is what interactive AI clients (Claude
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
3. The client fetches that metadata. Its `authorization_servers` names the MCP
   server itself (`https://domonda.app/api/mcp/`), **not** Auth0: the server
   proxies `/authorize`, `/token` and `/register` through to Auth0 so it can
   add [CIMD](#client-id-metadata-document-cimd) support.
4. The client fetches the authorization server metadata
   ([RFC 8414](https://www.rfc-editor.org/rfc/rfc8414)) from that identifier
   and completes a standard OAuth 2.1 authorization-code flow with PKCE.
5. The client retries with `Authorization: Bearer <access token>`.

To act on a client company other than your default, set the
`X-Selected-Client-Company-ID` header to the target company UUID. The server
verifies your admin or super-admin role at the selected company before
proceeding; super-admins count as admins.

#### Public endpoints

These paths answer without a token. Every other path under `/api/mcp/`
requires authentication:

| Path                                      | Document                                   |
|-------------------------------------------|--------------------------------------------|
| `/.well-known/oauth-protected-resource`   | Protected resource metadata (RFC 9728)     |
| `/.well-known/oauth-authorization-server` | Authorization server metadata (RFC 8414)   |
| `/.well-known/openid-configuration`       | The same document, OpenID Connect fallback |
| `/.well-known/mcp.json`                   | MCP server card                            |
| `/authorize`, `/token`, `/register`       | OAuth endpoints, proxied to Auth0          |

`/register` is public, but does not accept every callback — see
[Dynamic Client Registration](#dynamic-client-registration-register).
`/token` accepts only the grants the authorization server metadata lists — see
[Token endpoint](#token-endpoint-token).

The four well-known documents are also the only paths that answer a CORS
preflight and carry `Access-Control-Allow-Origin: *`, so a browser-based MCP
client (the MCP Inspector, for example) can read them from its own origin.
They are public and carry no credentials.

#### `GET` on the base path

`GET https://domonda.app/api/mcp/` is the Streamable HTTP server-to-client
stream, not a landing page. Without a token it answers `401` with the
`WWW-Authenticate` header above — **that is the healthy response**, and it is
what starts OAuth discovery. A connectivity check that expects `200` here is
checking the wrong thing; use the `POST`-based `scripts/health_check.sh` in
your client folder instead.

Only a request that explicitly asks for HTML gets the human-readable info
page:

```bash
curl -H 'Accept: text/html' https://domonda.app/api/mcp/
```

The stream is long-lived and kept alive with a 30-second heartbeat, so a
client that disappears without closing its connection is noticed and cleaned
up.

Only clients of protocol 2025-11-25 and earlier open this stream. An
authenticated `GET` that declares protocol 2026-07-28 in its
`Mcp-Protocol-Version` header gets `405`, because that revision has no
standing stream (see [Protocol versions](#protocol-versions)).

#### Token endpoint (`/token`)

`/token` takes its parameters form-encoded
(`application/x-www-form-urlencoded`), as
[RFC 6749](https://www.rfc-editor.org/rfc/rfc6749#section-4.1.3) requires, and
accepts the two grants the authorization server metadata lists in
`grant_types_supported`: `authorization_code` and `refresh_token`. A form with
any other `grant_type` is answered with `400 unsupported_grant_type`, and one
without a `grant_type` with `400 invalid_request`. A form that does not parse,
or that repeats `grant_type` or `client_id`, is answered with
`400 invalid_request` as well: RFC 6749 allows each parameter only once
([section 3.2](https://www.rfc-editor.org/rfc/rfc6749#section-3.2)).

**Enterprise-managed authorization is not supported.** On Claude Team and
Enterprise plans, an admin can set up managed authorization for a connector,
so that members are signed in through the organization's identity provider
instead of signing in themselves. Claude then presents an identity assertion
at `/token` with the `urn:ietf:params:oauth:grant-type:jwt-bearer` grant, and
domonda refuses it with `400 unsupported_grant_type`. Leave managed
authorization off for the domonda connector and let each member sign in with
OAuth, or use an [API key](#api-key-flow).

### Client ID Metadata Document (CIMD)

The server supports the
[Client ID Metadata Document](https://modelcontextprotocol.io/specification/2025-11-25/basic/authorization)
mechanism from the 2025-11-25 MCP authorization spec, used by clients that
identify themselves with an `https://` URL as their `client_id` (e.g. Claude
and ChatGPT). No client registration step is required on your side — the server
advertises CIMD support and registers the client transparently.

The registration is stored server-side and outlives restarts and
redeployments, so a connector added once keeps working: you are not sent back
through sign-in and consent because the server forgot the client. The
authorization server metadata also lists `offline_access` in
`scopes_supported`, so a client that follows it asks for a **refresh token**
and renews an expired session on its own instead of opening the browser again.

**Claude Code direct OAuth caveat.** Claude Code declares
`http://localhost/callback` in its CIMD document and then listens on a
port picked at runtime. The proxy accepts that, but the callback registered
with Auth0 carries no port and Auth0 matches callbacks exactly, so the
authorization request is rejected one step later. Claude Code users can use a
Domonda connector signed in on claude.ai with the same account, or authenticate
a directly configured server with an [API key](#api-key-flow). See the
[Claude Code guide](./claude-code/).

#### Dynamic Client Registration (`/register`)

A client that does not use CIMD falls back to OAuth 2.0 Dynamic Client
Registration ([RFC 7591](https://www.rfc-editor.org/rfc/rfc7591))
automatically. That endpoint is unauthenticated — registration is what a
client does before it holds any credential — so it only accepts the callbacks
a **native** client receives on the user's own machine:

- a loopback address on any port: `http://localhost/…`, `http://127.0.0.1/…`,
  `http://[::1]/…`, or
- a private-use URI scheme carrying a path
  ([RFC 8252 section 7.1](https://www.rfc-editor.org/rfc/rfc8252#section-7.1)),
  e.g. `cursor://callback`.

Every other `redirect_uri` — `https://` included — is answered with
`400 invalid_redirect_uri` and the reason, before the registration reaches
Auth0. A client served from the web has a stable HTTPS origin and should
publish a CIMD document instead, which is what claude.ai and chatgpt.com do.
A registration listing more than 20 redirect URIs is rejected as well.

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
- [`chatgpt/`](./chatgpt/) — ChatGPT custom MCP plugin and the OpenAI Responses API / Agents SDK
- [`openclaw/`](./openclaw/) — OpenClaw skill and MCP connection

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
- **`execute_query`** — Execute a read-only SQL query (api schema only).
  Available to API-key callers and to **admin or super-admin** OAuth users;
  other OAuth users are refused. See
  [What `execute_query` rejects](#what-execute_query-rejects), and
  [Why `execute_query`: two measured examples](#why-execute_query-two-measured-examples)
  for what it answers that the filter tools cannot.

### Documents
- **`list_documents`** — List documents with filtering (archived, category, fulltext, dates, amounts, status, …)
- **`get_document`** — Get a single document by ID (excludes `fulltext`)
- **`get_document_fulltext`** — Get the extracted fulltext (OCR) of a document, only of its first 50 pages; the PDF has all pages. **Can be large.**
- **`search_documents`** — Fulltext search across the first 50 pages of all documents
- **`download_document_pdf`** — Download a document PDF (optional audit trail)
- **`get_document_thumbnail`** — Default JPEG thumbnail (256-wide)
- **`get_document_preview`** — Preview image (964×1364 JPEG)
- **`add_document`** — *(write)* Upload a document of up to 64 MiB into the authenticated company

### Invoices
- **`list_invoices`** — List invoices of non-archived documents with filtering (date range, partner, exact invoice number of at least 3 characters, amounts)
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

### What `execute_query` rejects

The SQL you pass is validated before it reaches the database. A query is
rejected — with the reason in the error message — when it:

- does not start with `SELECT` or `WITH`,
- contains a DML or DDL keyword (`INSERT`, `UPDATE`, `DELETE`, `DROP`, …),
- references a schema other than `api` (`public`, `pg_catalog`,
  `information_schema`, …). Unqualified table names resolve to `api`,
- uses a `pg_*` identifier outside the small allowlist of introspection
  helpers,
- calls one of the mutating `api` schema functions by name (`update_document`,
  `create_invoice_accounting_item`, `partner_company_mutation`, …), or
- uses Postgres Unicode-escape syntax (`U&"…"` / `U&'…'`).

The identifier scans (`pg_*` and the mutating `api` functions) run against the
raw SQL text, so a rejected name matches even inside a string literal or a
comment: `WHERE title ILIKE '%update_document%'` is rejected. Search for such
text with `search_documents` or the `search_text` filter of `list_documents`
instead.

A query that tries to escape the scope it was given is failed outright and its
rows discarded, rather than answered.

### Why `execute_query`: two measured examples

The filter tools (`list_invoices`, `list_documents`, `list_money_transactions`,
`get_invoice_payments`, …) return whole rows from one entity at a time. A
question that needs a join, an aggregate or an anti-join therefore makes the
agent fetch every candidate row into its context and do the database's work
itself. `execute_query` runs the question where the data is and returns only
the answer.

Both examples were measured on 2026-11-16 against production
(`https://domonda.app/api/mcp/`) on a small tenant: 967
invoices, 970 documents, 7,055 money transactions, 356 invoice–payment matches.
Bytes are the JSON-RPC response bodies the client receives; times are curl's
`time_total` per call, summed. The time the model spends reading the rows and
calling the next tool comes on top and is not counted. Each `execute_query`
was run three times; the times are the repeat runs.

| | Example 1: supplier spend in EUR | Example 2: payments without receipt |
|---|---|---|
| Filter tools: calls | 8 | 969, and the answer is still incomplete |
| Filter tools: bytes | 226,311 | 1,652,883 |
| Filter tools: time | 2.4 s | 198 s |
| `execute_query`: calls | 1 | 1 |
| `execute_query`: bytes | 514 | 842 |
| `execute_query`: time | 0.19 s | 0.20 s |
| Reduction | 440× fewer bytes, 8× fewer calls | 1,960× fewer bytes, 969× fewer calls |

The filter-tool bytes grow with the tenant; the `execute_query` bytes grow only
with the size of the answer. A tenant 100 times larger needs 100 times the
pages and the `get_invoice_payments` calls, while both queries still return
one row per partner or per month.

#### Example 1: supplier spend per partner in EUR for 2025

*"How much did we spend per supplier in 2025, in EUR?"*

```sql
WITH inv AS (
  SELECT i.partner_company_id,
         CASE WHEN i.credit_memo THEN -i.total ELSE i.total END
         / CASE WHEN i.currency = 'EUR' THEN 1
                ELSE COALESCE(i.conversion_rate, (
                  SELECT r.rate FROM api.currency_rate r
                  WHERE r.currency = i.currency AND r.date <= i.invoice_date
                  ORDER BY r.date DESC LIMIT 1))
           END AS total_eur
  FROM api.invoice i
    JOIN api.document d ON d.id = i.document_id
  WHERE d.type = 'INCOMING_INVOICE'
    AND i.invoice_date >= '2025-01-01' AND i.invoice_date < '2026-01-01'
)
SELECT p.name AS partner, count(*) AS invoices, round(sum(inv.total_eur)::numeric, 2) AS total_eur
FROM inv
  LEFT JOIN api.partner_company p ON p.id = inv.partner_company_id
GROUP BY p.name
ORDER BY total_eur DESC NULLS LAST
LIMIT 10
```

No single filter tool can answer this:

- `list_invoices` has no direction field, so incoming and outgoing invoices
  come back mixed (160 rows for 2025, 15 of them outgoing).
- `list_documents` filters by `document_types: ["INCOMING_INVOICE"]` but
  returns no amounts, currencies or partners (145 rows).
- Neither list tool returns the conversion rate. The 2025 incoming invoices
  include 6 in USD or GBP, and only `get_invoice` returns their EUR amounts.

The cheapest route is therefore `list_invoices` for 2025 (98,178 bytes),
`list_documents` for the incoming IDs (85,178 bytes), and one `get_invoice` per
foreign-currency invoice (6 calls, 42,955 bytes). The agent then intersects 145
IDs with 160 rows, subtracts credit memos, divides the foreign totals by their
rates and sums by partner, all in its context. Done with `jq` for the
measurement, the route returned the same total as the query (12,682.85 EUR). An
LLM doing the same arithmetic over 145 rows by hand is the error-prone part
that the query removes.

#### Example 2: outgoing payments with neither a receipt nor a payment category

*"Which outgoing payments in 2025 have no receipt, per month?"* This is the
monthly closing question: a bank debit that is matched to no document and has
no payment category (`money_category`, such as bank fees or salaries) is a
missing receipt.

```sql
SELECT
  to_char(mt.booking_date, 'YYYY-MM') AS month,
  count(*) AS payments,
  round(sum(abs(mt.amount))::numeric, 2) AS amount
FROM api.money_transaction mt
WHERE mt.type = 'OUTGOING'
  AND mt.booking_date >= '2025-01-01' AND mt.booking_date < '2026-01-01'
  AND mt.money_category_id IS NULL
  AND NOT EXISTS (
    SELECT FROM api.document_money_transaction dmt
    WHERE dmt.money_transaction_id = mt.id
  )
GROUP BY 1
ORDER BY 1
```

The filter tools cannot answer this at all:

- **The category condition cannot be evaluated.** No tool returns
  `money_category_id`; `list_money_transactions` omits it and has no filter
  for it.
- **The match condition is only visible from the invoice side.**
  `list_money_transactions` has no "matched" flag. The matches show up only in
  `get_invoice_payments`, one invoice per call, so the agent must call it for
  every invoice of the tenant: a 2025 payment can settle an invoice from any
  earlier year. Matches to documents that are not invoices are not reachable
  through any tool.

The closest the filter tools get is `list_invoices` for all invoice IDs
(575,179 bytes), `list_money_transactions` for the 828 outgoing 2025
transactions (476,405 bytes), and `get_invoice_payments` for all 967 invoices
(601,299 bytes, 194 s), followed by an anti-join of 828 transaction IDs against
356 matched ones in the agent's context. At 1.65 MB that is more than a
200k-token context window holds, so in practice an agent cannot even load it.
On this tenant the route found the same 775 unmatched payments as the query,
but only because the tenant uses no payment categories and has no
non-invoice matches. On a tenant that books bank fees and salaries through
payment categories, the result would list every one of them as a missing
receipt.

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
JSON-RPC body. Most examples below are plain JSON-RPC calls without `_meta` or
a session ID, which the server accepts without a handshake (see
[Protocol versions](#protocol-versions)). Set your API key first:

```bash
 export DOMONDA_API_KEY="eyJhbGciOiJIUzI1NiIs..."   # leading space keeps this out of shell history
```

Discover the OAuth protected-resource metadata (public, no auth):

```bash
curl -s "https://domonda.app/api/mcp/.well-known/oauth-protected-resource" | jq .
```

Check that the auth challenge is wired — this is expected to answer `401` with
a `WWW-Authenticate` header, see [`GET` on the base path](#get-on-the-base-path):

```bash
curl -s -D- -o /dev/null "https://domonda.app/api/mcp/"
```

Initialize a session / verify a token:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"curl","version":"1.0"}}}' | jq .
```

Discover the server the 2026-07-28 way, with the protocol version, client info
and capabilities in `_meta` and the matching `Mcp-Protocol-Version` and
`Mcp-Method` headers. The result lists `supportedVersions` and carries the
cache hints `ttlMs` and `cacheScope`:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -H "Mcp-Protocol-Version: 2026-07-28" \
  -H "Mcp-Method: server/discover" \
  -d '{"jsonrpc":"2.0","id":1,"method":"server/discover","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientInfo":{"name":"curl","version":"1.0"},"io.modelcontextprotocol/clientCapabilities":{}}}}' | jq .
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
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"execute_query","arguments":{"sql":"SELECT id, title, imported_at FROM api.document ORDER BY imported_at DESC LIMIT 5"}}}' | jq .
```

List recent invoices:

```bash
curl -s -X POST "https://domonda.app/api/mcp/" \
  -H "Authorization: Bearer $DOMONDA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_invoices","arguments":{"from_date":"2025-01-01","until_date":"2025-12-31","limit":10}}}' | jq .
```

Upload a document (up to 64 MiB). `jq` builds the request body and pipes it to
curl, because a base64-encoded file on the command line exceeds the argument
size limit for all but small files; `tr` removes the line breaks GNU `base64`
inserts. Use `list_document_categories` to find a `category_id`:

```bash
base64 < invoice.pdf | tr -d '\n' |
  jq -Rs '{jsonrpc:"2.0",id:1,method:"tools/call",params:{name:"add_document",arguments:{filename:"invoice.pdf",file_data:.,category_id:"CATEGORY_UUID"}}}' |
  curl -s -X POST "https://domonda.app/api/mcp/" \
    -H "Authorization: Bearer $DOMONDA_API_KEY" \
    -H "Content-Type: application/json" \
    --data-binary @- | jq .
```

A request body larger than a 64 MiB file in base64 plus 1 MiB (about 86 MiB)
gets `413`.

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
  status) and paginate when you need more. A document upload through
  `add_document` can be up to **64 MiB**.
