import type { QuestionType } from "./types"

export function formatDate(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return ""
  return new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "short",
    year: "numeric",
  }).format(date)
}

export const questionTypeLabels: Record<QuestionType, string> = {
  single: "Один ответ",
  multiple: "Несколько ответов",
  text: "Текстовый ответ",
}

/** Ссылка на публичную страницу голосования. */
export function pollPublicUrl(shortId: string): string {
  if (typeof window === "undefined") return `/p/${shortId}`
  return `${window.location.origin}/p/${shortId}`
}

/** Преобразует код страны (ISO) в эмодзи-флаг. */
export function countryFlag(code: string): string {
  if (!code || code.length !== 2) return "🌐"
  const upper = code.toUpperCase()
  const A = 0x1f1e6
  return String.fromCodePoint(
    A + (upper.charCodeAt(0) - 65),
    A + (upper.charCodeAt(1) - 65),
  )
}
