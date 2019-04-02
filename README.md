# DOMONDA API

The general Domonda API is implemented using the GraphQL protocol: <http://graphql.org/>

The only exception are file uploads which are implemented as Multipart MIME HTTP POST requests.

A static GraphQL documentation is available at: <https://domonda.github.io/api/>

For interactive access and documentation use the webbased GraphiQL: <https://app.domonda.com/api/public/graphiql>

You can provide a authentication token for authenticated access to your Domonda data.
Without a valid token, demo data is provived by the API.

Alternatively you can use the desktop client Altair (<https://altair.sirmuel.design/>).


## Authentication

Authentication is implemented via Bearer token. Replace `API_KEY` with your companies's API key. If you don't have one, request it from api@domonda.com

```http
POST https://app.domonda.com/api/public/graphql

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
`https://app.domonda.com/api/public/upload`

Note that currently all document processing is done synchronously.
OCR and PDF processing may take up to 5 seconds per page, so adjust timeouts accordingly.

To identify the category of the uploaded document, either the form filed `documentCategory`
must be provided or alternatively the form field `documentType` with the additional fields
`bookingType` and `bookingCategory` if their value for the category is non null.
A combination of `documentType`, `bookingType`, `bookingCategory` uniquely identifies
a document category and may be easier to use than querieng document category IDs upfront.

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

The form field `document` must contain a file that serves as the visual representation of the document
and must be one of the following formats: PDF, PNG, JPEG, TIFF

Feature coming soon:

The optional form field `ebInterface` contains an XML file in the ebInterface 5.0 format as specified at:
<https://www.wko.at/service/netzwerke/ebinterface-aktuelle-version-xml-rechnungsstandard.html>

Reference XML files can be created online at: <https://formular.ebinterface.at/>

Example using the CURL command line tool with a `documentCategory` ID:

```sh
curl -X POST \
  -H "Authorization: Bearer API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F "documentCategory=fe110406-e38d-416a-a8d8-29f0a20f1c8d" \
  -F "document=@invoice.pdf" \
  -F "ebInterface=@invoice.xml" \
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
  -F "document=@invoice.pdf" \
  -F "ebInterface=@invoice.xml" \
  https://app.domonda.com/api/public/upload
```

Example with just `documentType` (`bookingType` and `bookingCategory` would be null in the GraphQL query for the document category):

```sh
curl -X POST \
  -H "Authorization: Bearer API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F "documentType=OUTGOING_INVOICE" \
  -F "document=@invoice.pdf" \
  -F "ebInterface=@invoice.xml" \
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
