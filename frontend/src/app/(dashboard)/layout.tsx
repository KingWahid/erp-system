"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import {
  Avatar,
  Dropdown,
  Layout,
  Menu,
  Spin,
} from "antd";
import {
  AppstoreOutlined,
  FunnelPlotOutlined,
  HolderOutlined,
  LogoutOutlined,
  QuestionCircleOutlined,
  SettingOutlined,
  ShopOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { useAuth } from "@/hooks/useAuth";

const { Sider, Content } = Layout;
const SIDER_WIDTH = 220;

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, isLoading, user, logout } = useAuth();
  const [openKeys, setOpenKeys] = useState<string[]>(["supplier-management"]);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isLoading, isAuthenticated, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <Spin size="large" />
      </div>
    );
  }

  if (!isAuthenticated) return null;

  // derive active key
  let selectedKey = "dashboard";
  if (pathname === "/suppliers") selectedKey = "supplier-list";
  else if (pathname.startsWith("/suppliers/reviews")) selectedKey = "review-approvals";
  else if (pathname.startsWith("/suppliers/configurations")) selectedKey = "configurations";
  else if (pathname.startsWith("/suppliers")) selectedKey = "supplier-list";
  else if (pathname.startsWith("/funnel")) selectedKey = "funnel";

  const menuItems = [
    {
      key: "dashboard",
      icon: <AppstoreOutlined />,
      label: "Dashboard",
      onClick: () => router.push("/dashboard"),
    },
    {
      key: "supplier-management",
      icon: <ShopOutlined />,
      label: "Supplier Management",
      children: [
        {
          key: "supplier-dashboard",
          label: "Dashboard",
          onClick: () => router.push("/suppliers/dashboard"),
        },
        {
          key: "supplier-list",
          label: "Supplier List",
          onClick: () => router.push("/suppliers"),
        },
        {
          key: "review-approvals",
          label: "Review & Approvals",
          onClick: () => router.push("/suppliers/reviews"),
        },
        {
          key: "configurations",
          label: "Configurations",
          onClick: () => router.push("/suppliers/configurations"),
        },
      ],
    },
    {
      key: "funnel",
      icon: <FunnelPlotOutlined />,
      label: "Funnel Management",
      onClick: () => router.push("/funnel"),
    },
    { type: "divider" as const },
    {
      key: "help",
      icon: <QuestionCircleOutlined />,
      label: "Help & Support",
    },
    {
      key: "settings",
      icon: <SettingOutlined />,
      label: "Settings",
    },
  ];

  const userMenu = [
    {
      key: "profile",
      icon: <UserOutlined />,
      label: user?.name ?? "Profile",
    },
    { type: "divider" as const },
    {
      key: "logout",
      icon: <LogoutOutlined />,
      label: "Log Out",
      danger: true,
      onClick: logout,
    },
  ];

  return (
    <Layout className="min-h-screen">
      {/* ── Sidebar ───────────────────────────────────────── */}
      <Sider
        width={SIDER_WIDTH}
        className="!fixed left-0 top-0 bottom-0 overflow-auto flex flex-col !bg-white border-r border-gray-200 shadow-md"
        style={{ background: "#fff" }}
      >
        {/* Branding */}
        <div className="flex items-center gap-2.5 px-5 py-4 border-b border-gray-100">
          <div className="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center flex-shrink-0">
            <HolderOutlined className="text-white text-base" />
          </div>
          <div className="min-w-0">
            <p className="text-sm font-bold text-gray-800 leading-tight truncate">
              ALISA
            </p>
            <p className="text-[10px] text-gray-400 leading-tight uppercase tracking-wider">
              Enterprise Portal
            </p>
          </div>
        </div>

        {/* Navigation */}
        <div className="flex-1 overflow-y-auto py-2">
          <Menu
            mode="inline"
            selectedKeys={[selectedKey]}
            openKeys={openKeys}
            onOpenChange={setOpenKeys}
            items={menuItems}
            className="!border-none !bg-transparent [&_.ant-menu-item]:rounded-lg [&_.ant-menu-item]:mx-2 [&_.ant-menu-submenu-title]:rounded-lg [&_.ant-menu-submenu-title]:mx-2"
            style={{ background: "transparent" }}
          />
        </div>

        {/* User at bottom */}
        <div className="border-t border-gray-100 p-3">
          <Dropdown menu={{ items: userMenu }} placement="topLeft" trigger={["click"]}>
            <div className="flex items-center gap-2.5 px-2 py-2 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors">
              <Avatar
                size={36}
                icon={<UserOutlined />}
                className="!bg-blue-100 !text-blue-600 flex-shrink-0"
              />
              <div className="min-w-0 flex-1">
                <p className="text-sm font-semibold text-gray-800 leading-tight truncate">
                  {user?.name ?? "John Doe"}
                </p>
                <p className="text-[11px] text-gray-400 leading-tight uppercase tracking-wide truncate">
                  {user?.roles?.[0] ?? "Senior Manager"}
                </p>
              </div>
              <LogoutOutlined
                className="text-gray-400 text-xs flex-shrink-0"
                onClick={(e) => {
                  e.stopPropagation();
                  logout();
                }}
                title="Log Out"
              />
            </div>
          </Dropdown>
        </div>
      </Sider>

      {/* ── Main ──────────────────────────────────────────── */}
      <Layout style={{ marginLeft: SIDER_WIDTH }}>
        {/* Page content */}
        <Content className="p-6 min-h-screen bg-[#eef0f5]">
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}
