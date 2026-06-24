import type { ReactNode } from "react"
import { AppHeader } from "@/components/app-header"
import { RequireAuth } from "@/components/require-auth"

export default function DashboardLayout({ children }: { children: ReactNode }) {
  return (
    <RequireAuth>
      <div className="flex min-h-svh flex-col">
        <AppHeader />
        <main className="mx-auto w-full max-w-5xl flex-1 px-4 py-8">{children}</main>
      </div>
    </RequireAuth>
  )
}
