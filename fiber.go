package eznodefiber

import (
	"strconv"
	"time"

	"github.com/amovah/eznode"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// RegisterRouter register a router on instance of fiber for proxy request to eznode and send response to user
// it registers a route with prefix which it accepts request and routes it to eznode
// for example: RegisterRoute(app, e, "/test-eznode", logrus.DebugLevel), it does register
// `/test-eznode/:chainId` on fiber app instance. request can be sent to desired chain with this route
// for example: GET `/test-eznode/ethereum`
// IMPORTANT NOTE: chainId is what you set in your eznode instance
func RegisterRoute(app *fiber.App, e *eznode.EzNode, routePrefix string, logLevel logrus.Level) {
	logger := logrus.New()
	logger.SetLevel(logLevel)

	app.All(routePrefix+"/:chainId/*", func(c *fiber.Ctx) error {
		return handleRequest(c, c.Method(), e, logger)
	})
}

// DisableNodeRequest disable node request body
type DisableNodeRequest struct {
	ChainId  string `json:"chain_id"`
	NodeName string `json:"node_name"`
	WithTime int    `json:"with_time"`
}

// EnableNodeRequest enable node request body
type EnableNodeRequest struct {
	ChainId  string `json:"chain_id"`
	NodeName string `json:"node_name"`
}

// StartFiber start a fiber app with a port to listen on
// default routes:
//   - ANY `/:chainId`: route to connect to eznode with chainId, you request to your specific node
//     with the chainId which you have set in the eznode instance. for example: GET `/ethereum`
//     IMPORTANT NOTE: chainId is what you set in your eznode instance
//   - POST `/manage/disable-node`: disable a node, accept DisableNodeRequest as body
//   - POST `/manage/enable-node`: enable a node, accept EnableNodeRequest as body
func StartFiber(port int, e *eznode.EzNode, logLevel logrus.Level) {
	logger := logrus.New()
	logger.SetLevel(logLevel)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler(logger),
	})

	app.Post("/manage/disable-node", func(c *fiber.Ctx) error {
		body := &DisableNodeRequest{}
		if err := c.BodyParser(body); err != nil {
			logger.Debugln("failed to parse body", err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if body.WithTime > 0 {
			e.DisableNodeWithTime(body.ChainId, body.NodeName, time.Duration(body.WithTime)*time.Minute)
		} else {
			e.DisableNode(body.ChainId, body.NodeName)
		}

		return c.SendString("done")
	})

	app.Post("/manage/enable-node", func(c *fiber.Ctx) error {
		body := &EnableNodeRequest{}
		if err := c.BodyParser(body); err != nil {
			logger.Debugln("failed to parse body", err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		e.EnableNode(body.ChainId, body.NodeName)

		return c.SendString("done")
	})

	app.All("/:chainId/*", func(c *fiber.Ctx) error {
		return handleRequest(c, c.Method(), e, logger)
	})

	logger.Fatal(app.Listen(":" + strconv.Itoa(port)))
}
