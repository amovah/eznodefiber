package eznodefiber

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/amovah/eznode"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func createRequest(
	uri []byte,
	method string,
	body *[]byte,
	header *fasthttp.RequestHeader,
) (*http.Request, error) {
	requestUrl := "/" + strings.Join(strings.Split(string(uri), "/")[2:], "/")
	request, err := http.NewRequest(method, requestUrl, bytes.NewBuffer(*body))
	if err != nil {
		return new(http.Request), err
	}

	request.Header.Set("Accept", string(header.Peek("Accept")))
	request.Header.Set("Accept-Encoding", string(header.Peek("Accept-Encoding")))
	request.Header.Set("Accept-Language", string(header.Peek("Accept-Language")))
	request.Header.Set("Content-Type", string(header.Peek("Content-Type")))

	return request, nil
}

func handleRequest(c *fiber.Ctx, httpMethod string, e *eznode.EzNode, logger *logrus.Logger) error {
	body := c.Request().Body()
	request, err := createRequest(c.Request().RequestURI(), httpMethod, &body, &c.Request().Header)
	if err != nil {
		return err
	}

	response, err := e.SendRequest(c.Params("chainId"), request)
	if err != nil {
		if nodeError, ok := err.(eznode.EzNodeError); ok {
			statusCode := fiber.StatusInternalServerError
			message := "Internal Server Error"

			for _, trace := range nodeError.Metadata.Trace {
				logger.
					WithField("request id", c.Context().ID()).
					WithField("chain id", nodeError.Metadata.ChainId).
					WithField("url", nodeError.Metadata.RequestedUrl).
					WithField("node name", trace.NodeName).
					WithField("status code", trace.StatusCode).
					Info(trace.Err)

				if trace.StatusCode > 0 {
					statusCode = trace.StatusCode
					if err != nil {
						message = trace.Err.Error()
					} else {
						message = http.StatusText(trace.StatusCode)
					}
				}
			}

			return fiber.NewError(statusCode, message)
		}

		return fiber.NewError(fiber.StatusNotAcceptable, fmt.Sprintf("%s is not valid", c.Params("chainId")))
	}

	for key, value := range *response.Headers {
		c.Response().Header.Set(key, strings.Join(value, ", "))
	}

	for _, trace := range response.Metadata.Trace {
		createdLog := logger.
			WithField("request id", c.Context().ID()).
			WithField("chain id", response.Metadata.ChainId).
			WithField("url", response.Metadata.RequestedUrl).
			WithField("node name", trace.NodeName).
			WithField("status code", trace.StatusCode)
		if trace.Err != nil {
			createdLog.Debug(trace.Err)
		} else {
			createdLog.Debug("respond successfully")
		}
	}

	c.Status(response.StatusCode)
	return c.Send(response.Body)
}
