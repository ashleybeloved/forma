"use client";

import { useMemo, useState } from "react";
import useSWR from "swr";
import Link from "next/link";
import { CheckCircle2, Loader2, Lock, Send } from "lucide-react";
import { toast } from "sonner";
import { api, ApiError } from "@/lib/api";
import type { Answer, Poll } from "@/lib/types";
import { useAuth } from "@/components/auth-provider";
import { Logo } from "@/components/logo";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Skeleton } from "@/components/ui/skeleton";

interface VoteFormProps {
  shortId: string;
}

// Хранилище ответов: ключ конкретного инстанса вопроса -> выбранные опции / текст
type AnswersState = Record<string, string[]>;

export function VoteForm({ shortId }: VoteFormProps) {
  const { isAuthenticated } = useAuth();
  const [answers, setAnswers] = useState<AnswersState>({});
  const [submitting, setSubmitting] = useState(false);
  const [justVoted, setJustVoted] = useState(false);

  const pollReq = useSWR<Poll>(
    ["public-poll", shortId],
    () => api.getPoll(shortId),
    {
      shouldRetryOnError: false,
      revalidateOnFocus: false,
    },
  );
  const checkReq = useSWR(
    ["check-vote", shortId],
    () => api.checkVote(shortId),
    { shouldRetryOnError: false, revalidateOnFocus: false },
  );

  const poll = pollReq.data;
  const questions = useMemo(() => poll?.config?.questions ?? [], [poll]);
  const alreadyVoted = checkReq.data?.voted ?? false;

  function getQuestionKey(questionId: number, index: number) {
    return `${questionId}-${index}`;
  }

  function setSingle(questionKey: string, value: string) {
    setAnswers((prev) => ({ ...prev, [questionKey]: [value] }));
  }

  function toggleMultiple(
    questionKey: string,
    value: string,
    checked: boolean,
  ) {
    setAnswers((prev) => {
      const current = prev[questionKey] ?? [];
      return {
        ...prev,
        [questionKey]: checked
          ? [...current, value]
          : current.filter((v) => v !== value),
      };
    });
  }

  function setText(questionKey: string, value: string) {
    setAnswers((prev) => ({ ...prev, [questionKey]: value ? [value] : [] }));
  }

  async function handleSubmit() {
    // Проверяем, что на каждый вопрос есть ответ.
    for (const [index, q] of questions.entries()) {
      const selected = answers[getQuestionKey(q.id, index)] ?? [];
      if (
        selected.length === 0 ||
        (q.type === "text" && !selected[0]?.trim())
      ) {
        toast.error("Ответьте на все вопросы");
        return;
      }
    }

    const payload: Answer[] = questions.map((q, index) => {
      const questionKey = getQuestionKey(q.id, index);
      return {
        question_id: q.id,
        options:
          q.type === "text"
            ? [(answers[questionKey]?.[0] ?? "").trim()]
            : (answers[questionKey] ?? []),
      };
    });

    setSubmitting(true);
    try {
      await api.vote(shortId, { answers: payload });
      setJustVoted(true);
      toast.success("Ваш голос учтён");
      checkReq.mutate();
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "Не удалось проголосовать";
      toast.error(message);
    } finally {
      setSubmitting(false);
    }
  }

  const loading = pollReq.isLoading || checkReq.isLoading;

  return (
    <div className="flex min-h-svh flex-col">
      <header className="border-b border-border">
        <div className="mx-auto flex h-16 w-full max-w-2xl items-center justify-between px-4">
          <Logo />
          {!isAuthenticated && (
            <Button variant="ghost" size="sm" render={<Link href="/login" />}>
              Войти
            </Button>
          )}
        </div>
      </header>

      <main className="mx-auto w-full max-w-2xl flex-1 px-4 py-8">
        {loading && (
          <div className="flex flex-col gap-4">
            <Skeleton className="h-8 w-3/4" />
            <Skeleton className="h-40 w-full" />
            <Skeleton className="h-40 w-full" />
          </div>
        )}

        {!loading && pollReq.error && (
          <Card className="flex flex-col items-center gap-3 p-12 text-center">
            <h1 className="text-lg font-medium">Опрос не найден</h1>
            <p className="text-sm text-muted-foreground">
              Возможно, ссылка устарела или опрос был удалён.
            </p>
            <Button variant="outline" render={<Link href="/" />}>
              На главную
            </Button>
          </Card>
        )}

        {!loading && poll && (justVoted || alreadyVoted) && (
          <Card className="flex flex-col items-center gap-3 p-12 text-center">
            <span className="flex size-14 items-center justify-center rounded-full bg-accent text-accent-foreground">
              <CheckCircle2 className="size-7" />
            </span>
            <h1 className="text-xl font-medium">Спасибо за участие!</h1>
            <p className="max-w-sm text-sm text-muted-foreground">
              Ваш голос в опросе «{poll.title}» уже учтён.
            </p>
            <Button variant="outline" render={<Link href="/" />}>
              На главную
            </Button>
          </Card>
        )}

        {!loading && poll && !justVoted && !alreadyVoted && (
          <div className="flex flex-col gap-6">
            <div className="flex flex-col gap-2">
              <h1 className="text-pretty text-2xl font-semibold tracking-tight">
                {poll.title}
              </h1>
              {poll.description && (
                <p className="text-pretty leading-relaxed text-muted-foreground">
                  {poll.description}
                </p>
              )}
              {poll.auth_only && !isAuthenticated && (
                <div className="flex items-center gap-2 rounded-lg border border-border bg-muted/40 px-3 py-2 text-sm text-muted-foreground">
                  <Lock className="size-4 shrink-0" />
                  Для голосования в этом опросе нужно{" "}
                  <Link
                    href="/login"
                    className="font-medium text-primary hover:underline"
                  >
                    войти в аккаунт
                  </Link>
                </div>
              )}
            </div>

            {questions.map((q, index) => {
              const questionKey = getQuestionKey(q.id, index);
              return (
                <Card key={questionKey} className="flex flex-col gap-4 p-5">
                  <div className="flex flex-col gap-1">
                    <span className="text-xs font-medium text-muted-foreground">
                      Вопрос {index + 1} из {questions.length}
                    </span>
                    <h2 className="text-pretty font-medium leading-snug">
                      {q.title}
                    </h2>
                    {q.description && (
                      <p className="text-sm text-muted-foreground">
                        {q.description}
                      </p>
                    )}
                  </div>

                  {q.type === "single" && (
                    <RadioGroup
                      value={answers[questionKey]?.[0] ?? ""}
                      onValueChange={(v) => setSingle(questionKey, v)}
                      className="flex flex-col gap-2"
                    >
                      {q.options.map((option, i) => (
                        <Label
                          key={`${questionKey}-${i}`}
                          htmlFor={`q${questionKey}-o${i}`}
                          className="flex cursor-pointer items-center gap-3 rounded-lg border border-border px-4 py-3 font-normal transition-colors hover:bg-accent has-data-[state=checked]:border-primary has-data-[state=checked]:bg-accent"
                        >
                          <RadioGroupItem
                            id={`q${questionKey}-o${i}`}
                            value={option}
                          />
                          {option}
                        </Label>
                      ))}
                    </RadioGroup>
                  )}

                  {q.type === "multiple" && (
                    <div className="flex flex-col gap-2">
                      {q.options.map((option, i) => {
                        const checked = (answers[questionKey] ?? []).includes(
                          option,
                        );
                        return (
                          <Label
                            key={`${questionKey}-${i}`}
                            htmlFor={`q${questionKey}-o${i}`}
                            className="flex cursor-pointer items-center gap-3 rounded-lg border border-border px-4 py-3 font-normal transition-colors hover:bg-accent has-data-[state=checked]:border-primary has-data-[state=checked]:bg-accent"
                          >
                            <Checkbox
                              id={`q${questionKey}-o${i}`}
                              checked={checked}
                              onCheckedChange={(c) =>
                                toggleMultiple(questionKey, option, Boolean(c))
                              }
                            />
                            {option}
                          </Label>
                        );
                      })}
                    </div>
                  )}

                  {q.type === "text" && (
                    <Textarea
                      value={answers[questionKey]?.[0] ?? ""}
                      onChange={(e) => setText(questionKey, e.target.value)}
                      placeholder="Ваш ответ"
                      rows={3}
                    />
                  )}
                </Card>
              );
            })}

            <Button
              size="lg"
              onClick={handleSubmit}
              disabled={submitting || (poll.auth_only && !isAuthenticated)}
            >
              {submitting ? (
                <Loader2 className="size-4 animate-spin" />
              ) : (
                <Send className="size-4" />
              )}
              Отправить ответы
            </Button>
          </div>
        )}
      </main>
    </div>
  );
}
