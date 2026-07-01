package handler

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

type UploadHandler struct {
	uploadSvc *services.UploadService
}

func NewUploadHandler(uploadSvc *services.UploadService) *UploadHandler {
	return &UploadHandler{uploadSvc: uploadSvc}
}

func (h *UploadHandler) UploadFile(ctx *gin.Context) {
	// Limit size of the uploaded file to 32 MB
	ctx.Request.ParseMultipartForm(32 << 20) // 32 MB

	file, header, err := ctx.Request.FormFile("file") // "file" is name of <input type="file">
	if err != nil {
		logger.Log.WithError(err).Error("UploadFile: Failed to get form file")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "")
		return
	}
	defer file.Close()

	// file this is io.Reader (multipart.File)
	var reader io.Reader = file

	keyName, err := h.uploadSvc.UploadFile(ctx.Request.Context(), header.Filename, reader, header.Header.Get("Content-Type"))
	if err != nil {
		logger.Log.WithError(err).Error("UploadFile: Failed to upload file")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, "upload failed")
		return
	}
	data := &domain.UploadFileResponse{
		FileKey: keyName,
		URL:     h.uploadSvc.PublicURL(keyName, ""),
	}

	HandleSuccess(ctx, data, "Upload file successfully!")
}

func (h *UploadHandler) DeleteFile(ctx *gin.Context) {
	fileKey := ctx.Query("fileKey")
	if fileKey == "" {
		logger.Log.Error("DeleteFile: fileKey is required")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "fileKey is required")
		return
	}

	err := h.uploadSvc.DeleteFile(ctx.Request.Context(), fileKey)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteFile: failed to delete file")
		HandleError(ctx, domain.ErrorCodeInternalServerError, nil, "delete failed")
		return
	}

	HandleSuccess(ctx, nil, "Delete file successfully!")
}

func (h *UploadHandler) UploadFileTinyEditor(ctx *gin.Context) {
	// Limit size of the uploaded file to 32 MB
	ctx.Request.ParseMultipartForm(32 << 20) // 32 MB

	file, header, err := ctx.Request.FormFile("file") // "file" is name of <input type="file">
	if err != nil {
		logger.Log.WithError(err).Error("UploadFileTinyEditor: Failed to get form file")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "")
		return
	}
	defer file.Close()

	// file this is io.Reader (multipart.File)
	var reader io.Reader = file

	keyName, err := h.uploadSvc.UploadFileWithBucket(ctx.Request.Context(), config.LoadConfig().TinyEditor.Bucket, header.Filename, reader, header.Header.Get("Content-Type"))
	if err != nil {
		logger.Log.WithError(err).Error("UploadFileTinyEditor: Failed to upload file")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "upload failed")
		return
	}
	data := &domain.UploadFileResponseTinyEditor{
		FileKey: keyName,
		URL:     h.uploadSvc.PublicURL(keyName, ""),
	}

	HandleSuccess(ctx, data, "Upload file successfully!")
}

func (h *UploadHandler) DeleteFileTinyEditor(ctx *gin.Context) {
	fileKey := ctx.Query("fileKey")
	if fileKey == "" {
		logger.Log.Error("DeleteFileTinyEditor: fileKey is required")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "fileKey is required")
		return
	}

	err := h.uploadSvc.DeleteFileWithBucket(ctx.Request.Context(), config.LoadConfig().TinyEditor.Bucket, fileKey)
	if err != nil {
		logger.Log.WithError(err).Error("DeleteFileTinyEditor: failed to delete file")
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, "delete failed")
		return
	}

	HandleSuccess(ctx, nil, "Delete file successfully!")
}
