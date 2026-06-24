"use client"

import Link from "next/link"
import { ArrowRight } from "lucide-react"
import { useAuth } from "@/components/auth-provider"
import { Logo } from "@/components/logo"
import { Button } from "@/components/ui/button"

export function MarketingHeader() {
  const { isAuthenticated, isLoading } = useAuth()

  return (
    <header className="sticky top-0 z-40 border-b border-border/60 bg-background/80 backdrop-blur-sm">
      <div className="mx-auto flex h-16 w-full max-w-6xl items-center justify-between px-4">
        <Logo />
        <nav className="hidden items-center gap-8 text-sm text-muted-foreground md:flex">
          <a href="#features" className="transition-colors hover:text-foreground">
            Возможности
          </a>
          <a href="#how" className="transition-colors hover:text-foreground">
            Как это работает
          </a>
        </nav>
        <div className="flex items-center gap-2">
          {isLoading ? (
            <div className="h-9 w-32" />
          ) : isAuthenticated ? (
            <Button size="sm" render={<Link href="/dashboard" />}>
              Дашборд
              <ArrowRight className="size-4" />
            </Button>
          ) : (
            <>
              <Button variant="ghost" size="sm" render={<Link href="/login" />}>
                Войти
              </Button>
              <Button size="sm" render={<Link href="/register" />}>
                Начать бесплатно
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  )
}
