import { cookies } from "next/headers";
import { DashboardLayout } from "@/shared/components";
import { RAIL_PINNED_COOKIE } from "@/shared/constants/navigation";

export default async function DashboardRootLayout({ children }) {
  const railPinned = (await cookies()).get(RAIL_PINNED_COOKIE)?.value !== "false";
  return <DashboardLayout railPinned={railPinned}>{children}</DashboardLayout>;
}
