import apiClient from "@/lib/axios";
import type {
  AdvanceStageRequest,
  AddressRequest,
  BlockSupplierRequest,
  BlockSupplierResponse,
  ContactRequest,
  CreateSupplierRequest,
  GroupRequest,
  MaterialItem,
  MessageResponse,
  RatingRequest,
  SupplierDetailResponse,
  SupplierListParams,
  SupplierListResponse,
  SupplierStatsResponse,
  UpdateSupplierRequest,
  AddressResponse,
  ContactResponse,
  GroupResponse,
  RatingListResponse,
  AddRatingResponse,
  StageHistoryListResponse,
  AddressListResponse,
  ContactListResponse,
  GroupListResponse,
  OutstandingListResponse,
  OutstandingListParams,
} from "@/types/api";

export const getSupplierStats = async (): Promise<SupplierStatsResponse> => {
  const { data } = await apiClient.get<SupplierStatsResponse>("/suppliers/stats");
  return data;
};

export const listSuppliers = async (
  params: SupplierListParams = {}
): Promise<SupplierListResponse> => {
  const { data } = await apiClient.get<SupplierListResponse>("/suppliers", { params });
  return data;
};

export const getSupplier = async (id: string): Promise<SupplierDetailResponse> => {
  const { data } = await apiClient.get<SupplierDetailResponse>(`/suppliers/${id}`);
  return data;
};

export const createSupplier = async (
  payload: CreateSupplierRequest
): Promise<SupplierDetailResponse> => {
  const { data } = await apiClient.post<SupplierDetailResponse>("/suppliers", payload);
  return data;
};

export const updateSupplier = async (
  id: string,
  payload: UpdateSupplierRequest
): Promise<SupplierDetailResponse> => {
  const { data } = await apiClient.put<SupplierDetailResponse>(`/suppliers/${id}`, payload);
  return data;
};

export const deleteSupplier = async (id: string): Promise<void> => {
  await apiClient.delete(`/suppliers/${id}`);
};

export const blockSupplier = async (
  id: string,
  payload: BlockSupplierRequest
): Promise<BlockSupplierResponse> => {
  const { data } = await apiClient.post<BlockSupplierResponse>(
    `/suppliers/${id}/block`,
    payload
  );
  return data;
};

export const advanceSupplierStage = async (
  id: string,
  payload?: AdvanceStageRequest
): Promise<SupplierDetailResponse> => {
  const { data } = await apiClient.post<SupplierDetailResponse>(
    `/suppliers/${id}/next-stage`,
    payload ?? {}
  );
  return data;
};

export const listAddresses = async (supplierId: string): Promise<AddressListResponse> => {
  const { data } = await apiClient.get<AddressListResponse>(
    `/suppliers/${supplierId}/addresses`
  );
  return data;
};

export const addAddress = async (
  supplierId: string,
  payload: AddressRequest
): Promise<AddressResponse> => {
  const { data } = await apiClient.post<AddressResponse>(
    `/suppliers/${supplierId}/addresses`,
    payload
  );
  return data;
};

export const updateAddress = async (
  supplierId: string,
  addressId: string,
  payload: AddressRequest
): Promise<AddressResponse> => {
  const { data } = await apiClient.put<AddressResponse>(
    `/suppliers/${supplierId}/addresses/${addressId}`,
    payload
  );
  return data;
};

export const deleteAddress = async (supplierId: string, addressId: string): Promise<void> => {
  await apiClient.delete(`/suppliers/${supplierId}/addresses/${addressId}`);
};

export const listContacts = async (supplierId: string): Promise<ContactListResponse> => {
  const { data } = await apiClient.get<ContactListResponse>(
    `/suppliers/${supplierId}/contacts`
  );
  return data;
};

export const addContact = async (
  supplierId: string,
  payload: ContactRequest
): Promise<ContactResponse> => {
  const { data } = await apiClient.post<ContactResponse>(
    `/suppliers/${supplierId}/contacts`,
    payload
  );
  return data;
};

export const updateContact = async (
  supplierId: string,
  contactId: string,
  payload: ContactRequest
): Promise<ContactResponse> => {
  const { data } = await apiClient.put<ContactResponse>(
    `/suppliers/${supplierId}/contacts/${contactId}`,
    payload
  );
  return data;
};

export const deleteContact = async (supplierId: string, contactId: string): Promise<void> => {
  await apiClient.delete(`/suppliers/${supplierId}/contacts/${contactId}`);
};

export const listGroups = async (supplierId: string): Promise<GroupListResponse> => {
  const { data } = await apiClient.get<GroupListResponse>(`/suppliers/${supplierId}/groups`);
  return data;
};

export const addGroup = async (
  supplierId: string,
  payload: GroupRequest
): Promise<GroupResponse> => {
  const { data } = await apiClient.post<GroupResponse>(
    `/suppliers/${supplierId}/groups`,
    payload
  );
  return data;
};

export const deleteGroup = async (supplierId: string, groupId: string): Promise<void> => {
  await apiClient.delete(`/suppliers/${supplierId}/groups/${groupId}`);
};

export const updateMaterials = async (
  supplierId: string,
  materials: Omit<MaterialItem, "id">[]
): Promise<MessageResponse> => {
  const { data } = await apiClient.put<MessageResponse>(
    `/suppliers/${supplierId}/materials`,
    { materials }
  );
  return data;
};

export const listRatings = async (supplierId: string): Promise<RatingListResponse> => {
  const { data } = await apiClient.get<RatingListResponse>(
    `/suppliers/${supplierId}/ratings`
  );
  return data;
};

export const addRating = async (
  supplierId: string,
  payload: RatingRequest
): Promise<AddRatingResponse> => {
  const { data } = await apiClient.post<AddRatingResponse>(
    `/suppliers/${supplierId}/ratings`,
    payload
  );
  return data;
};

export const listStageHistory = async (supplierId: string): Promise<StageHistoryListResponse> => {
  const { data } = await apiClient.get<StageHistoryListResponse>(
    `/suppliers/${supplierId}/stage-history`
  );
  return data;
};

export const listOutstandings = async (
  supplierId: string,
  params: OutstandingListParams = {}
): Promise<OutstandingListResponse> => {
  const { data } = await apiClient.get<OutstandingListResponse>(
    `/suppliers/${supplierId}/outstandings`,
    { params }
  );
  return data;
};
