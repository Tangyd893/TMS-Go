package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/file/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	fileService service.FileService
}

func NewHandler(svc service.FileService) *Handler {
	return &Handler{fileService: svc}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Fail(w, r, 400, 400001, "上传文件过大", nil); return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Fail(w, r, 400, 400001, "请选择文件", nil); return
	}
	defer file.Close()

	bizType := r.FormValue("bizType")
	bizID := r.FormValue("bizId")
	userID := middleware.UserIDFromContext(r.Context())

	fileObj, err := h.fileService.Upload(r.Context(), file, header.Filename, bizType, bizID, userID)
	if err != nil {
		apErr := unwrap(err)
		response.Fail(w, r, 400, apErr.Code, apErr.Message, nil); return
	}

	response.Success(w, r, fileObj)
}

func (h *Handler) ListByBiz(w http.ResponseWriter, r *http.Request) {
	bizType := r.URL.Query().Get("bizType")
	bizID := r.URL.Query().Get("bizId")
	if bizType == "" || bizID == "" {
		response.Fail(w, r, 400, 400001, "缺少bizType或bizId参数", nil); return
	}
	files, err := h.fileService.ListByBiz(r.Context(), bizType, bizID)
	if err != nil { response.Fail(w, r, 500, 500001, err.Error(), nil); return }
	response.Success(w, r, files)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" { response.Fail(w, r, 400, 400001, "缺少文件ID", nil); return }
	if err := h.fileService.Delete(r.Context(), id); err != nil {
		response.Fail(w, r, 500, unwrap(err).Code, err.Error(), nil); return
	}
	response.Success(w, r, nil)
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" { response.Fail(w, r, 400, 400001, "缺少文件ID", nil); return }

	reader, fileObj, err := h.fileService.GetDownloadURL(r.Context(), id)
	if err != nil {
		response.Fail(w, r, 404, unwrap(err).Code, err.Error(), nil); return
	}
	if reader != nil {
		defer reader.Close()
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileObj.OriginalName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	if reader != nil {
		io.Copy(w, reader)
	}
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /api/v1/files/upload", handler.Upload)
	mux.HandleFunc("GET /api/v1/files", handler.ListByBiz)
	mux.HandleFunc("GET /api/v1/files/{id}/download", handler.Download)
	mux.HandleFunc("DELETE /api/v1/files/{id}", handler.Delete)
}

func extractID(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx < 0 || idx == len(path)-1 { return "" }
	return path[idx+1:]
}

func unwrap(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok { return e }
	return errors.New(500001, err.Error())
}
