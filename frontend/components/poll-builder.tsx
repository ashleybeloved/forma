"use client";

import { useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, Loader2, Plus, Save } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import { api, ApiError } from "@/lib/api";
import type {
  NewPollRequest,
  Poll,
  Question,
  UpdatePollRequest,
} from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { QuestionEditor } from "@/components/question-editor";

interface PollBuilderProps {
  mode?: "create" | "edit";
  initialPoll?: Poll;
}

type EditableQuestion = Question & {
  clientId: number;
};

function createEmptyQuestion(clientId: number): EditableQuestion {
  return {
    clientId,
    id: clientId,
    type: "single",
    title: "",
    description: "",
    image: "",
    options: ["", ""],
  };
}

function cloneQuestions(questions: Question[]): EditableQuestion[] {
  return questions.map((question, index) => ({
    ...question,
    clientId: index + 1,
    options: [...question.options],
  }));
}

function getInitialQuestions(initialPoll?: Poll): EditableQuestion[] {
  const questions = initialPoll?.config?.questions ?? [];
  if (questions.length > 0) {
    return cloneQuestions(questions);
  }

  return [createEmptyQuestion(1)];
}

export function PollBuilder({
  mode = "create",
  initialPoll,
}: PollBuilderProps) {
  const router = useRouter();
  const isEditMode = mode === "edit" && Boolean(initialPoll);
  const initialQuestions = getInitialQuestions(initialPoll);
  const [title, setTitle] = useState(initialPoll?.title ?? "");
  const [description, setDescription] = useState(
    initialPoll?.description ?? "",
  );
  const [secured, setSecured] = useState(initialPoll?.secured ?? true);
  const [authOnly, setAuthOnly] = useState(initialPoll?.auth_only ?? false);
  const [questions, setQuestions] =
    useState<EditableQuestion[]>(initialQuestions);
  const [saving, setSaving] = useState(false);

  const nextQuestionId = useRef(
    Math.max(0, ...initialQuestions.map((question) => question.clientId)) + 1,
  );

  const backHref = initialPoll
    ? `/dashboard/polls/${initialPoll.short_id}/stats`
    : "/dashboard";
  const pageTitle = isEditMode ? "Редактирование опроса" : "Новый опрос";
  const submitLabel = isEditMode ? "Сохранить изменения" : "Сохранить опрос";

  function updateQuestion(clientId: number, q: EditableQuestion) {
    setQuestions((prev) =>
      prev.map((item) => (item.clientId === clientId ? q : item)),
    );
  }

  function removeQuestion(clientId: number) {
    setQuestions((prev) => prev.filter((item) => item.clientId !== clientId));
  }

  function addQuestion() {
    const id = nextQuestionId.current;
    nextQuestionId.current += 1;
    setQuestions((prev) => [...prev, createEmptyQuestion(id)]);
  }

  function validate(): string | null {
    if (!title.trim()) return "Введите название опроса";
    if (questions.length === 0) return "Добавьте хотя бы один вопрос";

    for (let i = 0; i < questions.length; i++) {
      const q = questions[i];
      if (!q.title.trim()) return `Введите текст вопроса №${i + 1}`;

      if (q.type !== "text") {
        const filled = q.options.filter((o) => o.trim());
        if (filled.length < 2) {
          return `В вопросе №${i + 1} нужно минимум 2 варианта ответа`;
        }
      }
    }

    return null;
  }

  function buildPayload(): NewPollRequest {
    return {
      title: title.trim(),
      description: description.trim(),
      secured,
      auth_only: authOnly,
      config: {
        questions: questions.map((q, index) => ({
          id: index + 1,
          type: q.type,
          title: q.title.trim(),
          description: q.description.trim(),
          image: q.image,
          options:
            q.type === "text"
              ? []
              : q.options.map((o) => o.trim()).filter(Boolean),
        })),
      },
    };
  }

  async function handleSave() {
    const errorMsg = validate();
    if (errorMsg) {
      toast.error(errorMsg);
      return;
    }

    const payload = buildPayload();

    setSaving(true);
    try {
      if (isEditMode && initialPoll) {
        const updatePayload: UpdatePollRequest = payload;

        await api.updatePoll(initialPoll.short_id, updatePayload);
        toast.success("Опрос обновлён");
        router.push(`/dashboard/polls/${initialPoll.short_id}/stats`);
      } else {
        const poll = await api.createPoll(payload);
        toast.success("Опрос создан");
        router.push(`/dashboard/polls/${poll.short_id}/stats`);
      }
    } catch (err) {
      const message =
        err instanceof ApiError
          ? err.message
          : isEditMode
            ? "Не удалось обновить опрос"
            : "Не удалось создать опрос";
      toast.error(message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="icon"
            className="size-9"
            render={<Link href={backHref} aria-label="Назад" />}
          >
            <ArrowLeft className="size-4" />
          </Button>
          <h1 className="text-2xl font-semibold tracking-tight">{pageTitle}</h1>
        </div>
        <Button onClick={handleSave} disabled={saving}>
          {saving ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <Save className="size-4" />
          )}
          {isEditMode ? "Сохранить" : "Создать"}
        </Button>
      </div>

      <Card className="flex flex-col gap-4 p-5">
        <div className="flex flex-col gap-2">
          <Label htmlFor="poll-title">Название</Label>
          <Input
            id="poll-title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Например: Опрос об инструментах разработки"
          />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="poll-description">
            Описание{" "}
            <span className="text-muted-foreground">(необязательно)</span>
          </Label>
          <Textarea
            id="poll-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Кратко расскажите, о чём этот опрос"
            rows={3}
          />
        </div>
      </Card>

      <Card className="flex flex-col gap-4 p-5">
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-col gap-0.5">
            <Label htmlFor="secured">Защита от повторного голосования</Label>
            <p className="text-sm text-muted-foreground">
              Один голос с одного IP-адреса и устройства
            </p>
          </div>
          <Switch id="secured" checked={secured} onCheckedChange={setSecured} />
        </div>
        <div className="h-px bg-border" />
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-col gap-0.5">
            <Label htmlFor="auth-only">Только для авторизованных</Label>
            <p className="text-sm text-muted-foreground">
              Голосовать смогут только вошедшие пользователи
            </p>
          </div>
          <Switch
            id="auth-only"
            checked={authOnly}
            onCheckedChange={setAuthOnly}
          />
        </div>
      </Card>

      <div className="flex flex-col gap-4">
        {questions.map((q, index) => (
          <QuestionEditor
            key={q.clientId}
            question={q}
            index={index}
            instanceId={q.clientId}
            onChange={(updated) =>
              updateQuestion(q.clientId, updated as EditableQuestion)
            }
            onRemove={() => removeQuestion(q.clientId)}
            canRemove={questions.length > 1}
          />
        ))}
      </div>

      <Button variant="outline" onClick={addQuestion} className="self-start">
        <Plus className="size-4" />
        Добавить вопрос
      </Button>

      <div className="flex justify-end border-t border-border pt-6">
        <Button onClick={handleSave} disabled={saving} size="lg">
          {saving ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <Save className="size-4" />
          )}
          {submitLabel}
        </Button>
      </div>
    </div>
  );
}
