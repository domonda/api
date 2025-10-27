package domonda

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ungerik/go-fs"

	"github.com/domonda/go-types/uu"
)

// UploadDocument uploads a document file (PDF, PNG, JPEG, or TIFF) to create a new document in Domonda.
// The document will be processed synchronously (creating/fixing PDF, rendering page images).
// Invoice data extraction happens asynchronously by default.
//
// Arguments:
//   - ctx:              Context for the HTTP request (for cancellation and timeouts)
//   - apiKey:           API key (bearer token) for authentication
//   - documentCategory: UUID of the document category (query via GraphQL allDocumentCategories)
//   - documentFile:     Document file to upload (PDF, PNG, JPEG, or TIFF format)
//   - invoiceFile:      Optional JSON file with structured invoice data (can be nil)
//   - tags:             Optional tags to attach to the document
//
// Returns:
//   - documentID: UUID of the created document
//   - err:        Error if upload fails, including HTTP status errors (409 for duplicates)
//
// The function uses a multipart form POST request to https://domonda.app/api/public/upload
//
// Note: Basic document processing may take up to 5 seconds per page.
// For synchronous invoice extraction, use the REST API directly with waitForExtraction=true.
func UploadDocument(ctx context.Context, apiKey string, documentCategory uu.ID, documentFile, invoiceFile fs.FileReader, tags ...string) (documentID uu.ID, err error) {
	body := bytes.NewBuffer(nil)
	form := multipart.NewWriter(body)

	err = form.WriteField("documentCategory", documentCategory.String())
	if err != nil {
		return uu.IDNil, err
	}

	for _, tag := range tags {
		err = form.WriteField("tag", tag)
		if err != nil {
			return uu.IDNil, err
		}
	}

	documentWriter, err := form.CreateFormFile("document", documentFile.Name())
	if err != nil {
		return uu.IDNil, err
	}
	_, err = documentFile.WriteTo(documentWriter)
	if err != nil {
		return uu.IDNil, err
	}

	if invoiceFile != nil {
		invoiceWriter, err := form.CreateFormFile("invoice", invoiceFile.Name())
		if err != nil {
			return uu.IDNil, err
		}
		_, err = invoiceFile.WriteTo(invoiceWriter)
		if err != nil {
			return uu.IDNil, err
		}
	}

	err = form.Close()
	if err != nil {
		return uu.IDNil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", BaseURL+"/upload", body)
	if err != nil {
		return uu.IDNil, err
	}
	request.Header.Add("Content-Type", form.FormDataContentType())
	request.Header.Add("Authorization", "Bearer "+apiKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return uu.IDNil, err
	}

	if response.StatusCode != 200 {
		return uu.IDNil, fmt.Errorf("%d: %s", response.StatusCode, response.Status)
	}

	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return uu.IDNil, err
	}

	return uu.IDFromBytes(responseBody)
}
