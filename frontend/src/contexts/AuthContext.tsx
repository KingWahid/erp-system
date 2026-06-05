"use client";

import {
  createContext,
  useCallback,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import Cookies from "js-cookie";
import { authService } from "@/services/auth.service";
import type { LoginRequest, UserPayload } from "@/types/api";

interface AuthContextValue {
  user: UserPayload | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (payload: LoginRequest) => Promise<void>;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextValue | null>(null);

const TOKEN_COOKIE = process.env.NEXT_PUBLIC_TOKEN_COOKIE ?? "erp_token";
const COOKIE_OPTIONS = { expires: 1, sameSite: "strict" as const };

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserPayload | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Restore session from cookie on mount
  useEffect(() => {
    const storedToken = Cookies.get(TOKEN_COOKIE);
    if (!storedToken) {
      setIsLoading(false);
      return;
    }

    setToken(storedToken);
    authService
      .getProfile()
      .then((res) => {
        // getProfile returns {id, name, email, roles} — build UserPayload shape
        setUser({
          id: res.data.id,
          name: res.data.name,
          email: res.data.email,
          roles: res.data.roles,
          permissions: [], // populated via PermissionContext
        });
      })
      .catch(() => {
        // Token expired or invalid — clean up
        Cookies.remove(TOKEN_COOKIE);
        setToken(null);
      })
      .finally(() => setIsLoading(false));
  }, []);

  const login = useCallback(async (payload: LoginRequest) => {
    const res = await authService.login(payload);
    const { token: newToken, user: newUser } = res.data;

    Cookies.set(TOKEN_COOKIE, newToken, COOKIE_OPTIONS);
    setToken(newToken);
    setUser(newUser);
  }, []);

  const logout = useCallback(() => {
    Cookies.remove(TOKEN_COOKIE);
    setToken(null);
    setUser(null);
    window.location.href = "/login";
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      isAuthenticated: !!token && !!user,
      isLoading,
      login,
      logout,
    }),
    [user, token, isLoading, login, logout]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
