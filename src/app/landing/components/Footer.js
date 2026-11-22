"use client";

export default function Footer() {
  return (
    <footer className="border-t border-border bg-bg-alt pt-16 pb-8 px-6">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-8 mb-16">
          {/* Brand */}
          <div className="col-span-2 lg:col-span-2">
            <div className="flex items-center gap-3 mb-6">
              <div className="size-6 rounded bg-primary flex items-center justify-center text-on-primary">
                <span className="material-symbols-outlined text-[16px]">hub</span>
              </div>
              <h3 className="text-text-main text-lg font-bold">OMEG4MON ROUTER</h3>
            </div>
            <p className="text-muted text-sm max-w-xs mb-6">
              Connect, route, and manage your AI providers from one endpoint.
            </p>
            <div className="flex gap-4">
              <a className="text-muted hover:text-text-main transition-colors" href="https://github.com/omeg4monx/router-gateway" target="_blank" rel="noopener noreferrer">
                <span className="material-symbols-outlined">code</span>
              </a>
            </div>
          </div>
          
          {/* Product */}
          <div className="flex flex-col gap-4">
            <h4 className="font-bold text-text-main">Product</h4>
            <a className="text-muted hover:text-[var(--color-brand-500)] text-sm transition-colors" href="#features">Features</a>
            <a className="text-muted hover:text-[var(--color-brand-500)] text-sm transition-colors" href="/dashboard">Dashboard</a>
            <a className="text-muted hover:text-[var(--color-brand-500)] text-sm transition-colors" href="https://github.com/omeg4monx/router-gateway" target="_blank" rel="noopener noreferrer">Changelog</a>
          </div>
          
          {/* Resources */}
          <div className="flex flex-col gap-4">
            <h4 className="font-bold text-text-main">Resources</h4>
            <a className="text-muted hover:text-[var(--color-brand-500)] text-sm transition-colors" href="https://github.com/omeg4monx/router-gateway/tree/main/docs" target="_blank" rel="noopener noreferrer">Documentation</a>
            <a className="text-muted hover:text-[var(--color-brand-500)] text-sm transition-colors" href="https://github.com/omeg4monx/router-gateway" target="_blank" rel="noopener noreferrer">GitHub</a>
          </div>
          
          {/* Legal */}
          <div className="flex flex-col gap-4">
            <h4 className="font-bold text-text-main">Legal</h4>
            <a className="text-muted hover:text-[var(--color-brand-500)] text-sm transition-colors" href="https://github.com/omeg4monx/router-gateway/blob/main/LICENSE" target="_blank" rel="noopener noreferrer">MIT License</a>
          </div>
        </div>
        
        {/* Bottom */}
        <div className="border-t border-border pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-muted text-sm">© 2025 OMEG4MON ROUTER. All rights reserved.</p>
          <div className="flex gap-6">
            <a className="text-muted hover:text-text-main text-sm transition-colors" href="https://github.com/omeg4monx/router-gateway" target="_blank" rel="noopener noreferrer">GitHub</a>
          </div>
        </div>
      </div>
    </footer>
  );
}

