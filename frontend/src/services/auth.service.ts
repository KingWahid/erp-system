import apiClient from "@/lib/axios";
import type {
  AuthResponse,
  LoginRequest,
  ProfileResponse,
  RegisterRequest,
} from "@/types/api";

export const authService = {
  login: async (payload: LoginRequest): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>("/auth/login", payload);
    return data;
  },

  register: async (payload: RegisterRequest): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>("/auth/register", payload);
    return data;
  },

  getProfile: async (): Promise<ProfileResponse> => {
    const { data } = await apiClient.get<ProfileResponse>("/auth/me");
    return data;
  },
};
