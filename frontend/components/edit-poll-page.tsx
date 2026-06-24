"use client"

import Link from "next/link"
import useSWR from "swr"
import { ArrowLeft, Loader2, ShieldAlert } from "lucide-react"
import { api, ApiError } from "@/lib/api"
import type { Poll } from "@/lib/types"
import { useAuth } from "@/components/auth-provider"
import { PollBuilder } from "@/components/poll-builder"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"

interface EditPollPageProps {
  shortId: string
}

export function EditPollPage({ shortId }: EditPollPageProps) {
  const { user, isLoading: authLoading } = useAuth()
  const pollReq = useSWR<Poll>(["edit-poll", shortId], () => api.getPoll(shortId), {
    shouldRetryOnError: false,
    revalidateOnFocus: false,
  })

  if (authLoading || pollReq.isLoading) {
    return (
      <div className="flex flex-col gap-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-36 w-full" />
        <Skeleton className="h-36 w-full" />
        <Skeleton className="h-52 w-full" />
      </div>
    )
  }

  if (pollReq.error || !pollReq.data) {
    const message =
      pollReq.error instanceof ApiError ? pollReq.error.message : "Не удалось загрузить опрос"

    return (
      <Card className="flex flex-col items-center gap-4 p-12 text-center">
        <p className="text-muted-foreground">{message}</p>
        <Button variant="outline" render={<Link href="/dashboard" />}>
          <ArrowLeft className="size-4" />
          К дашборду
        </Button>
      </Card>
    )
  }

  if (pollReq.data.creator_id !== user?.id) {
    return (
      <Card className="flex flex-col items-center gap-4 p-12 text-center">
        <span className="flex size-14 items-center justify-center rounded-full bg-accent text-accent-foreground">
          <ShieldAlert className="size-6" />
        </span>
        <div className="flex flex-col gap-1">
          <h1 className="text-lg font-medium">Нет доступа к редактированию</h1>
          <p className="text-sm text-muted-foreground">
            Редактировать опрос может только его создатель.
          </p>
        </div>
        <Button variant="outline" render={<Link href="/dashboard" />}>
          <ArrowLeft className="size-4" />
          К дашборду
        </Button>
      </Card>
    )
  }

  return <PollBuilder mode="edit" initialPoll={pollReq.data} />
}
