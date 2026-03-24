package v1

import (
	"embed"

	"github.com/gofiber/fiber/v2"
)

var (
	//go:embed web/index.html
	webFiles embed.FS
)

func (r *V1) showUI(ctx *fiber.Ctx) error {
	file, err := webFiles.ReadFile("web/index.html")
	if err != nil {
		r.l.Error(err, "restapi - v1 - showUI")

		return errorResponse(ctx, fiber.StatusInternalServerError, "problems with load UI")
	}

	ctx.Set(fiber.HeaderContentType, "text/html")

	return ctx.Send(file)
}
