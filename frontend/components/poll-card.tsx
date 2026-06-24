"use client";

import { useState } from "react";
import Link from "next/link";
import {
  BarChart3,
  Check,
  Copy,
  ExternalLink,
  Loader2,
  Lock,
  MoreVertical,
  Pencil,
  Trash2,
  Users,
} from "lucide-react";
import { toast } from "sonner";
import { copyText } from "@/lib/clipboard";
import type { Poll } from "@/lib/types";
import { formatDate, pollPublicUrl } from "@/lib/poll-utils";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface PollCardProps {
  poll: Poll;
  onDelete: (shortId: string) => Promise<void>;
}

export function PollCard({ poll, onDelete }: PollCardProps) {
  const [copied, setCopied] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const questionsCount = poll.config?.questions?.length ?? 0;

  async function copyLink() {
    try {
      const copied = await copyText(pollPublicUrl(poll.short_id));
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

  async function handleDelete() {
    if (deleting) return;
    setDeleting(true);
    try {
      await onDelete(poll.short_id);
    } finally {
      setDeleting(false);
    }
  }

  return (
    <Card className="relative flex h-full flex-col gap-4 p-5 transition-colors hover:border-primary/40">
      <DropdownMenu>
        <DropdownMenuTrigger
          render={
            <Button
              variant="ghost"
              size="icon"
              className="absolute top-3 right-3 z-10 size-8 shrink-0"
              aria-label="Действия с опросом"
              title="Действия"
            />
          }
        >
          {deleting ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <MoreVertical className="size-4" />
          )}
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-44">
          <DropdownMenuItem
            render={<Link href={`/dashboard/polls/${poll.short_id}/stats`} />}
          >
            <BarChart3 className="size-4" />
            Статистика
          </DropdownMenuItem>
          <DropdownMenuItem
            render={<Link href={`/dashboard/polls/${poll.short_id}/edit`} />}
          >
            <Pencil className="size-4" />
            Редактировать
          </DropdownMenuItem>
          <DropdownMenuItem
            render={
              <Link href={pollPublicUrl(poll.short_id)} target="_blank" />
            }
          >
            <ExternalLink className="size-4" />
            Открыть опрос
          </DropdownMenuItem>
          <DropdownMenuItem onClick={copyLink}>
            <Copy className="size-4" />
            Копировать ссылку
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem variant="destructive" onClick={handleDelete}>
            <Trash2 className="size-4" />
            Удалить
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <Link
        href={`/dashboard/polls/${poll.short_id}/stats`}
        className="flex flex-col gap-3 pr-10"
      >
        <div className="flex flex-col gap-1.5">
          <h3 className="text-pretty font-medium leading-snug">{poll.title}</h3>
          {poll.description && (
            <p className="line-clamp-2 text-sm text-muted-foreground">
              {poll.description}
            </p>
          )}
        </div>

        {(poll.secured || poll.auth_only) && (
          <div className="flex flex-wrap items-center gap-2">
            {poll.secured && (
              <Badge variant="secondary" className="gap-1 text-xs">
                <Lock className="size-3" />
                Защищён
              </Badge>
            )}
            {poll.auth_only && (
              <Badge variant="secondary" className="text-xs">
                Только для вошедших
              </Badge>
            )}
          </div>
        )}
      </Link>

      <div className="mt-auto flex items-center justify-end gap-1.5">
        <Button
          variant="outline"
          size="icon-xs"
          render={<Link href={`/dashboard/polls/${poll.short_id}/stats`} />}
          aria-label="Статистика"
          title="Статистика"
        >
          <BarChart3 className="size-3.5" />
        </Button>
        <Button
          variant="outline"
          size="icon-xs"
          render={<Link href={`/dashboard/polls/${poll.short_id}/edit`} />}
          aria-label="Редактировать"
          title="Редактировать"
        >
          <Pencil className="size-3.5" />
        </Button>
        <Button
          variant="outline"
          size="icon-xs"
          onClick={copyLink}
          aria-label="Скопировать ссылку"
          title="Скопировать ссылку"
        >
          {copied ? (
            <Check className="size-3.5" />
          ) : (
            <Copy className="size-3.5" />
          )}
        </Button>
      </div>

      <div className="flex items-center justify-between border-t border-border pt-4 text-xs text-muted-foreground">
        <span className="flex items-center gap-1.5">
          <Users className="size-3.5" />
          {questionsCount}{" "}
          {questionsCount === 1
            ? "вопрос"
            : questionsCount < 5
              ? "вопроса"
              : "вопросов"}
        </span>
        <span>{formatDate(poll.created_at)}</span>
      </div>
    </Card>
  );
}
