package handler

import (
	"erp-system/internal/generated"
	supplierusecase "erp-system/internal/supplier/usecase"
	apperrors "erp-system/pkg/errors"
	"erp-system/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// SupplierAdapter implements the supplier portion of generated.ServerInterface.
type SupplierAdapter struct {
	usecase supplierusecase.SupplierUsecase
	logger  *zap.Logger
}

func NewSupplierAdapter(uc supplierusecase.SupplierUsecase, logger *zap.Logger) *SupplierAdapter {
	return &SupplierAdapter{usecase: uc, logger: logger}
}

// GetSupplierStats implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplierStats(ctx echo.Context, params generated.GetSupplierStatsParams) error {
	stats, err := a.usecase.GetStats(ctx.Request().Context())
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, stats)
}

// ListSuppliers implements generated.ServerInterface
func (a *SupplierAdapter) ListSuppliers(ctx echo.Context, params generated.ListSuppliersParams) error {
	req := supplierusecase.ListRequest{
		Page:    1,
		PerPage: 20,
	}
	if params.Search != nil {
		req.Search = *params.Search
	}
	if params.Status != nil {
		req.Status = string(*params.Status)
	}
	if params.Page != nil {
		req.Page = *params.Page
	}
	if params.Limit != nil {
		req.PerPage = *params.Limit
	}

	res, err := a.usecase.ListSuppliers(ctx.Request().Context(), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.SuccessWithMeta(ctx, res.Items, response.Meta{
		Page:       res.Page,
		PerPage:    res.PerPage,
		Total:      res.Total,
		TotalPages: int((res.Total + int64(res.PerPage) - 1) / int64(res.PerPage)),
	})
}

// CreateSupplier implements generated.ServerInterface
func (a *SupplierAdapter) CreateSupplier(ctx echo.Context, params generated.CreateSupplierParams) error {
	var body generated.CreateSupplierJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.CreateSupplierRequest{
		Name: body.Name,
		Code: body.Code,
	}
	if body.Address != nil {
		req.Address = *body.Address
	}
	if body.Alias != nil {
		req.Alias = *body.Alias
	}
	if body.City != nil {
		req.City = *body.City
	}
	if body.Country != nil {
		req.Country = *body.Country
	}
	if body.Phone != nil {
		req.Phone = *body.Phone
	}
	if body.Email != nil {
		req.Email = string(*body.Email)
	}
	if body.Website != nil {
		req.Website = *body.Website
	}
	if body.Notes != nil {
		req.Notes = *body.Notes
	}

	res, err := a.usecase.CreateSupplier(ctx.Request().Context(), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Created(ctx, res)
}

// GetSupplier implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplier(ctx echo.Context, supplierID generated.SupplierID, params generated.GetSupplierParams) error {
	res, err := a.usecase.GetSupplierByID(ctx.Request().Context(), uuid.UUID(supplierID))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}

// UpdateSupplier implements generated.ServerInterface
func (a *SupplierAdapter) UpdateSupplier(ctx echo.Context, supplierID generated.SupplierID, params generated.UpdateSupplierParams) error {
	var body generated.UpdateSupplierJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.UpdateSupplierRequest{}
	if body.Name != nil {
		req.Name = *body.Name
	}
	if body.Alias != nil {
		req.Alias = *body.Alias
	}
	if body.Address != nil {
		req.Address = *body.Address
	}
	if body.City != nil {
		req.City = *body.City
	}
	if body.Country != nil {
		req.Country = *body.Country
	}
	if body.Phone != nil {
		req.Phone = *body.Phone
	}
	if body.Email != nil {
		req.Email = string(*body.Email)
	}
	if body.Website != nil {
		req.Website = *body.Website
	}
	if body.Notes != nil {
		req.Notes = *body.Notes
	}

	res, err := a.usecase.UpdateSupplier(ctx.Request().Context(), uuid.UUID(supplierID), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}

// DeleteSupplier implements generated.ServerInterface
func (a *SupplierAdapter) DeleteSupplier(ctx echo.Context, supplierID generated.SupplierID, params generated.DeleteSupplierParams) error {
	if err := a.usecase.DeleteSupplier(ctx.Request().Context(), uuid.UUID(supplierID)); err != nil {
		return response.Error(ctx, err)
	}
	return response.NoContent(ctx)
}

// BlockSupplier implements generated.ServerInterface
func (a *SupplierAdapter) BlockSupplier(ctx echo.Context, supplierID generated.SupplierID, params generated.BlockSupplierParams) error {
	var body generated.BlockSupplierJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	var err error
	if body.Block {
		reason := ""
		if body.Reason != nil {
			reason = *body.Reason
		}
		err = a.usecase.BlockSupplier(ctx.Request().Context(), uuid.UUID(supplierID), reason)
	} else {
		err = a.usecase.UnblockSupplier(ctx.Request().Context(), uuid.UUID(supplierID))
	}

	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, map[string]bool{"blocked": body.Block})
}

// AdvanceSupplierStage implements generated.ServerInterface
func (a *SupplierAdapter) AdvanceSupplierStage(ctx echo.Context, supplierID generated.SupplierID, params generated.AdvanceSupplierStageParams) error {
	var body generated.AdvanceSupplierStageJSONRequestBody
	_ = ctx.Bind(&body)

	notes := ""
	if body.Notes != nil {
		notes = *body.Notes
	}
	changedBy, _ := ctx.Get("user_id").(string)

	res, err := a.usecase.AdvanceStage(ctx.Request().Context(), uuid.UUID(supplierID), notes, changedBy)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}

// GetSupplierMaterials implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplierMaterials(ctx echo.Context, supplierID generated.SupplierID, params generated.GetSupplierMaterialsParams) error {
	detail, err := a.usecase.GetSupplierByID(ctx.Request().Context(), uuid.UUID(supplierID))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, detail.Materials)
}

// UpdateSupplierMaterials implements generated.ServerInterface
func (a *SupplierAdapter) UpdateSupplierMaterials(ctx echo.Context, supplierID generated.SupplierID, params generated.UpdateSupplierMaterialsParams) error {
	var body generated.UpdateSupplierMaterialsJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := make([]supplierusecase.MaterialRequest, 0, len(body))
	for _, m := range body {
		active := true
		if m.IsActive != nil {
			active = *m.IsActive
		}
		req = append(req, supplierusecase.MaterialRequest{
			MaterialGroup: m.MaterialGroup,
			MaterialID:    m.MaterialId,
			IsActive:      active,
		})
	}

	if err := a.usecase.UpdateMaterials(ctx.Request().Context(), uuid.UUID(supplierID), req); err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, map[string]string{"message": "materials updated"})
}

// GetSupplierRatings implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplierRatings(ctx echo.Context, supplierID generated.SupplierID, params generated.GetSupplierRatingsParams) error {
	ratings, err := a.usecase.GetPerformanceRatings(ctx.Request().Context(), uuid.UUID(supplierID))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, ratings)
}

// AddSupplierRating implements generated.ServerInterface
func (a *SupplierAdapter) AddSupplierRating(ctx echo.Context, supplierID generated.SupplierID, params generated.AddSupplierRatingParams) error {
	var body generated.AddSupplierRatingJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.RatingRequest{
		PriceRating:    body.PriceRating,
		DeliveryRating: body.DeliveryRating,
	}
	if body.Notes != nil {
		req.Notes = *body.Notes
	}
	if body.ReviewedBy != nil {
		req.ReviewedBy = *body.ReviewedBy
	}

	if err := a.usecase.AddPerformanceRating(ctx.Request().Context(), uuid.UUID(supplierID), req); err != nil {
		return response.Error(ctx, err)
	}
	return response.Created(ctx, map[string]string{"message": "rating added"})
}

// GetSupplierAddresses implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplierAddresses(ctx echo.Context, supplierID generated.SupplierID, params generated.GetSupplierAddressesParams) error {
	addrs, err := a.usecase.GetSupplierAddresses(ctx.Request().Context(), uuid.UUID(supplierID))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, addrs)
}

// AddSupplierAddress implements generated.ServerInterface
func (a *SupplierAdapter) AddSupplierAddress(ctx echo.Context, supplierID generated.SupplierID, params generated.AddSupplierAddressParams) error {
	var body generated.AddressRequest
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.AddressRequest{
		Name:    body.Name,
		Address: body.Address,
	}
	if body.City != nil {
		req.City = *body.City
	}
	if body.Province != nil {
		req.Province = *body.Province
	}
	if body.Country != nil {
		req.Country = *body.Country
	}
	if body.PostalCode != nil {
		req.PostalCode = *body.PostalCode
	}
	if body.IsMain != nil {
		req.IsMain = *body.IsMain
	}

	res, err := a.usecase.AddSupplierAddress(ctx.Request().Context(), uuid.UUID(supplierID), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Created(ctx, res)
}

// UpdateSupplierAddress implements generated.ServerInterface
func (a *SupplierAdapter) UpdateSupplierAddress(ctx echo.Context, supplierID generated.SupplierID, addressID generated.AddressID, params generated.UpdateSupplierAddressParams) error {
	var body generated.AddressRequest
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.AddressRequest{
		Name:    body.Name,
		Address: body.Address,
	}
	if body.City != nil {
		req.City = *body.City
	}
	if body.Province != nil {
		req.Province = *body.Province
	}
	if body.Country != nil {
		req.Country = *body.Country
	}
	if body.PostalCode != nil {
		req.PostalCode = *body.PostalCode
	}
	if body.IsMain != nil {
		req.IsMain = *body.IsMain
	}

	res, err := a.usecase.UpdateSupplierAddress(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(addressID), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}

// DeleteSupplierAddress implements generated.ServerInterface
func (a *SupplierAdapter) DeleteSupplierAddress(ctx echo.Context, supplierID generated.SupplierID, addressID generated.AddressID, params generated.DeleteSupplierAddressParams) error {
	if err := a.usecase.DeleteSupplierAddress(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(addressID)); err != nil {
		return response.Error(ctx, err)
	}
	return response.NoContent(ctx)
}

// SetMainSupplierAddress implements generated.ServerInterface
func (a *SupplierAdapter) SetMainSupplierAddress(ctx echo.Context, supplierID generated.SupplierID, addressID generated.AddressID, params generated.SetMainSupplierAddressParams) error {
	if err := a.usecase.SetMainAddress(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(addressID)); err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, map[string]string{"message": "main address set"})
}

// GetSupplierContacts implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplierContacts(ctx echo.Context, supplierID generated.SupplierID, params generated.GetSupplierContactsParams) error {
	contacts, err := a.usecase.GetSupplierContacts(ctx.Request().Context(), uuid.UUID(supplierID))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, contacts)
}

// AddSupplierContact implements generated.ServerInterface
func (a *SupplierAdapter) AddSupplierContact(ctx echo.Context, supplierID generated.SupplierID, params generated.AddSupplierContactParams) error {
	var body generated.ContactRequest
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.ContactRequest{
		Name: body.Name,
	}
	if body.Position != nil {
		req.Position = *body.Position
	}
	if body.Phone != nil {
		req.Phone = *body.Phone
	}
	if body.Mobile != nil {
		req.Mobile = *body.Mobile
	}
	if body.Email != nil {
		req.Email = string(*body.Email)
	}
	if body.IsPrimary != nil {
		req.IsPrimary = *body.IsPrimary
	}

	res, err := a.usecase.AddSupplierContact(ctx.Request().Context(), uuid.UUID(supplierID), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Created(ctx, res)
}

// UpdateSupplierContact implements generated.ServerInterface
func (a *SupplierAdapter) UpdateSupplierContact(ctx echo.Context, supplierID generated.SupplierID, contactID generated.ContactID, params generated.UpdateSupplierContactParams) error {
	var body generated.ContactRequest
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.ContactRequest{
		Name: body.Name,
	}
	if body.Position != nil {
		req.Position = *body.Position
	}
	if body.Phone != nil {
		req.Phone = *body.Phone
	}
	if body.Mobile != nil {
		req.Mobile = *body.Mobile
	}
	if body.Email != nil {
		req.Email = string(*body.Email)
	}
	if body.IsPrimary != nil {
		req.IsPrimary = *body.IsPrimary
	}

	res, err := a.usecase.UpdateSupplierContact(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(contactID), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}

// DeleteSupplierContact implements generated.ServerInterface
func (a *SupplierAdapter) DeleteSupplierContact(ctx echo.Context, supplierID generated.SupplierID, contactID generated.ContactID, params generated.DeleteSupplierContactParams) error {
	if err := a.usecase.DeleteSupplierContact(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(contactID)); err != nil {
		return response.Error(ctx, err)
	}
	return response.NoContent(ctx)
}

// SetPrimarySupplierContact implements generated.ServerInterface
func (a *SupplierAdapter) SetPrimarySupplierContact(ctx echo.Context, supplierID generated.SupplierID, contactID generated.ContactID, params generated.SetPrimarySupplierContactParams) error {
	if err := a.usecase.SetPrimaryContact(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(contactID)); err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, map[string]string{"message": "primary contact set"})
}

// GetSupplierGroups implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplierGroups(ctx echo.Context, supplierID generated.SupplierID, params generated.GetSupplierGroupsParams) error {
	groups, err := a.usecase.GetSupplierGroups(ctx.Request().Context(), uuid.UUID(supplierID))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, groups)
}

// AddSupplierGroup implements generated.ServerInterface
func (a *SupplierAdapter) AddSupplierGroup(ctx echo.Context, supplierID generated.SupplierID, params generated.AddSupplierGroupParams) error {
	var body generated.GroupRequest
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.GroupRequest{
		GroupName: body.GroupName,
		Value:     body.Value,
	}
	if body.IsActive != nil {
		req.IsActive = *body.IsActive
	}

	res, err := a.usecase.AddSupplierGroup(ctx.Request().Context(), uuid.UUID(supplierID), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Created(ctx, res)
}

// UpdateSupplierGroup implements generated.ServerInterface
func (a *SupplierAdapter) UpdateSupplierGroup(ctx echo.Context, supplierID generated.SupplierID, groupID generated.GroupID, params generated.UpdateSupplierGroupParams) error {
	var body generated.GroupRequest
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	req := supplierusecase.GroupRequest{
		GroupName: body.GroupName,
		Value:     body.Value,
	}
	if body.IsActive != nil {
		req.IsActive = *body.IsActive
	}

	res, err := a.usecase.UpdateSupplierGroup(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(groupID), req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}

// DeleteSupplierGroup implements generated.ServerInterface
func (a *SupplierAdapter) DeleteSupplierGroup(ctx echo.Context, supplierID generated.SupplierID, groupID generated.GroupID, params generated.DeleteSupplierGroupParams) error {
	if err := a.usecase.DeleteSupplierGroup(ctx.Request().Context(), uuid.UUID(supplierID), uuid.UUID(groupID)); err != nil {
		return response.Error(ctx, err)
	}
	return response.NoContent(ctx)
}

// GetStageHistory implements generated.ServerInterface
func (a *SupplierAdapter) GetStageHistory(ctx echo.Context, supplierID generated.SupplierID, params generated.GetStageHistoryParams) error {
	histories, err := a.usecase.GetStageHistories(ctx.Request().Context(), uuid.UUID(supplierID))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, histories)
}

// GetSupplierOutstandings implements generated.ServerInterface
func (a *SupplierAdapter) GetSupplierOutstandings(ctx echo.Context, supplierID generated.SupplierID, params generated.GetSupplierOutstandingsParams) error {
	page := 1
	perPage := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		perPage = *params.Limit
	}

	invoices, total, err := a.usecase.GetOutstandingInvoices(ctx.Request().Context(), uuid.UUID(supplierID), page, perPage)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.SuccessWithMeta(ctx, invoices, response.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	})
}
