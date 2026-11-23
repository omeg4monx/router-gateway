# OMEG4MON ROUTER — Design Direction

Direction ini berasal dari pemilik produk (joshz), ditulis apa adanya. Ini **bukan** panduan gaya buatan agen: warna, logo, dan password default adalah keputusan pemilik.

## Identity

- **Produk:** OMEG4MON ROUTER, gateway LLM (OpenAI + Anthropic compatible) yang di-host sendiri.
- **Karakter:** tool infrastruktur yang hangat, bukan dashboard korporat dingin. Tegas, presisi, tapi tidak steril.
- **Motif identitas:** garis aksen tipis (`from-transparent via-[var(--accent-primary)] to-transparent`) yang mengikuti aksen **Ruby Rose** dari custom palette `main` (`#b5173a` di Light, `#eb6583` di Dark). Satu gesture, dipakai konsisten.
- **Logo:** caduceus (tongkat Hermes). Keputusan pemilik, jangan diganti.

## Palette (Parity dengan `main` + Custom Accent `rose` / Ruby Rose)

Palet warna mengikuti sistem tema **Omegamon** dari branch `main` dengan aksen **Ruby Rose** (`preset("rose", "Ruby Rose", ["#b5173a", "#951330"], ["#eb6583", "#ef8099"])`):

### Light — Omegamon X-Antibody + Ruby Rose
| Peran | Nilai | Keterangan |
|---|---|---|
| Primary (Brand / Ruby Rose) | `#b5173a` (hover `#951330`) | Custom palette `rose` light; `on-primary` `#ffffff` (kontras >= 6:1) |
| Secondary | `#E8A317` (bg `#fdf4e0`, fg `#8a5200`) | Grey Sword gold |
| Accent | `#C8102E` (bg `#fce8eb`, fg `#a50d26`) | Red armor trim |
| Highlight | `#00B8E6` | All Delete cyan (decorative / glow) |
| Surfaces | `#F5F6F8` canvas, `#ffffff` surface, `#eef1f5` raised | Pearl-white armor |
| Borders | `#e8ebf0` subtle, `#d9dee6` strong | Cool pearl seams |
| Text | `#0B0D12` primary, `#4a5260` secondary, `#5f6775` muted | Black armor plates |

### Dark — Omegamon Zwart (Zwart X) + Ruby Rose
| Peran | Nilai | Keterangan |
|---|---|---|
| Primary (Brand / Ruby Rose) | `#eb6583` (hover `#ef8099`) | Custom palette `rose` dark; `on-primary` `#050608` (kontras >= 6:1) |
| Secondary | `#C0C6CF` (fg `#d5dae1`) | Silver trim |
| Accent | `#FF2D3A` (fg `#ff7a82`) | Burning red eyes / energy core |
| Highlight | `#22D3EE` | Electric cyan glow |
| Surfaces | `#050608` canvas, `#101217` surface, `#171a20` raised, `#0a0c10` input | Obsidian armor & gunmetal plates |
| Borders | `#171a20` subtle, `#262a33` strong | Gunmetal seams |
| Text | `#E6E9EE` primary, `#a3aab6` secondary, `#7e8592` muted | Silver |

## Typography

- **Sans:** Geist Sans. Alasan: produk developer, dan Geist punya bentuk numerik dan tabular yang rapi untuk tabel usage dan key. Dipilih sadar, bukan karena default.
- **Mono:** Geist Mono, khusus untuk nilai teknis: key, hash, endpoint, token, ID. Mono dipakai sebagai penanda data, bukan sebagai kostum "terminal".
- **Skala:** heading besar tegas untuk fokus, body tenang, label kecil untuk metadata.

## Surfaces

- **Solid, tanpa glassmorphism.** Blur/translucency sudah dihapus dari seluruh UI (R-10). Panel, sidebar, header, dan modal semuanya opak.
- **Elevation minimal.** Shadow hanya untuk modal dan dropdown yang benar-benar mengangkat diri dari halaman (R-12).
- **Glow dose cap:** maksimal 1 elemen sebagai focus accent (R-13).

## Dials

Dial: **ENERGY 2 / RHYTHM 2 / MOTION 1**

- ENERGY 2: hadir dan hangat, tapi tidak berteriak. Aksen oranye muncul di momen kunci.
- RHYTHM 2: konsisten dengan beberapa jeda. Login punya komposisi split yang berbeda dari dashboard shell.
- MOTION 1: hover dan transisi state saja. Tidak ada loop tanpa henti, tidak ada parallax.

## Surface archetypes (per layar)

- **Login:** Configure. User sedang menyiapkan akses. Dekorasi rendah, state validasi jelas, satu fokus.
- **Dashboard:** Monitor. User memantau state. Kepadatan dan hierarki sekali-pandang menang atas dekorasi.

## Non-negotiables

- Tanpa emoji di teks UI. Ikon SVG bila perlu, dengan alasan (R-04).
- Semua fitur tetap ada. Perubahan dibatasi ke lapisan tampilan (UI, GUI, UX).
- Kontras teks memenuhi WCAG AA (R-25).
- Setiap kontrol punya perilaku nyata (R-26).
