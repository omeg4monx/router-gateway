"use client";

import React from "react";

export default function LogoMark({ className = "size-full" }) {
  return (
    <span
      className={`relative inline-flex items-center justify-center ${className}`.trim()}
      style={{
        color: "var(--color-primary, #E56A4A)",
      }}
    >
      <svg
        aria-hidden="true"
        viewBox="0 0 100 100"
        width="100%"
        height="100%"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="size-full overflow-hidden block"
      >
        {/* Squircle container background */}
        <rect
          x="4"
          y="4"
          width="92"
          height="92"
          rx="22"
          className="fill-surface-2 dark:fill-[#18191C] stroke-border dark:stroke-[#26282D]"
          strokeWidth="2"
        />

        <path d="M 29 50 H 39 M 60 43 L 69 34 H 76 M 60 57 L 69 66 H 76" stroke="currentColor" strokeWidth="7" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx="50" cy="50" r="11" stroke="currentColor" strokeWidth="7" />
      </svg>
    </span>
  );
}
