"use client";

import { GripVertical, Plus, Trash2, X } from "lucide-react";
import type { Question, QuestionType } from "@/lib/types";
import { questionTypeLabels } from "@/lib/poll-utils";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";

interface QuestionEditorProps {
  question: Question;
  index: number;
  instanceId: number;
  onChange: (q: Question) => void;
  onRemove: () => void;
  canRemove: boolean;
}

const types: QuestionType[] = ["single", "multiple", "text"];

export function QuestionEditor({
  question,
  index,
  instanceId,
  onChange,
  onRemove,
  canRemove,
}: QuestionEditorProps) {
  const isChoice = question.type === "single" || question.type === "multiple";

  function update(patch: Partial<Question>) {
    onChange({ ...question, ...patch });
  }

  function setType(type: QuestionType) {
    if (type === "text") {
      update({ type, options: [] });
    } else {
      update({
        type,
        options: question.options.length ? question.options : ["", ""],
      });
    }
  }

  function updateOption(i: number, value: string) {
    const options = [...question.options];
    options[i] = value;
    update({ options });
  }

  function addOption() {
    update({ options: [...question.options, ""] });
  }

  function removeOption(i: number) {
    update({ options: question.options.filter((_, idx) => idx !== i) });
  }

  return (
    <Card className="flex flex-col gap-4 p-5">
      <div className="flex items-center justify-between">
        <span className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
          <GripVertical className="size-4" />
          Вопрос {index + 1}
        </span>
        <Button
          variant="ghost"
          size="icon"
          className="size-8 text-muted-foreground hover:text-destructive"
          onClick={onRemove}
          disabled={!canRemove}
          aria-label="Удалить вопрос"
        >
          <Trash2 className="size-4" />
        </Button>
      </div>

      <div className="flex flex-col gap-2">
        <Label htmlFor={`q-title-${instanceId}`}>Текст вопроса</Label>
        <Input
          id={`q-title-${instanceId}`}
          value={question.title}
          onChange={(e) => update({ title: e.target.value })}
          placeholder="Например: Какой ваш любимый язык программирования?"
        />
      </div>

      <div className="flex flex-col gap-2">
        <Label htmlFor={`q-desc-${instanceId}`}>
          Описание{" "}
          <span className="text-muted-foreground">(необязательно)</span>
        </Label>
        <Input
          id={`q-desc-${instanceId}`}
          value={question.description}
          onChange={(e) => update({ description: e.target.value })}
          placeholder="Дополнительное пояснение к вопросу"
        />
      </div>

      <div className="flex flex-col gap-2">
        <Label>Тип ответа</Label>
        <div className="flex flex-wrap gap-2">
          {types.map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => setType(t)}
              className={
                "rounded-lg border px-3 py-1.5 text-sm transition-colors " +
                (question.type === t
                  ? "border-primary bg-primary text-primary-foreground"
                  : "border-border bg-background text-foreground hover:bg-accent")
              }
            >
              {questionTypeLabels[t]}
            </button>
          ))}
        </div>
      </div>

      {isChoice && (
        <div className="flex flex-col gap-2">
          <Label>Варианты ответов</Label>
          <div className="flex flex-col gap-2">
            {question.options.map((option, i) => (
              <div
                key={`${instanceId}-${i}`}
                className="flex items-center gap-2"
              >
                <Input
                  value={option}
                  onChange={(e) => updateOption(i, e.target.value)}
                  placeholder={`Вариант ${i + 1}`}
                />
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-9 shrink-0 text-muted-foreground hover:text-destructive"
                  onClick={() => removeOption(i)}
                  disabled={question.options.length <= 2}
                  aria-label="Удалить вариант"
                >
                  <X className="size-4" />
                </Button>
              </div>
            ))}
          </div>
          <Button
            variant="outline"
            size="sm"
            className="self-start"
            onClick={addOption}
          >
            <Plus className="size-4" />
            Добавить вариант
          </Button>
        </div>
      )}

      {question.type === "text" && (
        <div className="rounded-lg border border-dashed border-border p-3">
          <Textarea
            disabled
            placeholder="Участник введёт свой ответ здесь"
            className="resize-none bg-muted/40"
            rows={2}
          />
        </div>
      )}
    </Card>
  );
}
