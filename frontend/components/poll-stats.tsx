"use client";

import Link from "next/link";
import { useState } from "react";
import useSWR from "swr";
import {
  ArrowLeft,
  BarChart3,
  Check,
  Copy,
  ExternalLink,
  Globe2,
  Loader2,
  Pencil,
  Users,
} from "lucide-react";
import { toast } from "sonner";
import { api, ApiError } from "@/lib/api";
import { copyText } from "@/lib/clipboard";
import type { Poll, Stats } from "@/lib/types";
import { countryFlag, pollPublicUrl } from "@/lib/poll-utils";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";

interface PollStatsProps {
  shortId: string;
}

export function PollStats({ shortId }: PollStatsProps) {
  const [copied, setCopied] = useState(false);

  const pollReq = useSWR<Poll>(["poll", shortId], () => api.getPoll(shortId), {
    shouldRetryOnError: false,
  });
  const statsReq = useSWR<Stats>(
    ["stats", shortId],
    () => api.getPollStats(shortId),
    {
      shouldRetryOnError: false,
    },
  );

  const poll = pollReq.data;
  const stats = statsReq.data;
  const isLoading = pollReq.isLoading || statsReq.isLoading;
  const error = pollReq.error || statsReq.error;

  async function copyLink() {
    try {
      const copied = await copyText(pollPublicUrl(shortId));
      if (!copied) {
        throw new Error("copy failed");
      }
      setCopied(true);
      toast.success("Ссылка скопирована");
      setTimeout(() => setCopied(false), 2000);
    } catch {
      toast.error("Не удалось скопировать ссылку");
    }
  }

  if (isLoading) {
    return (
      <div className="flex flex-col gap-6">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-4 sm:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-24 w-full" />
          ))}
        </div>
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (error) {
    const message =
      error instanceof ApiError
        ? error.message
        : "Не удалось загрузить статистику";
    return (
      <Card className="flex flex-col items-center gap-4 p-12 text-center">
        <p className="text-muted-foreground">{message}</p>
        <Button variant="outline" render={<Link href="/dashboard" />}>
          <ArrowLeft className="size-4" />К дашборду
        </Button>
      </Card>
    );
  }

  const questions = poll?.config?.questions ?? [];
  const totalVotes = stats?.total_votes ?? 0;

  // Сопоставляем результаты статистики с вопросами по id.
  function questionTitle(id: number): string {
    return questions.find((q) => q.id === id)?.title ?? `Вопрос #${id}`;
  }

  return (
    <div className="flex flex-col gap-6">
      {/* Шапка */}
      <div className="flex flex-col gap-4">
        <Button
          variant="ghost"
          size="sm"
          className="-ml-2 self-start"
          render={<Link href="/dashboard" />}
        >
          <ArrowLeft className="size-4" />К дашборду
        </Button>
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-pretty text-2xl font-semibold tracking-tight">
              {poll?.title}
            </h1>
            {poll?.description && (
              <p className="text-muted-foreground">{poll.description}</p>
            )}
          </div>
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" size="sm" onClick={copyLink}>
              {copied ? (
                <Check className="size-4" />
              ) : (
                <Copy className="size-4" />
              )}
              Ссылка
            </Button>
            <Button
              variant="outline"
              size="sm"
              render={<Link href={`/dashboard/polls/${shortId}/edit`} />}
            >
              <Pencil className="size-4" />
              Редактировать
            </Button>
            <Button
              size="sm"
              render={<Link href={pollPublicUrl(shortId)} target="_blank" />}
            >
              <ExternalLink className="size-4" />
              Открыть
            </Button>
          </div>
        </div>
      </div>

      {/* Сводка */}
      <div className="grid gap-4 sm:grid-cols-3">
        <Card className="flex flex-col gap-1 p-5">
          <span className="flex items-center gap-1.5 text-sm text-muted-foreground">
            <Users className="size-4" />
            Всего голосов
          </span>
          <span className="text-3xl font-semibold">{totalVotes}</span>
        </Card>
        <Card className="flex flex-col gap-1 p-5">
          <span className="flex items-center gap-1.5 text-sm text-muted-foreground">
            <BarChart3 className="size-4" />
            Вопросов
          </span>
          <span className="text-3xl font-semibold">{questions.length}</span>
        </Card>
        <Card className="flex flex-col gap-1 p-5">
          <span className="flex items-center gap-1.5 text-sm text-muted-foreground">
            <Globe2 className="size-4" />
            Стран
          </span>
          <span className="text-3xl font-semibold">
            {stats?.top_countries?.length ?? 0}
          </span>
        </Card>
      </div>

      {totalVotes === 0 ? (
        <Card className="flex flex-col items-center gap-3 border-dashed p-12 text-center">
          <span className="flex size-14 items-center justify-center rounded-full bg-accent text-accent-foreground">
            <BarChart3 className="size-6" />
          </span>
          <h2 className="text-lg font-medium">Пока нет голосов</h2>
          <p className="max-w-sm text-sm text-muted-foreground">
            Поделитесь ссылкой на опрос, чтобы начать собирать ответы.
          </p>
          <Button onClick={copyLink} variant="outline">
            <Copy className="size-4" />
            Скопировать ссылку
          </Button>
        </Card>
      ) : (
        <div className="flex flex-col gap-6">
          {/* Результаты по вопросам */}
          {(stats?.results ?? []).map((result, resultIndex) => {
            const maxVotes = Math.max(...result.options.map((o) => o.votes), 1);
            return (
              <Card
                key={`${result.id}-${resultIndex}`}
                className="flex flex-col gap-4 p-5"
              >
                <h3 className="font-medium">{questionTitle(result.id)}</h3>
                <div className="flex flex-col gap-4">
                  {result.options.map((opt, i) => {
                    const isTop = opt.votes === maxVotes;
                    return (
                      <div
                        key={`${result.id}-${i}`}
                        className="flex flex-col gap-1.5"
                      >
                        <div className="flex items-center justify-between gap-3 text-sm">
                          <span className="font-medium">
                            {opt.option || "—"}
                          </span>
                          <span className="shrink-0 text-muted-foreground">
                            {opt.votes} · {opt.percentage.toFixed(1)}%
                          </span>
                        </div>
                        <div className="h-2.5 w-full overflow-hidden rounded-full bg-secondary">
                          <div
                            className={
                              "h-full rounded-full " +
                              (isTop ? "bg-primary" : "bg-primary/45")
                            }
                            style={{
                              width: `${Math.max(opt.percentage, 1.5)}%`,
                            }}
                          />
                        </div>
                      </div>
                    );
                  })}
                </div>
              </Card>
            );
          })}

          {/* Топ стран */}
          {stats && stats.top_countries.length > 0 && (
            <Card className="flex flex-col gap-4 p-5">
              <h3 className="flex items-center gap-2 font-medium">
                <Globe2 className="size-4 text-muted-foreground" />
                География голосов
              </h3>
              <div className="flex flex-col gap-2">
                {stats.top_countries.map((country) => {
                  const pct = totalVotes
                    ? (country.votes / totalVotes) * 100
                    : 0;
                  return (
                    <div
                      key={country.country_code}
                      className="flex items-center justify-between gap-3 rounded-lg border border-border px-3 py-2"
                    >
                      <span className="flex items-center gap-2 text-sm font-medium">
                        <span className="text-lg leading-none">
                          {countryFlag(country.country_code)}
                        </span>
                        {country.country_code || "Неизвестно"}
                      </span>
                      <Badge variant="secondary">
                        {country.votes} ({pct.toFixed(0)}%)
                      </Badge>
                    </div>
                  );
                })}
              </div>
            </Card>
          )}
        </div>
      )}
    </div>
  );
}
