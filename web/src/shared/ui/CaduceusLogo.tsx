import type { SVGProps } from "react";

interface CaduceusLogoProps extends Omit<SVGProps<SVGSVGElement>, "width" | "height"> {
  size?: number;
  decorative?: boolean;
}

export function CaduceusLogo({
  size = 28,
  decorative = false,
  className = "",
  ...rest
}: CaduceusLogoProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 32 32"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      aria-hidden={decorative ? true : undefined}
      role={decorative ? undefined : "img"}
      aria-label={decorative ? undefined : "OMEG4MON ROUTER caduceus mark"}
      {...rest}
    >
      {/* Central staff */}
      <line
        x1="16"
        y1="4"
        x2="16"
        y2="29"
        stroke="#E56A4A"
        strokeWidth="2"
        strokeLinecap="round"
      />
      {/* Top finial sphere */}
      <circle cx="16" cy="4.5" r="2" fill="#E56A4A" />
      {/* Geometric Hermes wings */}
      <path
        d="M16 9.5L8.5 6.5L5 9.5L11.5 11.5L16 10.5"
        stroke="#E56A4A"
        strokeWidth="1.65"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M16 9.5L23.5 6.5L27 9.5L20.5 11.5L16 10.5"
        stroke="#E56A4A"
        strokeWidth="1.65"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      {/* Intertwined left serpent curve */}
      <path
        d="M12.2 12.6C10.5 13.8 11.2 16 16 17.5C20.5 18.9 21 21.4 16 23.2C13.6 24.1 13.4 25.8 15.2 26.8"
        stroke="#E56A4A"
        strokeWidth="1.75"
        strokeLinecap="round"
      />
      {/* Intertwined right serpent curve */}
      <path
        d="M19.8 12.6C21.5 13.8 20.8 16 16 17.5C11.5 18.9 11 21.4 16 23.2C18.4 24.1 18.6 25.8 16.8 26.8"
        stroke="#E56A4A"
        strokeWidth="1.75"
        strokeLinecap="round"
      />
    </svg>
  );
}
