package apicommon

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/render"

	"github.com/CognyHub/evmd-golib/logger"
)

type ErrResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	StatusText     string `json:"status"`          // user-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	logger.Errorln(e.Err)
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func Err(err error, statusCode int) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: statusCode,
		StatusText:     http.StatusText(statusCode),
		ErrorText:      err.Error(),
	}
}

func RenderError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	errResp := &ErrResponse{
		Err:            err,
		HTTPStatusCode: statusCode,
		StatusText:     http.StatusText(statusCode),
		ErrorText:      err.Error(),
	}

	logger.Errorln(err)

	switch render.GetAcceptedContentType(r) {
	case render.ContentTypeJSON:
		render.DefaultResponder(w, r, errResp)
	default:
		workDir, _ := os.Getwd()
		htmlfile := filepath.Join(workDir, "html/error.html")

		tmpl, err := template.ParseFiles(htmlfile)
		if err != nil {
			render.DefaultResponder(w, r, Err(fmt.Errorf("template.ParseFiles(): %w", err), http.StatusInternalServerError))
			return
		}

		if err := tmpl.Execute(w, errResp); err != nil {
			render.DefaultResponder(w, r, Err(fmt.Errorf("tmpl.Execute(): %w", err), http.StatusInternalServerError))
			return
		}
	}
}
