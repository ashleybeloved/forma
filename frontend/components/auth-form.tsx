"use client";

import type React from "react";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { useAuth } from "@/components/auth-provider";
import { Logo } from "@/components/logo";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface AuthFormProps {
  mode: "login" | "register";
}

export function AuthForm({ mode }: AuthFormProps) {
  const router = useRouter();
  const { refresh, isAuthenticated, isLoading } = useAuth();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const isLogin = mode === "login";

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.replace("/dashboard");
    }
  }, [isAuthenticated, isLoading, router]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (submitting) return;

    if (username.trim().length < 3) {
      toast.error("Имя пользователя должно быть не короче 3 символов");
      return;
    }
    if (password.length < 6) {
      toast.error("Пароль должен быть не короче 6 символов");
      return;
    }

    setSubmitting(true);
    try {
      if (isLogin) {
        await api.login(username.trim(), password);
      } else {
        await api.register(username.trim(), password);
      }
      await refresh();
      toast.success(isLogin ? "С возвращением!" : "Аккаунт создан");
      router.push("/dashboard");
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "Что-то пошло не так";
      toast.error(message);
    } finally {
      setSubmitting(false);
    }
  }

  if (isLoading || isAuthenticated) {
    return (
      <div className="flex min-h-svh items-center justify-center px-4 py-12">
        <Loader2 className="size-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="flex min-h-svh flex-col items-center justify-center px-4 py-12">
      <div className="mb-8">
        <Logo />
      </div>
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">
            {isLogin ? "Вход в аккаунт" : "Создать аккаунт"}
          </CardTitle>
          <CardDescription>
            {isLogin
              ? "Введите данные, чтобы продолжить работу с опросами"
              : "Зарегистрируйтесь, чтобы создавать опросы"}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="username">Имя пользователя</Label>
              <Input
                id="username"
                autoComplete="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="ivan"
                disabled={submitting}
                required
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="password">Пароль</Label>
              <Input
                id="password"
                type="password"
                autoComplete={isLogin ? "current-password" : "new-password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                disabled={submitting}
                required
              />
            </div>
            <Button type="submit" className="mt-2 w-full" disabled={submitting}>
              {submitting && <Loader2 className="size-4 animate-spin" />}
              {isLogin ? "Войти" : "Зарегистрироваться"}
            </Button>
          </form>

          <p className="mt-6 text-center text-sm text-muted-foreground">
            {isLogin ? "Ещё нет аккаунта? " : "Уже есть аккаунт? "}
            <Link
              href={isLogin ? "/register" : "/login"}
              className="font-medium text-primary hover:underline"
            >
              {isLogin ? "Зарегистрироваться" : "Войти"}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
