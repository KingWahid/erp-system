"use client";

import { useRouter } from "next/navigation";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button, Form, Input, App, Card } from "antd";
import { ArrowLeftOutlined, SaveOutlined } from "@ant-design/icons";
import { createSupplier } from "@/services/supplier.service";

const schema = z.object({
  name: z.string().min(2, "Nama minimal 2 karakter"),
  code: z.string().min(2, "Kode minimal 2 karakter").max(10, "Kode maksimal 10 karakter"),
  alias: z.string().optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  country: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().email("Format email tidak valid").optional().or(z.literal("")),
  website: z.string().optional(),
  notes: z.string().optional(),
});

type FormValues = z.infer<typeof schema>;

export default function NewSupplierPage() {
  const router = useRouter();
  const qc = useQueryClient();
  const { message } = App.useApp();

  const {
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      country: "Indonesia",
    },
  });

  const mutation = useMutation({
    mutationFn: createSupplier,
    onSuccess: (res) => {
      qc.invalidateQueries({ queryKey: ["suppliers"] });
      qc.invalidateQueries({ queryKey: ["supplier-stats"] });
      message.success("Supplier berhasil dibuat");
      router.push(`/suppliers/${res.data.id}`);
    },
    onError: () => message.error("Gagal membuat supplier. Cek kode apakah sudah digunakan."),
  });

  const onSubmit = (values: FormValues) => {
    mutation.mutate({
      name: values.name,
      code: values.code,
      alias: values.alias || undefined,
      address: values.address || undefined,
      city: values.city || undefined,
      country: values.country || undefined,
      phone: values.phone || undefined,
      email: values.email || undefined,
      website: values.website || undefined,
      notes: values.notes || undefined,
    });
  };

  return (
    <div>
      {/* Header */}
      <div className="flex items-center gap-3 mb-5">
        <Button
          icon={<ArrowLeftOutlined />}
          onClick={() => router.back()}
          className="rounded-md"
        />
        <h1 className="text-xl font-bold text-gray-800 m-0">New Supplier</h1>
      </div>

      <Form layout="vertical" onFinish={handleSubmit(onSubmit)} requiredMark="optional">
        <div className="grid grid-cols-3 gap-4">
          {/* Main info */}
          <div className="col-span-2 space-y-4">
            <Card
              title={<span className="text-sm font-semibold">Informasi Dasar</span>}
              className="rounded-lg"
              size="small"
            >
              <div className="grid grid-cols-2 gap-4">
                {/* Name */}
                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Nama Perusahaan *</span>}
                  validateStatus={errors.name ? "error" : ""}
                  help={errors.name?.message}
                  className="col-span-2"
                >
                  <Controller
                    name="name"
                    control={control}
                    render={({ field }) => (
                      <Input {...field} placeholder="Contoh: PT Setroom Indonesia" className="rounded-md" />
                    )}
                  />
                </Form.Item>

                {/* Code */}
                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Kode *</span>}
                  validateStatus={errors.code ? "error" : ""}
                  help={errors.code?.message}
                >
                  <Controller
                    name="code"
                    control={control}
                    render={({ field }) => (
                      <Input
                        {...field}
                        placeholder="Contoh: STRM"
                        maxLength={10}
                        className="rounded-md uppercase"
                        onChange={(e) => field.onChange(e.target.value.toUpperCase())}
                      />
                    )}
                  />
                </Form.Item>

                {/* Alias */}
                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Alias</span>}
                  validateStatus={errors.alias ? "error" : ""}
                  help={errors.alias?.message}
                >
                  <Controller
                    name="alias"
                    control={control}
                    render={({ field }) => (
                      <Input {...field} placeholder="Contoh: Setroom" className="rounded-md" />
                    )}
                  />
                </Form.Item>

                {/* Phone */}
                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Telepon</span>}
                >
                  <Controller
                    name="phone"
                    control={control}
                    render={({ field }) => (
                      <Input {...field} placeholder="021-123456" className="rounded-md" />
                    )}
                  />
                </Form.Item>

                {/* Email */}
                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Email</span>}
                  validateStatus={errors.email ? "error" : ""}
                  help={errors.email?.message}
                >
                  <Controller
                    name="email"
                    control={control}
                    render={({ field }) => (
                      <Input {...field} type="email" placeholder="info@perusahaan.com" className="rounded-md" />
                    )}
                  />
                </Form.Item>

                {/* Website */}
                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Website</span>}
                  className="col-span-2"
                >
                  <Controller
                    name="website"
                    control={control}
                    render={({ field }) => (
                      <Input {...field} placeholder="https://perusahaan.com" className="rounded-md" />
                    )}
                  />
                </Form.Item>
              </div>
            </Card>

            {/* Address */}
            <Card
              title={<span className="text-sm font-semibold">Alamat</span>}
              className="rounded-lg"
              size="small"
            >
              <div className="grid grid-cols-2 gap-4">
                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Alamat</span>}
                  className="col-span-2"
                >
                  <Controller
                    name="address"
                    control={control}
                    render={({ field }) => (
                      <Input.TextArea
                        {...field}
                        rows={2}
                        placeholder="Jl. Sudirman No. 1"
                        className="rounded-md"
                      />
                    )}
                  />
                </Form.Item>

                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Kota</span>}
                >
                  <Controller
                    name="city"
                    control={control}
                    render={({ field }) => (
                      <Input {...field} placeholder="Jakarta" className="rounded-md" />
                    )}
                  />
                </Form.Item>

                <Form.Item
                  label={<span className="text-sm font-medium text-gray-700">Negara</span>}
                >
                  <Controller
                    name="country"
                    control={control}
                    render={({ field }) => (
                      <Input {...field} placeholder="Indonesia" className="rounded-md" />
                    )}
                  />
                </Form.Item>
              </div>
            </Card>

            {/* Notes */}
            <Card
              title={<span className="text-sm font-semibold">Catatan</span>}
              className="rounded-lg"
              size="small"
            >
              <Form.Item className="mb-0">
                <Controller
                  name="notes"
                  control={control}
                  render={({ field }) => (
                    <Input.TextArea
                      {...field}
                      rows={3}
                      placeholder="Catatan tambahan tentang supplier ini..."
                      className="rounded-md"
                    />
                  )}
                />
              </Form.Item>
            </Card>
          </div>

          {/* Right sidebar — actions */}
          <div>
            <Card
              title={<span className="text-sm font-semibold">Aksi</span>}
              className="rounded-lg sticky top-4"
              size="small"
            >
              <div className="space-y-2">
                <Button
                  type="primary"
                  htmlType="submit"
                  icon={<SaveOutlined />}
                  block
                  loading={isSubmitting || mutation.isPending}
                  className="rounded-md"
                >
                  Simpan Supplier
                </Button>
                <Button
                  block
                  onClick={() => router.back()}
                  className="rounded-md"
                >
                  Batal
                </Button>
              </div>

              <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                <p className="text-xs text-blue-600 font-medium mb-1">Info</p>
                <p className="text-xs text-blue-500">
                  Supplier baru akan dibuat dengan status <strong>Draft</strong>.
                  Anda dapat mengubah stage setelah supplier dibuat.
                </p>
              </div>
            </Card>
          </div>
        </div>
      </Form>
    </div>
  );
}
