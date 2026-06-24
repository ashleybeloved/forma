"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { LogOut, Plus, User } from "lucide-react";
import { toast } from "sonner";
import { useAuth } from "@/components/auth-provider";
import { Logo } from "@/components/logo";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function AppHeader() {
  const { user, logout } = useAuth();
  const router = useRouter();

  async function handleLogout() {
    try {
      await logout();
      toast.success("Вы вышли из аккаунта");
      router.push("/");
    } catch {
      toast.error("Не удалось выйти");
    }
  }

  const initial = user?.username?.charAt(0).toUpperCase() ?? "?";

  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur-sm">
      <div className="mx-auto flex h-16 w-full max-w-5xl items-center justify-between px-4">
        <Logo />
        <div className="flex items-center gap-2">
          <Button size="sm" render={<Link href="/dashboard/new" />}>
            <Plus className="size-4" />
            <span className="hidden sm:inline">Новый опрос</span>
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button variant="ghost" size="icon" className="rounded-full" />
              }
            >
              <Avatar className="size-8">
                <AvatarFallback className="bg-accent text-accent-foreground text-sm">
                  {initial}
                </AvatarFallback>
              </Avatar>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <div className="flex items-center gap-2 px-2 py-1.5 text-sm">
                <User className="size-4 text-muted-foreground" />
                <span className="truncate">{user?.username ?? "Аккаунт"}</span>
              </div>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={handleLogout} variant="destructive">
                <LogOut className="size-4" />
                Выйти
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  );
}
