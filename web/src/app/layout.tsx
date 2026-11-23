import type { Metadata } from "next";
import { JetBrains_Mono, Plus_Jakarta_Sans } from "next/font/google";
import type { ReactNode } from "react";
import { QueryProvider } from "@/shared/api/QueryProvider";
import "./globals.css";

const plusJakartaSans = Plus_Jakarta_Sans({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-ui-sans",
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  variable: "--font-ui-mono",
  display: "swap",
});

export const metadata: Metadata = {
  title: "OMEG4MON ROUTER",
  description: "Self-hosted OpenAI and Anthropic compatible LLM routing gateway.",
};

const UI_BOOTSTRAP_SCRIPT = `(function(){try{var t=localStorage.getItem("omeg4mon-theme");if(t==="light"||t==="dark"){document.documentElement.setAttribute("data-theme",t);}var l=localStorage.getItem("omeg4mon-locale");if(l==="en"||l==="id"){document.documentElement.setAttribute("lang",l);document.documentElement.setAttribute("data-locale",l);}}catch(e){}})();`;

interface RootLayoutProps {
  children: ReactNode;
}

export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <html
      lang="en"
      data-theme="light"
      data-locale="en"
      className={`${plusJakartaSans.variable} ${jetbrainsMono.variable}`}
      suppressHydrationWarning
    >
      <head>
        {/* biome-ignore lint/security/noDangerouslySetInnerHtml: Static theme and locale bootstrap script with zero user input to prevent FOUC */}
        <script dangerouslySetInnerHTML={{ __html: UI_BOOTSTRAP_SCRIPT }} />
      </head>
      <body className="min-h-screen bg-[var(--bg-canvas)] font-sans text-[var(--text-primary)]">
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
