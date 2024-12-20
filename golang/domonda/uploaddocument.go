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

// UploadDocument uploads a file for a new document
// with a multipart form POST request to the endpoint https://domonda.app/api/public/upload
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
