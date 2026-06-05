"use client";

import { Card } from "antd";
import { SettingOutlined } from "@ant-design/icons";

export default function ConfigurationsPage() {
  return (
    <div>
      <h1 className="text-xl font-bold text-gray-800 mb-6">Configurations</h1>

      <Card className="rounded-lg">
        <div className="flex flex-col items-center justify-center py-16 gap-4 text-center">
          <div className="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center">
            <SettingOutlined className="text-3xl text-gray-300" />
          </div>
          <div>
            <h3 className="text-base font-semibold text-gray-600">
              Konfigurasi Supplier
            </h3>
            <p className="text-sm text-gray-400 mt-1 max-w-xs">
              Halaman ini akan menampilkan pengaturan grup supplier, kategori material, dan konfigurasi workflow.
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
