"use client";

import { type ReactNode } from "react";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import { App, ConfigProvider } from "antd";
import idID from "antd/locale/id_ID";
import { queryClient } from "@/lib/queryClient";
import { AuthProvider } from "@/contexts/AuthContext";
import { PermissionProvider } from "@/contexts/PermissionContext";

/**
 * Top-level provider tree.
 * Order matters:
 *   QueryClient → Antd → Auth → Permission
 */
export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <AntdRegistry>
        <ConfigProvider
          locale={idID}
          theme={{
            token: {
              colorPrimary: "#1677ff",
              borderRadius: 6,
            },
          }}
          // Suppress React 19 compatibility warning
          warning={{ strict: false }}
        >
          <App>
            <AuthProvider>
              <PermissionProvider>
                {children}
              </PermissionProvider>
            </AuthProvider>
          </App>
        </ConfigProvider>
      </AntdRegistry>
      {process.env.NODE_ENV === "development" && (
        <ReactQueryDevtools initialIsOpen={false} />
      )}
    </QueryClientProvider>
  );
}
