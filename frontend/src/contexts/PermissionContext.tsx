"use client";

import {
  createContext,
  useContext,
  useMemo,
  type ReactNode,
} from "react";
import { AuthContext } from "./AuthContext";

// e.g. "supplier:read", "supplier:create", "workflow:advance"

type Permission = string;

interface PermissionContextValue {
  permissions: Permission[];
  roles: string[];
  /** Check if the user has ALL of the given permissions */
  can: (...perms: Permission[]) => boolean;
  /** Check if the user has ANY of the given permissions */
  canAny: (...perms: Permission[]) => boolean;
  /** Check if the user has a specific role */
  hasRole: (role: string) => boolean;
  isAdmin: boolean;
}

export const PermissionContext = createContext<PermissionContextValue | null>(null);

export function PermissionProvider({ children }: { children: ReactNode }) {
  const auth = useContext(AuthContext);

  const permissions: Permission[] = auth?.user?.permissions ?? [];
  const roles: string[] = auth?.user?.roles ?? [];

  const value = useMemo<PermissionContextValue>(() => {
    const permSet = new Set(permissions);
    const roleSet = new Set(roles);
    const isAdmin = roleSet.has("admin");

    return {
      permissions,
      roles,
      can: (...perms) => isAdmin || perms.every((p) => permSet.has(p)),
      canAny: (...perms) => isAdmin || perms.some((p) => permSet.has(p)),
      hasRole: (role) => roleSet.has(role),
      isAdmin,
    };
  }, [permissions, roles]);

  return (
    <PermissionContext.Provider value={value}>
      {children}
    </PermissionContext.Provider>
  );
}
