"use client";

import { useQuery } from "@tanstack/react-query";
import { Table, Tag, Card, Button } from "antd";
import { useRouter } from "next/navigation";
import { listSuppliers } from "@/services/supplier.service";
import type { SupplierListItem } from "@/types/api";

export default function ReviewsPage() {
  const router = useRouter();

  const { data, isLoading } = useQuery({
    queryKey: ["suppliers-review"],
    queryFn: () => listSuppliers({ status: "in_progress", limit: 50 }),
  });

  return (
    <div>
      <div className="flex items-center justify-between mb-5">
        <h1 className="text-xl font-bold text-gray-800 m-0">Review &amp; Approvals</h1>
      </div>

      <Card className="rounded-lg !p-0" styles={{ body: { padding: 0 } }}>
        <Table<SupplierListItem>
          rowKey="id"
          loading={isLoading}
          dataSource={data?.data ?? []}
          size="middle"
          className="[&_.ant-table-thead_th]:bg-gray-50 [&_.ant-table-thead_th]:font-semibold [&_.ant-table-row:hover_td]:bg-blue-50 [&_.ant-table-row]:cursor-pointer"
          columns={[
            {
              title: "Supplier",
              key: "name",
              render: (_: unknown, r: SupplierListItem) => (
                <div>
                  <p className="font-semibold text-gray-800 text-[13px] m-0">{r.name}</p>
                  <p className="text-xs text-gray-400 m-0">{r.code} · {r.supplier_no}</p>
                </div>
              ),
            },
            {
              title: "Alamat",
              dataIndex: "address",
              key: "address",
              render: (v: string) => v || "-",
            },
            {
              title: "Kontak",
              dataIndex: "contact",
              key: "contact",
              render: (v: string) => v || "-",
            },
            {
              title: "Status",
              dataIndex: "status",
              key: "status",
              render: () => <Tag color="processing">In Progress</Tag>,
            },
            {
              title: "Aksi",
              key: "action",
              render: (_: unknown, r: SupplierListItem) => (
                <Button
                  size="small"
                  type="primary"
                  className="rounded"
                  onClick={() => router.push(`/suppliers/${r.id}`)}
                >
                  Review
                </Button>
              ),
            },
          ]}
          pagination={{
            pageSize: 20,
            showTotal: (t) => `Total ${t} supplier`,
          }}
          onRow={(r) => ({ onClick: () => router.push(`/suppliers/${r.id}`) })}
          locale={{ emptyText: "Tidak ada supplier yang perlu direview" }}
        />
      </Card>
    </div>
  );
}
