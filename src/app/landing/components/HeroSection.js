"use client";

export default function HeroSection() {
  return (
    <section className="relative pt-32 pb-20 px-6 min-h-[90vh] flex flex-col items-center justify-center overflow-hidden">
      {/* Glow effect */}
      <div className="absolute inset-x-0 top-0 h-[420px] pointer-events-none" style={{ background: "radial-gradient(ellipse 65% 100% at 50% 0%, color-mix(in srgb, var(--color-primary) 14%, transparent), transparent 70%)" }}></div>
      
      <div className="relative z-10 max-w-4xl w-full text-center flex flex-col items-center gap-8">
        {/* Version badge */}
        <div className="inline-flex items-center gap-2 rounded-full border border-border bg-surface-2/50 px-3 py-1 text-xs font-medium text-primary">
          <span className="flex h-2 w-2 rounded-full bg-[var(--color-brand-500)] animate-pulse"></span>
          v1.0 is now live
        </div>

        {/* Main heading */}
        <h1 className="text-5xl md:text-7xl font-black leading-[1.1] tracking-tight">
          One Endpoint for <br/>
          <span className="text-[var(--color-brand-500)]">All AI Providers</span>
        </h1>

        {/* Description */}
        <p className="text-lg md:text-xl text-muted max-w-2xl mx-auto font-light">
          AI endpoint proxy with web dashboard - A JavaScript port of CLIProxyAPI. Works seamlessly with Claude Code, OpenAI Codex, Cline, RooCode, and other CLI tools.
        </p>

        {/* CTA Buttons */}
        <div className="flex flex-wrap items-center justify-center gap-4 w-full">
          <button className="h-12 px-8 rounded-lg bg-primary hover:bg-primary-hover text-on-primary text-base font-bold transition-all shadow-[0_0_15px_color-mix(in_srgb,var(--color-primary)_40%,transparent)] flex items-center gap-2">
            <span className="material-symbols-outlined">rocket_launch</span>
            Get Started
          </button>
          <a 
            href="https://github.com/omeg4monx/router-gateway"
            target="_blank" 
            rel="noopener noreferrer"
            className="h-12 px-8 rounded-lg border border-border bg-surface-2 hover:bg-surface-3 text-text-main text-base font-bold transition-all flex items-center gap-2"
          >
            <span className="material-symbols-outlined">code</span>
            View on GitHub
          </a>
        </div>
      </div>
    </section>
  );
}

