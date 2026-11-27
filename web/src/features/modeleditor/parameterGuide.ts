import type { UILocale } from "@/shared/store/uiStore";

export type ParameterGuideId =
  | "model_pattern"
  | "target_model"
  | "temperature"
  | "max_tokens"
  | "strip_params"
  | "extra_params"
  | "dry_run";

export interface ParameterGuideEntry {
  id: ParameterGuideId;
  shortLabel: Record<UILocale, string>;
  title: Record<UILocale, string>;
  purpose: Record<UILocale, string>;
  defaultValue: Record<UILocale, string>;
  whenIncreased: {
    label: Record<UILocale, string>;
    effect: Record<UILocale, string>;
  };
  whenDecreased: {
    label: Record<UILocale, string>;
    effect: Record<UILocale, string>;
  };
  example: string;
}

export const PARAMETER_GUIDE_ENTRIES: readonly [ParameterGuideEntry, ...ParameterGuideEntry[]] = [
  {
    id: "temperature",
    shortLabel: { en: "Temperature", id: "Temperature" },
    title: { en: "Temperature (0.0 – 2.0)", id: "Temperature (0.0 – 2.0)" },
    purpose: {
      en: "Controls randomness and creativity when the model chooses the next token.",
      id: "Mengatur tingkat keacakan dan kreativitas saat model memilih token berikutnya.",
    },
    defaultValue: {
      en: "Empty = keep client/provider default (typically 0.7 – 1.0).",
      id: "Kosong = ikut default klien/provider (umumnya 0.7 – 1.0).",
    },
    whenIncreased: {
      label: { en: "Higher (0.8 – 1.5+)", id: "Dinaikkan (0.8 – 1.5+)" },
      effect: {
        en: "More creative and varied wording, but higher risk of hallucinations, drift, or broken JSON.",
        id: "Jawaban lebih kreatif & variatif, tetapi risiko halusinasi, melenceng, atau JSON rusak lebih tinggi.",
      },
    },
    whenDecreased: {
      label: { en: "Lower (0.0 – 0.3)", id: "Diturunkan (0.0 – 0.3)" },
      effect: {
        en: "Deterministic, focused, and repeatable output. Best for coding, SQL, JSON, and strict extraction.",
        id: "Output lebih deterministik, fokus, dan konsisten. Ideal untuk coding, SQL, JSON, dan ekstraksi data.",
      },
    },
    example: "0.2",
  },
  {
    id: "extra_params",
    shortLabel: { en: "Extra JSON", id: "Extra JSON" },
    title: { en: "Extra Params JSON", id: "Extra Params JSON" },
    purpose: {
      en: "Injects or overwrites arbitrary root-level fields in the outgoing JSON request body.",
      id: "Menyuntikkan atau menimpa field JSON bebas di level root body request sebelum dikirim ke provider.",
    },
    defaultValue: {
      en: "{} (empty object — no extra fields injected).",
      id: "{} (objek kosong — tidak ada field tambahan yang disuntikkan).",
    },
    whenIncreased: {
      label: { en: "Added keys", id: "Ditambahkan key" },
      effect: {
        en: "Enables provider-specific controls (e.g. seed, reasoning_effort). Unknown keys may trigger HTTP 400 on strict providers.",
        id: "Mengaktifkan fitur khusus provider (mis. seed, reasoning_effort). Key yang tidak dikenali provider bisa memicu error HTTP 400.",
      },
    },
    whenDecreased: {
      label: { en: "Empty ({})", id: "Dikosongkan ({})" },
      effect: {
        en: "Keeps the client request payload clean with zero custom field injection.",
        id: "Payload request dari klien tetap bersih tanpa injeksi field tambahan.",
      },
    },
    example: '{"seed": 42, "reasoning_effort": "medium"}',
  },
  {
    id: "max_tokens",
    shortLabel: { en: "Max Tokens", id: "Max Tokens" },
    title: { en: "Max Tokens (1 – 2,000,000)", id: "Max Tokens (1 – 2.000.000)" },
    purpose: {
      en: "Caps the maximum number of output tokens the model is allowed to generate in one response.",
      id: "Membatasi jumlah maksimum token output yang boleh dihasilkan model dalam satu respons.",
    },
    defaultValue: {
      en: "Empty = use client/provider limit (commonly 4096 – 8192).",
      id: "Kosong = ikut batas klien/provider (umumnya 4096 – 8192).",
    },
    whenIncreased: {
      label: { en: "Higher", id: "Dinaikkan" },
      effect: {
        en: "Allows longer, complete responses and code outputs, at the cost of higher token spend and latency.",
        id: "Memungkinkan jawaban atau kode yang lebih panjang & tuntas, namun biaya token dan latensi ikut naik.",
      },
    },
    whenDecreased: {
      label: { en: "Lower", id: "Diturunkan" },
      effect: {
        en: "Faster replies and lower cost ceiling, but long answers may be cut off mid-sentence.",
        id: "Respons lebih cepat dan batas biaya lebih aman, tetapi jawaban panjang berisiko terpotong di tengah.",
      },
    },
    example: "4096",
  },
  {
    id: "strip_params",
    shortLabel: { en: "Strip Params", id: "Strip Params" },
    title: { en: "Strip Params (CSV)", id: "Strip Params (CSV)" },
    purpose: {
      en: "Removes comma-separated fields from the incoming request before forwarding upstream.",
      id: "Menghapus daftar field (dipisah koma) dari request masuk sebelum diteruskan ke provider upstream.",
    },
    defaultValue: {
      en: "Empty = forward all client parameters as-is.",
      id: "Kosong = seluruh parameter dari klien diteruskan apa adanya.",
    },
    whenIncreased: {
      label: { en: "Added fields", id: "Ditambahkan field" },
      effect: {
        en: "Prevents HTTP 400 errors on strict/reasoning models (e.g. o3) that reject temperature or penalty fields.",
        id: "Mencegah error HTTP 400 pada model ketat/reasoning (mis. o3) yang menolak temperature atau penalty.",
      },
    },
    whenDecreased: {
      label: { en: "Removed / Empty", id: "Dikurangkan / Kosong" },
      effect: {
        en: "Preserves every client parameter; unsupported fields will reach the upstream provider.",
        id: "Semua parameter klien tetap dikirim; field yang tidak didukung bisa ditolak provider.",
      },
    },
    example: "temperature, presence_penalty",
  },
  {
    id: "model_pattern",
    shortLabel: { en: "Model Pattern", id: "Pola Model" },
    title: { en: "Model Pattern (Wildcard)", id: "Pola Model (Wildcard)" },
    purpose: {
      en: "Matches incoming request model names to decide whether this override rule should run.",
      id: "Mencocokkan nama model pada request masuk untuk menentukan apakah aturan override ini aktif.",
    },
    defaultValue: {
      en: "Required — exact model name or glob wildcard (*).",
      id: "Wajib diisi — nama model persis atau wildcard (*).",
    },
    whenIncreased: {
      label: { en: "Broader (*)", id: "Lebih luas (*)" },
      effect: {
        en: "Applies to many or all models at once. Fast for global defaults, but can affect unintended models.",
        id: "Berlaku untuk banyak atau semua model sekaligus. Cepat untuk aturan global, tapi bisa mengenai model lain.",
      },
    },
    whenDecreased: {
      label: { en: "More specific", id: "Lebih spesifik" },
      effect: {
        en: "Targets only an exact model or prefix (e.g. coding-*), keeping other models untouched.",
        id: "Hanya menyentuh model atau prefix tertentu (mis. coding-*), sehingga model lain tetap aman.",
      },
    },
    example: "coding-*",
  },
  {
    id: "target_model",
    shortLabel: { en: "Target Model", id: "Model Target" },
    title: { en: "Target Model (Rewrite)", id: "Model Target (Rewrite)" },
    purpose: {
      en: "Rewrites the outgoing model name before the gateway routes the request to a provider.",
      id: "Mengganti nama model keluar sebelum gateway meneruskan request ke provider.",
    },
    defaultValue: {
      en: "Empty = keep the original requested model name.",
      id: "Kosong = tetap memakai nama model asli dari request.",
    },
    whenIncreased: {
      label: { en: "Filled in", id: "Diisi" },
      effect: {
        en: "Transparently redirects alias/legacy names (e.g. coding-fast) to a real upstream model.",
        id: "Mengalihkan alias/nama lama (mis. coding-fast) secara transparan ke model upstream yang sebenarnya.",
      },
    },
    whenDecreased: {
      label: { en: "Left empty", id: "Dikosongkan" },
      effect: {
        en: "Only parameter overrides apply; the model name itself stays unchanged.",
        id: "Hanya parameter yang diubah; nama model tidak diganti.",
      },
    },
    example: "claude-sonnet-4",
  },
  {
    id: "dry_run",
    shortLabel: { en: "Dry-Run", id: "Dry-Run" },
    title: { en: "Dry-Run Preview", id: "Dry-Run Preview" },
    purpose: {
      en: "Simulates how all enabled overrides transform a sample JSON payload without calling any AI provider.",
      id: "Mensimulasikan bagaimana semua override aktif mengubah payload JSON contoh tanpa memanggil provider AI.",
    },
    defaultValue: {
      en: "Uses openai-chat protocol with a sample request body (zero token cost).",
      id: "Memakai protokol openai-chat dengan payload contoh (tanpa biaya token).",
    },
    whenIncreased: {
      label: { en: "More active overrides", id: "Banyak override aktif" },
      effect: {
        en: "Shows the exact final merged payload and rewritten model so you can spot conflicts early.",
        id: "Menampilkan hasil gabungan payload akhir dan model hasil rewrite untuk mendeteksi konflik sejak awal.",
      },
    },
    whenDecreased: {
      label: { en: "No matching override", id: "Tanpa override cocok" },
      effect: {
        en: "Returns the original request payload unchanged.",
        id: "Mengembalikan payload request asli tanpa perubahan.",
      },
    },
    example: '{"model":"coding-fast","temperature":0.9}',
  },
] as const;

export function getParameterGuideEntry(id: ParameterGuideId): ParameterGuideEntry {
  return PARAMETER_GUIDE_ENTRIES.find((entry) => entry.id === id) ?? PARAMETER_GUIDE_ENTRIES[0];
}

export const MODEL_EDITOR_COPY = {
  en: {
    subtitle: "Override model parameters and test mutations without spending tokens.",
    guideButton: "Guide",
    refreshButton: "Refresh",
    tabOverride: "New Override",
    tabRules: "Active",
    tabDryRun: "Dry-Run",
    newOverrideTitle: "New Model Override",
    activeOverridesTitle: "Active Overrides",
    dryRunTitle: "Dry-Run Preview",
    nameLabel: "Name",
    namePlaceholder: "Reasoning Safe",
    patternLabel: "Model Pattern",
    patternPlaceholder: "coding-* or o3-*",
    targetLabel: "Target Model",
    targetPlaceholder: "optional rewrite",
    tempLabel: "Temperature",
    tempPlaceholder: "0.2",
    maxTokensLabel: "Max Tokens",
    maxTokensPlaceholder: "4096",
    stripLabel: "Strip Params (CSV)",
    stripPlaceholder: "temperature, presence_penalty",
    extraLabel: "Extra Params JSON",
    extraPlaceholder: '{"seed": 42}',
    createButton: "Create Override",
    creatingButton: "Creating...",
    loadingOverrides: "Loading overrides...",
    emptyOverrides: "No overrides yet. Create one above to shape outgoing model requests.",
    enableLabel: "Enable",
    disableLabel: "Disable",
    deleteAriaLabel: "Delete override",
    previewModelLabel: "Model",
    previewBodyLabel: "Request JSON",
    runPreviewButton: "Run Preview",
    runningPreviewButton: "Running...",
    mutatedOutputLabel: "Mutated Output",
    previewPlaceholder: "Run preview to inspect transformed payload and matched model.",
    errorTitle: "Model Editor Error",
    dismissLabel: "Dismiss",
    guideTitle: "Parameter Guide",
    guideSubtitle: "Purpose, default behavior, and impact when increased or reduced.",
    guidePurpose: "What it does",
    guideDefault: "Standard / Default",
    guideExample: "Example",
    closeLabel: "Close",
    helpAriaPrefix: "Explain",
  },
  id: {
    subtitle: "Ubah parameter model dan uji hasil mutasi tanpa memakan token.",
    guideButton: "Panduan",
    refreshButton: "Muat Ulang",
    tabOverride: "Override Baru",
    tabRules: "Aktif",
    tabDryRun: "Dry-Run",
    newOverrideTitle: "New Model Override",
    activeOverridesTitle: "Override Aktif",
    dryRunTitle: "Dry-Run Preview",
    nameLabel: "Nama",
    namePlaceholder: "Reasoning Safe",
    patternLabel: "Model Pattern",
    patternPlaceholder: "coding-* atau o3-*",
    targetLabel: "Target Model",
    targetPlaceholder: "opsional rewrite",
    tempLabel: "Temperature",
    tempPlaceholder: "0.2",
    maxTokensLabel: "Max Tokens",
    maxTokensPlaceholder: "4096",
    stripLabel: "Strip Params (CSV)",
    stripPlaceholder: "temperature, presence_penalty",
    extraLabel: "Extra Params JSON",
    extraPlaceholder: '{"seed": 42}',
    createButton: "Simpan Override",
    creatingButton: "Menyimpan...",
    loadingOverrides: "Memuat daftar override...",
    emptyOverrides:
      "Belum ada override. Buat aturan di atas untuk mengatur parameter request model.",
    enableLabel: "Aktifkan",
    disableLabel: "Nonaktifkan",
    deleteAriaLabel: "Hapus override",
    previewModelLabel: "Model",
    previewBodyLabel: "Request JSON",
    runPreviewButton: "Jalankan Preview",
    runningPreviewButton: "Memproses...",
    mutatedOutputLabel: "Hasil Mutasi",
    previewPlaceholder: "Jalankan preview untuk melihat hasil payload dan model yang cocok.",
    errorTitle: "Error Model Editor",
    dismissLabel: "Tutup",
    guideTitle: "Panduan Parameter",
    guideSubtitle: "Kegunaan, nilai standar/default, serta dampak saat dinaikkan atau diturunkan.",
    guidePurpose: "Kegunaan",
    guideDefault: "Standar / Default",
    guideExample: "Contoh",
    closeLabel: "Tutup",
    helpAriaPrefix: "Jelaskan",
  },
} as const;
