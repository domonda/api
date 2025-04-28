![image](img/domonda_logo_schriftzug.png)

# DOMONDA API

The Domonda API consists of a [GraphQL](http://graphql.org/)
and a [REST](https://en.wikipedia.org/wiki/REST) part.

[**GraphQL**](#graphql-api) is the main API for querying data and changing single data items with mutations.

[**REST**](#rest-api) is used for file up- and downloads. This includes PDFs for invoices
but also importing bulk master data for client companies.

[**Authentication**](#authentication) of requests is identical for GraphQL and REST
and is based on API keys that give access to the data of a client company.

## Table of Content

1. [**Authentication**](#authentication)
2. [**GraphQL API**](#graphql-api)
   * [Graph*i*QL interactive access and documentation](#graphiql-interactive-access-and-documentation)
   * [GraphQL query examples](#graphQL-query-examples)
3. [**REST API**](#rest-api)
   * [Document PDF download](#document-pdf-download)
   * [File uploads](#file-uploads)
   * [Upload structured invoice data as JSON](#upload-structured-invoice-data-as-json)
   * [Upload company master data as JSON](#upload-company-master-data-as-json)

## Authentication

Authentication for [GraphQL](#graphql-api) and [REST](#rest-api) requests
is implemented via `Bearer` token that holds an API key specific to a client-company. 

In the following examples, replace `${DOMONDA_API_KEY}` with your actual API key.

```http
POST https://domonda.app/api/public/graphql

Authorization: Bearer ${DOMONDA_API_KEY}
```

Minimal example:

```sh
curl -X POST \
  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
  -H "Content-Type: application/graphql" \
  --data "{ allDocuments{ totalCount } }" \
  -o - \
  https://domonda.app/api/public/graphql
```


## GraphQL API

### Graph*i*QL interactive access and documentation

The URL for the **Domonda API Graph*i*QL tool** is:
https://domonda.app/api/public/graphiql

You can provide an authentication token to authenticate access to your Domonda data.
Without a valid token, demo data is provided by the API.

Alternatively, you can use the desktop client Altair (<https://altair.sirmuel.design/>).

For an intro on how to work With GraphQL and the interactive Graph*i*GL tool, we recommend watching this YouTube video:
https://www.youtube.com/watch?v=L5bzHzqUn9M


### Basic GraphQL usage

You can find all GraphQL query types in the generated documentation:
<https://domonda.github.io/api/doc/schema/query.doc.html>

Referenced fields always have an extra field for querying the actual data behind it.

Query own company data:

```gql
{
  currentClientCompany {
    companyRowId
    hasAccounting
    companyByCompanyRowId {
      name
      brandName
      alternativeNames
    }
    mainLocation {
      street
      city
      zip
      country
      vatNo
    }
  }
}
```

If you want to query all documents with the additional information of the import user, you can achieve this by using the
field `userByImportedBy`, which gets you the associated user.

```gql
{
  query {
    allDocuments {
      nodes {
        rowId
        userByImportedBy {
          firstName
          lastName
        }
      }
    }
  }
}
```

### GraphQL query examples

#### Query all document categories:

```gql
{
  allDocumentCategories {
    nodes {
      rowId
      documentType
      bookingType
      bookingCategory
      description
      emailAlias
      createdAt
    }
  }
}
```

#### Query all documents:

```gql
{
  allDocuments {
    nodes {
      rowId
      categoryRowId
      workflowStepRowId
      name
      title
      language
      tags
      numPages
      numAttachPages
      version
      importedBy
      updatedAt
      createdAt
    }
  }
}
```

#### Query all invoices:

```gql
{
  allInvoices {
    nodes {
      documentRowId
      partnerName
      invoiceNumber
      invoiceDate
      dueDate
      orderNumber
      orderDate
      creditMemo
      net
      total
      vatPercent
      vatPercentages
      discountPercent
      discountUntil
      currency
      conversionRate
      conversionRateDate
      conversionRateSource
      goodsServices
      deliveredFrom
      deliveredUntil
      iban
      bic
    }
  }
}
```

#### All incoming invoices of a month, including payment status

```gql
{
  filterDocuments(dateFilterType: INVOICE_DATE, fromDate: "2022-06-01", untilDate: "2022-07-01", documentTypes: [INCOMING_INVOICE], orderBys: [INVOICE_DATE_ASC]) {
    nodes {
      rowId
      documentCategoryByCategoryRowId {
        documentType
      }
      tags
      paymentStatus
      documentMoneyTransactionsByDocumentRowId {
        nodes {
          moneyTransactionByMoneyTransactionRowId {
            bookingDate
            amount
          }
        }
      }
      invoiceByDocumentRowId {
        invoiceDate
        invoiceNumber
        net
        total
        paymentStatus
        paidDate
        invoiceCostCentersByInvoiceDocumentRowId {
          nodes {
            clientCompanyCostCenterByClientCompanyCostCenterRowId {
              number
              description
            }
            amount
          }
        }
        partnerCompanyByPartnerCompanyRowId {
          name
          alternativeNames
          derivedName
        }
        companyLocationByPartnerCompanyLocationRowId {
          main
          street
          city
          zip
          country
          email
          website
          vatNo
          registrationNo
        }
      }
    }
  }
}
```

#### Query invoice booking lines (accounting items):

```gql
{
  invoiceByDocumentRowId(
    documentRowId: "035bda2e-a5a1-445d-a712-6943e803f108"
  ) {
    invoiceAccountingItemsByInvoiceDocumentRowId {
      nodes {
        title
        amount
        amountType
        bookingType
        generalLedgerAccountRowId
        valueAddedTaxRowId
        valueAddedTaxPercentageRowId
      }
    }
  }
}
```

#### Query all delivery notes:

```gql
{
  allDeliveryNotes {
    nodes {
      documentRowId
      partnerCompanyRowId
      invoiceDocumentRowId
      invoiceNumber
      deliveryNoteNumber
      deliveryDate
      createdAt
    }
  }
}
```

#### Query all delivery note items:

```gql
{
  allDeliveryNoteItems {
    nodes {
      deliveryNoteDocumentRowId
      posNo
      quantity
      productNo
      gtinNo
      eanNo
      description
    }
  }
}
```

#### Find money transactions:

```gql
{
  filterMoneyTransactions(searchText: "My Transaction Reference") {
    nodes {
      rowId
      accountRowId
      type
      partnerName
      partnerIban
      partnerCompanyRowId
      amount
      foreignCurrency
      foreignAmount
      purpose
      bookingDate
      valueDate
      importDocumentRowId
      moneyCategoryRowId
      updatedAt
      createdAt
    }
  }
}
```

#### Load last month's invoices in Google Sheets

1. Open Google Sheets
1. In the app bar, click on Extensions -> Apps Script
1. Add [Moment](https://momentjs.com)
   1. Click on plus in Libraries
   1. Use script ID: `15hgNOjKHUG4UtyZl9clqBbl23sDvWMS8pfDJOyIapZk5RBqwL3i-rlCo`
   1. Click on `Look up`
   1. Use v9 and `Moment` as an identifier
   1. Click on Add
1. Use script:

   ```js
   const DOMONDA_API = "https://domonda.app/api/public/graphql";
   const API_TOKEN = ""; // your personal access token

   function lastMonthImportedInvoices() {
     var lastMonth = Moment.moment(new Date()).subtract(1, "months");
     var options = {
       method: "POST",
       headers: { Authorization: "Bearer " + API_TOKEN },
       contentType: "application/json",
       payload: JSON.stringify({
         // Use our interactive API explorer here https://domonda.app/api/public/graphiql to get the most out of your data.
         query: `query lastMonthImportedInvoices($from: Date!, $until: Date!) {
            filterDocuments(dateFilterType: IMPORT_DATE, fromDate: $from, untilDate: $until) {
              nodes {
                invoice: invoiceByDocumentRowId {
                  partnerName
                  partnerVatID: partnerVatRowIdNo
                  invoiceDate
                  invoiceNumber
                  totalInEur
                }
              }
            }
          }`,
         variables: {
           from: lastMonth.startOf("month").format("YYYY-MM-DD"),
           until: lastMonth.endOf("month").format("YYYY-MM-DD"),
         },
       }),
     };

     var response = UrlFetchApp.fetch(DOMONDA_API, options);

     var json = JSON.parse(response.getContentText());
     var invoices = json.data.filterDocuments.nodes.map(
       ({ invoice }) => invoice
     );

     var rows = [];
     for (const invoice of invoices) {
       // skip non-invoices
       if (!invoice) continue;

       const row = [];

       // when rows is empty, start by creating the header
       if (rows.length === 0) {
         for (const key of Object.keys(invoice)) {
           row.push(key);
         }
         rows.push(row);
         continue;
       }

       // first row is always header
       const header = rows[0];

       // add data rows following the header
       for (const key of header) {
         row.push(invoice[key]);
       }
       rows.push(row);
     }

     const active = SpreadsheetApp.getActive();

     // get or create sheet named YYYY-MM
     const sheetName = lastMonth.format("YYYY-MM");
     let sheet = active.getSheetByName(sheetName);
     if (!sheet) {
       sheet = active.insertSheet();
       sheet.setName(sheetName);
     }

     if (rows.length === 0) {
       sheet.getRange(1, 1).setValue("No data");
     } else {
       sheet.getRange(1, 1, rows.length, rows[0].length).setValues(rows);
     }

     active.setActiveSheet(sheet);
   }
   ```

1. Click on Run
1. A new sheet with data named `YYYY-MM` should be added to the Google Sheet

## REST API

### Document PDF download

To request the PDF file for the document with the ID `00000000-0000-0000-0000-000000000000`
(replace zeros with actual UUID hex-code) make the following GET request:

```sh
curl \
  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
  --fail \
  --remote-name \
  https://domonda.app/api/public/document/00000000-0000-0000-0000-000000000000.pdf
```

To add an audit trail to the PDF or request only the audit trail,
use the following values for the query parameter `auditTrail`:

  - `append`: append the audit trail to the end of the document
  - `prepend`: prepend the audit trail before the first page of the document
  - `configured`: use the audit trail position configured for the client company
  - `only`: download only the audit trail as PDF

The default language of the audit trail is German.
To request the audit trail in a different language,
use the query parameter `auditTrailLang` with a language code.

Example URL to download only the audit trail in English:

```
https://domonda.app/api/public/document/00000000-0000-0000-0000-000000000000.pdf?auditTrail=only&auditTrailLang=en
```

### File uploads

File uploads are not using GraphQL, but Multipart MIME HTTP POST requests to the following URL:
`https://domonda.app/api/public/upload`

Note that basic document processing like creating or fixing a PDF file
and rendering page images is done synchronously and
may take up to 5 seconds per page, so adjust timeouts accordingly.

Extracting invoice data is done asynchronously by default.
Set the form field `waitForExtraction` to `true` for synchronous
extraction, where results will be available via GraphQL directly after
the upload request returns (add another 30 seconds to timeouts).

To identify the category of the uploaded document, either the form field `documentCategory`
must be provided, or the form field `documentType` with the additional fields
`bookingType` and `bookingCategory` if their value for the category is non-null/empty.
A combination of `documentType`, `bookingType`, `bookingCategory` uniquely identifies
a document category and may be easier to use than querying document category IDs upfront.

Valid values for `documentType` are:

```txt
INCOMING_INVOICE
OUTGOING_INVOICE
INCOMING_DUNNING_LETTER
OUTGOING_DUNNING_LETTER
INCOMING_DELIVERY_NOTE
OUTGOING_DELIVERY_NOTE
BANK_STATEMENT
CREDITCARD_STATEMENT
FACTORING_STATEMENT
OTHER_DOCUMENT
```

Valid values for `bookingType` are either an empty string (or not provided at all) or:

```txt
CASH_BOOK
CLEARING_ACCOUNT
```

`bookingCategory` is a generic string that may be empty or not provided at all.

Document categories can be queried via GraphQL:
<https://domonda.github.io/api/doc/schema/documentcategory.doc.html>

Example GraphQL query:

```gql
{
  allDocumentCategories {
    nodes {
      rowId
      documentType
      bookingType
      bookingCategory
      description
      emailAlias
    }
  }
}
```

The form field `tag` can be used multiple times to add multiple tags to the document.

The form field `document` must contain a file that serves as the visual representation of the document
and must be one of the following formats: PDF, PNG, JPEG, TIFF

If there is already another document with an identical content hash of the uploaded file
then an error response with the HTTP status code `409: Conflict` is returned.
The body of the response will be a JSON object with the following format:

```json
{
  "error": "Duplicate document content",
  "detail": {
    "documentFileHash": "<SOME_HASH_STRING>",
    "duplicateDocumentIDs": ["00000000-0000-0000-0000-000000000000"]
  }
}
```

If content-hash duplicates are uploaded as overlapping requests where the first request
hasn't had enough time yet to create a document with an ID in the database
then the filename of the first still processing upload is returned as `processingFileName``:

```json
{
  "error": "Duplicate document content",
  "detail": {
    "documentFileHash": "<SOME_HASH_STRING>",
    "processingFileName": "ORIGINAL_FILENAME.pdf"
  }
}
```

By default, the content hash conflict check also includes documents that have been marked as deleted.
To exclude deleted documents from the check, add the form field `allowDuplicateDeleted` with the value `true`.
With `allowDuplicateDeleted=true`, a duplicate of an already deleted document can be uploaded
as a new document.

### File upload with Go SDK

See Go SDK with the package
[github.com/domonda/api/golang/domonda](https://pkg.go.dev/github.com/domonda/api/golang/domonda).

Using this Go package to make API requests has the benefit of basic
client side validation of the data before sending the request
and the client package will always be kept up to date
with the API server implementation.

Go function: https://pkg.go.dev/github.com/domonda/api/golang/domonda#UploadDocument


### Upload structured invoice data as JSON

The optional form field `invoice` contains a [JSON file](example/invoice.jsonc)
with the schema [invoice.schema.json](invoice.schema.json) (see https://json-schema.org/).

The JSONC format variant with comments is supported.

Example:

```jsonc
{
  // Optional string or null
  "confirmedBy": "My Custom CRM",
  // Optional string or null
  "partnerName": "Muster I AG",
  // Optional string or null
  "partnerVatId": "ATU10223006",
  // Optional string or null
  "invoiceNumber": "17",
  // Optional string or null
  "internalNumber": null,
  // Optional string with format "YYYY-MM-DD" or null
  "invoiceDate": "2020-10-07",
  // Optional string with format "YYYY-MM-DD" or null
  "dueDate": "2020-10-21",
  // Optional string or null
  "orderNumber": "Auftrag 2020/1234",
  // Optional string with format "YYYY-MM-DD" or null
  "orderDate": "2020-08-15",
  // Optional boolean, false is used when missing
  "creditMemo": false,
  // Optional number or null
  "net": 5610.5,
  // Optional number or null
  "total": 6732.6,
  // Optional number or null
  "vatPercent": 20,
  // Optional array of numbers or null
  "vatPercentages": [20, 20, 20],
  // Optional array of numbers or null
  "vatAmounts": [210, 900, 12.1],
  // Optional number or null
  "discountPercent": 0,
  // Optional string with format "YYYY-MM-DD" or null
  "discountUntil": null,
  // Optional object or null
  "costCenters": {
    // Cost-center "number" as key with net amount as value
    "1000": 1050,
    // Cost-center "number" as key with net amount as value
    "2000": 4500,
    // Cost-center "number" as key with net amount as value
    "9000": 60.5
  },
  // Optional string or null, 3 character ISO 4217 alphabetic code
  "currency": "EUR",
  // Optional number greater zero or null
  "conversionRate": 1,
  // Optional string with format "YYYY-MM-DD" or null
  "conversionRateDate": "2020-10-07",
  // Optional string or null
  "goodsServices": "Website Design",
  // Optional string with format "YYYY-MM-DD" or null, use for performance period
  "deliveredFrom": "2020-09-01",
  // Optional string with format "YYYY-MM-DD" or null, use as single delivery date
  "deliveredUntil": "2020-09-30",
  // Optional string array
  "deliveryNoteNumbers": ["D12345"],
  // Optional string or null
  "iban": "DE02120300000000202051",
  // Optional string or null
  "bic": "BYLADEM1001",
  // Optional array of accounting-items or null
  "accountingItems": [
    {
      // Required string
      "title": "Test",
      // Required string
      "generalLedgerAccountNumber": 1000,
      // Required enum
      "bookingType": "DEBIT", // or CREDIT
      // Required enum
      "amountType": "NET", // or TOTAL
      // Required number
      "amount": 5000,
      // Optional UUID or null
      "valueAddedTax": "e77d686e-92f2-4c96-a5c1-b7c912327e90", // See: vat-codes-and-percentages.csv
      // Optional number or null
      "valueAddedTaxPercentageAmount": 20
    }
  ]
}
```

If `confirmedBy` is set to a non-empty string, then all values from the JSON
will be marked as confirmed and not overwritten by values from Domonda's automated invoice data extraction.
Upload API confirmations can be overwritten by users of the domonda app if they have sufficient rights.

The optional form field `ebInterface` contains an XML file in the ebInterface 5.0 format as specified at:
<https://www.wko.at/service/netzwerke/ebinterface-aktuelle-version-xml-rechnungsstandard.html>

Reference XML files can be created online at: <https://formular.ebinterface.at/>

Example using the CURL command-line tool with a `documentCategory` ID and multiple `tag` fields and a user-defined `uuid` for the document that must not exist in domonda yet:

```sh
curl -X POST \
  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
  -H "Content-Type: multipart/form-data" \
  -F "uuid=01505320-42f7-4cff-a930-4669eeb5e999"
  -F "documentCategory=fe110406-e38d-416a-a8d8-29f0a20f1c8d" \
  -F "document=@example/invoice.pdf" \
  -F "invoice=@example/invoice.jsonc" \
  -F "tag=TagA" \
  -F "tag=TagB" \
  -F "allowDuplicateDeleted=true" \
  -o - \
  https://domonda.app/api/public/upload
```

Example with `documentType`, `bookingType`, `bookingCategory`, and `waitForExtraction`:

```sh
curl -X POST \
  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
  -H "Content-Type: multipart/form-data" \
  -F "documentType=INCOMING_INVOICE" \
  -F "bookingType=CLEARING_ACCOUNT" \
  -F "bookingCategory=VKxx" \
  -F "document=@example/invoice.pdf" \
  -F "ebInterface=@example/invoice.xml" \
  -F "allowDuplicateDeleted=false" \
  -F "waitForExtraction=true" \
  -o - \
  https://domonda.app/api/public/upload
```

Example with just `documentType` (`bookingType` and `bookingCategory` would be null in the GraphQL query for the document category)
and `waitForExtraction`:

```sh
curl -X POST \
  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
  -H "Content-Type: multipart/form-data" \
  -F "documentType=OUTGOING_INVOICE" \
  -F "document=@example/invoice.pdf" \
  -F "ebInterface=@example/invoice.xml" \
  -F "allowDuplicateDeleted=false" \
  -F "waitForExtraction=true" \
  -o - \
  https://domonda.app/api/public/upload
```

The response will be an HTTP status code 200 message with the created document's UUID in plaintext format:

```txt
HTTP/1.1 200 OK
Content-Type: text/plain

ef059fa4-7288-4b77-8017-adce142e29a8
```

The document ID will be a new v4 random UUID, except if a user-defined `uuid` that does not exist yet in domonda is passed as form field. 

This UUID can be used in the GraphQL document API:
<https://domonda.github.io/api/doc/schema/document.doc.html>

### File upload HTTP response status codes

In case of errors, standard 4xx and 5xx HTTP status code responses will be returned with plaintext error messages in the body.

- 200 OK: Upload was successful
- 400 Bad Request: The request contains invalid data
- 401 Unauthorized: Invalid API key
- 402 Payment Required: The client company is not active or was blocked because of missing payments
- 409 Conflict: A document with the same UUID or file contents already exists
- 500 Internal Server Error: Something unexpected went wrong

### Querying uploaded document data with GraphQL

Just checking if the automated extraction of the uploaded document data has finished
with the UUID returned as plaintext body from the upload request as `rowId`:

```gql
{
  documentByRowId(rowId: "cbd03cbe-5d2f-4f97-bf12-03f1481d6c41") {
    extracted
  }
}
```

Querying invoice data:

```gql
{
  documentByRowId(rowId: "cbd03cbe-5d2f-4f97-bf12-03f1481d6c41") {
    numPages
    tags
    extracted
    invoiceByDocumentRowId {
      partnerName
      invoiceNumber
      invoiceNumberConfirmedBy
      invoiceDate
      invoiceDateConfirmedBy
      dueDate
      dueDateConfirmedBy
      orderNumber
      orderNumberConfirmedBy
      orderDate
      orderDateConfirmedBy
      creditMemo
      net
      netConfirmedBy
      total
      totalConfirmedBy
      vatPercent
      vatPercentConfirmedBy
      vatPercentages
      discountPercent
      discountPercentConfirmedBy
      discountUntil
      discountUntilConfirmedBy
      currency
      currencyConfirmedBy
      conversionRate
      conversionRateDate
      conversionRateSource
      goodsServices
      goodsServicesConfirmedBy
      deliveredFrom
      deliveredFromConfirmedBy
      deliveredUntil
      deliveredUntilConfirmedBy
      iban
      bic
    }
  }
}
```

### Upload company master data as JSON

Updating the master data of a client-company can be done in bulk using the REST API
with JSON data. In general, those API endpoints follow an `upsert` logic,
meaning that if data records can be identified by an ID or name,
they will be updated with the provided data,
else new records will be inserted.

For those endpoints, we also provide a Go SDK with the package
[github.com/domonda/api/golang/domonda](https://pkg.go.dev/github.com/domonda/api/golang/domonda).

Using this Go package to make API requests has the benefit of basic
client-side validation of the data before sending the request
and the client package will always be kept up to date
with the API server implementation.

The server is also implemented with Go and uses the standard library
JSON parser, which is documented here: <https://pkg.go.dev/encoding/json>.

The struct types in `github.com/domonda/api/golang/domonda` define the JSON API
where every struct is mapped to a JSON object
with the names of the struct fields used as JSON object value names.
The Go JSON parser detects object values names case insensitive
and object values that are left out in the JSON will have
their types' Go zero value after parsing.

So for example, the following JSON:

```json
{
  "number": "1/200",
  "name": "Account for object 66",
  "category": null,
  "objectNo": "66"
}
```

...can be parsed as `GLAccount` struct:

```go
type GLAccount struct {
	Number   account.Number
	Name     nullable.TrimmedString
	Category nullable.TrimmedString
	ObjectNo account.NullableNumber
}
```

The types used in those structs are from the package [github.com/domonda/go-types](https://pkg.go.dev/github.com/domonda/go-types) and should be self-explanatory.

Just note that `account.Number` is not an integer number,
but an alphanumeric string which has to match the regular expression
`^[0-9A-Za-z][0-9A-Za-z_\-\/:.;,]*$`



#### POST General Ledger Accounts

Go function: https://pkg.go.dev/github.com/domonda/api/golang/domonda#PostGLAccounts

API endpoint: https://domonda.app/api/public/masterdata/gl-accounts

Optional URL query parameters:
* `source`: string describing the source of the data; use your company or service name
* `objectSpecificAccountNos`: pass `true` if the account numbers are specific for objects, meaning that an account number can mean something different per object

Example: `https://domonda.app/api/public/masterdata/gl-accounts?source=MyCompany&objectSpecificAccountNos=true`

The request body is JSON array with objects matching the struct:

```go
type GLAccount struct {
	Number   account.Number         // Alphanumeric account number
	Name     nullable.TrimmedString // Name of the account
	Category nullable.TrimmedString // Higher level description of the account
	ObjectNo account.NullableNumber // Real estate object number
}
```

Example:

```json
[
  {
    "Number": "1/200",
    "Name": "Account for object 66",
    "Category": null,
    "ObjectNo": "66"
  },
  {
    "Number": "9/100",
    "Name": "Other costs",
  }
]
```

#### POST Partner Companies

Existing partner companies are identified by VATIDNo, VendorAccountNumber,
ClientAccountNumber, and Name and then updated.
Else, a new partner company is created.

Go function: https://pkg.go.dev/github.com/domonda/api/golang/domonda#PostPartners

API endpoint: https://domonda.app/api/public/masterdata/partner-companies

Optional URL query parameters:
  - `failOnInvalid`: If true, the function will fail if any partner data is invalid.
  - `useCleanedInvalid`: If true, the function will clean invalid partner data and use what's left. If false, invalid partner will be skipped.
  - `allOrNone`: If true, the function will import either all partners or none in case of any error.
  - `source`: The source of the partner data; use your company or service name.

The request body is JSON array with objects matching the struct:

```go
type Partner struct {
	Name             notnull.TrimmedString
	AlternativeNames notnull.StringArray // used when merging

	// main location
	Street    nullable.TrimmedString
	City      nullable.TrimmedString
	ZIP       nullable.TrimmedString
	Country   country.NullableCode
	Phone     nullable.TrimmedString
	Email     email.NullableAddress
	Website   nullable.TrimmedString
	CompRegNo nullable.TrimmedString
	TaxIDNo   nullable.TrimmedString
	VATIDNo   vat.NullableID

	// partner accounts
	VendorAccountNumber account.NullableNumber // "" means not set -> will not create a partner account
	ClientAccountNumber account.NullableNumber // "" means not set -> will not create a partner account

	// A single payment bank account for the partner.
	// IBAN and BIC are well suited as CSV columns.
	IBAN bank.NullableIBAN
	BIC  bank.NullableBIC
	// More payment bank accounts for the partner.
	// As struct better suited for JSON import.
	BankAccounts []bank.Account
}
```

The response is a JSON array with one object per request partner JSON object
with the following object members.

How the input from the request was normalized and used:
  - `NormalizedInput` (how the request input was normalized and used)
  - `InputWarnings` (warnings from validation and normalization)

Created or updated partner data, only present when data exists:
  - `PartnerCompany`
  - `PartnerLocations`
  - `VendorAccount`
  - `ClientAccount`
  - `PaymentPresets`

Result on partner level:
  - `State` (one of `"UNCHANGED"`, `"UPDATED"`, `"CREATED"`, `"ERROR"`)
  - `Error` (error message in case of `"State": "ERROR"` )

Example:

```sh
curl -X POST \
  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
  -H "Content-Type: application/json" \
  --data '[{"Name":"Beispiel GmbH","Street":"Ringstrasse","ZIP":"10115","City":"Berlin","Country":"DE","Email":"contact@example.com","ClientAccountNumber":"K1","VendorAccountNumber":"L1"}]' \
  -o - \
  https://domonda.app/api/public/masterdata/partner-companies
```


#### POST Bank Accounts

Existing bank accounts are identified by their IBAN and then updated instead of inserted.

Go function: https://pkg.go.dev/github.com/domonda/api/golang/domonda#PostBankAccounts

API endpoint: https://domonda.app/api/public/masterdata/bank-accounts

Optional URL query parameters:
* `source`: string describing the source of the data; use your company or service name

The request body is JSON array with objects matching the struct:

```go
type BankAccount struct {
	IBAN     bank.IBAN
	BIC      bank.BIC
	Currency money.Currency
	Holder   notnull.TrimmedString

	// Optional
	AccountNumber nullable.TrimmedString
	Name          nullable.TrimmedString
	Description   nullable.TrimmedString
}
```

#### POST Real Estate Objects

The following endpoint updates or inserts instances of the `RealEstateObject` class
using the `Number` property as the identifier for objects.

Go function: https://pkg.go.dev/github.com/domonda/api/golang/domonda#PostRealEstateObjects

API endpoint: https://domonda.app/api/public/masterdata/real-estate-objects

Optional URL query parameters:
* `source`: string describing the source of the data; use your company or service name

The request body is a JSON array of objects matching the struct:

```go
type RealEstateObject struct {
	Type                 RealEstateObjectType
	Number               account.Number
	AccountingArea       account.NullableNumber
	UserAccount          account.NullableNumber
	Description          nullable.TrimmedString
	StreetAddress        notnull.TrimmedString
	AlternativeAddresses nullable.StringArray
	ZipCode              nullable.TrimmedString
	City                 nullable.TrimmedString
	Country              country.Code
	BankAccounts         []bank.Account
	Active               bool
}
```

`RealEstateObjectType` is an enum with the following string values:

```go
type RealEstateObjectType string //#enum

const (
	RealEstateObjectTypeHI      RealEstateObjectType = "HI"
	RealEstateObjectTypeWEG     RealEstateObjectType = "WEG"
	RealEstateObjectTypeSUB     RealEstateObjectType = "SUB"
	RealEstateObjectTypeKREIS   RealEstateObjectType = "KREIS"
	RealEstateObjectTypeMANDANT RealEstateObjectType = "MANDANT"
)
```

#### POST Custom Object Class Instances

The following endpoint updates or inserts instances of a custom object class
using the property `idPropName` as the identifier for the objects.

Go function: https://pkg.go.dev/github.com/domonda/api/golang/domonda#PostObjectInstancesWithIDProp

API endpoint schema: https://domonda.app/api/public/masterdata/upsert-objects/{className}/id-prop/{idPropName}

Replace `{className}` with the name of the custom object class and
`{idPropName}` with the name of the property that is used as the identifier for the objects.

Optional URL query parameters:
* `source`: string describing the source of the data; use your company or service name

The request body is a JSON array of objects with keys matching the property names of the class.

The properties of a class can be queried with the GraphQL API using the following query
(replace `RealEstateObject` with the name of your custom class):

```graphql
{
  objectClassProps(className: "RealEstateObject") {
    nodes {
      rowId
      name
      type
      required
      options
      description
    }
  }
}
```
