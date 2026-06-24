import { cn } from "@/lib/utils"
import Link from "next/link"

export function Logo({ className }: { className?: string }) {
  return (
    <Link
      href="/"
      className={cn("inline-flex items-center gap-2 font-semibold tracking-tight", className)}
    >
      <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-primary text-primary-foreground">
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <path d="M4 5h16" />
          <path d="M4 12h10" />
          <path d="M4 19h6" />
        </svg>
      </span>
      <span className="text-lg">Forma</span>
    </Link>
  )
}
