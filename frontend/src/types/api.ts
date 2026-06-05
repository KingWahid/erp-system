/**
 * API types — mirrored from backend OpenAPI spec (openapi.yaml)
 * Keep in sync with backend/api/openapi.yaml
 */

export type SupplierStatus = "draft" | "active" | "in_progress" | "blocked" | "inactive";
export type SupplierStage = "draft" | "in_review" | "in_assessment" | "active";
export type UserRole = "admin" | "manager" | "viewer" | "supplier";

export interface ErrorBody {
  code: string;
  message: string;
  message_id?: string;
}

export interface ErrorResponse {
  success: false;
  error: ErrorBody;
}

export interface Meta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface MessageResponse {
  success: boolean;
  data: {
    message: string;
  };
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  role?: UserRole;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface UserPayload {
  id: string;
  name: string;
  email: string;
  roles: string[];
  permissions: string[];
}

export interface AuthData {
  token: string;
  user: UserPayload;
}

export interface AuthResponse {
  success: boolean;
  data: AuthData;
}

export interface ProfileResponse {
  success: boolean;
  data: {
    id: string;
    name: string;
    email: string;
    roles: string[];
  };
}

export interface SupplierStats {
  total_supplier: number;
  new_supplier: number;
  avg_cost_supplier: number;
  blocked_supplier: number;
}

export interface SupplierStatsResponse {
  success: boolean;
  data: SupplierStats;
}

export interface SupplierListItem {
  id: string;
  code: string;
  supplier_no: string;
  name: string;
  alias?: string;
  address?: string;
  contact?: string;
  status: SupplierStatus;
}

export interface SupplierListResponse {
  success: boolean;
  data: SupplierListItem[];
  meta: Meta;
}

export interface SupplierListParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: SupplierStatus;
}

export interface AddressItem {
  id: string;
  name: string;
  address: string;
  city?: string;
  province?: string;
  country?: string;
  postal_code?: string;
  is_main: boolean;
}

export interface AddressRequest {
  name: string;
  address: string;
  city?: string;
  province?: string;
  country?: string;
  postal_code?: string;
  is_main?: boolean;
}

export interface ContactItem {
  id: string;
  name: string;
  position?: string;
  phone?: string;
  mobile?: string;
  email?: string;
  is_primary: boolean;
}

export interface ContactRequest {
  name: string;
  position?: string;
  phone?: string;
  mobile?: string;
  email?: string;
  is_primary?: boolean;
}

export interface GroupItem {
  id: string;
  group_name: string;
  value: string;
  is_active: boolean;
}

export interface GroupRequest {
  group_name: string;
  value: string;
  is_active?: boolean;
}

export interface MaterialItem {
  id: string;
  material_group: string;
  material_id: string;
  is_active: boolean;
}

export interface MaterialRequest {
  material_group: string;
  material_id: string;
  is_active?: boolean;
}

export interface RatingItem {
  id: string;
  price_rating: number;
  delivery_rating: number;
  notes?: string;
  reviewed_by?: string;
  reviewed_at: string;
}

export interface RatingRequest {
  price_rating: number;
  delivery_rating: number;
  notes?: string;
  reviewed_by?: string;
}

export interface StageHistoryItem {
  id: string;
  from_stage?: SupplierStage;
  to_stage: SupplierStage;
  notes?: string;
  changed_by?: string;
  elapsed_ms?: number;
  created_at: string;
}

export interface SupplierDetail {
  id: string;
  code: string;
  supplier_no: string;
  name: string;
  alias?: string;
  logo_url?: string;
  address?: string;
  city?: string;
  country?: string;
  phone?: string;
  email?: string;
  website?: string;
  status: SupplierStatus;
  stage: SupplierStage;
  sla_hours: number;
  is_blocked: boolean;
  block_reason?: string;
  notes?: string;
  addresses: AddressItem[];
  contacts: ContactItem[];
  groups: GroupItem[];
  materials: MaterialItem[];
  stage_histories: StageHistoryItem[];
  created_at: string;
  updated_at: string;
}

export interface SupplierDetailResponse {
  success: boolean;
  data: SupplierDetail;
}

export interface CreateSupplierRequest {
  name: string;
  code: string;
  alias?: string;
  address?: string;
  city?: string;
  country?: string;
  phone?: string;
  email?: string;
  website?: string;
  notes?: string;
}

export interface UpdateSupplierRequest {
  name?: string;
  alias?: string;
  address?: string;
  city?: string;
  country?: string;
  phone?: string;
  email?: string;
  website?: string;
  notes?: string;
}

export interface BlockSupplierRequest {
  block: boolean;
  reason?: string;
}

export interface BlockSupplierResponse {
  success: boolean;
  data: {
    blocked: boolean;
  };
}

export interface AdvanceStageRequest {
  notes?: string;
}

export interface AddressResponse {
  success: boolean;
  data: AddressItem;
}

export interface AddressListResponse {
  success: boolean;
  data: AddressItem[];
}

export interface ContactResponse {
  success: boolean;
  data: ContactItem;
}

export interface ContactListResponse {
  success: boolean;
  data: ContactItem[];
}

export interface GroupResponse {
  success: boolean;
  data: GroupItem;
}

export interface GroupListResponse {
  success: boolean;
  data: GroupItem[];
}

export interface RatingListResponse {
  success: boolean;
  data: RatingItem[];
}

export interface AddRatingResponse {
  success: boolean;
  data: { message: string };
}

export interface StageHistoryListResponse {
  success: boolean;
  data: StageHistoryItem[];
}

// ── Outstandings ──────────────────────────────────────────────

export type InvoiceStatus = "unpaid" | "partial" | "paid" | "overdue";

export interface InvoiceItem {
  id: string;
  invoice_number: string;
  project_name?: string;
  amount: number;
  currency: string;
  invoice_date: string;
  due_date: string;
  aging_days: number;
  status: InvoiceStatus;
  paid_amount: number;
}

export interface OutstandingListResponse {
  success: boolean;
  data: InvoiceItem[];
  meta: Meta;
}

export interface OutstandingListParams {
  page?: number;
  limit?: number;
}
