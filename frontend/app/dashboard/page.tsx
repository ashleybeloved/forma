"use client";

import Link from "next/link";
import useSWR from "swr";
import { useState } from "react";
import { FileBarChart, Loader2, Plus } from "lucide-react";
import { toast } from "sonner";
import { api, ApiError } from "@/lib/api";
import { config } from "@/lib/config";
import type { Poll } from "@/lib/types";
import { useAuth } from "@/components/auth-provider";
import { PollCard } from "@/components/poll-card";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

const PAGE_SIZE = config.pollsPageSize;

export default function DashboardPage() {
  const { user } = useAuth();
  const [page, setPage] = useState(0);

  const { data, error, isLoading, mutate } = useSWR<Poll[]>(
    ["my-polls", page],
    () => api.getMyPolls(PAGE_SIZE, page * PAGE_SIZE),
    { keepPreviousData: true, shouldRetryOnError: false },
  );

  // Бэкенд возвращает 404, если опросов нет — трактуем как пустой список.
  const notFound = error instanceof ApiError && error.status === 404;
  const polls = data ?? [];
  const isEmpty = !isLoading && (notFound || polls.length === 0) && page === 0;
  const realError = error && !notFound;

  async function handleDelete(shortId: string) {
    try {
      await api.deletePoll(shortId);
      toast.success("Опрос удалён");
      await mutate();
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "Не удалось удалить опрос";
      toast.error(message);
    }
  }

  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-col gap-1">
        <h1 className="text-2xl font-semibold tracking-tight">
          Привет, {user?.username ?? "пользователь"}
        </h1>
        <p className="text-muted-foreground">Ваши опросы и их результаты</p>
      </div>

      {isLoading && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i} className="flex flex-col gap-4 p-5">
              <Skeleton className="h-5 w-20" />
              <Skeleton className="h-5 w-3/4" />
              <Skeleton className="h-4 w-full" />
              <Skeleton className="mt-4 h-9 w-full" />
            </Card>
          ))}
        </div>
      )}

      {realError && (
        <Card className="flex flex-col items-center gap-3 p-10 text-center">
          <p className="text-muted-foreground">
            Не удалось загрузить опросы. Проверьте, что бэкенд запущен.
          </p>
          <Button variant="outline" onClick={() => mutate()}>
            Повторить
          </Button>
        </Card>
      )}

      {isEmpty && (
        <Card className="flex flex-col items-center gap-4 border-dashed p-12 text-center">
          <span className="flex size-14 items-center justify-center rounded-full bg-accent text-accent-foreground">
            <FileBarChart className="size-6" />
          </span>
          <div className="flex flex-col gap-1">
            <h2 className="text-lg font-medium">Пока нет опросов</h2>
            <p className="max-w-sm text-sm text-muted-foreground">
              Создайте первый опрос, чтобы начать собирать ответы и смотреть
              статистику.
            </p>
          </div>
          <Button render={<Link href="/dashboard/new" />}>
            <Plus className="size-4" />
            Создать опрос
          </Button>
        </Card>
      )}

      {!isLoading && !realError && polls.length > 0 && (
        <>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {polls.map((poll) => (
              <PollCard key={poll.id} poll={poll} onDelete={handleDelete} />
            ))}
          </div>

          <div className="flex items-center justify-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={page === 0}
              onClick={() => setPage((p) => Math.max(0, p - 1))}
            >
              Назад
            </Button>
            <span className="text-sm text-muted-foreground">
              Страница {page + 1}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={polls.length < PAGE_SIZE}
              onClick={() => setPage((p) => p + 1)}
            >
              {data === undefined ? (
                <Loader2 className="size-4 animate-spin" />
              ) : (
                "Вперёд"
              )}
            </Button>
          </div>
        </>
      )}
    </div>
  );
}
