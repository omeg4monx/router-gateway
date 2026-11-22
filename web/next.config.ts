import type { NextConfig } from "next";

const isStaticExport = process.env.NEXT_OUTPUT_EXPORT === "true";

const nextConfig: NextConfig = {
  output: isStaticExport ? "export" : undefined,
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
  ...(isStaticExport
    ? {}
    : {
        async rewrites() {
          return [
            {
              source: "/api/:path*",
              destination: "http://127.0.0.1:20128/api/:path*",
            },
            {
              source: "/v1/:path*",
              destination: "http://127.0.0.1:20128/v1/:path*",
            },
            {
              source: "/v1beta/:path*",
              destination: "http://127.0.0.1:20128/v1beta/:path*",
            },
          ];
        },
      }),
};

export default nextConfig;
