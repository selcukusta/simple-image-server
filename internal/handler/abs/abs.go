package abs

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/model"
	"github.com/valyala/fasthttp"
)

//Handler is using connect to Azure Blob Storage and get the image
func Handler(ctx *fasthttp.RequestCtx, vars map[string]string) {
	path := vars["path"]

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	credential, err := azblob.NewSharedKeyCredential(connection.AccountName, connection.AccountKey)
	if err != nil {
		customError := model.CustomError{Message: "Unable to create shared keys for Azure Blob Storage connection", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	URL, err := url.Parse(fmt.Sprintf("https://%s.%s/", connection.AccountName, connection.AzureURI))
	if err != nil {
		customError := model.CustomError{Message: "Unable to parse  Azure Blob Storage URI", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	containerURL := azblob.NewContainerURL(*URL, azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			TryTimeout: 10 * time.Second,
		},
	}))
	blobURL := containerURL.NewBlockBlobURL(path)

	downloadResponse, err := blobURL.Download(context, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		customError := model.CustomError{Message: "Blob cannot be downloaded", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	var buf bytes.Buffer
	reader := downloadResponse.Body(azblob.RetryReaderOptions{})
	_, err = buf.ReadFrom(reader)
	if err != nil {
		customError := model.CustomError{Message: "Buffer can not be read", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}
	defer reader.Close()

	headers := make(map[string]string)
	headers["ETag"] = string(downloadResponse.ETag())
	headers["Last-Modified"] = downloadResponse.LastModified().Format(time.RFC1123)
	model.SucceededFinalizer{ResponseWriter: ctx, ContentType: downloadResponse.ContentType()}.Finalize(vars, buf.Bytes())
}
