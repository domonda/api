# DOMONDA API

The Domonda API is implemented using the GraphQL protocol: http://graphql.org/

For interactive access and documentation use web based GraphiQL.app (<https://github.com/skevy/graphiql-app>) or 
the desktop application Altair (<https://altair.sirmuel.design/>).


## Authentication

Authentication is implemented via Bearer token. Replace `API_KEY` with your companies's API key. If you don't have one, request it from erik@domonda.com

```http
POST https://app.domonda.com/api/public/graphql

Authorization: Bearer API_KEY
```

For a full API specification see:

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
