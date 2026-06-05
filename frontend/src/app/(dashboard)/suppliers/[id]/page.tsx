"use client";

import { use, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Avatar,
  Breadcrumb,
  Skeleton,
  Tag,
  Tabs,
  Table,
  Input,
  App,
  Rate,
} from "antd";
import {
  ArrowLeftOutlined,
  StopOutlined,
  CheckCircleOutlined,
  ArrowRightOutlined,
  EnvironmentOutlined,
  PhoneOutlined,
  MailOutlined,
  GlobalOutlined,
  EditOutlined,
  ShopOutlined,
  TagsOutlined,
  AppstoreOutlined,
  UserOutlined,
  FileTextOutlined,
  StarOutlined,
  AuditOutlined,
  ProjectOutlined,
  ToolOutlined,
  HistoryOutlined,
} from "@ant-design/icons";
import { useRouter } from "next/navigation";
import {
  getSupplier,
  blockSupplier,
  advanceSupplierStage,
  listRatings,
  listOutstandings,
} from "@/services/supplier.service";
import { usePermission } from "@/hooks/usePermission";
import type { SupplierStage, SupplierStatus, RatingItem, InvoiceItem } from "@/types/api";

const { TextArea } = Input;

const STATUS_CFG: Record<SupplierStatus, { label: string; pill: string; dot: string }> = {
  active:      { label: "Active",      pill: "bg-green-50 text-green-700",   dot: "bg-green-500"  },
  draft:       { label: "Draft",       pill: "bg-gray-100 text-gray-600",    dot: "bg-gray-400"   },
  in_progress: { label: "In Progress", pill: "bg-orange-50 text-orange-700", dot: "bg-orange-400" },
  blocked:     { label: "Blocked",     pill: "bg-red-50 text-red-700",       dot: "bg-red-500"    },
  inactive:    { label: "Inactive",    pill: "bg-gray-100 text-gray-600",    dot: "bg-gray-400"   },
};

const STAGE_STEPS: SupplierStage[] = ["draft", "in_review", "in_assessment", "active"];
const STAGE_LABEL: Record<SupplierStage, string> = {
  draft:         "Draft",
  in_review:     "In Review",
  in_assessment: "In Assessment",
  active:        "Active",
};

function StatusPill({ status }: { status: SupplierStatus }) {
  const c = STATUS_CFG[status];
  return (
    <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold ${c.pill}`}>
      <span className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${c.dot}`} />
      {c.label}
    </span>
  );
}

function SupplierLogo({ status }: { status: SupplierStatus }) {
  const bg: Record<SupplierStatus, string> = {
    active:      "bg-blue-50   border-blue-200",
    draft:       "bg-gray-50   border-gray-200",
    in_progress: "bg-orange-50 border-orange-200",
    blocked:     "bg-red-50    border-red-200",
    inactive:    "bg-gray-50   border-gray-200",
  };
  return (
    <div className={`w-14 h-14 rounded-xl border-2 flex items-center justify-center flex-shrink-0 ${bg[status]}`}>
      <div className="flex flex-col gap-[3px]">
        <div className="w-6 h-[2px] bg-gray-400 rounded-full" />
        <div className="w-5 h-[2px] bg-gray-400 rounded-full" />
        <div className="w-6 h-[2px] bg-gray-400 rounded-full" />
      </div>
    </div>
  );
}

function DetailRow({
  label,
  value,
  icon,
}: {
  label: string;
  value?: React.ReactNode;
  icon?: React.ReactNode;
}) {
  if (!value) return null;
  return (
    <div className="flex items-start gap-3 py-2.5 border-b border-gray-100 last:border-0">
      {icon && (
        <span className="text-gray-400 mt-0.5 text-sm w-4 flex-shrink-0">{icon}</span>
      )}
      <span className="text-xs text-gray-400 w-32 shrink-0 pt-0.5 leading-tight">{label}</span>
      <span className="text-sm text-gray-700 flex-1 leading-tight">{value}</span>
    </div>
  );
}

function StageStepper({
  steps,
  current,
}: {
  steps: { key: string; label: string }[];
  current: number;
}) {
  return (
    <div className="flex items-center w-full">
      {steps.map((step, i) => {
        const done    = i < current;
        const active  = i === current;
        const pending = i > current;
        return (
          <div key={step.key} className="flex items-center flex-1 last:flex-none">
            <div className="flex flex-col items-center gap-1">
              <div
                className={`w-6 h-6 rounded-full border-2 flex items-center justify-center transition-colors ${
                  done
                    ? "border-green-500 bg-green-500"
                    : active
                    ? "border-blue-500 bg-white"
                    : "border-gray-300 bg-white"
                }`}
              >
                {done ? (
                  <svg className="w-3 h-3 text-white" viewBox="0 0 12 12" fill="none" stroke="currentColor" strokeWidth={2}>
                    <path d="M2 6 L5 9 L10 3" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                ) : (
                  <span
                    className={`w-2 h-2 rounded-full ${
                      active ? "bg-blue-500" : "bg-gray-300"
                    }`}
                  />
                )}
              </div>
              <span
                className={`text-[10px] font-medium whitespace-nowrap ${
                  done ? "text-green-600" : active ? "text-blue-600" : "text-gray-400"
                }`}
              >
                {step.label}
              </span>
            </div>
            {i < steps.length - 1 && (
              <div
                className={`flex-1 h-[2px] mx-1 -mt-4 rounded-full transition-colors ${
                  i < current ? "bg-green-400" : "bg-gray-200"
                }`}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}

function StarDisplay({ value, max = 5 }: { value: number; max?: number }) {
  return (
    <span className="inline-flex gap-0.5">
      {Array.from({ length: max }).map((_, i) => (
        <svg
          key={i}
          className={`w-3.5 h-3.5 ${i < value ? "text-amber-400" : "text-gray-200"}`}
          viewBox="0 0 24 24"
          fill="currentColor"
        >
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
      ))}
    </span>
  );
}

function RatingCard({ rating }: { rating: RatingItem }) {
  const date = new Date(rating.reviewed_at).toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
  return (
    <div className="bg-gray-50 border border-gray-100 rounded-xl p-3 space-y-1.5">
      <div className="flex items-center gap-2">
        <span className="text-[11px] text-gray-400 w-24 shrink-0">Price</span>
        <StarDisplay value={rating.price_rating} />
      </div>
      <div className="flex items-center gap-2">
        <span className="text-[11px] text-gray-400 w-24 shrink-0">Delivery Time</span>
        <StarDisplay value={rating.delivery_rating} />
      </div>
      {rating.notes && (
        <div className="flex items-start gap-2">
          <span className="text-[11px] text-gray-400 w-24 shrink-0">Notes</span>
          <span className="text-[11px] text-gray-600">{rating.notes}</span>
        </div>
      )}
      <p className="text-[10px] text-gray-400 pt-1">
        {date}
        {rating.reviewed_by ? ` by ${rating.reviewed_by}` : ""}
      </p>
    </div>
  );
}

const TABLE_CLS = [
  "[&_.ant-table-thead>tr>th]:bg-gray-50",
  "[&_.ant-table-thead>tr>th]:text-xs",
  "[&_.ant-table-thead>tr>th]:font-semibold",
  "[&_.ant-table-thead>tr>th]:text-gray-500",
  "[&_.ant-table-thead>tr>th]:uppercase",
  "[&_.ant-table-thead>tr>th]:py-3",
  "[&_.ant-table-tbody>tr>td]:py-3",
  "[&_.ant-table-tbody>tr>td]:border-b",
  "[&_.ant-table-tbody>tr>td]:border-gray-50",
  "[&_.ant-table]:border-none",
  "[&_.ant-table-container]:border-none",
].join(" ");

export default function SupplierDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();
  const qc = useQueryClient();
  const { can } = usePermission();
  const { message, modal } = App.useApp();

  const [stageNotes, setStageNotes] = useState("");

  const { data, isLoading } = useQuery({
    queryKey: ["supplier", id],
    queryFn: () => getSupplier(id),
    enabled: !!id,
  });

  const { data: ratingsData } = useQuery({
    queryKey: ["supplier-ratings", id],
    queryFn: () => listRatings(id),
    enabled: !!id,
  });

  const { data: outstandingsData, isLoading: outstandingsLoading } = useQuery({
    queryKey: ["supplier-outstandings", id],
    queryFn: () => listOutstandings(id, { page: 1, limit: 50 }),
    enabled: !!id,
  });

  const blockMut = useMutation({
    mutationFn: ({ block, reason }: { block: boolean; reason?: string }) =>
      blockSupplier(id, { block, reason }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["supplier", id] });
      message.success("Status supplier berhasil diperbarui");
    },
    onError: () => message.error("Gagal memperbarui status"),
  });

  const stageMut = useMutation({
    mutationFn: (notes?: string) => advanceSupplierStage(id, { notes }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["supplier", id] });
      setStageNotes("");
      message.success("Stage berhasil dimajukan");
    },
    onError: () => message.error("Gagal memajukan stage"),
  });

  const handleBlock = () => {
    const isBlocked = supplier?.is_blocked ?? false;
    if (isBlocked) {
      modal.confirm({
        title: "Buka blokir supplier?",
        content: "Supplier akan dapat diakses kembali.",
        okText: "Buka Blokir",
        onOk: () => blockMut.mutate({ block: false }),
      });
    } else {
      let reason = "";
      modal.confirm({
        title: "Blokir Supplier",
        content: (
          <div>
            <p className="mb-2 text-sm text-gray-600">Masukkan alasan pemblokiran:</p>
            <TextArea
              rows={3}
              placeholder="Contoh: Payment default"
              onChange={(e) => { reason = e.target.value; }}
            />
          </div>
        ),
        okText: "Blokir",
        okButtonProps: { danger: true },
        onOk: () => blockMut.mutate({ block: true, reason }),
      });
    }
  };

  const handleNextStage = () => {
    const nextStage = STAGE_STEPS[STAGE_STEPS.indexOf(supplier!.stage) + 1];
    modal.confirm({
      title: `Majukan ke stage "${STAGE_LABEL[nextStage]}"?`,
      content: stageNotes
        ? <p className="text-sm text-gray-600">Catatan: <em>{stageNotes}</em></p>
        : <p className="text-sm text-gray-400">Tidak ada catatan.</p>,
      okText: "Konfirmasi",
      onOk: () => stageMut.mutate(stageNotes || undefined),
    });
  };

  const supplierRaw = data?.data;
  // Defensive defaults — backend may return null for sub-arrays
  const supplier = supplierRaw
    ? {
        ...supplierRaw,
        addresses:      supplierRaw.addresses      ?? [],
        contacts:       supplierRaw.contacts       ?? [],
        groups:         supplierRaw.groups         ?? [],
        materials:      supplierRaw.materials      ?? [],
        stage_histories: supplierRaw.stage_histories ?? [],
      }
    : null;
  const ratings   = ratingsData?.data ?? [];
  const outstandings = outstandingsData?.data ?? [];
  const stageIdx  = supplier ? STAGE_STEPS.indexOf(supplier.stage) : 0;
  const canAdvance = can("workflow:advance") && !!supplier && stageIdx < STAGE_STEPS.length - 1;
  const location  = supplier
    ? [supplier.address, supplier.city, supplier.country].filter(Boolean).join(", ")
    : "";

  if (isLoading) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-6">
        <Skeleton active paragraph={{ rows: 16 }} />
      </div>
    );
  }

  if (!supplier) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 p-12 flex flex-col items-center gap-4">
        <ShopOutlined className="text-5xl text-gray-300" />
        <p className="text-gray-400">Supplier tidak ditemukan</p>
        <button
          onClick={() => router.back()}
          className="px-4 h-9 rounded-xl border border-gray-200 text-sm text-gray-600 hover:border-blue-500 hover:text-blue-500 transition-colors"
        >
          Kembali
        </button>
      </div>
    );
  }

  const overviewSubTabs = [
    {
      key: "general",
      label: "General",
      children: (
        <div className="p-4 space-y-4">
          {/* General info */}
          <div className="bg-gray-50 border border-gray-200 rounded-xl p-4">
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
              Informasi Umum
            </p>
            <DetailRow label="Nama Perusahaan" value={supplier.name} />
            <DetailRow label="Alias" value={supplier.alias} />
            <DetailRow label="Kode" value={supplier.code} />
            <DetailRow label="No. Supplier" value={supplier.supplier_no} />
            <DetailRow label="SLA" value={`${supplier.sla_hours} jam`} />
            <DetailRow label="Status" value={<StatusPill status={supplier.status} />} />
            {supplier.block_reason && (
              <DetailRow
                label="Alasan Blokir"
                value={<span className="text-red-500">{supplier.block_reason}</span>}
              />
            )}
            {supplier.notes && (
              <DetailRow label="Catatan" value={supplier.notes} />
            )}
          </div>

          {/* Contact & location */}
          <div className="bg-gray-50 border border-gray-200 rounded-xl p-4">
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
              Kontak &amp; Lokasi
            </p>
            <DetailRow
              icon={<MailOutlined />}
              label="Email"
              value={
                supplier.email
                  ? <a href={`mailto:${supplier.email}`} className="text-blue-500">{supplier.email}</a>
                  : undefined
              }
            />
            <DetailRow icon={<PhoneOutlined />} label="Telepon" value={supplier.phone} />
            <DetailRow
              icon={<GlobalOutlined />}
              label="Website"
              value={
                supplier.website
                  ? <a href={supplier.website} target="_blank" rel="noreferrer" className="text-blue-500">{supplier.website}</a>
                  : undefined
              }
            />
            <DetailRow
              icon={<EnvironmentOutlined />}
              label="Alamat"
              value={location || undefined}
            />
          </div>
        </div>
      ),
    },
    {
      key: "address",
      label: `Address (${supplier.addresses.length})`,
      children: (
        <div className="p-4">
          <Table
            rowKey="id"
            size="small"
            dataSource={supplier.addresses}
            pagination={false}
            className={TABLE_CLS}
            columns={[
              { title: "Nama",          dataIndex: "name",    key: "name" },
              {
                title: "Alamat Lengkap",
                key: "addr",
                render: (_: unknown, r: typeof supplier.addresses[0]) =>
                  [r.address, r.city, r.province, r.country, r.postal_code]
                    .filter(Boolean)
                    .join(", "),
              },
              {
                title: "Utama",
                dataIndex: "is_main",
                key: "is_main",
                width: 80,
                render: (v: boolean) =>
                  v ? <Tag color="blue">Utama</Tag> : null,
              },
            ]}
          />
        </div>
      ),
    },
    {
      key: "contacts",
      label: `Contacts (${supplier.contacts.length})`,
      children: (
        <div className="p-4">
          <Table
            rowKey="id"
            size="small"
            dataSource={supplier.contacts}
            pagination={false}
            className={TABLE_CLS}
            columns={[
              { title: "Nama",     dataIndex: "name",     key: "name"     },
              { title: "Jabatan",  dataIndex: "position", key: "position" },
              { title: "Telepon",  dataIndex: "phone",    key: "phone"    },
              { title: "Mobile",   dataIndex: "mobile",   key: "mobile"   },
              { title: "Email",    dataIndex: "email",    key: "email"    },
              {
                title: "Primary",
                dataIndex: "is_primary",
                key: "is_primary",
                width: 80,
                render: (v: boolean) =>
                  v ? <Tag color="blue">Utama</Tag> : null,
              },
            ]}
          />
        </div>
      ),
    },
    {
      key: "groups",
      label: `Groups (${supplier.groups.length})`,
      children: (
        <div className="p-4">
          {supplier.groups.length === 0 ? (
            <p className="text-sm text-gray-400 py-4 text-center">Belum ada grup</p>
          ) : (
            <div className="flex flex-wrap gap-2 py-1">
              {supplier.groups.map((g) => (
                <Tag
                  key={g.id}
                  color={g.is_active ? "blue" : "default"}
                  className="rounded-full px-3 py-1 text-xs"
                >
                  {g.group_name}: {g.value}
                </Tag>
              ))}
            </div>
          )}
        </div>
      ),
    },
    {
      key: "material-list",
      label: `Material List (${supplier.materials.length})`,
      children: (
        <div className="p-4">
          <Table
            rowKey="id"
            size="small"
            dataSource={supplier.materials}
            pagination={false}
            className={TABLE_CLS}
            columns={[
              { title: "Grup Material", dataIndex: "material_group", key: "material_group" },
              { title: "Material ID",   dataIndex: "material_id",    key: "material_id"    },
              {
                title: "Status",
                dataIndex: "is_active",
                key: "is_active",
                width: 100,
                render: (v: boolean) => (
                  <Tag color={v ? "success" : "default"}>{v ? "Aktif" : "Nonaktif"}</Tag>
                ),
              },
            ]}
          />
        </div>
      ),
    },
  ];

  const mainTabs = [
    {
      key: "overview",
      label: (
        <span className="flex items-center gap-1.5">
          <AppstoreOutlined />Overview
        </span>
      ),
      children: (
        <div>
          {/* Sub-tabs */}
          <Tabs
            size="small"
            items={overviewSubTabs}
            className={[
              "[&_.ant-tabs-nav]:!px-4",
              "[&_.ant-tabs-nav]:!mb-0",
              "[&_.ant-tabs-nav]:border-b",
              "[&_.ant-tabs-nav]:border-gray-100",
              "[&_.ant-tabs-tab]:!text-xs",
              "[&_.ant-tabs-tab]:!py-2",
            ].join(" ")}
          />

          {/* Outstandings — always visible below sub-tabs */}
          <div className="p-4 border-t border-gray-100">
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
              Outstandings
            </p>
            <Table<InvoiceItem>
              rowKey="id"
              size="small"
              dataSource={outstandings}
              loading={outstandingsLoading}
              pagination={false}
              className={TABLE_CLS}
              locale={{ emptyText: "Tidak ada outstanding invoice" }}
              columns={[
                {
                  title: "#",
                  key: "no",
                  width: 48,
                  align: "center" as const,
                  render: (_: unknown, __: InvoiceItem, idx: number) => idx + 1,
                },
                {
                  title: "Invoice Number",
                  dataIndex: "invoice_number",
                  key: "invoice_number",
                  render: (v: string, r: InvoiceItem) => (
                    <span className="font-medium text-gray-800">{v}</span>
                  ),
                },
                { title: "Project Name", dataIndex: "project_name", key: "project_name",
                  render: (v?: string) => v || <span className="text-gray-400">-</span> },
                {
                  title: "Amount",
                  dataIndex: "amount",
                  key: "amount",
                  render: (v: number, r: InvoiceItem) =>
                    <span className="font-medium">
                      {v.toLocaleString("id-ID")}
                    </span>,
                },
                {
                  title: "Status",
                  dataIndex: "status",
                  key: "status",
                  width: 100,
                  render: (v: string) => {
                    const cfg: Record<string, { color: string; label: string }> = {
                      unpaid:  { color: "orange", label: "Unpaid"  },
                      partial: { color: "blue",   label: "Partial" },
                      overdue: { color: "red",    label: "Overdue" },
                      paid:    { color: "green",  label: "Paid"    },
                    };
                    const c = cfg[v] ?? { color: "default", label: v };
                    return <Tag color={c.color}>{c.label}</Tag>;
                  },
                },
                {
                  title: "Aging (days)",
                  dataIndex: "aging_days",
                  key: "aging_days",
                  width: 110,
                  align: "center" as const,
                  render: (v: number) => (
                    <span className={v > 30 ? "text-red-600 font-semibold" : "text-gray-700"}>
                      {v}
                    </span>
                  ),
                },
              ]}
            />
          </div>
        </div>
      ),
    },
    {
      key: "assessment",
      label: (
        <span className="flex items-center gap-1.5">
          <AuditOutlined />Assessment
        </span>
      ),
      children: (
        <div className="p-5 flex flex-col items-center justify-center py-12 gap-3 text-center">
          <AuditOutlined className="text-4xl text-gray-200" />
          <p className="text-sm text-gray-400">Modul assessment belum tersedia</p>
        </div>
      ),
    },
    {
      key: "material-catalog",
      label: (
        <span className="flex items-center gap-1.5">
          <TagsOutlined />Material Catalog
        </span>
      ),
      children: (
        <div className="p-5">
          <Table
            rowKey="id"
            size="small"
            dataSource={supplier.materials}
            pagination={false}
            className={TABLE_CLS}
            columns={[
              { title: "Grup Material", dataIndex: "material_group", key: "material_group" },
              { title: "Material ID",   dataIndex: "material_id",    key: "material_id"    },
              {
                title: "Status",
                dataIndex: "is_active",
                key: "is_active",
                width: 100,
                render: (v: boolean) => (
                  <Tag color={v ? "success" : "default"}>{v ? "Aktif" : "Nonaktif"}</Tag>
                ),
              },
            ]}
          />
        </div>
      ),
    },
    {
      key: "orders",
      label: (
        <span className="flex items-center gap-1.5">
          <FileTextOutlined />Orders
        </span>
      ),
      children: (
        <div className="p-5 flex flex-col items-center justify-center py-12 gap-3 text-center">
          <FileTextOutlined className="text-4xl text-gray-200" />
          <p className="text-sm text-gray-400">Belum ada data orders</p>
        </div>
      ),
    },
    {
      key: "invoices",
      label: (
        <span className="flex items-center gap-1.5">
          <FileTextOutlined />Invoices
        </span>
      ),
      children: (
        <div className="p-5">
          {/* Outstandings */}
          <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
            Outstandings
          </p>
          <Table<InvoiceItem>
            rowKey="id"
            size="small"
            dataSource={outstandings}
            loading={outstandingsLoading}
            pagination={false}
            className={TABLE_CLS}
            locale={{ emptyText: "Tidak ada outstanding invoice" }}
            columns={[
              {
                title: "#",
                key: "no",
                width: 48,
                align: "center" as const,
                render: (_: unknown, __: InvoiceItem, idx: number) => idx + 1,
              },
              {
                title: "Invoice Number",
                dataIndex: "invoice_number",
                key: "invoice_number",
                render: (v: string) => <span className="font-medium text-gray-800">{v}</span>,
              },
              { title: "Project Name", dataIndex: "project_name", key: "project_name",
                render: (v?: string) => v || <span className="text-gray-400">-</span> },
              {
                title: "Amount",
                dataIndex: "amount",
                key: "amount",
                render: (v: number, r: InvoiceItem) => (
                  <span className="font-medium">{v.toLocaleString("id-ID")}</span>
                ),
              },
              {
                title: "Status",
                dataIndex: "status",
                key: "status",
                width: 100,
                render: (v: string) => {
                  const cfg: Record<string, { color: string; label: string }> = {
                    unpaid:  { color: "orange", label: "Unpaid"  },
                    partial: { color: "blue",   label: "Partial" },
                    overdue: { color: "red",    label: "Overdue" },
                    paid:    { color: "green",  label: "Paid"    },
                  };
                  const c = cfg[v] ?? { color: "default", label: v };
                  return <Tag color={c.color}>{c.label}</Tag>;
                },
              },
              {
                title: "Due Date",
                dataIndex: "due_date",
                key: "due_date",
                render: (v: string) =>
                  new Date(v).toLocaleDateString("id-ID", {
                    day: "2-digit", month: "short", year: "numeric",
                  }),
              },
              {
                title: "Aging (days)",
                dataIndex: "aging_days",
                key: "aging_days",
                width: 110,
                align: "center" as const,
                render: (v: number) => (
                  <span className={v > 30 ? "text-red-600 font-semibold" : "text-gray-700"}>
                    {v}
                  </span>
                ),
              },
            ]}
          />
        </div>
      ),
    },
    {
      key: "projects",
      label: (
        <span className="flex items-center gap-1.5">
          <ProjectOutlined />Projects
        </span>
      ),
      children: (
        <div className="p-5 flex flex-col items-center justify-center py-12 gap-3 text-center">
          <ProjectOutlined className="text-4xl text-gray-200" />
          <p className="text-sm text-gray-400">Belum ada data projects</p>
        </div>
      ),
    },
    {
      key: "services",
      label: (
        <span className="flex items-center gap-1.5">
          <ToolOutlined />Services
        </span>
      ),
      children: (
        <div className="p-5 flex flex-col items-center justify-center py-12 gap-3 text-center">
          <ToolOutlined className="text-4xl text-gray-200" />
          <p className="text-sm text-gray-400">Belum ada data services</p>
        </div>
      ),
    },
    {
      key: "history",
      label: (
        <span className="flex items-center gap-1.5">
          <HistoryOutlined />History
        </span>
      ),
      children: (
        <div className="p-5">
          <Table
            rowKey="id"
            size="small"
            dataSource={supplier.stage_histories}
            pagination={false}
            locale={{ emptyText: "Belum ada histori stage" }}
            className={TABLE_CLS}
            columns={[
              {
                title: "Dari",
                dataIndex: "from_stage",
                key: "from_stage",
                render: (v: SupplierStage) =>
                  v ? <Tag>{STAGE_LABEL[v]}</Tag> : <span className="text-gray-400">-</span>,
              },
              {
                title: "Ke",
                dataIndex: "to_stage",
                key: "to_stage",
                render: (v: SupplierStage) => <Tag color="blue">{STAGE_LABEL[v]}</Tag>,
              },
              { title: "Catatan",    dataIndex: "notes",      key: "notes"      },
              { title: "Diubah oleh", dataIndex: "changed_by", key: "changed_by" },
              {
                title: "Elapsed",
                dataIndex: "elapsed_ms",
                key: "elapsed_ms",
                render: (v: number) => (v ? `${Math.round(v / 1000)}s` : "-"),
              },
              {
                title: "Tanggal",
                dataIndex: "created_at",
                key: "created_at",
                render: (v: string) => new Date(v).toLocaleString("id-ID"),
              },
            ]}
          />
        </div>
      ),
    },
  ];

  return (
    <div>
      {/* Breadcrumb — di atas, terpisah dari card stack */}
      <Breadcrumb
        className="mb-10"
        items={[
          {
            title: (
              <button
                className="text-blue-500 hover:underline bg-transparent border-none cursor-pointer p-0 text-sm"
                onClick={() => router.push("/suppliers")}
              >
                Supplier List
              </button>
            ),
          },
          { title: <span className="text-gray-600 text-sm">{supplier.name}</span> },
        ]}
      />

      <div className="space-y-4">

      {/* ── Header card ─────────────────────────────────────── */}
      <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-5">
        <div className="flex items-center justify-between gap-4">
          {/* Left: back + logo + info */}
          <div className="flex items-center gap-4 min-w-0">
            <button
              onClick={() => router.back()}
              className="w-9 h-9 rounded-xl border border-gray-200 flex items-center justify-center text-gray-500 hover:border-blue-500 hover:text-blue-500 transition-colors flex-shrink-0"
            >
              <ArrowLeftOutlined className="text-sm" />
            </button>

            <SupplierLogo status={supplier.status} />

            <div className="min-w-0">
              <div className="flex items-center gap-2.5 flex-wrap">
                <h2 className="text-lg font-bold text-gray-800 m-0 leading-tight truncate">
                  {supplier.name}
                </h2>
                <StatusPill status={supplier.status} />
                {supplier.is_blocked && (
                  <span className="inline-flex items-center gap-1 text-xs text-red-600 bg-red-50 px-2 py-0.5 rounded-full">
                    <span className="w-1.5 h-1.5 rounded-full bg-red-500" />
                    Diblokir
                  </span>
                )}
              </div>

              {/* Sub info line */}
              {location && (
                <div className="flex items-center gap-1 mt-1">
                  <span className="text-xs text-gray-400 flex items-center gap-1">
                    <EnvironmentOutlined className="text-[11px]" />
                    {location}
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Right: actions */}
          <div className="flex items-center gap-2 flex-shrink-0">
            {can("supplier:update") && (
              <button className="flex items-center gap-1.5 px-3.5 h-9 rounded-xl border border-gray-200 text-sm font-medium text-gray-600 hover:border-blue-500 hover:text-blue-500 transition-colors">
                <EditOutlined className="text-xs" />
                Edit
              </button>
            )}
            {can("supplier:block") && (
              <button
                onClick={handleBlock}
                disabled={blockMut.isPending}
                className={`flex items-center gap-1.5 px-3.5 h-9 rounded-xl text-sm font-semibold transition-colors ${
                  supplier.is_blocked
                    ? "bg-green-600 hover:bg-green-700 text-white"
                    : "bg-red-600 hover:bg-red-700 text-white"
                } disabled:opacity-60`}
              >
                {supplier.is_blocked
                  ? <><CheckCircleOutlined className="text-xs" /> Buka Blokir</>
                  : <><StopOutlined className="text-xs" /> Block / Unblock</>
                }
              </button>
            )}
          </div>
        </div>
      </div>

      {/* ── Body: two-column ────────────────────────────────── */}
      <div className="grid grid-cols-3 gap-4 items-start">

        {/* ── Left 2/3 — tabs ─────────────────────────────── */}
        <div className="col-span-2 bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
          <Tabs
            defaultActiveKey="overview"
            items={mainTabs}
            className={[
              "[&_.ant-tabs-nav]:!px-5",
              "[&_.ant-tabs-nav]:!mb-0",
              "[&_.ant-tabs-nav]:border-b",
              "[&_.ant-tabs-nav]:border-gray-100",
              "[&_.ant-tabs-tab]:!text-sm",
              "[&_.ant-tabs-tab-active_.ant-tabs-tab-btn]:!font-semibold",
              "[&_.ant-tabs-content-holder]:min-h-[300px]",
            ].join(" ")}
          />
        </div>

        {/* ── Right 1/3 — sidebar ─────────────────────────── */}
        <div className="space-y-4">

          {/* Stage card */}
          <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-5">
            <div className="flex items-start justify-between mb-4">
              <div>
                <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide leading-none mb-0.5">
                  Stage
                </p>
                <p className="text-sm font-bold text-gray-800 m-0">
                  {STAGE_LABEL[supplier.stage]}
                </p>
              </div>
              <span className="text-[11px] font-medium text-gray-400 bg-gray-50 border border-gray-200 rounded-lg px-2 py-1 whitespace-nowrap">
                SLA: {supplier.sla_hours} hour(s)
              </span>
            </div>

            {/* Stepper */}
            <StageStepper
              steps={STAGE_STEPS.map((s) => ({ key: s, label: STAGE_LABEL[s] }))}
              current={stageIdx}
            />

            {/* Elapsed */}
            {supplier.stage_histories.length > 0 && (() => {
              const last = supplier.stage_histories[supplier.stage_histories.length - 1];
              if (!last.elapsed_ms) return null;
              const totalSec = Math.round(last.elapsed_ms / 1000);
              const h = Math.floor(totalSec / 3600);
              const m = Math.floor((totalSec % 3600) / 60);
              const s = totalSec % 60;
              return (
                <p className="text-[11px] text-gray-400 mt-3 mb-0">
                  Elapsed{" "}
                  <span className="font-medium text-gray-600 font-mono">
                    {String(h).padStart(2, "0")}:{String(m).padStart(2, "0")}:{String(s).padStart(2, "0")}
                  </span>
                </p>
              );
            })()}

            {/* Notes + Next Stage */}
            {canAdvance && (
              <div className="mt-4 space-y-2">
                <TextArea
                  rows={2}
                  placeholder="Notes (opsional)…"
                  value={stageNotes}
                  onChange={(e) => setStageNotes(e.target.value)}
                  className="!rounded-xl !text-sm !resize-none"
                />
                <button
                  onClick={handleNextStage}
                  disabled={stageMut.isPending}
                  className="w-full flex items-center justify-center gap-1.5 h-9 rounded-xl bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold transition-colors disabled:opacity-60"
                >
                  Next Stage
                  <ArrowRightOutlined className="text-xs" />
                </button>
              </div>
            )}
          </div>

          {/* Performance Ratings card */}
          <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-5">
            <div className="flex items-center justify-between mb-3">
              <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide m-0 flex items-center gap-1.5">
                <StarOutlined />
                Performance Ratings
              </p>
              <span className="text-[11px] text-gray-400">{ratings.length} review{ratings.length !== 1 ? "s" : ""}</span>
            </div>

            {ratings.length === 0 ? (
              <div className="flex flex-col items-center gap-2 py-6">
                <StarOutlined className="text-3xl text-gray-200" />
                <p className="text-xs text-gray-400">Belum ada penilaian</p>
              </div>
            ) : (
              <div className="space-y-2">
                {ratings.map((r) => (
                  <RatingCard key={r.id} rating={r} />
                ))}
              </div>
            )}
          </div>

        </div>
        {/* end sidebar */}
      </div>
    </div>
    </div>
  );
}
