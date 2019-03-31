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

## File uploads

File uploads are not using GraphQL, but Multipart MIME HTTP POST requests to the following URL:
https://app.domonda.com/api/public/upload

The form field `documentCategory` must contain the id of a document category that can be requested at:
https://domonda.github.io/api/doc/schema/documentcategory.doc.html

The form field `document` must contain a file that serves as the visual representation of the document
and must be one of the following formats: PDF, PNG, JPEG, TIFF

The optional form field `ebInterface` contains an XML file in the ebInterface 5.0 format as specified at:
https://www.wko.at/service/netzwerke/ebinterface-aktuelle-version-xml-rechnungsstandard.html

Example using the CURL command line tool:

```sh
curl -X POST \
  -H "Authorization: Bearer API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F "documentCategory=fe110406-e38d-416a-a8d8-29f0a20f1c8d" \
  -F "document=@invoice.pdf" \
  -F "ebInterface=@invoice.xml" \
  https://app.domonda.com/api/public/upload
```

The response will be a HTTP status code 200 message with the created document's UUID in plaintext format:

```http
HTTP/1.1 200 OK
Content-Type: text/plain

ef059fa4-7288-4b77-8017-adce142e29a8
```

This UUID can be used in the GraphQL document API:
https://domonda.github.io/api/doc/schema/document.doc.html

In case of an error, standard 4xx and 5xx HTTP status code responses will be returned with plaintext error messages in the body.


## GraphQL API specification

* [domonda-api.gql](domonda-api.gql)
* [domonda-api.json](domonda-api.json)


## Basic usage

You can find all GraphQL query types in the generated documentation: <https://domonda.github.io/api/doc/schema/query.doc.html>.

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


## Example queries

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
