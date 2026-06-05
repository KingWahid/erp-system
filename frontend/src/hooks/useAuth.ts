import { useContext } from "react";
import { AuthContext } from "@/contexts/AuthContext";

/**
 * Access the current auth state and actions.
 * Must be used inside <AuthProvider>.
 */
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within <AuthProvider>");
  }
  return ctx;
}
