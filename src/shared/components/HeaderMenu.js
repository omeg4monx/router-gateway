"use client";

import { useState, useRef, useEffect } from "react";
import Link from "next/link";
import PropTypes from "prop-types";
import ChangelogModal from "./ChangelogModal";
import { ConfirmModal } from "./Modal";
import ThemeModal from "./ThemeModal";

function MenuItem({ icon, label, onClick, trailing, danger }) {
  return (
    <button
      type="button"
      role="menuitem"
      onClick={onClick}
      className={`flex items-center gap-2.5 w-full rounded-[8px] px-2.5 py-2 text-[12.5px] font-medium transition-colors duration-150 cursor-pointer active:scale-[0.98] ${
        danger
          ? "text-danger hover:bg-danger-bg"
          : "text-text-main hover:bg-surface-2"
      }`}
    >
      <span className={`material-symbols-outlined text-[18px] leading-none shrink-0 ${danger ? "" : "text-text-muted"}`}>
        {icon}
      </span>
      <span className="flex-1 text-left truncate">{label}</span>
      {trailing && <span className="text-base">{trailing}</span>}
    </button>
  );
}

MenuItem.propTypes = {
  icon: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
  trailing: PropTypes.node,
  danger: PropTypes.bool,
};

export default function HeaderMenu({ onLogout }) {
  const [isOpen, setIsOpen] = useState(false);
  const [changelogOpen, setChangelogOpen] = useState(false);
  const [themeModalOpen, setThemeModalOpen] = useState(false);
  const [shutdownOpen, setShutdownOpen] = useState(false);
  const [isShuttingDown, setIsShuttingDown] = useState(false);
  const menuRef = useRef(null);
  const triggerRef = useRef(null);

  const handleShutdown = async () => {
    setIsShuttingDown(true);
    try {
      await fetch("/api/version/shutdown", { method: "POST" });
    } catch (e) {
      // Expected to fail as server shuts down; ignore error
    }
    setIsShuttingDown(false);
    setShutdownOpen(false);
  };

  useEffect(() => {
    if (!isOpen) return undefined;
    const handleClickOutside = (e) => {
      if (menuRef.current && !menuRef.current.contains(e.target)) {
        setIsOpen(false);
      }
    };
    const handleKeyDown = (e) => {
      if (e.key === "Escape") {
        setIsOpen(false);
        triggerRef.current?.focus();
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [isOpen]);

  const close = () => setIsOpen(false);

  return (
    <>
      <div className="relative" ref={menuRef}>
        <button
          ref={triggerRef}
          type="button"
          onClick={() => setIsOpen((v) => !v)}
          aria-haspopup="menu"
          aria-expanded={isOpen}
          className={`flex items-center justify-center size-8 rounded-[10px] transition-all duration-150 cursor-pointer active:scale-[0.97] ${
            isOpen
              ? "bg-primary/12 text-primary border border-primary/20"
              : "text-text-muted hover:text-text-main hover:bg-surface-2 border border-transparent"
          }`}
          title="Account menu"
          aria-label="Account menu"
        >
          <span className="material-symbols-outlined text-[20px] leading-none">account_circle</span>
        </button>

        {isOpen && (
          <div
            className="absolute left-0 bottom-full mb-2.5 w-52 rounded-xl shell-popover popover-in z-[70] p-1.5"
            role="menu"
          >
            <MenuItem
              icon="history"
              label="Change Log"
              onClick={() => { close(); setChangelogOpen(true); }}
            />
            <MenuItem
              icon="palette"
              label="Theme"
              onClick={() => { close(); setThemeModalOpen(true); }}
            />
            <Link
              href="/dashboard/profile"
              role="menuitem"
              onClick={close}
              className="flex items-center gap-2.5 w-full rounded-[8px] px-2.5 py-2 text-[12.5px] font-medium text-text-main hover:bg-surface-2 transition-colors duration-150 cursor-pointer active:scale-[0.98]"
            >
              <span className="material-symbols-outlined text-[18px] leading-none shrink-0 text-text-muted">settings</span>
              <span className="flex-1 text-left truncate">Settings</span>
            </Link>
            <div className="my-1 border-t border-border-subtle/80" role="separator" />
            <MenuItem
              icon="power_settings_new"
              label="Shutdown"
              danger
              onClick={() => { close(); setShutdownOpen(true); }}
            />
            <MenuItem
              icon="logout"
              label="Logout"
              danger
              onClick={() => { close(); onLogout(); }}
            />
          </div>
        )}
      </div>

      <ChangelogModal isOpen={changelogOpen} onClose={() => setChangelogOpen(false)} />
      <ThemeModal isOpen={themeModalOpen} onClose={() => setThemeModalOpen(false)} />
      <ConfirmModal
        isOpen={shutdownOpen}
        onClose={() => setShutdownOpen(false)}
        onConfirm={handleShutdown}
        title="Close Proxy"
        message="Are you sure you want to close the proxy server?"
        confirmText="Close"
        cancelText="Cancel"
        variant="danger"
        loading={isShuttingDown}
      />
    </>
  );
}

HeaderMenu.propTypes = {
  onLogout: PropTypes.func.isRequired,
};
