"use client";

import { useRouter } from "next/navigation";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button, Form, Input, App } from "antd";
import { LockOutlined, MailOutlined } from "@ant-design/icons";
import { useAuth } from "@/hooks/useAuth";

const loginSchema = z.object({
  email: z.string().email("Format email tidak valid"),
  password: z.string().min(1, "Password wajib diisi"),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const { message } = App.useApp();

  const {
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = async (values: LoginFormValues) => {
    try {
      await login(values);
      router.replace("/suppliers");
    } catch {
      message.error("Email atau password salah. Silakan coba lagi.");
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* ── Left branding panel ───────────────────────────── */}
      <div className="flex-1 flex flex-col items-center justify-center px-12 bg-gradient-to-br from-[#1e2640] to-[#2d3a5c]">
        {/* Logo */}
        <div className="w-16 h-16 rounded-full bg-blue-600 flex items-center justify-center text-white text-2xl font-black mb-4">
          A
        </div>
        <h1 className="text-4xl font-bold text-white tracking-widest mb-2">ALISA</h1>
        <p className="text-white/60 text-base mb-8">ERP Supplier Management System</p>

        <div className="w-64 h-px bg-white/15 my-6" />

        <p className="text-white/40 text-sm text-center max-w-xs leading-relaxed">
          Kelola seluruh ekosistem supplier Anda dalam satu platform terintegrasi
        </p>

        {/* Feature bullets */}
        <div className="mt-10 space-y-3 w-full max-w-xs">
          {[
            "Manajemen supplier end-to-end",
            "Workflow approval otomatis",
            "Dashboard & analytics real-time",
          ].map((f) => (
            <div key={f} className="flex items-center gap-3">
              <div className="w-1.5 h-1.5 rounded-full bg-blue-400" />
              <span className="text-white/55 text-sm">{f}</span>
            </div>
          ))}
        </div>
      </div>

      {/* ── Right login form ──────────────────────────────── */}
      <div className="w-[420px] bg-white flex flex-col justify-center px-10 py-12 shadow-2xl">
        <h2 className="text-2xl font-bold text-gray-800 mb-1">Masuk ke Akun</h2>
        <p className="text-gray-400 text-sm mb-8">
          Gunakan akun ERP Anda untuk melanjutkan
        </p>

        <Form layout="vertical" onFinish={handleSubmit(onSubmit)} requiredMark={false}>
          <Form.Item
            label={<span className="font-semibold text-gray-700 text-sm">Email</span>}
            validateStatus={errors.email ? "error" : ""}
            help={errors.email?.message}
          >
            <Controller
              name="email"
              control={control}
              render={({ field }) => (
                <Input
                  {...field}
                  prefix={<MailOutlined className="text-gray-300" />}
                  placeholder="admin@erp.local"
                  size="large"
                  autoComplete="email"
                  className="rounded-lg"
                />
              )}
            />
          </Form.Item>

          <Form.Item
            label={<span className="font-semibold text-gray-700 text-sm">Password</span>}
            validateStatus={errors.password ? "error" : ""}
            help={errors.password?.message}
            className="mb-7"
          >
            <Controller
              name="password"
              control={control}
              render={({ field }) => (
                <Input.Password
                  {...field}
                  prefix={<LockOutlined className="text-gray-300" />}
                  placeholder="••••••••"
                  size="large"
                  autoComplete="current-password"
                  className="rounded-lg"
                />
              )}
            />
          </Form.Item>

          <Form.Item className="mb-0">
            <Button
              type="primary"
              htmlType="submit"
              size="large"
              block
              loading={isSubmitting}
              className="rounded-lg h-11 font-semibold text-base"
            >
              Masuk
            </Button>
          </Form.Item>
        </Form>

        <p className="text-center text-xs text-gray-400 mt-5">
          Default: <span className="font-mono">admin@erp.local</span> /{" "}
          <span className="font-mono">Admin@123</span>
        </p>
      </div>
    </div>
  );
}
