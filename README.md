![image](img/domonda_logo_schriftzug.png)

# DOMONDA API

The general Domonda API is implemented using the GraphQL protocol: <http://graphql.org/>

The only exception are file uploads which are implemented as Multipart MIME HTTP POST requests.

A static GraphQL documentation is available at: <https://domonda.github.io/api/>

For interactive access and documentation use the webbased GraphiQL: <https://domonda.app/api/public/graphiql>

You can provide a authentication token for authenticated access to your Domonda data.
Without a valid token, demo data is provived by the API.

Alternatively you can use the desktop client Altair (<https://altair.sirmuel.design/>).


## Authentication

Authentication is implemented via Bearer token. Replace `API_KEY` with your companies's API key. If you don't have one, request it from api@domonda.com

```http
POST https://domonda.app/api/public/graphql

Authorization: Bearer API_KEY
```


## Basic usage

You can find all GraphQL query types in the generated documentation:
<https://domonda.github.io/api/doc/schema/query.doc.html>

Referenced fields always have an extra field for querying the actual data behind it.

If you want to query all documents with the additional information of the import user you can achieve this by using the
field `userByImportedBy` which gets you the associated user.

```gql
{
  query{
    allDocuments{
      nodes{
        id
        userByImportedBy{
          firstName
          lastName
        }
      }
    }
  }
}
```


## File uploads

File uploads are not using GraphQL, but Multipart MIME HTTP POST requests to the following URL:
`https://domonda.app/api/public/upload`

Note that currently all document processing is done synchronously.
OCR and PDF processing may take up to 5 seconds per page, so adjust timeouts accordingly.

To identify the category of the uploaded document, either the form field `documentCategory`
must be provided or alternatively the form field `documentType` with the additional fields
`bookingType` and `bookingCategory` if their value for the category is non null/empty.
A combination of `documentType`, `bookingType`, `bookingCategory` uniquely identifies
a document category and may be easier to use than querying document category IDs upfront.

Valid values for documentType are:

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

Document categories can be queried with via GraphQL:
<https://domonda.github.io/api/doc/schema/documentcategory.doc.html>

Example GraphQL query:

```gql
{
  allDocumentCategories{
    nodes{
      id,
      documentType,
      bookingType,
      bookingCategory,
      description,
      emailAlias
    }
  }
}
```

The form field `tag` can be used multiple times to add multiple tags to the document.

The form field `document` must contain a file that serves as the visual representation of the document
and must be one of the following formats: PDF, PNG, JPEG, TIFF

The optional form field `invoice` contains a [JSON file](example/invoice.jsonc) with the following fields
(the JSONC format variant with comments is supported): 

```jsonc
{
    "confirmedBy": "My Custom CRM",     // Optional string or null
    "partnerName": "Muster I AG",       // Optional string or null
    "partnerVatId": "ATU10223006",      // Optional string or null
    "invoiceNumber": "17",              // Optional string or null
    "internalNumber": null,             // Optional string or null
    "invoiceDate": "2020-10-07",        // Optional string with format "YYYY-MM-DD" or null
    "dueDate": "2020-10-21",            // Optional string with format "YYYY-MM-DD" or null
    "creditMemo": false,                // Optional boolean, false is used when missing
    "net": 5610,                        // Optional number or null
    "total": 6732,                      // Optional number or null
    "vatPercent": 20,                   // Optional number or null
    "vatPercentages": [20, 20, 20],     // Optional array of numbers or null
    "vatAmounts": [210, 900, 12],       // Optional array of numbers or null
    "discountPercent": 0,               // Optional number or null
    "discountUntil": null,              // Optional string with format "YYYY-MM-DD" or null
    "currency": "EUR",                  // Optional string or null, 3 character ISO 4217 alphabetic code
    "conversionRate": 1,                // Optional number greater zero or null
    "conversionRateDate": "2020-10-07", // Optional string with format "YYYY-MM-DD" or null
    "goodsServices": "Website Design",  // Optional string or null
    "deliveredFrom": "2020-09-01",      // Optional string with format "YYYY-MM-DD" or null
    "deliveredUntil": "2020-09-30",     // Optional string with format "YYYY-MM-DD" or null
    "iban": "DE02120300000000202051",   // Optional string or null
    "bic": "BYLADEM1001"                // Optional string or null
}
```

If `confirmedBy` is set to a non-empty string then all values from the JSON
will be marked as confirmed and not overwritten by values from domonda's automated incoice data extraction.
The `confirmedBy` string will be shown next to confirmed values in the domonda web app,
so a representive, human readable string should be chosen.
Upload API confirmations can be changed by users of the domonda app, if they have sufficient rights.


The optional form field `ebInterface` contains a XML file in the ebInterface 5.0 format as specified at:
<https://www.wko.at/service/netzwerke/ebinterface-aktuelle-version-xml-rechnungsstandard.html>

Reference XML files can be created online at: <https://formular.ebinterface.at/>

Example using the CURL command line tool with a `documentCategory` ID and multiple `tag` fields:

```sh
curl -X POST \
  -H "Authorization: Bearer API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F "documentCategory=fe110406-e38d-416a-a8d8-29f0a20f1c8d" \
  -F "document=@example/invoice.pdf" \
  -F "tag=TagA" \
  -F "tag=TagB" \
  -F "ebInterface=@example/invoice.xml" \
  https://app.domonda.com/api/public/upload
```

Example with `documentType`, `bookingType`, `bookingCategory`:

```sh
curl -X POST \
  -H "Authorization: Bearer API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F "documentType=INCOMING_INVOICE" \
  -F "bookingType=CLEARING_ACCOUNT" \
  -F "bookingCategory=VKxx" \
  -F "document=@example/invoice.pdf" \
  -F "ebInterface=@example/invoice.xml" \
  https://app.domonda.com/api/public/upload
```

Example with just `documentType` (`bookingType` and `bookingCategory` would be null in the GraphQL query for the document category):

```sh
curl -X POST \
  -H "Authorization: Bearer API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F "documentType=OUTGOING_INVOICE" \
  -F "document=@example/invoice.pdf" \
  -F "ebInterface=@example/invoice.xml" \
  https://app.domonda.com/api/public/upload
```

The response will be a HTTP status code 200 message with the created document's UUID in plaintext format:

```txt
HTTP/1.1 200 OK
Content-Type: text/plain

ef059fa4-7288-4b77-8017-adce142e29a8
```

This UUID can be used in the GraphQL document API:
<https://domonda.github.io/api/doc/schema/document.doc.html>

In case of an error, standard 4xx and 5xx HTTP status code responses will be returned with plaintext error messages in the body.


## Example GraphQL queries

Query all document categories:

```gql
{
  allDocumentCategories{
    nodes{
      id,
      documentType,
      bookingType,
      bookingCategory,
      description,
      emailAlias
    }
  }
}
```

Query all documents:

```gql
{
  allDocuments {
    nodes {
        id,
        type,
        categoryId,
        workflowStepId,
        importDate,
        periodDate,
        name,
        title,
        language,
        version,
        numPages,
        numAttachPages
    }
  }
}
```

Query all invoices:

```gql
{
  allInvoices {
    nodes {
      documentId,
      type,
      title,
      version,
      invoiceNumber,
      invoiceDate,
      net,
      total,
      vatPercent,
      vatId,
      currency,
      iban,
      bic,
      dueDate,
      paymentStatus,
      paidDate,
    }
  }
}
```

Query all delivery notes:

```gql
{
  allDeliveryNotes {
    nodes {
      documentId,
      type,
      version,
      partnerCompanyId,
      invoiceId,
      invoiceNumber,
      deliveryNoteNr,
      deliveryDate,
    }
  }
}
```

Quary all delivery note items:

```gql
{
  allDeliveryNoteItems {
    nodes {
      documentId,
      deliveryNoteNr,
      deliveryDate,
      posNr,
      quantity,
      productNr,
      eanNr,
      gtinNr,
      description,
    }
  }
}
```
