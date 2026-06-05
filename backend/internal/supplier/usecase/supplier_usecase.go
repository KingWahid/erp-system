package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"erp-system/internal/supplier/domain"
	apperrors "erp-system/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type supplierUsecase struct {
	repo   domain.SupplierRepository
	logger *zap.Logger
}

func NewSupplierUsecase(repo domain.SupplierRepository, logger *zap.Logger) SupplierUsecase {
	return &supplierUsecase{repo: repo, logger: logger}
}

func (uc *supplierUsecase) ListSuppliers(ctx context.Context, req ListRequest) (*ListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}

	filter := domain.ListFilter{
		Search:  req.Search,
		Status:  domain.SupplierStatus(req.Status),
		Page:    req.Page,
		PerPage: req.PerPage,
	}

	suppliers, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	items := make([]*SupplierListItem, 0, len(suppliers))
	for _, s := range suppliers {
		item := &SupplierListItem{
			ID:         s.ID,
			Code:       s.Code,
			SupplierNo: s.SupplierNo,
			Name:       s.Name,
			Alias:      s.Alias,
			Address:    fmt.Sprintf("%s, %s", s.City, s.Country),
			Status:     string(s.Status),
		}
		for _, c := range s.Contacts {
			if c.IsPrimary {
				item.Contact = c.Name
				break
			}
		}
		items = append(items, item)
	}

	uc.logger.Info("suppliers listed",
		zap.Int("count", len(items)),
		zap.Int64("total", total),
	)

	return &ListResponse{
		Items:   items,
		Total:   total,
		Page:    req.Page,
		PerPage: req.PerPage,
	}, nil
}

func (uc *supplierUsecase) GetStats(ctx context.Context) (*domain.SupplierStats, error) {
	stats, err := uc.repo.GetStats(ctx)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier stats fetched")
	return stats, nil
}

func (uc *supplierUsecase) CreateSupplier(ctx context.Context, req CreateSupplierRequest) (*SupplierDetailResponse, error) {
	// Check duplicate code
	existing, _ := uc.repo.FindByCode(ctx, req.Code)
	if existing != nil {
		return nil, apperrors.NewCustomError("supplier code already exists").
			WithErrorCode(apperrors.ErrCodeSupplierCodeExists).
			WithMessageID("error_supplier_code_exists").
			WithHTTPCode(http.StatusConflict)
	}

	supplier := &domain.Supplier{
		Name:     req.Name,
		Code:     req.Code,
		Alias:    req.Alias,
		Address:  req.Address,
		City:     req.City,
		Country:  req.Country,
		Phone:    req.Phone,
		Email:    req.Email,
		Website:  req.Website,
		Notes:    req.Notes,
		Status:   domain.StatusDraft,
		Stage:    domain.StageDraft,
		SLAHours: 72,
	}
	supplier.SupplierNo = generateSupplierNo()

	if err := uc.repo.Create(ctx, supplier); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	uc.logger.Info("supplier created",
		zap.String("supplier_id", supplier.ID.String()),
		zap.String("code", supplier.Code),
	)

	return toDetailResponse(supplier), nil
}

func (uc *supplierUsecase) GetSupplierByID(ctx context.Context, id uuid.UUID) (*SupplierDetailResponse, error) {
	supplier, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	return toDetailResponse(supplier), nil
}

func (uc *supplierUsecase) UpdateSupplier(ctx context.Context, id uuid.UUID, req UpdateSupplierRequest) (*SupplierDetailResponse, error) {
	supplier, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.Alias != "" {
		supplier.Alias = req.Alias
	}
	if req.Address != "" {
		supplier.Address = req.Address
	}
	if req.City != "" {
		supplier.City = req.City
	}
	if req.Country != "" {
		supplier.Country = req.Country
	}
	if req.Phone != "" {
		supplier.Phone = req.Phone
	}
	if req.Email != "" {
		supplier.Email = req.Email
	}
	if req.Website != "" {
		supplier.Website = req.Website
	}
	if req.Notes != "" {
		supplier.Notes = req.Notes
	}

	if err := uc.repo.Update(ctx, supplier); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	uc.logger.Info("supplier updated",
		zap.String("supplier_id", supplier.ID.String()),
	)

	return toDetailResponse(supplier), nil
}

func (uc *supplierUsecase) DeleteSupplier(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}
	if err := uc.repo.SoftDelete(ctx, id); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier deleted", zap.String("supplier_id", id.String()))
	return nil
}

func (uc *supplierUsecase) BlockSupplier(ctx context.Context, id uuid.UUID, reason string) error {
	supplier, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	supplier.IsBlocked = true
	supplier.Status = domain.StatusBlocked
	supplier.BlockReason = reason

	if err := uc.repo.Update(ctx, supplier); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	uc.logger.Info("supplier blocked",
		zap.String("supplier_id", id.String()),
		zap.String("reason", reason),
	)
	return nil
}

func (uc *supplierUsecase) UnblockSupplier(ctx context.Context, id uuid.UUID) error {
	supplier, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	supplier.IsBlocked = false
	supplier.Status = domain.StatusActive
	supplier.BlockReason = ""

	if err := uc.repo.Update(ctx, supplier); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	uc.logger.Info("supplier unblocked", zap.String("supplier_id", id.String()))
	return nil
}

// Stage transition map: defines allowed next stage
var stageTransitions = map[domain.SupplierStage]domain.SupplierStage{
	domain.StageDraft:      domain.StageInReview,
	domain.StageInReview:   domain.StageAssessment,
	domain.StageAssessment: domain.StageActive,
}

func (uc *supplierUsecase) AdvanceStage(ctx context.Context, id uuid.UUID, notes, changedBy string) (*SupplierDetailResponse, error) {
	supplier, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	nextStage, ok := stageTransitions[supplier.Stage]
	if !ok {
		return nil, apperrors.NewCustomError("supplier is already in final stage").
			WithErrorCode(apperrors.ErrCodeSupplierInvalidStage).
			WithMessageID("error_supplier_invalid_stage").
			WithHTTPCode(http.StatusBadRequest)
	}

	history := &domain.SupplierStageHistory{
		SupplierID: supplier.ID,
		FromStage:  supplier.Stage,
		ToStage:    nextStage,
		Notes:      notes,
		ChangedBy:  changedBy,
		CreatedAt:  time.Now(),
	}

	supplier.Stage = nextStage
	if nextStage == domain.StageActive {
		supplier.Status = domain.StatusActive
	} else {
		supplier.Status = domain.StatusInProgress
	}

	if err := uc.repo.Update(ctx, supplier); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	_ = uc.repo.AddStageHistory(ctx, history)

	uc.logger.Info("supplier stage advanced",
		zap.String("supplier_id", id.String()),
		zap.String("from_stage", string(history.FromStage)),
		zap.String("to_stage", string(nextStage)),
		zap.String("changed_by", changedBy),
	)

	return toDetailResponse(supplier), nil
}

func (uc *supplierUsecase) UpdateMaterials(ctx context.Context, id uuid.UUID, req []MaterialRequest) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	materials := make([]domain.SupplierMaterial, 0, len(req))
	for _, m := range req {
		materials = append(materials, domain.SupplierMaterial{
			SupplierID:    id,
			MaterialGroup: m.MaterialGroup,
			MaterialID:    m.MaterialID,
			IsActive:      m.IsActive,
		})
	}

	if err := uc.repo.DeleteMaterials(ctx, id); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	if err := uc.repo.UpsertMaterials(ctx, id, materials); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	uc.logger.Info("supplier materials updated",
		zap.String("supplier_id", id.String()),
		zap.Int("count", len(materials)),
	)
	return nil
}

func (uc *supplierUsecase) AddPerformanceRating(ctx context.Context, id uuid.UUID, req RatingRequest) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	rating := &domain.PerformanceRating{
		SupplierID:     id,
		PriceRating:    req.PriceRating,
		DeliveryRating: req.DeliveryRating,
		Notes:          req.Notes,
		ReviewedBy:     req.ReviewedBy,
		ReviewedAt:     time.Now(),
	}

	if err := uc.repo.AddPerformanceRating(ctx, rating); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}

	uc.logger.Info("performance rating added", zap.String("supplier_id", id.String()))
	return nil
}

func (uc *supplierUsecase) GetPerformanceRatings(ctx context.Context, id uuid.UUID) ([]*RatingResponse, error) {
	ratings, err := uc.repo.GetPerformanceRatings(ctx, id)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	result := make([]*RatingResponse, 0, len(ratings))
	for _, r := range ratings {
		result = append(result, &RatingResponse{
			ID:             r.ID,
			PriceRating:    r.PriceRating,
			DeliveryRating: r.DeliveryRating,
			Notes:          r.Notes,
			ReviewedBy:     r.ReviewedBy,
			ReviewedAt:     r.ReviewedAt,
		})
	}
	return result, nil
}

func (uc *supplierUsecase) GetSupplierAddresses(ctx context.Context, supplierID uuid.UUID) ([]AddressResponse, error) {
	addrs, err := uc.repo.GetAddresses(ctx, supplierID)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	result := make([]AddressResponse, 0, len(addrs))
	for _, a := range addrs {
		result = append(result, AddressResponse{
			ID:         a.ID,
			Name:       a.Name,
			Address:    a.Address,
			City:       a.City,
			Province:   a.Province,
			Country:    a.Country,
			PostalCode: a.PostalCode,
			IsMain:     a.IsMain,
		})
	}
	return result, nil
}

func (uc *supplierUsecase) AddSupplierAddress(ctx context.Context, supplierID uuid.UUID, req AddressRequest) (*AddressResponse, error) {
	addr := &domain.SupplierAddress{
		SupplierID: supplierID,
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		Province:   req.Province,
		Country:    req.Country,
		PostalCode: req.PostalCode,
		IsMain:     req.IsMain,
	}
	if err := uc.repo.AddAddress(ctx, addr); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier address added", zap.String("supplier_id", supplierID.String()))
	return &AddressResponse{
		ID:         addr.ID,
		Name:       addr.Name,
		Address:    addr.Address,
		City:       addr.City,
		Province:   addr.Province,
		Country:    addr.Country,
		PostalCode: addr.PostalCode,
		IsMain:     addr.IsMain,
	}, nil
}

func (uc *supplierUsecase) UpdateSupplierAddress(ctx context.Context, supplierID, addressID uuid.UUID, req AddressRequest) (*AddressResponse, error) {
	addr := &domain.SupplierAddress{
		ID:         addressID,
		SupplierID: supplierID,
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		Province:   req.Province,
		Country:    req.Country,
		PostalCode: req.PostalCode,
		IsMain:     req.IsMain,
	}
	if err := uc.repo.UpdateAddress(ctx, addr); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier address updated", zap.String("address_id", addressID.String()))
	return &AddressResponse{
		ID:         addr.ID,
		Name:       addr.Name,
		Address:    addr.Address,
		City:       addr.City,
		Province:   addr.Province,
		Country:    addr.Country,
		PostalCode: addr.PostalCode,
		IsMain:     addr.IsMain,
	}, nil
}

func (uc *supplierUsecase) DeleteSupplierAddress(ctx context.Context, supplierID, addressID uuid.UUID) error {
	if err := uc.repo.DeleteAddress(ctx, supplierID, addressID); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier address deleted", zap.String("address_id", addressID.String()))
	return nil
}

func (uc *supplierUsecase) SetMainAddress(ctx context.Context, supplierID, addressID uuid.UUID) error {
	if err := uc.repo.SetMainAddress(ctx, supplierID, addressID); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier main address set", zap.String("address_id", addressID.String()))
	return nil
}

func (uc *supplierUsecase) GetSupplierContacts(ctx context.Context, supplierID uuid.UUID) ([]ContactResponse, error) {
	contacts, err := uc.repo.GetContacts(ctx, supplierID)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	result := make([]ContactResponse, 0, len(contacts))
	for _, c := range contacts {
		result = append(result, ContactResponse{
			ID:        c.ID,
			Name:      c.Name,
			Position:  c.Position,
			Phone:     c.Phone,
			Mobile:    c.Mobile,
			Email:     c.Email,
			IsPrimary: c.IsPrimary,
		})
	}
	return result, nil
}

func (uc *supplierUsecase) AddSupplierContact(ctx context.Context, supplierID uuid.UUID, req ContactRequest) (*ContactResponse, error) {
	contact := &domain.SupplierContact{
		SupplierID: supplierID,
		Name:       req.Name,
		Position:   req.Position,
		Phone:      req.Phone,
		Mobile:     req.Mobile,
		Email:      req.Email,
		IsPrimary:  req.IsPrimary,
	}
	if err := uc.repo.AddContact(ctx, contact); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier contact added", zap.String("supplier_id", supplierID.String()))
	return &ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Position:  contact.Position,
		Phone:     contact.Phone,
		Mobile:    contact.Mobile,
		Email:     contact.Email,
		IsPrimary: contact.IsPrimary,
	}, nil
}

func (uc *supplierUsecase) UpdateSupplierContact(ctx context.Context, supplierID, contactID uuid.UUID, req ContactRequest) (*ContactResponse, error) {
	contact := &domain.SupplierContact{
		ID:         contactID,
		SupplierID: supplierID,
		Name:       req.Name,
		Position:   req.Position,
		Phone:      req.Phone,
		Mobile:     req.Mobile,
		Email:      req.Email,
		IsPrimary:  req.IsPrimary,
	}
	if err := uc.repo.UpdateContact(ctx, contact); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier contact updated", zap.String("contact_id", contactID.String()))
	return &ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Position:  contact.Position,
		Phone:     contact.Phone,
		Mobile:    contact.Mobile,
		Email:     contact.Email,
		IsPrimary: contact.IsPrimary,
	}, nil
}

func (uc *supplierUsecase) DeleteSupplierContact(ctx context.Context, supplierID, contactID uuid.UUID) error {
	if err := uc.repo.DeleteContact(ctx, supplierID, contactID); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier contact deleted", zap.String("contact_id", contactID.String()))
	return nil
}

func (uc *supplierUsecase) SetPrimaryContact(ctx context.Context, supplierID, contactID uuid.UUID) error {
	if err := uc.repo.SetPrimaryContact(ctx, supplierID, contactID); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier primary contact set", zap.String("contact_id", contactID.String()))
	return nil
}

func (uc *supplierUsecase) GetSupplierGroups(ctx context.Context, supplierID uuid.UUID) ([]GroupResponse, error) {
	groups, err := uc.repo.GetGroups(ctx, supplierID)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	result := make([]GroupResponse, 0, len(groups))
	for _, g := range groups {
		result = append(result, GroupResponse{
			ID:        g.ID,
			GroupName: g.GroupName,
			Value:     g.Value,
			IsActive:  g.IsActive,
		})
	}
	return result, nil
}

func (uc *supplierUsecase) AddSupplierGroup(ctx context.Context, supplierID uuid.UUID, req GroupRequest) (*GroupResponse, error) {
	group := &domain.SupplierGroup{
		SupplierID: supplierID,
		GroupName:  req.GroupName,
		Value:      req.Value,
		IsActive:   req.IsActive,
	}
	if err := uc.repo.AddGroup(ctx, group); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier group added", zap.String("supplier_id", supplierID.String()))
	return &GroupResponse{
		ID:        group.ID,
		GroupName: group.GroupName,
		Value:     group.Value,
		IsActive:  group.IsActive,
	}, nil
}

func (uc *supplierUsecase) UpdateSupplierGroup(ctx context.Context, supplierID, groupID uuid.UUID, req GroupRequest) (*GroupResponse, error) {
	group := &domain.SupplierGroup{
		ID:         groupID,
		SupplierID: supplierID,
		GroupName:  req.GroupName,
		Value:      req.Value,
		IsActive:   req.IsActive,
	}
	if err := uc.repo.UpdateGroup(ctx, group); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier group updated", zap.String("group_id", groupID.String()))
	return &GroupResponse{
		ID:        group.ID,
		GroupName: group.GroupName,
		Value:     group.Value,
		IsActive:  group.IsActive,
	}, nil
}

func (uc *supplierUsecase) DeleteSupplierGroup(ctx context.Context, supplierID, groupID uuid.UUID) error {
	if err := uc.repo.DeleteGroup(ctx, supplierID, groupID); err != nil {
		return apperrors.NewCustomErrorFromError(err)
	}
	uc.logger.Info("supplier group deleted", zap.String("group_id", groupID.String()))
	return nil
}

func (uc *supplierUsecase) GetStageHistories(ctx context.Context, supplierID uuid.UUID) ([]StageHistoryResponse, error) {
	histories, err := uc.repo.GetStageHistories(ctx, supplierID)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}
	result := make([]StageHistoryResponse, 0, len(histories))
	for _, h := range histories {
		result = append(result, StageHistoryResponse{
			ID:        h.ID,
			FromStage: string(h.FromStage),
			ToStage:   string(h.ToStage),
			Notes:     h.Notes,
			ChangedBy: h.ChangedBy,
			ElapsedMs: h.ElapsedMs,
			CreatedAt: h.CreatedAt,
		})
	}
	return result, nil
}

func (uc *supplierUsecase) GetOutstandingInvoices(ctx context.Context, supplierID uuid.UUID, page, perPage int) ([]InvoiceOutstandingResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	invoices, total, err := uc.repo.GetOutstandingInvoices(ctx, supplierID, page, perPage)
	if err != nil {
		return nil, 0, apperrors.NewCustomErrorFromError(err)
	}

	result := make([]InvoiceOutstandingResponse, 0, len(invoices))
	now := time.Now()
	for _, inv := range invoices {
		agingDays := 0
		if inv.DueDate.Before(now) {
			agingDays = int(now.Sub(inv.DueDate).Hours() / 24)
		}
		result = append(result, InvoiceOutstandingResponse{
			ID:            inv.ID,
			InvoiceNumber: inv.InvoiceNumber,
			ProjectName:   inv.ProjectName,
			Amount:        inv.Amount,
			Currency:      inv.Currency,
			InvoiceDate:   inv.InvoiceDate,
			DueDate:       inv.DueDate,
			AgingDays:     agingDays,
			Status:        inv.Status,
			PaidAmount:    inv.PaidAmount,
		})
	}

	uc.logger.Info("outstanding invoices fetched",
		zap.String("supplier_id", supplierID.String()),
		zap.Int("count", len(result)),
	)
	return result, total, nil
}

func generateSupplierNo() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%100000000+10000000)
}

func toDetailResponse(s *domain.Supplier) *SupplierDetailResponse {
	res := &SupplierDetailResponse{
		ID:          s.ID,
		Code:        s.Code,
		SupplierNo:  s.SupplierNo,
		Name:        s.Name,
		Alias:       s.Alias,
		LogoURL:     s.LogoURL,
		Address:     s.Address,
		City:        s.City,
		Country:     s.Country,
		Phone:       s.Phone,
		Email:       s.Email,
		Website:     s.Website,
		Status:      string(s.Status),
		Stage:       string(s.Stage),
		SLAHours:    s.SLAHours,
		IsBlocked:   s.IsBlocked,
		BlockReason: s.BlockReason,
		Notes:       s.Notes,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	for _, c := range s.Contacts {
		res.Contacts = append(res.Contacts, ContactResponse{
			ID:        c.ID,
			Name:      c.Name,
			Position:  c.Position,
			Phone:     c.Phone,
			Mobile:    c.Mobile,
			Email:     c.Email,
			IsPrimary: c.IsPrimary,
		})
	}

	for _, a := range s.Addresses {
		res.Addresses = append(res.Addresses, AddressResponse{
			ID:         a.ID,
			Name:       a.Name,
			Address:    a.Address,
			City:       a.City,
			Province:   a.Province,
			Country:    a.Country,
			PostalCode: a.PostalCode,
			IsMain:     a.IsMain,
		})
	}

	for _, g := range s.Groups {
		res.Groups = append(res.Groups, GroupResponse{
			ID:        g.ID,
			GroupName: g.GroupName,
			Value:     g.Value,
			IsActive:  g.IsActive,
		})
	}

	for _, m := range s.Materials {
		res.Materials = append(res.Materials, MaterialResponse{
			ID:            m.ID,
			MaterialGroup: m.MaterialGroup,
			MaterialID:    m.MaterialID,
			IsActive:      m.IsActive,
		})
	}

	for _, h := range s.StageHistories {
		res.StageHistories = append(res.StageHistories, StageHistoryResponse{
			ID:        h.ID,
			FromStage: string(h.FromStage),
			ToStage:   string(h.ToStage),
			Notes:     h.Notes,
			ChangedBy: h.ChangedBy,
			ElapsedMs: h.ElapsedMs,
			CreatedAt: h.CreatedAt,
		})
	}

	return res
}
