package domonda

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/domonda/go-types/uu"
	"github.com/ungerik/go-fs"
)

const (
	UploadURL = "https://domonda.app/api/public/upload"
)

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

	request, err := http.NewRequestWithContext(ctx, "POST", UploadURL, body)
	if err != nil {
		return uu.IDNil, err
	}

	request.Header.Add("Authorization", "Bearer "+apiKey)
	request.Header.Add("Content-Type", "multipart/form-data")

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
