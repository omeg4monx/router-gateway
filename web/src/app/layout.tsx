import type { Metadata } from "next";
import { GeistMono } from "geist/font/mono";
import { GeistSans } from "geist/font/sans";
import type { ReactNode } from "react";
import { QueryProvider } from "@/shared/api/QueryProvider";
import "./globals.css";

export const metadata: Metadata = {
  title: "OMEG4MON ROUTER",
  description: "Self-hosted OpenAI and Anthropic compatible LLM routing gateway.",
};

const THEME_INIT_SCRIPT = `(function(){try{var t=localStorage.getItem("omeg4mon-theme");if(t==="light"||t==="dark"){document.documentElement.setAttribute("data-theme",t);}}catch(e){}})();`;

interface RootLayoutProps {
  children: ReactNode;
}

export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <html
      lang="en"
      data-theme="dark"
      className={`${GeistSans.variable} ${GeistMono.variable}`}
      suppressHydrationWarning
    >
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME_INIT_SCRIPT }} />
      </head>
      <body className="min-h-screen bg-[var(--bg-canvas)] font-sans text-[var(--text-primary)]">
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
