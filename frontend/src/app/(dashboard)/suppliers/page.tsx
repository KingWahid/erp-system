"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Button, Dropdown, Select, Table, App, Tooltip, Input } from "antd";
import type { MenuProps } from "antd";
import {
  UserAddOutlined,
  AppstoreOutlined,
  UnorderedListOutlined,
  EllipsisOutlined,
  EyeOutlined,
  EditOutlined,
  StopOutlined,
  SearchOutlined,
} from "@ant-design/icons";
import type { ColumnsType } from "antd/es/table";
import { useRouter } from "next/navigation";
import {
  getSupplierStats,
  listSuppliers,
  blockSupplier,
} from "@/services/supplier.service";
import { usePermission } from "@/hooks/usePermission";
import type { SupplierListItem, SupplierStatus } from "@/types/api";

// Status pill

const STATUS_CFG: Record<
  SupplierStatus,
  { label: string; dot: string; pill: string }
> = {
  active:      { label: "Active",      dot: "bg-green-500",  pill: "bg-green-50  text-green-700"  },
  draft:       { label: "Draft",       dot: "bg-gray-400",   pill: "bg-gray-100  text-gray-600"   },
  in_progress: { label: "In Progress", dot: "bg-orange-400", pill: "bg-orange-50 text-orange-700" },
  blocked:     { label: "Blocked",     dot: "bg-red-500",    pill: "bg-red-50    text-red-700"    },
  inactive:    { label: "Inactive",    dot: "bg-gray-400",   pill: "bg-gray-100  text-gray-600"   },
};

function StatusPill({ status }: { status: SupplierStatus }) {
  const c = STATUS_CFG[status];
  return (
    <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium ${c.pill}`}>
      <span className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${c.dot}`} />
      {c.label}
    </span>
  );
}

// Stat card — matches wireframe: icon top-left, trend top-right,
//             label small uppercase, value large bold, subtitle

interface StatCardProps {
  label: string;
  value: string | number;
  subValue?: string;
  trend: number;
  trendSuffix?: string;
  icon: React.ReactNode;
  iconBg: string;
  trendPositive?: boolean; // override direction
}

function StatCard({ label, value, subValue, trend, trendSuffix = "vs last year", icon, iconBg, trendPositive }: StatCardProps) {
  const up = trendPositive ?? trend >= 0;
  return (
    <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-5 flex flex-col gap-2">
      {/* top row */}
      <div className="flex items-start justify-between">
        <div className={`w-10 h-10 rounded-xl flex items-center justify-center text-lg ${iconBg}`}>
          {icon}
        </div>
        <span className={`flex items-center gap-0.5 text-xs font-semibold ${up ? "text-green-500" : "text-red-500"}`}>
          {/* tiny arrow */}
          <svg className="w-2.5 h-2.5" viewBox="0 0 10 10" fill="currentColor">
            {up
              ? <path d="M5 2 L9 8 L1 8 Z" />
              : <path d="M5 8 L9 2 L1 2 Z" />}
          </svg>
          {up ? "+" : ""}{trend}%
        </span>
      </div>

      {/* label */}
      <p className="text-[11px] text-gray-400 uppercase tracking-widest leading-none">{label}</p>

      {/* value */}
      <div>
        <p className="text-[2rem] font-extrabold text-gray-900 leading-none">{value}</p>
        {subValue && (
          <p className="text-[1.25rem] font-extrabold text-gray-900 leading-tight -mt-0.5">{subValue}</p>
        )}
      </div>

      {/* trend label */}
      <p className="text-xs text-gray-400">{trendSuffix}</p>
    </div>
  );
}

// Supplier logo placeholder

function SupplierLogo({ status }: { status: SupplierStatus }) {
  const border: Record<SupplierStatus, string> = {
    active:      "border-blue-200 bg-blue-50",
    draft:       "border-gray-200 bg-gray-50",
    in_progress: "border-orange-200 bg-orange-50",
    blocked:     "border-red-200 bg-red-50",
    inactive:    "border-gray-200 bg-gray-50",
  };
  return (
    <div className={`w-9 h-9 rounded-lg border flex-shrink-0 flex items-center justify-center ${border[status]}`}>
      {/* simple "document" icon drawn with divs */}
      <div className="flex flex-col gap-[3px]">
        <div className="w-5 h-[2px] bg-gray-400 rounded-full" />
        <div className="w-4 h-[2px] bg-gray-400 rounded-full" />
        <div className="w-5 h-[2px] bg-gray-400 rounded-full" />
      </div>
    </div>
  );
}

// Format helpers

function fmtMn(v: number) {
  if (v === 0) return { main: "0", sub: undefined };
  const mn = v / 1_000_000;
  const fixed = mn.toFixed(1);
  const [int, dec] = fixed.split(".");
  return { main: `Rp ${int},${dec}`, sub: "Mn" };
}

// Page

export default function SuppliersPage() {
  const router = useRouter();
  const qc = useQueryClient();
  const { can } = usePermission();
  const { modal, message } = App.useApp();

  const PAGE_SIZE = 10;
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<SupplierStatus | undefined>(undefined);
  const [viewMode, setViewMode] = useState<"list" | "grid">("list");

  const { data: statsData } = useQuery({
    queryKey: ["supplier-stats"],
    queryFn: getSupplierStats,
  });

  const { data: listData, isLoading } = useQuery({
    queryKey: ["suppliers", page, statusFilter, search],
    queryFn: () => listSuppliers({ page, limit: PAGE_SIZE, status: statusFilter, search: search || undefined }),
  });

  const blockMut = useMutation({
    mutationFn: ({ id, block }: { id: string; block: boolean }) =>
      blockSupplier(id, { block }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["suppliers"] });
      qc.invalidateQueries({ queryKey: ["supplier-stats"] });
      message.success("Status berhasil diperbarui");
    },
  });

  const stats = statsData?.data;
  const total = listData?.meta?.total ?? 0;
  const rows = listData?.data ?? [];
  const totalPages = Math.ceil(total / PAGE_SIZE);

  // Per-row dropdown
  const rowMenu = (r: SupplierListItem): MenuProps["items"] => [
    {
      key: "view",
      icon: <EyeOutlined />,
      label: "Lihat Detail",
      onClick: () => router.push(`/suppliers/${r.id}`),
    },
    {
      key: "edit",
      icon: <EditOutlined />,
      label: "Edit",
      onClick: () => router.push(`/suppliers/${r.id}`),
    },
    { type: "divider" },
    {
      key: "block",
      icon: <StopOutlined />,
      label: r.status === "blocked" ? "Buka Blokir" : "Blokir",
      danger: r.status !== "blocked",
      onClick: () => {
        const blocked = r.status === "blocked";
        modal.confirm({
          title: blocked ? "Buka blokir supplier?" : "Blokir supplier ini?",
          content: r.name,
          okButtonProps: { danger: !blocked },
          okText: blocked ? "Buka Blokir" : "Blokir",
          onOk: () => blockMut.mutate({ id: r.id, block: !blocked }),
        });
      },
    },
  ];

  // Table columns
  const columns: ColumnsType<SupplierListItem> = [
    {
      title: <span className="text-xs font-semibold text-gray-400 uppercase">#</span>,
      key: "no",
      width: 52,
      align: "center" as const,
      render: (_: unknown, __: unknown, i: number) => (
        <span className="text-sm text-gray-500 font-medium">{(page - 1) * PAGE_SIZE + i + 1}</span>
      ),
    },
    {
      title: <span className="text-xs font-semibold text-gray-400 uppercase">Name</span>,
      key: "name",
      render: (_: unknown, r: SupplierListItem) => (
        <div className="flex items-center gap-3">
          <SupplierLogo status={r.status} />
          <div className="min-w-0">
            <button
              className="font-semibold text-blue-600 hover:underline text-sm leading-tight bg-transparent border-none cursor-pointer p-0 text-left block"
              onClick={(e) => { e.stopPropagation(); router.push(`/suppliers/${r.id}`); }}
            >
              {r.name}
            </button>
            <p className="text-[11px] text-gray-400 mt-0.5 m-0 leading-tight">
              <span className="text-gray-500 font-medium mr-0.5">{r.code}</span>
              <button
                className="text-blue-500 hover:underline bg-transparent border-none cursor-pointer p-0 text-[11px] mr-0.5"
                onClick={(e) => { e.stopPropagation(); router.push(`/suppliers/${r.id}`); }}
              >
                {r.supplier_no}
              </button>
              {r.alias && <span className="text-gray-400">[{r.alias}]</span>}
            </p>
          </div>
        </div>
      ),
    },
    {
      title: <span className="text-xs font-semibold text-gray-400 uppercase">Address</span>,
      dataIndex: "address",
      key: "address",
      render: (v: string) => <span className="text-sm text-gray-600">{v || "-"}</span>,
    },
    {
      title: <span className="text-xs font-semibold text-gray-400 uppercase">Contact</span>,
      dataIndex: "contact",
      key: "contact",
      render: (v: string) => <span className="text-sm text-gray-600">{v || "-"}</span>,
    },
    {
      title: <span className="text-xs font-semibold text-gray-400 uppercase">Status</span>,
      dataIndex: "status",
      key: "status",
      width: 140,
      render: (s: SupplierStatus) => <StatusPill status={s} />,
    },
    {
      title: <span className="text-xs font-semibold text-gray-400 uppercase">Actions</span>,
      key: "actions",
      width: 80,
      align: "center" as const,
      render: (_: unknown, r: SupplierListItem) => (
        <Dropdown menu={{ items: rowMenu(r) }} trigger={["click"]} placement="bottomRight">
          <button
            className="w-8 h-8 rounded-lg hover:bg-gray-100 flex items-center justify-center text-gray-400 hover:text-gray-700 transition-colors mx-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <EllipsisOutlined className="text-base" />
          </button>
        </Dropdown>
      ),
    },
  ];

  // Pagination pages array
  const pagesArr = (): (number | "...")[] => {
    if (totalPages <= 5) return Array.from({ length: totalPages }, (_, i) => i + 1);
    const arr: (number | "...")[] = [1];
    if (page > 3) arr.push("...");
    for (let i = Math.max(2, page - 1); i <= Math.min(totalPages - 1, page + 1); i++) arr.push(i);
    if (page < totalPages - 2) arr.push("...");
    arr.push(totalPages);
    return arr;
  };

  const { main: costMain, sub: costSub } = fmtMn(stats?.avg_cost_supplier ?? 0);

  return (
    <div className="min-h-full">
      {/* ── Header ─────────────────────────────────────── */}
      <div className="flex items-start justify-between mb-6">
        <div>
          <h1 className="text-xl font-bold text-gray-800 m-0 leading-tight">Supplier List</h1>
          <p className="text-sm text-gray-400 mt-1 m-0">
            Overview of your enterprise supply network
          </p>
        </div>
        {can("supplier:create") && (
          <Button
            type="primary"
            icon={<UserAddOutlined />}
            onClick={() => router.push("/suppliers/new")}
            className="!rounded-xl !h-9 !px-4 !font-semibold !text-sm"
          >
            New Supplier
          </Button>
        )}
      </div>

      {/* ── Stat cards ─────────────────────────────────── */}
      <div className="grid grid-cols-4 gap-4 mb-6">
        <StatCard
          label="Total Supplier"
          value={(stats?.total_supplier ?? 0).toLocaleString("id-ID")}
          trend={8}
          trendSuffix="vs last year"
          icon={
            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.8}>
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
          }
          iconBg="bg-blue-50 text-blue-500"
        />
        <StatCard
          label="New Supplier"
          value={stats?.new_supplier ?? 0}
          trend={1}
          trendSuffix="vs Last Year"
          icon={
            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.8}>
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <line x1="19" y1="8" x2="19" y2="14"/>
              <line x1="22" y1="11" x2="16" y2="11"/>
            </svg>
          }
          iconBg="bg-cyan-50 text-cyan-500"
        />
        <StatCard
          label="Avg Cost per Supplier"
          value={costMain}
          subValue={costSub}
          trend={-1}
          trendPositive={false}
          trendSuffix="vs Last Year"
          icon={
            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.8}>
              <rect x="2" y="5" width="20" height="14" rx="2"/>
              <line x1="2" y1="10" x2="22" y2="10"/>
            </svg>
          }
          iconBg="bg-orange-50 text-orange-500"
        />
        <StatCard
          label="Blocked Supplier"
          value={stats?.blocked_supplier ?? 0}
          trend={-4}
          trendPositive={false}
          trendSuffix="vs Last Year"
          icon={
            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.8}>
              <circle cx="12" cy="12" r="10"/>
              <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
            </svg>
          }
          iconBg="bg-red-50 text-red-400"
        />
      </div>

      {/* ── Table card ─────────────────────────────────── */}
      <div className="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">

        {/* Toolbar */}
        <div className="flex items-center justify-between px-5 py-3 border-b border-gray-100">
          <div className="flex items-center gap-2">
            <Input
              placeholder="Search supplier…"
              prefix={<SearchOutlined className="text-gray-400" />}
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1); }}
              allowClear
              className="!w-56 !rounded-xl !text-sm"
            />
            <Select
              value={statusFilter}
              placeholder="All Status"
              allowClear
              onChange={(val) => { setStatusFilter(val as SupplierStatus | undefined); setPage(1); }}
              className="w-36"
              options={[
                { value: "active",      label: "Active"      },
                { value: "draft",       label: "Draft"       },
                { value: "in_progress", label: "In Progress" },
                { value: "blocked",     label: "Blocked"     },
                { value: "inactive",    label: "Inactive"    },
              ]}
            />
          </div>

          <div className="flex items-center gap-1.5">
            <Tooltip title="Export CSV">
              <button
                onClick={() => {
                  if (!rows.length) return;
                  const headers = ["#", "Code", "Supplier No", "Name", "Alias", "Address", "Contact", "Status"];
                  const csvRows = rows.map((r, i) => [
                    (page - 1) * PAGE_SIZE + i + 1,
                    r.code,
                    r.supplier_no,
                    r.name,
                    r.alias ?? "",
                    r.address ?? "",
                    r.contact ?? "",
                    r.status,
                  ].map((v) => `"${String(v).replace(/"/g, '""')}"`).join(","));
                  const csv = [headers.join(","), ...csvRows].join("\n");
                  const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
                  const url = URL.createObjectURL(blob);
                  const a = document.createElement("a");
                  a.href = url;
                  a.download = `suppliers-${new Date().toISOString().slice(0, 10)}.csv`;
                  a.click();
                  URL.revokeObjectURL(url);
                }}
                className="flex items-center gap-1.5 px-3 h-8 rounded-lg border border-gray-200 text-xs font-medium text-gray-600 hover:border-blue-500 hover:text-blue-600 transition-colors"
              >
                <svg className="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2}>
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" strokeLinecap="round" strokeLinejoin="round" />
                  <polyline points="7 10 12 15 17 10" strokeLinecap="round" strokeLinejoin="round" />
                  <line x1="12" y1="15" x2="12" y2="3" strokeLinecap="round" />
                </svg>
                Export
              </button>
            </Tooltip>
            <div className="w-px h-5 bg-gray-200 mx-0.5" />
            <Tooltip title="Grid view">
              <button
                onClick={() => setViewMode("grid")}
                className={`w-8 h-8 rounded-lg flex items-center justify-center transition-colors text-sm ${
                  viewMode === "grid" ? "bg-blue-600 text-white" : "text-gray-400 hover:bg-gray-100"
                }`}
              >
                <AppstoreOutlined />
              </button>
            </Tooltip>
            <Tooltip title="List view">
              <button
                onClick={() => setViewMode("list")}
                className={`w-8 h-8 rounded-lg flex items-center justify-center transition-colors text-sm ${
                  viewMode === "list" ? "bg-blue-600 text-white" : "text-gray-400 hover:bg-gray-100"
                }`}
              >
                <UnorderedListOutlined />
              </button>
            </Tooltip>
          </div>
        </div>

        {/* Table */}
        <div style={{ "--ant-table-border-color": "#e5e7eb" } as React.CSSProperties}>
        <Table<SupplierListItem>
          rowKey="id"
          columns={columns}
          dataSource={rows}
          loading={isLoading}
          size="middle"
          pagination={false}
          className={[
            "[&_.ant-table-thead>tr>th]:bg-gray-50",
            "[&_.ant-table-thead>tr>th]:text-gray-500",
            "[&_.ant-table-thead>tr>th]:font-semibold",
            "[&_.ant-table-thead>tr>th]:py-3",
            "[&_.ant-table-tbody>tr>td]:py-3.5",
            "[&_.ant-table-row:hover>td]:!bg-blue-50/30",
            "[&_.ant-table-row]:cursor-pointer",
            "[&_.ant-table]:border-none",
            "[&_.ant-table-container]:border-none",
          ].join(" ")}
          onRow={(r) => ({ onClick: () => router.push(`/suppliers/${r.id}`) })}
        />
        </div>

        {/* Footer pagination */}
        <div className="flex items-center justify-between px-5 py-3.5 border-t border-gray-100">
          <p className="text-sm text-gray-500 m-0">
            Showing{" "}
            <span className="font-semibold text-gray-800">{Math.min(page * PAGE_SIZE, total)}</span>
            {" "}of{" "}
            <span className="font-semibold text-gray-800">{total.toLocaleString("id-ID")}</span>
            {" "}suppliers
          </p>

          <div className="flex items-center gap-1">
            {/* prev */}
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="w-8 h-8 rounded-lg border border-gray-200 flex items-center justify-center text-gray-500 text-base hover:border-blue-500 hover:text-blue-600 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              ‹
            </button>

            {pagesArr().map((p, i) =>
              p === "..." ? (
                <span key={`e${i}`} className="w-8 text-center text-gray-400 text-sm select-none">…</span>
              ) : (
                <button
                  key={p}
                  onClick={() => setPage(p as number)}
                  className={`w-8 h-8 rounded-lg text-sm font-medium transition-colors ${
                    page === p
                      ? "bg-blue-600 text-white border border-blue-600"
                      : "border border-gray-200 text-gray-600 hover:border-blue-500 hover:text-blue-600"
                  }`}
                >
                  {p}
                </button>
              )
            )}

            {/* next */}
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="w-8 h-8 rounded-lg border border-gray-200 flex items-center justify-center text-gray-500 text-base hover:border-blue-500 hover:text-blue-600 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              ›
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
