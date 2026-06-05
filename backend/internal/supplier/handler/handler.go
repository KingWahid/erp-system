package handler

import (
	"fmt"
	"net/http"

	supplierusecase "erp-system/internal/supplier/usecase"
	"erp-system/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SupplierHandler struct {
	usecase supplierusecase.SupplierUsecase
}

func NewSupplierHandler(uc supplierusecase.SupplierUsecase) *SupplierHandler {
	return &SupplierHandler{usecase: uc}
}

// ListSuppliers godoc
// @Summary      List suppliers
// @Tags         suppliers
// @Security     BearerAuth
// @Produce      json
// @Param        search   query  string  false  "Search by name, code"
// @Param        status   query  string  false  "Filter by status"
// @Param        page     query  int     false  "Page number"
// @Param        per_page query  int     false  "Items per page"
// @Success      200  {object}  response.Response
// @Router       /suppliers [get]
func (h *SupplierHandler) ListSuppliers(c echo.Context) error {
	req := supplierusecase.ListRequest{
		Search:  c.QueryParam("search"),
		Status:  c.QueryParam("status"),
		Page:    queryInt(c, "page", 1),
		PerPage: queryInt(c, "per_page", 20),
	}

	res, err := h.usecase.ListSuppliers(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.SuccessWithMeta(c, res.Items, response.Meta{
		Page:       res.Page,
		PerPage:    res.PerPage,
		Total:      res.Total,
		TotalPages: int((res.Total + int64(res.PerPage) - 1) / int64(res.PerPage)),
	})
}

// GetStats godoc
// @Summary      Get supplier dashboard stats
// @Tags         suppliers
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /suppliers/stats [get]
func (h *SupplierHandler) GetStats(c echo.Context) error {
	stats, err := h.usecase.GetStats(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, stats)
}

// CreateSupplier godoc
// @Summary      Create new supplier
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  supplierusecase.CreateSupplierRequest  true  "Supplier payload"
// @Success      201   {object}  response.Response
// @Router       /suppliers [post]
func (h *SupplierHandler) CreateSupplier(c echo.Context) error {
	var req supplierusecase.CreateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	res, err := h.usecase.CreateSupplier(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, res)
}

// GetSupplier godoc
// @Summary      Get supplier by ID
// @Tags         suppliers
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Supplier UUID"
// @Success      200  {object}  response.Response
// @Router       /suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	res, err := h.usecase.GetSupplierByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, res)
}

// UpdateSupplier godoc
// @Summary      Update supplier
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                                 true  "Supplier UUID"
// @Param        body  body  supplierusecase.UpdateSupplierRequest  true  "Update payload"
// @Success      200   {object}  response.Response
// @Router       /suppliers/{id} [put]
func (h *SupplierHandler) UpdateSupplier(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	var req supplierusecase.UpdateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	res, err := h.usecase.UpdateSupplier(c.Request().Context(), id, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, res)
}

// DeleteSupplier godoc
// @Summary      Soft delete supplier
// @Tags         suppliers
// @Security     BearerAuth
// @Param        id  path  string  true  "Supplier UUID"
// @Success      204
// @Router       /suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	if err := h.usecase.DeleteSupplier(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// BlockSupplier godoc
// @Summary      Block or unblock supplier
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  string  true  "Supplier UUID"
// @Success      200
// @Router       /suppliers/{id}/block [post]
func (h *SupplierHandler) BlockSupplier(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	var req struct {
		Block  bool   `json:"block"`
		Reason string `json:"reason"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	var ucErr error
	if req.Block {
		ucErr = h.usecase.BlockSupplier(c.Request().Context(), id, req.Reason)
	} else {
		ucErr = h.usecase.UnblockSupplier(c.Request().Context(), id)
	}

	if ucErr != nil {
		return response.Error(c, ucErr)
	}

	return response.Success(c, map[string]bool{"blocked": req.Block})
}

// AdvanceStage godoc
// @Summary      Advance supplier to next workflow stage
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  string  true  "Supplier UUID"
// @Success      200   {object}  response.Response
// @Router       /suppliers/{id}/next-stage [post]
func (h *SupplierHandler) AdvanceStage(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	var req struct {
		Notes string `json:"notes"`
	}
	_ = c.Bind(&req)

	changedBy, _ := c.Get("user_id").(string)

	res, err := h.usecase.AdvanceStage(c.Request().Context(), id, req.Notes, changedBy)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, res)
}

// UpdateMaterials godoc
// @Summary      Update supplier material list
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  string  true  "Supplier UUID"
// @Success      200
// @Router       /suppliers/{id}/materials [put]
func (h *SupplierHandler) UpdateMaterials(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	var req []supplierusecase.MaterialRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.usecase.UpdateMaterials(c.Request().Context(), id, req); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{"message": "materials updated"})
}

// AddPerformanceRating godoc
// @Summary      Add performance rating
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  string  true  "Supplier UUID"
// @Success      201
// @Router       /suppliers/{id}/ratings [post]
func (h *SupplierHandler) AddPerformanceRating(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	var req supplierusecase.RatingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.usecase.AddPerformanceRating(c.Request().Context(), id, req); err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, map[string]string{"message": "rating added"})
}

// GetPerformanceRatings godoc
// @Summary      Get supplier performance ratings
// @Tags         suppliers
// @Security     BearerAuth
// @Param        id  path  string  true  "Supplier UUID"
// @Success      200  {object}  response.Response
// @Router       /suppliers/{id}/ratings [get]
func (h *SupplierHandler) GetPerformanceRatings(c echo.Context) error {
	id, err := parseUUID(c, "id")
	if err != nil {
		return err
	}

	ratings, err := h.usecase.GetPerformanceRatings(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, ratings)
}

// ---- Helpers ----

func parseUUID(c echo.Context, param string) (uuid.UUID, error) {
	raw := c.Param(param)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid: " + param})
	}
	return id, nil
}

func queryInt(c echo.Context, key string, defaultVal int) int {
	val := c.QueryParam(key)
	if val == "" {
		return defaultVal
	}
	var n int
	if _, err := fmt.Sscanf(val, "%d", &n); err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
