# DOMONDA API

The Domonda API is implemented using the GraphQL protocol: http://graphql.org/

For interactive access and documentation use GraphiQL.app: https://github.com/skevy/graphiql-app

## Authentication

Authentication is implemented via Bearer token. Replace `API_KEY` with your companies's API key. If you don't have one, request it from erik@domonda.com

```http
POST https://app.domonda.com/api/public/graphql

Authorization: Bearer API_KEY
```

For a full API specification see:

* [domonda-api.gql](domonda-api.gql)
* [domonda-api.json](domonda-api.json)

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

Query delivery note with note items:

```gql
{
  deliveryNoteByDocumentId(id: "00d3d910-d6e8-47eb-b6f4-9aae5daf44b1") {
    documentId,
    type,
    version,
    partnerCompanyId,
    invoiceId,
    invoiceNumber,
    deliveryNoteNr,
    deliveryDate,
    noteItems {
      nodes {
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
}
```