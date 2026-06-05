"use client";

import { useRouter } from "next/navigation";
import { Button } from "antd";
import { ShopOutlined, ArrowRightOutlined } from "@ant-design/icons";

export default function DashboardPage() {
  const router = useRouter();

  return (
    <div>
      <h1 className="text-xl font-bold text-gray-800 mb-6">Dashboard</h1>

      <div className="grid grid-cols-3 gap-4">
        {/* Quick access card */}
        <div
          className="bg-white rounded-lg border border-gray-200 p-6 flex flex-col gap-4 cursor-pointer hover:border-blue-300 hover:shadow-sm transition-all"
          onClick={() => router.push("/suppliers")}
        >
          <div className="w-12 h-12 rounded-lg bg-blue-50 flex items-center justify-center">
            <ShopOutlined className="text-2xl text-blue-500" />
          </div>
          <div>
            <h3 className="font-semibold text-gray-800 text-base mb-1">
              Supplier Management
            </h3>
            <p className="text-sm text-gray-400">
              Kelola data supplier, workflow approval, dan performance rating
            </p>
          </div>
          <Button
            type="link"
            icon={<ArrowRightOutlined />}
            className="self-start p-0 font-medium"
          >
            Buka Supplier List
          </Button>
        </div>
      </div>
    </div>
  );
}
