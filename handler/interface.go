package handler

import (
	"context"
	"net/http"
)

type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
}
