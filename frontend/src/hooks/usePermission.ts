import { useContext } from "react";
import { PermissionContext } from "@/contexts/PermissionContext";

/**
 * Access permission-checking utilities.
 * Must be used inside <PermissionProvider>.
 *
 * @example
 * const { can, canAny, hasRole } = usePermission();
 * if (can("supplier:create")) { ... }
 * if (canAny("supplier:update", "supplier:delete")) { ... }
 * if (hasRole("admin")) { ... }
 */
export function usePermission() {
  const ctx = useContext(PermissionContext);
  if (!ctx) {
    throw new Error("usePermission must be used within <PermissionProvider>");
  }
  return ctx;
}
