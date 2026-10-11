import {
  Stethoscope,
  Database,
  History,
  ShieldCheck,
  BookOpen,
  Receipt,
  Building2,
} from "@lucide/svelte";

export interface SubMenuItem {
  id: string;
  label: string;
  badge?: string | null;
  permission?: string;
}

export interface MenuItem {
  id: string;
  label: string;
  shortLabel: string;
  icon: any;
  children: SubMenuItem[];
}

export const navigationStructure: MenuItem[] = [
  {
    id: "clinical",
    label: "Pelayanan Medis",
    shortLabel: "Klinis",
    icon: Stethoscope,
    children: [
      {
        id: "physical",
        label: "Status Lokalis & Anatomi",
        badge: "Body",
        permission: "patient:read",
      },
      {
        id: "anamnesis",
        label: "Anamnesis & SOAP",
        badge: null,
        permission: "patient:read",
      },
    ],
  },
  {
    id: "master",
    label: "Master Data RS",
    shortLabel: "Master",
    icon: Database,
    children: [
      {
        id: "master-patient",
        label: "Pasien (Patient)",
        badge: null,
        permission: "patient:read",
      },
      {
        id: "master-practitioner",
        label: "Tenaga Medis (Practitioner)",
        badge: null,
        permission: "practitioner:read",
      },
      {
        id: "master-departement",
        label: "Departemen & Instalasi",
        badge: null,
        permission: "department:read",
      },
      {
        id: "master-service-unit",
        label: "Unit Layanan & Poli",
        badge: null,
        permission: "service_unit:read",
      },
      {
        id: "master-room",
        label: "Ruangan & Bed (Room)",
        badge: null,
        permission: "room:read",
      },
      {
        id: "master-customer",
        label: "Debitur & Penjamin (Customer)",
        badge: null,
        permission: "customer:read",
      },
      {
        id: "master-referal",
        label: "Faskes Rujukan (Referal)",
        badge: null,
        permission: "referal:read",
      },
      {
        id: "master-tariff-class",
        label: "Kelas Tarif (Tariff Class)",
        badge: null,
        permission: "tariff_class:read",
      },
      {
        id: "master-item",
        label: "Katalog Item & Tarif",
        badge: "Universal",
        permission: "item:read",
      },
    ],
  },
  {
    id: "accounting",
    label: "Chart of Account",
    shortLabel: "Coa",
    icon: BookOpen,
    children: [
      {
        id: "accounting-coa",
        label: "Bagan Akun (COA)",
        badge: "Aktif",
        permission: "accounting:read",
      },
      {
        id: "accounting-journals",
        label: "Jurnal Umum Transaksi",
        badge: "Tahap 2",
        permission: "accounting:read",
      },
      {
        id: "accounting-ledger",
        label: "Buku Besar & Kas/Bank",
        badge: "Tahap 3",
        permission: "accounting:read",
      },
      {
        id: "accounting-reports",
        label: "Laporan Keuangan & Neraca",
        badge: "Tahap 4",
        permission: "accounting:read",
      },
    ],
  },
  {
    id: "tariff",
    label: "Buku & Manajemen Tarif",
    shortLabel: "Tarif",
    icon: Receipt,
    children: [
      {
        id: "finance-tariff",
        label: "Tarif & Price Plan (Tariff)",
        badge: "Tariff",
        permission: "tariff_class:read",
      },
      {
        id: "finance-price-plan",
        label: "Buku Tarif (Price Plan)",
        badge: "Core",
        permission: "tariff_price_plan:read",
      },
      {
        id: "finance-tariff-matrix",
        label: "Editor Matriks Tarif",
        badge: "Batch",
        permission: "tariff_price_plan:update",
      },
      {
        id: "finance-tariff-lookup",
        label: "Kalkulator Lookup Tarif",
        badge: "Engine",
        permission: "tariff_price_plan:read",
      },
      {
        id: "finance-tariff-class",
        label: "Kelas Tarif (Class)",
        badge: null,
        permission: "tariff_class:read",
      },
      {
        id: "finance-tariff-component",
        label: "Komponen Tarif (Component)",
        badge: null,
        permission: "tariff_component:read",
      },
    ],
  },
  {
    id: "customer",
    label: "Debitur & Penjamin (Customer)",
    shortLabel: "Penjamin",
    icon: Building2,
    children: [
      {
        id: "finance-customer",
        label: "Master Debitur & Penjamin",
        badge: "Customer",
        permission: "customer:read",
      },
    ],
  },
  {
    id: "auth",
    label: "Pengguna & Akses",
    shortLabel: "Akses",
    icon: ShieldCheck,
    children: [
      {
        id: "auth-users",
        label: "Kelola Pengguna",
        badge: null,
        permission: "user:read",
      },
      {
        id: "auth-roles",
        label: "Peran & Hak Akses",
        badge: null,
        permission: "role:read",
      },
      {
        id: "auth-permissions",
        label: "Katalog Izin",
        badge: null,
        permission: "role:read",
      },
    ],
  },
  {
    id: "activity",
    label: "Aktivitas Pasien",
    shortLabel: "Aktivitas",
    icon: History,
    children: [
      {
        id: "history",
        label: "Riwayat Kunjungan RME",
        badge: null,
        permission: "patient:read",
      },
      {
        id: "schedule",
        label: "Jadwal Kontrol Poliklinik",
        badge: null,
        permission: "appointment:view",
      },
    ],
  },
];
