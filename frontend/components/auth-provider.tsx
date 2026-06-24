"use client"

import { createContext, useContext, type ReactNode } from "react"
import useSWR from "swr"
import { api, ApiError } from "@/lib/api"
import type { UserResponse } from "@/lib/types"

interface AuthContextValue {
  user: UserResponse | null
  isLoading: boolean
  isAuthenticated: boolean
  refresh: () => Promise<unknown>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const { data, error, isLoading, mutate } = useSWR<UserResponse>(
    "auth/me",
    () => api.me(),
    {
      shouldRetryOnError: false,
      revalidateOnFocus: false,
    },
  )

  // 401 — это нормальное состояние «не вошёл», не считаем ошибкой загрузки.
  const unauthenticated = error instanceof ApiError && error.status === 401

  const value: AuthContextValue = {
    user: data ?? null,
    isLoading,
    isAuthenticated: Boolean(data) && !unauthenticated,
    refresh: () => mutate(),
    logout: async () => {
      await api.logout()
      await mutate(undefined, { revalidate: false })
    },
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used within AuthProvider")
  return ctx
}
