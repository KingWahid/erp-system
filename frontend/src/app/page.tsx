import { redirect } from "next/navigation";

/**
 * Root route — redirect to dashboard.
 * Auth guard in the dashboard layout handles unauthenticated users.
 */
export default function RootPage() {
  redirect("/suppliers");
}
