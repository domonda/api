---
name: domonda
description: >
  Use when you need to access domonda financial data — documents, invoices,
  payments, money transactions, companies, GL accounts, partner companies,
  and document workflows. Connects to the domonda MCP server via Streamable
  HTTP with Bearer token authentication (API key or OAuth). All query tools
  are read-only; the only write tool is `add_document`, which uploads a file
  into the authenticated company's document pipeline.
env:
  - name: DOMONDA_API_KEY
    required: true
    secret: true
    description: API key (HS256 JWT) or Auth0 OAuth access token for Bearer authentication
  - name: DOMONDA_URL
    required: false
    secret: false
    description: "Override domonda base URL (default: https://domonda.app)"
---

# domonda MCP Skill

Connect to the [domonda](https://domonda.app) financial data platform via its MCP server.
Query documents, invoices, payments, money transactions, companies, GL accounts, partner companies, and document workflows. All query tools are read-only; the only write tool is `add_document`, which uploads a file into the authenticated company's document pipeline.

> **Language note.** The live MCP tool and parameter descriptions surfaced by
> the server are written in **German**, including a "Nutzen wenn …" block of
> trigger phrases. This matches the domonda-app UI and the natural vocabulary
> of end users, and improves discoverability for clients (e.g. ChatGPT) that
> match user prompts against tool descriptions by embedding similarity. The
> English descriptions in this file are developer-facing reference only —
> what your agent actually sees at runtime is the German text from the server.
>
> If you are unsure which tool to call, call `list_capabilities` first to get
> a topic-grouped catalog of all available tools.

## Agent Security Rules

- **Never** display, log, or transmit the API key in plain text.
- **Never** pass the API key as a CLI argument or embed it in URLs.
- Only use the API key for authenticating against the domonda MCP endpoint.
- Do not write, modify, or delete data via `execute_query` or any query tool — those tools are read-only and any DML will be rejected by the server.
- The only write tool is `add_document`. Only call it when the user has explicitly asked you to upload a document. Never call it on your own initiative, and never use it to write arbitrary data.

## Quick Start

Pick one of:

**a. Custom Connector UI (recommended, production, OAuth):**

1. In Claude Desktop: **Settings → Connectors → Add custom connector**.
2. URL: `https://domonda.app/api/mcp/`
3. Complete the Auth0 OAuth prompt.

**b. `mcp-remote` stdio bridge (localhost dev or raw JWT):**

Edit `claude_desktop_config.json` and add:

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

Restart Claude Desktop afterwards. Verify endpoint connectivity with:

```bash
./scripts/health_check.sh
```

See `README.md` for when to choose which path.

## Available Tools

> The description paragraph for each tool below is written in **German**
> (matching what the running MCP server actually surfaces to clients),
> because domonda end users chat in German and tools are picked by
> similarity to the user's prompt. The parameter tables and technical
> notes stay in English.

### Discovery

#### `list_capabilities`

Liefert einen nach Thema gruppierten Katalog ALLER verfügbaren Domonda-Tools
mit einer kurzen deutschen Einzeiler-Beschreibung pro Tool. Dient als
Einstiegspunkt: bei Unsicherheit, welches Tool passt, zuerst dieses
aufrufen und dann das passende Tool direkt nutzen. Nutzen wenn der Benutzer
nach "was kannst du", "welche Tools", "welche Werkzeuge", "hilf mir",
"Übersicht", "Katalog", "was ist verfügbar", "Capabilities", "Hilfe" fragt.

**Parameters:** None

### Schema Tools

#### `list_tables`

Listet alle Tabellen und Views im `api`-Schema der Domonda-Datenbank (Name,
Typ Tabelle/View, Beschreibung). Nutzen für Datenmodell-Exploration, bevor
eine benutzerdefinierte Abfrage via `execute_query` gebaut wird, oder wenn
der Benutzer nach "welche Tabellen gibt es", "Datenmodell", "Schema", "was
ist in der DB" fragt.

**Parameters:** None

#### `describe_table`

Beschreibt die Spalten einer Tabelle oder View im `api`-Schema:
Spaltennamen, Datentypen, NULL-Zulässigkeit, Default-Werte. Nutzen bevor
eine `SELECT`-Abfrage via `execute_query` geschrieben wird, um die
verfügbaren Spalten zu kennen, oder wenn der Benutzer nach "welche
Spalten", "Struktur der Tabelle", "Felder" fragt.

| Parameter    | Type   | Required | Description                                                                                     |
|------------- |--------|----------|-------------------------------------------------------------------------------------------------|
| `table_name` | string | yes      | Name of the table or view in the `api` schema (e.g. `invoice`, `document`, `partner_company`)   |

### Query Tools

#### `execute_query`

Führt eine READ-ONLY SQL-Abfrage (`SELECT` oder `WITH`) gegen das
`api`-Schema der Domonda-Datenbank aus. Nur das `api`-Schema ist zugänglich;
`pg_catalog`, `information_schema` und andere Schemata werden abgelehnt.
Unqualifizierte Tabellennamen werden ins `api`-Schema aufgelöst. Ergebnis:
JSON-Array von Objekten, maximal 1000 Zeilen. Daten sind automatisch auf
den angemeldeten Mandanten eingeschränkt. Nutzen nur wenn die
spezialisierten Tools (`list_documents`, `list_invoices`,
`list_money_transactions`, …) nicht ausreichen und eine eigene Abfrage
nötig ist. Vorher `list_tables` und `describe_table` nutzen, um das Schema
zu erkunden.

| Parameter | Type   | Required | Description                                                                                        |
|-----------|--------|----------|----------------------------------------------------------------------------------------------------|
| `sql`     | string | yes      | The SQL query to execute. Must start with `SELECT` or `WITH`. Only the `api` schema is accessible. |

### Document Tools

#### `list_documents`

Sucht, filtert und listet Dokumente des Mandanten: Rechnungen, Verträge,
Belege, Lieferscheine, Kontoauszüge, sonstige Dokumente, DMS. Unterstützt
Volltextsuche, Datumsbereiche (mehrere Datumstypen), Beträge, Zahl- und
Exportstatus, Geschäftspartner, Kategorien, Dokumentarten,
Genehmigungsstatus, Duplikate, Kommentare, Vertragstypen und konfigurierbare
Sortierung. Nutzen wenn der Benutzer nach "Dokumente", "Belege",
"Rechnungen", "Verträge", "Lieferscheine", "Kontoauszüge", "Mahnungen",
"was liegt im Eingang", "Dokumentenliste", "welche Dokumente gibt es zu …",
"suche Beleg über …", "offene Rechnungen", "nicht freigegebene", "nicht
exportierte" fragt. Für reine Volltextsuche ohne Filter siehe
`search_documents`.

All filter parameters are optional; omitting a parameter lets the server default apply.

| Parameter                | Type            | Description                                                                                                               |
|--------------------------|-----------------|---------------------------------------------------------------------------------------------------------------------------|
| `limit`                  | number          | Maximum number of documents to return (default 50, max 1000)                                                              |
| `offset`                 | number          | Number of documents to skip for pagination                                                                                |
| `search_text`            | string          | Fulltext search over the extracted document content (OCR + structured data)                                               |
| `archived`               | string enum     | `exclude` (default, only active), `only` (only archived), or `include` (both)                                             |
| `has_warning`            | boolean         | `true`: only documents with a warning; `false`: only documents without                                                    |
| `internal_number`        | string          | Filter by domonda-internal sequential document number                                                                     |
| `date_filter_type`       | string enum     | Which date `from_date`/`until_date` apply to: `DOCUMENT_DATE`, `INVOICE_DATE` (default), `IMPORT_DATE`, `DUE_DATE`, `DISCOUNT_DUE_DATE`, `PAID_DATE`, `DELIVERED_DATE` |
| `from_date`              | string          | Lower bound (inclusive) for the date selected via `date_filter_type` (YYYY-MM-DD)                                         |
| `until_date`             | string          | Upper bound (inclusive) for the date selected via `date_filter_type` (YYYY-MM-DD)                                         |
| `min_total`              | number          | Only documents with invoice gross total ≥ this amount                                                                     |
| `max_total`              | number          | Only documents with invoice gross total ≤ this amount                                                                     |
| `export_status`          | string enum     | Accounting export status (exported / not exported, for "booking" or "ready for booking"; see tool schema for the full enum) |
| `paid_status`            | string enum     | Invoice payment status: paid, partially paid, not paid, not payable, or paid via a specific payment method                 |
| `document_types`         | string[] enum   | One or more document types (incoming/outgoing invoice, dunning letter, delivery note, bank/credit-card statement, other document, DMS, …) |
| `document_category_ids`  | string[] (UUID) | One or more document category IDs (from `list_document_categories`)                                                        |
| `partner_company_ids`    | string[] (UUID) | One or more partner company IDs (customer or supplier, from `list_partner_companies`)                                      |
| `is_credit_memo`         | boolean         | `true`: only credit memos (Gutschriften); `false`: exclude credit memos                                                    |
| `has_pain001_payment`    | boolean         | `true`: only documents with a generated SEPA credit transfer (PAIN.001); `false`: without                                  |
| `has_pain008_direct_debit` | boolean       | `true`: only documents with a generated SEPA direct debit (PAIN.008); `false`: without                                     |
| `is_payment_reminder_sent` | boolean       | `true`: only documents for which a payment reminder/dunning letter has already been sent                                   |
| `is_approved`            | boolean         | `true`: only approved documents; `false`: only not-yet-approved documents (waiting for workflow approval)                  |
| `is_duplicate`           | boolean         | `true`: only documents detected as duplicates; `false`: only non-duplicates                                                |
| `has_comments`           | boolean         | `true`: only documents with comments; `false`: without                                                                     |
| `is_recurring`           | boolean         | `true`: only recurring documents; `false`: only one-off documents                                                          |
| `contract_type`          | string enum     | Contract type filter for `OTHER_DOCUMENT`s (service contracts, rental contracts, insurance, …; see tool schema)            |
| `order_bys`              | string[] enum   | Ordering rules applied in sequence; default `IMPORT_DATE_DESC` (most recently imported first)                              |

> The full enum values for `export_status`, `paid_status`, `document_types`, `contract_type`, and `order_bys` are exposed on the live tool schema — your MCP client will surface them.

#### `get_document`

Liefert Metadaten eines einzelnen Dokuments anhand seiner UUID: Dokumentart,
Kategorie, Geschäftspartner, Datum, Status, Beträge (soweit vorhanden).
Enthält **nicht** den Volltext — siehe `get_document_fulltext` für den
extrahierten OCR-Text oder `get_invoice` für Rechnungs-Positionen. Nutzen
wenn der Benutzer nach "Details zum Dokument", "Info zu diesem Beleg", "was
ist das für ein Dokument", "wer hat das hochgeladen", "welche Kategorie"
fragt.

| Parameter     | Type   | Required | Description              |
|---------------|--------|----------|--------------------------|
| `document_id` | string | yes      | The UUID of the document |

#### `get_document_fulltext`

Liefert den extrahierten Volltext (OCR) eines Dokuments. Nutzen wenn der
Benutzer den kompletten Inhalt eines Dokuments lesen will — "Volltext",
"was steht im Beleg", "OCR-Text", "Inhalt dieses Dokuments".

**Warning:** fulltexts can be very large; long documents may approach or
exceed the MCP response byte cap (default 10 MiB)
and the call will fail. For pure search prefer `search_documents` and only
fetch the fulltext of a specific hit when really needed.

| Parameter     | Type   | Required | Description              |
|---------------|--------|----------|--------------------------|
| `document_id` | string | yes      | The UUID of the document |

#### `search_documents`

Volltextsuche über alle Dokumente des Mandanten; Treffer sind nach Relevanz
sortiert. Nutzen wenn der Benutzer etwas FINDEN will, statt gezielt zu
filtern — z. B. "finde Beleg über Büromöbel", "such mir alles zu Projekt
X", "gibt es eine Rechnung mit …", "wo taucht der Begriff … auf". Für
strukturierte Filter (Datum, Partner, Status, Beträge) lieber
`list_documents` verwenden.

| Parameter | Type   | Required | Description                                          |
|-----------|--------|----------|------------------------------------------------------|
| `query`   | string | yes      | Search term or phrase for the fulltext search        |
| `limit`   | number | no       | Maximum number of results (default: 20, max: 1000)   |

#### `download_document_pdf`

Lädt die Original-PDF-Datei eines Dokuments herunter — optional mit
angehängtem oder vorangestelltem Audit-Trail (Protokoll der
Bearbeitungsschritte, Freigaben, Änderungen). Standard: die PDF-Bytes
werden base64-kodiert inline zurückgegeben. Mit `inline: false` wird
stattdessen ein MCP `resource_link` auf `document://{document_id}/pdf?...`
geliefert, den ein resource-fähiger Client via `resources/read` nachladen
kann. Nutzen wenn der Benutzer "die PDF", "das Original", "den Beleg
herunterladen", "Rechnung als PDF", "mit Audit-Trail" sagt.

| Parameter          | Type    | Required | Description                                                                        |
|--------------------|---------|----------|------------------------------------------------------------------------------------|
| `document_id`      | string  | yes      | The UUID of the document                                                           |
| `audit_trail`      | string  | no       | Audit-trail position in the PDF: `append`, `prepend`, `configured` (mandant setting), or `only` (audit trail only). Empty = no audit trail. |
| `audit_trail_lang` | string  | no       | Language of the audit trail as ISO 639-1 code (default `de`)                       |
| `inline`           | boolean | no       | `true` (default): return the PDF bytes inline as a base64-encoded resource. `false`: return an MCP `resource_link` to `document://{document_id}/pdf?...` that a resource-aware client loads via `resources/read`. |

The PDF is also exposed as an MCP resource at `document://{document_id}/pdf{?audit_trail,audit_trail_lang}`, so resource-aware clients can fetch it via `resources/read`. The `configured` audit-trail value is resolved to `append`/`prepend` before a resource link is returned, so resource URIs are deterministic and cacheable.

#### `get_document_thumbnail`

Liefert die kleine Miniaturansicht (JPEG, 256 px breit, Seitenverhältnis
erhalten) eines Dokuments als base64-kodierten Blob. Auch als MCP-Ressource
unter `document://{document_id}/thumb_256x0.jpeg` verfügbar. Nutzen wenn
der Benutzer eine schnelle Vorschau, ein Miniaturbild oder ein "Icon" eines
Dokuments sehen will. Für eine größere Vorschau siehe `get_document_preview`.

| Parameter     | Type   | Required | Description              |
|---------------|--------|----------|--------------------------|
| `document_id` | string | yes      | The UUID of the document |

#### `get_document_preview`

Liefert die große Vorschau (JPEG, 964×1364) eines Dokuments als
base64-kodierten Blob — das größte vorgerenderte Vorschaubild. Auch als
MCP-Ressource unter `document://{document_id}/preview.jpeg` verfügbar.
Nutzen wenn der Benutzer "zeig mir das Dokument", "Vorschau", "wie sieht
der Beleg aus", "Bild des Dokuments" sagt. Für eine kleinere Miniaturansicht
siehe `get_document_thumbnail`; für die Original-PDF siehe
`download_document_pdf`.

| Parameter     | Type   | Required | Description              |
|---------------|--------|----------|--------------------------|
| `document_id` | string | yes      | The UUID of the document |

#### `add_document` (write)

Lädt ein neues Dokument (Rechnung, Beleg, Vertrag, Lieferschein, sonstige
Datei) für den angemeldeten Mandanten hoch. Dubletten-Erkennung und
Kategorie-Zuweisung laufen automatisch. Nutzen wenn der Benutzer "Dokument
hochladen", "Beleg importieren", "Rechnung hinzufügen", "upload", "neues
Dokument", "Datei einspielen" sagt. Vorher `list_document_categories`
aufrufen, um eine gültige `category_id` zu ermitteln.

Supported file types: PDF, PNG, JPEG, TIFF, XML, HTML, RTF, DOC, DOCX, ODT,
XLS, XLSX, ODS. This is the **only write operation** on the MCP server; the
file is processed through the same pipeline as the public upload endpoint.
Returns the new document UUID on success.

| Parameter                 | Type    | Required | Description                                                                                |
|---------------------------|---------|----------|--------------------------------------------------------------------------------------------|
| `file_data`               | string  | yes      | Base64-encoded file content (no URL, no path — the raw bytes base64-encoded)                |
| `filename`                | string  | yes      | Original filename including extension (e.g. `invoice_2026_04.pdf`)                         |
| `category_id`             | string  | yes      | Document category UUID (use `list_document_categories` to find valid IDs)                  |
| `tags`                    | string[]| no       | Optional document tags                                                                     |
| `allow_duplicate_deleted` | boolean | no       | If `true`, allow re-uploading a file whose hash matches a previously-deleted document. Default `false`. |
| `wait_for_extraction`     | boolean | no       | If `true`, wait synchronously for OCR/extraction before returning. Default `false` (asynchronous). |

> **Agent guidance:** Only call `add_document` when the user has explicitly asked you to upload a file. Do not use it on your own initiative.

### Invoice Tools

#### `list_invoices`

Listet Rechnungen des Mandanten (Eingangs- und Ausgangsrechnungen, inkl.
Gutschriften) mit Rechnungsnummer, Rechnungsdatum, Fälligkeitsdatum,
Geschäftspartner, Netto, Brutto, Währung, USt-Satz, Zahlstatus, offenem
Betrag und Buchungsstatus. Filter: Datumsbereich, Partner-Name,
Mindest-/Maximalbetrag. Nutzen wenn der Benutzer nach "Rechnungen",
"Eingangsrechnungen", "Ausgangsrechnungen", "ER", "AR", "Gutschriften",
"offene Posten", "Forderungen", "Verbindlichkeiten", "unbezahlte
Rechnungen", "Rechnungen von <Firma>" fragt. Für einzelne Rechnungsdetails
siehe `get_invoice`.

| Parameter      | Type   | Required | Description                                                                                                  |
|----------------|--------|----------|--------------------------------------------------------------------------------------------------------------|
| `limit`        | number | no       | Maximum number of invoices to return (default: 50, max: 1000)                                                |
| `offset`       | number | no       | Number of invoices to skip for pagination (default: 0)                                                       |
| `from_date`    | string | no       | Lower bound (inclusive) for the invoice date (YYYY-MM-DD)                                                    |
| `until_date`   | string | no       | Upper bound (inclusive) for the invoice date (YYYY-MM-DD)                                                    |
| `partner_name` | string | no       | Filter by partner company name (customer or supplier, case-insensitive partial match)                        |
| `min_total`    | number | no       | Only invoices with gross total ≥ this amount                                                                 |
| `max_total`    | number | no       | Only invoices with gross total ≤ this amount                                                                 |

#### `get_invoice`

Liefert alle Details zu einer einzelnen Rechnung: Kopfdaten, Positionen
(Line Items), Skonto, Fälligkeit, umgerechnete EUR-Beträge, Geschäftspartner,
UID/VAT-ID des Partners, zugeordnete Zahlungen. Nutzen wenn der Benutzer
eine konkrete Rechnung im Detail sehen will — "zeig mir Rechnung …",
"Details zu dieser Rechnung", "was steht auf der Rechnung", "Positionen",
"Line Items", "Skonto", "Fälligkeit". Für eine Liste mehrerer Rechnungen
siehe `list_invoices`.

| Parameter     | Type   | Required | Description                                                                   |
|---------------|--------|----------|-------------------------------------------------------------------------------|
| `document_id` | string | yes      | The UUID of the invoice's document (from `list_invoices` or `list_documents`) |

### Payment Tools

#### `list_money_accounts`

Listet alle Bankkonten und Zahlungskonten des Mandanten: Bankkonten,
Kreditkarten, PayPal, Wise, Kassa. Nutzen wenn der Benutzer nach
"Bankkonto", "Bankkonten", "Konten", "Kreditkarte", "PayPal", "Kassa",
"welche Konten haben wir", "IBAN" fragt — oder wenn eine Kontoauswahl für
weitere Filter benötigt wird.

**Parameters:** None

#### `list_money_transactions`

Listet Banktransaktionen und Kreditkartenbewegungen des Mandanten mit
Betrag, Währung, Buchungsdatum, Valutadatum, Empfänger/Zahler,
Verwendungszweck und Typ (EINGEHEND/AUSGEHEND). Filter: Konto, Typ,
Datumsbereich, Partner-Name. Nutzen wenn der Benutzer nach "Transaktionen",
"Banktransaktionen", "Buchungen auf dem Konto", "Kontoauszug", "Ein- und
Ausgänge", "Zahlungseingang", "Zahlungsausgang", "Umsätze" fragt.

| Parameter      | Type   | Required | Description                                                           |
|----------------|--------|----------|-----------------------------------------------------------------------|
| `limit`        | number | no       | Maximum number of transactions to return (default: 50, max: 1000)     |
| `offset`       | number | no       | Number of transactions to skip for pagination (default: 0)            |
| `account_id`   | string | no       | Filter to a specific money account (UUID from `list_money_accounts`)  |
| `type`         | string | no       | Filter by direction: `INCOMING` (credit) or `OUTGOING` (debit)        |
| `from_date`    | string | no       | Lower bound (inclusive) for the booking date (YYYY-MM-DD)             |
| `until_date`   | string | no       | Upper bound (inclusive) for the booking date (YYYY-MM-DD)             |
| `partner_name` | string | no       | Text filter on the counterparty name (case-insensitive partial match) |

#### `get_invoice_payments`

Liefert alle Zahlungen (Banktransaktionen), die einer bestimmten Rechnung
zugeordnet sind, inkl. offenem Betrag, bereits verrechneter Summe,
Zahlstatus und Zahldatum. Nutzen wenn der Benutzer nach "Zahlungen zur
Rechnung", "ist bezahlt", "offener Betrag", "welche Zahlung wurde
zugeordnet", "Zahlungsabgleich", "Payment Matching" für eine konkrete
Rechnung fragt.

| Parameter     | Type   | Required | Description                                                                   |
|---------------|--------|----------|-------------------------------------------------------------------------------|
| `document_id` | string | yes      | The UUID of the invoice's document (from `list_invoices` or `list_documents`) |

### Company Tools

#### `get_my_company`

Liefert die Stammdaten des eigenen Mandanten (angemeldete Firma):
Firmenname, Rechtsform, Hauptsitz/Adresse, Telefon, E-Mail, Website,
UID/VAT-Nummer, Steuernummer, Handelsregister-Nummer. Nutzen wenn der
Benutzer nach "meine Firma", "mein Mandant", "unsere Firma", "Firmendaten",
"Stammdaten", "UID", "VAT-ID", "Steuernummer", "Firmenadresse", "Hauptsitz",
"wer sind wir" fragt.

**Parameters:** None

#### `list_document_categories`

Listet alle Dokumentenkategorien des Mandanten (visuelle Kategorien für
die Ablage). Nutzen wenn der Benutzer nach "Kategorien",
"Dokumentenkategorien", "Ablagekategorien", "welche Kategorien gibt es",
"wo wird abgelegt" fragt oder bevor ein Dokument via `add_document`
hochgeladen wird, um eine gültige `category_id` zu ermitteln.

**Parameters:** None

#### `list_gl_accounts`

Listet alle Sachkonten (GL accounts, Hauptbuchkonten) des Mandanten mit
Nummer und Bezeichnung. Nutzen wenn der Benutzer nach "Sachkonto",
"Sachkonten", "Kontenplan", "Buchungskonto", "GL account", "general
ledger", "Hauptbuch", "Kontonummer" fragt — typischerweise für
Buchhaltung, Kontierung und Accounting-Export.

**Parameters:** None

#### `list_partner_companies`

Listet Geschäftspartner des Mandanten: Kunden, Lieferanten, Debitoren,
Kreditoren — mit Namen, abgeleitetem Namen, Debitoren- und
Kreditoren-Kontonummer. Optional nach Namen filterbar. Nutzen wenn der
Benutzer nach "Kunde", "Kunden", "Lieferant", "Lieferanten",
"Geschäftspartner", "Partner", "Debitor", "Kreditor", "vendor",
"supplier", "customer" fragt.

| Parameter     | Type   | Required | Description                                                                                    |
|---------------|--------|----------|------------------------------------------------------------------------------------------------|
| `search_text` | string | no       | Optional text filter: partner companies whose name contains this substring (case-insensitive)  |
| `limit`       | number | no       | Maximum number of results (default: 50, max: 1000)                                             |

### User Tools

#### `get_user`

Liefert den aktuell angemeldeten Benutzer (Name, E-Mail, Sprache, Mandant).
Nutzen wenn der Benutzer nach "wer bin ich", "mein Profil", "mein Konto",
"Benutzer", "User", "angemeldet" fragt oder wenn ein anderes Tool die
aktuelle Benutzer-ID benötigt. Funktioniert nur bei OAuth-Login —
API-Key-Zugriff liefert hier einen Fehler, weil der Key einen Mandanten
identifiziert, keinen Benutzer.

**Parameters:** None

### Workflow Tools

#### `list_document_workflows`

Listet alle im Mandanten konfigurierten Dokument-Workflows
(Freigabe-Prozesse) mit Namen und IDs. Ein Workflow ist ein Ablauf, den
ein Dokument durchläuft, bestehend aus mehreren Workflow-Schritten. Nutzen
wenn der Benutzer nach "Workflows", "Freigabe-Prozess",
"Genehmigungsprozess", "Approval-Flow", "welche Workflows gibt es" fragt.
Für die einzelnen Schritte eines Workflows siehe
`list_document_workflow_steps`.

| Parameter | Type   | Required | Description                                                |
|-----------|--------|----------|------------------------------------------------------------|
| `limit`   | number | no       | Maximum number of rows to return (default: 100, max: 1000) |
| `offset`  | number | no       | Number of rows to skip for pagination (default: 0)         |

#### `list_document_workflow_steps`

Listet Workflow-Schritte des Mandanten (einzelne Stationen eines
Freigabe-Workflows). Ein Dokument steht zu jedem Zeitpunkt in genau einem
Workflow-Schritt. Nutzen wenn der Benutzer nach "Workflow-Schritten",
"Freigabeschritten", "Genehmigungsstufen", "welche Schritte hat unser
Workflow", "Approval-Steps" fragt. Für den aktuellen Schritt konkreter
Dokumente siehe `list_document_workflow_states`.

| Parameter     | Type   | Required | Description                                                                                   |
|---------------|--------|----------|-----------------------------------------------------------------------------------------------|
| `workflow_id` | string | no       | Optional: UUID of a specific workflow to list only its steps (from `list_document_workflows`) |
| `limit`       | number | no       | Maximum number of rows to return (default: 100, max: 1000)                                    |
| `offset`      | number | no       | Number of rows to skip for pagination (default: 0)                                            |

#### `list_document_workflow_states`

Liefert den AKTUELLEN Workflow-Zustand von Dokumenten: in welchem
Workflow-Schritt steht ein Dokument gerade, und wie lautet sein Status.
Eine Zeile pro Dokument. Nutzen wenn der Benutzer nach "Wo steht das
Dokument gerade", "in welchem Schritt", "Freigabestatus", "wartet auf
Freigabe von …", "welche Dokumente hängen wo", "Workflow-Status" fragt.

| Parameter     | Type   | Required | Description                                                            |
|---------------|--------|----------|------------------------------------------------------------------------|
| `document_id` | string | no       | Optional: UUID of a single document to return only its workflow state  |
| `limit`       | number | no       | Maximum number of rows to return (default: 50, max: 1000)              |
| `offset`      | number | no       | Number of rows to skip for pagination (default: 0)                     |

### App URL Tools

#### `get_app_urls`

Liefert Deep-Links in die Domonda Web-App, zugeschnitten auf den
angemeldeten Mandanten: Dokumente-Übersicht, Bank (Banking),
Buchhaltungsansicht, Workflows/Freigaben, Academy (Hilfe/Doku),
API-Repository. Nutzen wenn der Benutzer "öffne in Domonda", "zeig mir in
der App", "Link zu …", "wo finde ich das", "geh in die Bank-Ansicht",
"geh in die Buchhaltung" sagt.

- `documents_overview` — main documents list
- `banking` — banking overview
- `accounting` — documents list with accounting view
- `workflows` — documents list with workflow view (unapproved)
- `documentation` — user-facing documentation (academy)
- `api` — public API repository

**Parameters:** None

#### `get_document_urls`

Liefert Deep-Links in die Domonda Web-App für EIN einzelnes Dokument: Tabs
Prüfung (`review`), Zahlung (`payment`), Buchhaltung (`accounting`) und
Zusammenarbeit (`collaboration`). Nutzen wenn der Benutzer "öffne dieses
Dokument in Domonda", "Link zum Beleg", "spring zu
Zahlung/Buchhaltung/Prüfung" für ein konkretes Dokument sagt. Das Dokument
muss zum angemeldeten Mandanten gehören.

| Parameter     | Type   | Required | Description              |
|---------------|--------|----------|--------------------------|
| `document_id` | string | yes      | The UUID of the document |

#### `get_document_selection_url`

Liefert Deep-Links in die Domonda Web-App, die eine Auswahl MEHRERER
Dokumente gemeinsam öffnen — je eine URL pro Mehrfach-Ansicht:
`thumb_view` (Miniaturansicht), `list_view` (Listenansicht),
`accounting_view` (Buchhaltungsansicht), `workflow_view` (Freigabenansicht).
Nutzen wenn der Benutzer "öffne diese Rechnungen in Domonda", "zeig mir
diese Belege in der App", "Link zu diesen Dokumenten",
"Buchhaltungsansicht für diese Dokumente" sagt. Alle Dokumente müssen zum
angemeldeten Mandanten gehören.

URLs use the `?documentRowIds[]=<uuid>&view=<VIEW>` schema consumed by the
domonda-app frontend routing; document ownership is verified per-id on the
server.

| Parameter      | Type            | Required | Description                                                              |
|----------------|-----------------|----------|--------------------------------------------------------------------------|
| `document_ids` | string[] (UUID) | yes      | List of document UUIDs to open together in the selection (at least one) |

## Authentication

The domonda MCP server uses JWT Bearer token authentication. The API key is a JWT where the `sub` claim contains the client company ID. All queries execute within a read-only, company-scoped database transaction.

Pass the token in the HTTP `Authorization` header:

```
Authorization: Bearer <your-api-key>
```

## Usage Tips

- Start with `get_my_company` to confirm connectivity and see which company you are authenticated as.
- Use `list_tables` and `describe_table` for schema discovery before writing queries.
- Use `execute_query` for ad-hoc SQL queries when the specialized tools don't cover your use case.
- All date parameters use `YYYY-MM-DD` format.
- UUID parameters must be valid UUID strings.
- Results are capped at 1000 rows — use `limit` and `offset` for pagination.
- `offset` is clamped to 100,000; stream beyond that by narrowing filters (date range, partner, status) instead of paginating deeper.
- `execute_query` rejects SQL longer than 10,000 characters.
