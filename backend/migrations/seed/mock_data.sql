-- =============================================================================
-- MOCK DATA SEED
-- Run: make seed-mock
-- =============================================================================

BEGIN;

-- =============================================================================
-- USERS (password for all: Password@123)
-- bcrypt hash of "Password@123"
-- =============================================================================
INSERT INTO users (id, name, email, password, is_active) VALUES
    ('a1b2c3d4-0001-0001-0001-000000000001', 'Admin System',    'admin@erp.local',    '$2a$10$6InfEA2HKH9PgEIepo3qfuJi0mf7T2.PNP9l/0BwzS/UAoVAIi44K', TRUE),
    ('a1b2c3d4-0002-0002-0002-000000000002', 'Budi Manager',    'manager@erp.local',  '$2a$10$5p0RHHumuDhMArNKZSSBiupnEtcvc/AzhY0y2vCZmzpG6ESHSOJyO', TRUE),
    ('a1b2c3d4-0003-0003-0003-000000000003', 'Citra Viewer',    'viewer@erp.local',   '$2a$10$5p0RHHumuDhMArNKZSSBiupnEtcvc/AzhY0y2vCZmzpG6ESHSOJyO', TRUE),
    ('a1b2c3d4-0004-0004-0004-000000000004', 'Dewi Reviewer',   'reviewer@erp.local', '$2a$10$5p0RHHumuDhMArNKZSSBiupnEtcvc/AzhY0y2vCZmzpG6ESHSOJyO', TRUE)
ON CONFLICT (email) DO NOTHING;

-- =============================================================================
-- ASSIGN ROLES TO USERS
-- =============================================================================
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u CROSS JOIN roles r
WHERE u.email = 'admin@erp.local'   AND r.name = 'admin'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u CROSS JOIN roles r
WHERE u.email = 'manager@erp.local' AND r.name = 'manager'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u CROSS JOIN roles r
WHERE u.email = 'viewer@erp.local'  AND r.name = 'viewer'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u CROSS JOIN roles r
WHERE u.email = 'reviewer@erp.local' AND r.name = 'viewer'
ON CONFLICT DO NOTHING;

-- =============================================================================
-- SUPPLIERS
-- =============================================================================
INSERT INTO suppliers (
    id, code, supplier_no, name, alias,
    address, city, province, country,
    phone, email, website,
    status, stage, sla_hours,
    is_blocked, notes,
    created_at, updated_at
) VALUES
    -- Active suppliers
    ('b1000001-0000-0000-0000-000000000001', 'STRM', '61000012',
     'PT Setroom Indonesia', 'Setroom',
     'Fatmawati Raya St, 33', 'Jakarta Selatan', 'DKI Jakarta', 'Indonesia',
     '+62211234567', 'info@setroom.co.id', 'https://setroom.co.id',
     'active', 'active', 72, FALSE, 'Supplier aktif untuk perangkat IT',
     NOW() - INTERVAL '120 days', NOW()),

    ('b1000002-0000-0000-0000-000000000002', 'SKSK', '41000013',
     'PT Suko Suko', 'Sukasuko',
     'Jl. Cihampelas No. 22', 'Bandung', 'Jawa Barat', 'Indonesia',
     '+62222345678', 'contact@sukosuko.com', 'https://sukosuko.com',
     'in_progress', 'in_review', 72, FALSE, 'Dalam proses review',
     NOW() - INTERVAL '30 days', NOW()),

    ('b1000003-0000-0000-0000-000000000003', 'ARIB', '41000014',
     'PT Angin Ribut', 'Angin',
     'Jl. Gatot Subroto No. 88', 'Denpasar', 'Bali', 'Indonesia',
     '+62361234567', 'ops@anginribut.id', NULL,
     'blocked', 'active', 72, TRUE, 'Diblokir karena keterlambatan pembayaran',
     NOW() - INTERVAL '90 days', NOW()),

    ('b1000004-0000-0000-0000-000000000004', 'MNDR', '71000015',
     'CV Mandiri Jaya', 'Mandiri',
     'Jl. Pemuda No. 10', 'Surabaya', 'Jawa Timur', 'Indonesia',
     '+62315678901', 'info@mandiri-jaya.com', NULL,
     'active', 'active', 48, FALSE, 'Supplier alat tulis kantor',
     NOW() - INTERVAL '200 days', NOW()),

    ('b1000005-0000-0000-0000-000000000005', 'TKNO', '81000016',
     'PT Teknologi Nusantara', 'Teknusa',
     'Jl. Sudirman Kav. 45', 'Jakarta Pusat', 'DKI Jakarta', 'Indonesia',
     '+62212223344', 'hello@teknusa.id', 'https://teknusa.id',
     'in_progress', 'in_assessment', 72, FALSE, 'Dalam tahap assessment',
     NOW() - INTERVAL '15 days', NOW()),

    ('b1000006-0000-0000-0000-000000000006', 'GRND', '91000017',
     'UD Garuda Mas', 'Garuda',
     'Jl. Malioboro No. 5', 'Yogyakarta', 'DI Yogyakarta', 'Indonesia',
     '+62274567890', 'garuda@mas.co.id', NULL,
     'draft', 'draft', 72, FALSE, 'Supplier baru, belum diverifikasi',
     NOW() - INTERVAL '3 days', NOW()),

    ('b1000007-0000-0000-0000-000000000007', 'BRKH', '51000018',
     'PT Berkah Abadi', 'Berkah',
     'Jl. Ahmad Yani No. 100', 'Medan', 'Sumatera Utara', 'Indonesia',
     '+62618901234', 'info@berkah-abadi.com', 'https://berkah-abadi.com',
     'active', 'active', 24, FALSE, 'Supplier furnitur kantor premium',
     NOW() - INTERVAL '365 days', NOW()),

    ('b1000008-0000-0000-0000-000000000008', 'SMRT', '31000019',
     'PT Smart Solution', 'Smart',
     'Jl. Diponegoro No. 77', 'Semarang', 'Jawa Tengah', 'Indonesia',
     '+62247654321', 'cs@smart-solution.co.id', 'https://smart-solution.co.id',
     'active', 'active', 72, FALSE, 'Supplier software dan lisensi',
     NOW() - INTERVAL '180 days', NOW()),

    ('b1000009-0000-0000-0000-000000000009', 'PNTI', '21000020',
     'CV Panti Karya', 'Panti',
     'Jl. Veteran No. 33', 'Makassar', 'Sulawesi Selatan', 'Indonesia',
     '+62411234567', NULL, NULL,
     'inactive', 'active', 72, FALSE, 'Tidak aktif - kontrak habis',
     NOW() - INTERVAL '400 days', NOW()),

    ('b1000010-0000-0000-0000-000000000010', 'DLTA', '11000021',
     'PT Delta Prima', 'Delta',
     'Jl. Raya Darmo No. 15', 'Surabaya', 'Jawa Timur', 'Indonesia',
     '+62318765432', 'delta@prima.co.id', 'https://deltaprima.id',
     'active', 'active', 48, FALSE, 'Supplier peralatan cleaning service',
     NOW() - INTERVAL '250 days', NOW())
ON CONFLICT (code) DO NOTHING;

-- =============================================================================
-- SUPPLIER CONTACTS
-- =============================================================================
INSERT INTO supplier_contacts (supplier_id, name, position, phone, mobile, email, is_primary) VALUES
    ('b1000001-0000-0000-0000-000000000001', 'Albert Einstein',   'CEO',                 '021.123456', '0811234567', 'albert@setroom.co.id',      TRUE),
    ('b1000001-0000-0000-0000-000000000001', 'Isaac Newton',      'Mgr Proc',            '021.654321', '0811765432', 'isaac@setroom.co.id',        FALSE),
    ('b1000001-0000-0000-0000-000000000001', 'Marie Curie',       'Finance Director',    '021.111222', '0812345678', 'marie@setroom.co.id',        FALSE),
    ('b1000002-0000-0000-0000-000000000002', 'James Lee',         'Account Manager',     '022.234567', '0823456789', 'james@sukosuko.com',         TRUE),
    ('b1000003-0000-0000-0000-000000000003', 'Mario Chen',        'Operations Manager',  '0361.234567','0834567890', 'mario@anginribut.id',        TRUE),
    ('b1000004-0000-0000-0000-000000000004', 'Siti Rahayu',       'Director',            '031.456789', '0845678901', 'siti@mandiri-jaya.com',      TRUE),
    ('b1000005-0000-0000-0000-000000000005', 'Reza Pratama',      'Sales Manager',       '021.567890', '0856789012', 'reza@teknusa.id',            TRUE),
    ('b1000006-0000-0000-0000-000000000006', 'Wahyu Santoso',     'Owner',               '0274.567890','0867890123', 'wahyu@garudamas.co.id',      TRUE),
    ('b1000007-0000-0000-0000-000000000007', 'Putri Handayani',   'VP Operations',       '061.678901', '0878901234', 'putri@berkah-abadi.com',     TRUE),
    ('b1000007-0000-0000-0000-000000000007', 'Doni Kusuma',       'Technical Lead',      '061.678902', '0878901235', 'doni@berkah-abadi.com',      FALSE),
    ('b1000008-0000-0000-0000-000000000008', 'Kevin Wijaya',      'Sales Director',      '024.789012', '0889012345', 'kevin@smart-solution.co.id', TRUE),
    ('b1000010-0000-0000-0000-000000000010', 'Lina Marlina',      'Procurement Lead',    '031.890123', '0890123456', 'lina@deltaprima.id',         TRUE);

-- =============================================================================
-- SUPPLIER ADDRESSES
-- =============================================================================
INSERT INTO supplier_addresses (supplier_id, name, address, city, province, country, postal_code, is_main) VALUES
    ('b1000001-0000-0000-0000-000000000001', 'Head Office',   'Fatmawati Raya St, 123', 'Jakarta Selatan', 'DKI Jakarta',    'Indonesia', '12420', TRUE),
    ('b1000001-0000-0000-0000-000000000001', 'Branch Office', 'Ciawi No. 99',           'Bogor',           'Jawa Barat',     'Indonesia', '16730', FALSE),
    ('b1000002-0000-0000-0000-000000000002', 'Head Office',   'Jl. Cihampelas No. 22',  'Bandung',         'Jawa Barat',     'Indonesia', '40131', TRUE),
    ('b1000004-0000-0000-0000-000000000004', 'Head Office',   'Jl. Pemuda No. 10',      'Surabaya',        'Jawa Timur',     'Indonesia', '60271', TRUE),
    ('b1000004-0000-0000-0000-000000000004', 'Warehouse',     'Jl. Rungkut No. 5',      'Surabaya',        'Jawa Timur',     'Indonesia', '60293', FALSE),
    ('b1000007-0000-0000-0000-000000000007', 'Head Office',   'Jl. Ahmad Yani No. 100', 'Medan',           'Sumatera Utara', 'Indonesia', '20111', TRUE),
    ('b1000008-0000-0000-0000-000000000008', 'Head Office',   'Jl. Diponegoro No. 77',  'Semarang',        'Jawa Tengah',    'Indonesia', '50131', TRUE),
    ('b1000010-0000-0000-0000-000000000010', 'Head Office',   'Jl. Raya Darmo No. 15',  'Surabaya',        'Jawa Timur',     'Indonesia', '60264', TRUE);

-- =============================================================================
-- SUPPLIER GROUPS
-- =============================================================================
INSERT INTO supplier_groups (supplier_id, group_name, value, is_active) VALUES
    ('b1000001-0000-0000-0000-000000000001', 'Industry',      'Manufacture',       TRUE),
    ('b1000001-0000-0000-0000-000000000001', 'Telkom Group',  'Non Telkom Group',  TRUE),
    ('b1000002-0000-0000-0000-000000000002', 'Industry',      'Technology',        TRUE),
    ('b1000002-0000-0000-0000-000000000002', 'Telkom Group',  'Non Telkom Group',  TRUE),
    ('b1000004-0000-0000-0000-000000000004', 'Industry',      'Trading',           TRUE),
    ('b1000004-0000-0000-0000-000000000004', 'Telkom Group',  'Telkom Group',      TRUE),
    ('b1000005-0000-0000-0000-000000000005', 'Industry',      'Technology',        TRUE),
    ('b1000005-0000-0000-0000-000000000005', 'Telkom Group',  'Non Telkom Group',  FALSE),
    ('b1000007-0000-0000-0000-000000000007', 'Industry',      'Manufacture',       TRUE),
    ('b1000007-0000-0000-0000-000000000007', 'Telkom Group',  'Telkom Group',      TRUE),
    ('b1000008-0000-0000-0000-000000000008', 'Industry',      'Technology',        TRUE),
    ('b1000010-0000-0000-0000-000000000010', 'Industry',      'Service',           TRUE);

-- =============================================================================
-- SUPPLIER MATERIALS
-- =============================================================================
INSERT INTO supplier_materials (supplier_id, material_group, material_id, is_active) VALUES
    -- PT Setroom Indonesia - IT Devices
    ('b1000001-0000-0000-0000-000000000001', 'IT - Device',    'Computer / Notebook',    TRUE),
    ('b1000001-0000-0000-0000-000000000001', 'IT - Device',    'Computer / PC',          TRUE),
    ('b1000001-0000-0000-0000-000000000001', 'IT - Device',    'Monitor',                TRUE),
    ('b1000001-0000-0000-0000-000000000001', 'IT - Peripheral','Keyboard & Mouse',       TRUE),
    ('b1000001-0000-0000-0000-000000000001', 'IT - Peripheral','External Storage',       FALSE),

    -- PT Suko Suko - Networking
    ('b1000002-0000-0000-0000-000000000002', 'IT - Network',   'Router / Switch',        TRUE),
    ('b1000002-0000-0000-0000-000000000002', 'IT - Network',   'Network Cable',          TRUE),
    ('b1000002-0000-0000-0000-000000000002', 'IT - Device',    'Server',                 TRUE),

    -- PT Angin Ribut - Maintenance
    ('b1000003-0000-0000-0000-000000000003', 'Maintenance',    'AC Service',             TRUE),
    ('b1000003-0000-0000-0000-000000000003', 'Maintenance',    'Electrical Service',     TRUE),

    -- CV Mandiri Jaya - Office Supplies
    ('b1000004-0000-0000-0000-000000000004', 'Office Supply',  'Paper / Stationery',     TRUE),
    ('b1000004-0000-0000-0000-000000000004', 'Office Supply',  'Printer Toner',          TRUE),
    ('b1000004-0000-0000-0000-000000000004', 'Office Supply',  'Office Equipment',       TRUE),

    -- PT Teknologi Nusantara - Software
    ('b1000005-0000-0000-0000-000000000005', 'Software',       'ERP License',            TRUE),
    ('b1000005-0000-0000-0000-000000000005', 'Software',       'Security Software',      TRUE),
    ('b1000005-0000-0000-0000-000000000005', 'IT - Service',   'IT Consulting',          TRUE),

    -- PT Berkah Abadi - Furniture
    ('b1000007-0000-0000-0000-000000000007', 'Furniture',      'Office Chair',           TRUE),
    ('b1000007-0000-0000-0000-000000000007', 'Furniture',      'Office Desk',            TRUE),
    ('b1000007-0000-0000-0000-000000000007', 'Furniture',      'Meeting Table',          TRUE),
    ('b1000007-0000-0000-0000-000000000007', 'Furniture',      'Storage Cabinet',        FALSE),

    -- PT Smart Solution - Software & Services
    ('b1000008-0000-0000-0000-000000000008', 'Software',       'Microsoft 365',          TRUE),
    ('b1000008-0000-0000-0000-000000000008', 'Software',       'Adobe Creative Suite',   TRUE),
    ('b1000008-0000-0000-0000-000000000008', 'IT - Service',   'Cloud Hosting',          TRUE),

    -- PT Delta Prima - Cleaning
    ('b1000010-0000-0000-0000-000000000010', 'Cleaning',       'Cleaning Equipment',     TRUE),
    ('b1000010-0000-0000-0000-000000000010', 'Cleaning',       'Cleaning Chemical',      TRUE),
    ('b1000010-0000-0000-0000-000000000010', 'Cleaning Service','Office Cleaning',       TRUE);

-- =============================================================================
-- PERFORMANCE RATINGS
-- =============================================================================
INSERT INTO supplier_performance_ratings (
    supplier_id, price_rating, delivery_rating, notes, reviewed_by, reviewed_at
) VALUES
    -- PT Setroom Indonesia
    ('b1000001-0000-0000-0000-000000000001', 4, 5, 'Pengiriman sangat cepat, harga kompetitif',    'User Legal',   NOW() - INTERVAL '30 days'),
    ('b1000001-0000-0000-0000-000000000001', 3, 4, 'Harga sedikit naik dibanding kontrak lama',    'User Legal',   NOW() - INTERVAL '90 days'),
    ('b1000001-0000-0000-0000-000000000001', 5, 5, 'Performa sangat baik di Q3',                   'Budi Manager', NOW() - INTERVAL '180 days'),

    -- CV Mandiri Jaya
    ('b1000004-0000-0000-0000-000000000004', 4, 4, 'Konsisten dan tepat waktu',                    'Budi Manager', NOW() - INTERVAL '20 days'),
    ('b1000004-0000-0000-0000-000000000004', 3, 3, 'Ada keterlambatan di bulan Agustus',            'User Legal',   NOW() - INTERVAL '60 days'),

    -- PT Berkah Abadi
    ('b1000007-0000-0000-0000-000000000007', 5, 4, 'Kualitas furnitur premium, worth the price',   'Admin System', NOW() - INTERVAL '10 days'),
    ('b1000007-0000-0000-0000-000000000007', 4, 5, 'Instalasi rapi dan profesional',               'Budi Manager', NOW() - INTERVAL '120 days'),

    -- PT Smart Solution
    ('b1000008-0000-0000-0000-000000000008', 3, 5, 'Support responsif, harga lisensi agak tinggi', 'Admin System', NOW() - INTERVAL '45 days'),

    -- PT Delta Prima
    ('b1000010-0000-0000-0000-000000000010', 4, 4, 'Hasil kebersihan memuaskan',                   'Citra Viewer', NOW() - INTERVAL '15 days');

-- =============================================================================
-- STAGE HISTORIES
-- =============================================================================
INSERT INTO supplier_stage_histories (supplier_id, from_stage, to_stage, notes, changed_by, elapsed_ms) VALUES
    -- PT Setroom Indonesia: full journey draft → active
    ('b1000001-0000-0000-0000-000000000001', 'draft',        'in_review',     'Dokumen lengkap',              'admin@erp.local',   86400000),
    ('b1000001-0000-0000-0000-000000000001', 'in_review',    'in_assessment', 'Review selesai',               'manager@erp.local', 172800000),
    ('b1000001-0000-0000-0000-000000000001', 'in_assessment','active',        'Assessment lulus, diaktifkan', 'admin@erp.local',   259200000),

    -- PT Suko Suko: draft → in_review (masih proses)
    ('b1000002-0000-0000-0000-000000000002', 'draft',        'in_review',     'Dokumen diterima',             'admin@erp.local',   43200000),

    -- PT Teknologi Nusantara: draft → in_review → in_assessment
    ('b1000005-0000-0000-0000-000000000005', 'draft',        'in_review',     'Pengajuan awal',               'admin@erp.local',   86400000),
    ('b1000005-0000-0000-0000-000000000005', 'in_review',    'in_assessment', 'Dokumen valid, lanjut assessment', 'manager@erp.local', 43200000),

    -- CV Mandiri Jaya: full journey
    ('b1000004-0000-0000-0000-000000000004', 'draft',        'in_review',     'Pengajuan diterima',           'admin@erp.local',   86400000),
    ('b1000004-0000-0000-0000-000000000004', 'in_review',    'in_assessment', 'Lulus review',                 'manager@erp.local', 86400000),
    ('b1000004-0000-0000-0000-000000000004', 'in_assessment','active',        'Semua kriteria terpenuhi',     'admin@erp.local',   86400000),

    -- PT Berkah Abadi: full journey (supplier lama)
    ('b1000007-0000-0000-0000-000000000007', 'draft',        'in_review',     'Supplier lama re-registrasi',  'admin@erp.local',   86400000),
    ('b1000007-0000-0000-0000-000000000007', 'in_review',    'in_assessment', 'Track record baik',            'manager@erp.local', 43200000),
    ('b1000007-0000-0000-0000-000000000007', 'in_assessment','active',        'Diaktifkan kembali',           'admin@erp.local',   86400000);

-- =============================================================================
-- SUPPLIER INVOICES (Outstandings)
-- Seeded for: PT Setroom Indonesia, CV Mandiri Jaya, PT Berkah Abadi,
--             PT Smart Solution, PT Delta Prima
-- aging_days is computed at runtime from NOW() - due_date
-- =============================================================================
INSERT INTO supplier_invoices (
    supplier_id, invoice_number, project_name,
    amount, currency,
    invoice_date, due_date, paid_date,
    status, paid_amount, notes
) VALUES
    -- PT Setroom Indonesia (b1000001) — 3 outstandings
    ('b1000001-0000-0000-0000-000000000001', 'INV-2026-001', 'Project ABC',
     123000000, 'IDR',
     NOW() - INTERVAL '64 days', NOW() - INTERVAL '34 days', NULL,
     'overdue', 0, 'Pengadaan laptop Q1'),

    ('b1000001-0000-0000-0000-000000000001', 'INV-2026-002', 'Project DEF',
     87500000, 'IDR',
     NOW() - INTERVAL '30 days', NOW() + INTERVAL '15 days', NULL,
     'unpaid', 0, 'Pengadaan monitor & aksesoris'),

    ('b1000001-0000-0000-0000-000000000001', 'INV-2026-003', 'Project GHI',
     215000000, 'IDR',
     NOW() - INTERVAL '90 days', NOW() - INTERVAL '60 days', NULL,
     'partial', 100000000, 'Server procurement — partial payment received'),

    -- CV Mandiri Jaya (b1000004) — 2 outstandings
    ('b1000004-0000-0000-0000-000000000004', 'INV-2026-010', 'Office Restocking Q2',
     32500000, 'IDR',
     NOW() - INTERVAL '45 days', NOW() - INTERVAL '15 days', NULL,
     'overdue', 0, 'ATK & toner printer'),

    ('b1000004-0000-0000-0000-000000000004', 'INV-2026-011', 'Furniture Upgrade',
     18750000, 'IDR',
     NOW() - INTERVAL '10 days', NOW() + INTERVAL '20 days', NULL,
     'unpaid', 0, 'Kursi dan meja kerja tambahan'),

    -- PT Berkah Abadi (b1000007) — 2 outstandings
    ('b1000007-0000-0000-0000-000000000007', 'INV-2026-020', 'Kantor Lantai 5',
     540000000, 'IDR',
     NOW() - INTERVAL '120 days', NOW() - INTERVAL '90 days', NULL,
     'overdue', 200000000, 'Furnitur meeting room + lounge'),

    ('b1000007-0000-0000-0000-000000000007', 'INV-2026-021', 'Ekspansi Gedung Baru',
     320000000, 'IDR',
     NOW() - INTERVAL '20 days', NOW() + INTERVAL '10 days', NULL,
     'unpaid', 0, 'Workstation area baru'),

    -- PT Smart Solution (b1000008) — 1 outstanding
    ('b1000008-0000-0000-0000-000000000008', 'INV-2026-030', 'Lisensi Microsoft 365',
     98000000, 'IDR',
     NOW() - INTERVAL '55 days', NOW() - INTERVAL '25 days', NULL,
     'overdue', 0, 'Renewal tahunan 50 seat'),

    -- PT Delta Prima (b1000010) — 2 outstandings
    ('b1000010-0000-0000-0000-000000000010', 'INV-2026-040', 'Cleaning Service Q2',
     24000000, 'IDR',
     NOW() - INTERVAL '35 days', NOW() - INTERVAL '5 days', NULL,
     'overdue', 0, 'Jasa kebersihan bulanan April-Juni'),

    ('b1000010-0000-0000-0000-000000000010', 'INV-2026-041', 'Cleaning Supplies',
     11500000, 'IDR',
     NOW() - INTERVAL '5 days', NOW() + INTERVAL '25 days', NULL,
     'unpaid', 0, 'Perlengkapan kebersihan semester II');

COMMIT;
