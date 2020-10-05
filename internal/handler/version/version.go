package version

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/logger"
	"github.com/valyala/fasthttp"
)

//Handler is using connect to MongoDB and get the image
func Handler(ctx *fasthttp.RequestCtx) {
	version := "Unknown"
	if os.Getenv("APP_VERSION") != "" {
		version = os.Getenv("APP_VERSION")
	}
	ctx.Response.Header.Set("Content-Type", "text/plain")
	_, err := ctx.Write([]byte(version))
	if err != nil {
		msg := fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error())
		logger.InitExceptionWithRequest(ctx, errors.WithStack(err)).Error(msg)
	}
}
