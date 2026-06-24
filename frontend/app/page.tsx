import Link from "next/link"
import {
  ArrowRight,
  BarChart3,
  Globe2,
  ListChecks,
  Lock,
  Share2,
  Zap,
} from "lucide-react"
import { MarketingHeader } from "@/components/marketing-header"
import { Logo } from "@/components/logo"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"

const features = [
  {
    icon: ListChecks,
    title: "Гибкий конструктор",
    description:
      "Одиночный и множественный выбор, текстовые ответы. Добавляйте вопросы и варианты в пару кликов.",
  },
  {
    icon: BarChart3,
    title: "Статистика в реальном времени",
    description:
      "Смотрите распределение голосов по каждому вопросу с наглядными процентами и счётчиками.",
  },
  {
    icon: Globe2,
    title: "География ответов",
    description:
      "Узнайте, из каких стран приходят голоса — топ-страны определяются автоматически.",
  },
  {
    icon: Lock,
    title: "Защита от накруток",
    description:
      "Режим «один голос» по IP и токену, а также опросы только для авторизованных.",
  },
  {
    icon: Share2,
    title: "Ссылка для шеринга",
    description:
      "Каждый опрос получает короткий адрес — делитесь им где угодно и собирайте ответы.",
  },
  {
    icon: Zap,
    title: "Быстро и просто",
    description:
      "Никаких лишних настроек. Создали опрос, отправили ссылку, получили результаты.",
  },
]

const steps = [
  {
    number: "01",
    title: "Создайте опрос",
    description: "Добавьте вопросы, варианты ответов и настройте параметры доступа.",
  },
  {
    number: "02",
    title: "Поделитесь ссылкой",
    description: "Отправьте короткую ссылку участникам — в чат, соцсети или по почте.",
  },
  {
    number: "03",
    title: "Анализируйте результаты",
    description: "Следите за голосами и географией ответов на странице статистики.",
  },
]

export default function HomePage() {
  return (
    <div className="flex min-h-svh flex-col">
      <MarketingHeader />

      <main className="flex-1">
        {/* Hero */}
        <section className="mx-auto w-full max-w-6xl px-4 pb-16 pt-16 md:pb-24 md:pt-24">
          <div className="grid items-center gap-12 lg:grid-cols-2">
            <div className="flex flex-col items-start gap-6">
              <Badge variant="secondary" className="rounded-full px-3 py-1 text-xs">
                Платформа для создания опросов
              </Badge>
              <h1 className="text-balance text-4xl font-semibold tracking-tight md:text-6xl">
                Собирайте мнения с помощью красивых опросов
              </h1>
              <p className="max-w-md text-pretty text-lg leading-relaxed text-muted-foreground">
                Forma помогает быстро создавать опросы, делиться ими по короткой
                ссылке и видеть результаты в реальном времени — с географией и
                защитой от накруток.
              </p>
              <div className="flex flex-col gap-3 sm:flex-row">
                <Button size="lg" render={<Link href="/register" />}>
                  Создать опрос
                  <ArrowRight className="size-4" />
                </Button>
                <Button size="lg" variant="outline" render={<Link href="/login" />}>
                  У меня уже есть аккаунт
                </Button>
              </div>
            </div>

            {/* Превью результатов */}
            <div className="relative">
              <div className="rounded-2xl border border-border bg-card p-6 shadow-sm">
                <div className="mb-5 flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium">Какой фреймворк вы используете?</p>
                    <p className="text-xs text-muted-foreground">248 голосов</p>
                  </div>
                  <Badge variant="secondary" className="text-xs">
                    Активен
                  </Badge>
                </div>
                <div className="flex flex-col gap-4">
                  {[
                    { label: "Next.js", value: 62 },
                    { label: "Nuxt", value: 21 },
                    { label: "SvelteKit", value: 11 },
                    { label: "Другое", value: 6 },
                  ].map((row) => (
                    <div key={row.label} className="flex flex-col gap-1.5">
                      <div className="flex items-center justify-between text-sm">
                        <span className="font-medium">{row.label}</span>
                        <span className="text-muted-foreground">{row.value}%</span>
                      </div>
                      <div className="h-2.5 w-full overflow-hidden rounded-full bg-secondary">
                        <div
                          className="h-full rounded-full bg-primary"
                          style={{ width: `${row.value}%` }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Возможности */}
        <section id="features" className="border-t border-border bg-card/40 py-16 md:py-24">
          <div className="mx-auto w-full max-w-6xl px-4">
            <div className="mx-auto mb-12 max-w-2xl text-center">
              <h2 className="text-balance text-3xl font-semibold tracking-tight md:text-4xl">
                Всё, что нужно для опросов
              </h2>
              <p className="mt-4 text-pretty text-muted-foreground">
                От создания до анализа результатов — без лишней сложности.
              </p>
            </div>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {features.map((feature) => (
                <div
                  key={feature.title}
                  className="flex flex-col gap-3 rounded-xl border border-border bg-background p-6"
                >
                  <span className="flex size-10 items-center justify-center rounded-lg bg-accent text-accent-foreground">
                    <feature.icon className="size-5" />
                  </span>
                  <h3 className="font-medium">{feature.title}</h3>
                  <p className="text-sm leading-relaxed text-muted-foreground">
                    {feature.description}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Как это работает */}
        <section id="how" className="py-16 md:py-24">
          <div className="mx-auto w-full max-w-6xl px-4">
            <div className="mx-auto mb-12 max-w-2xl text-center">
              <h2 className="text-balance text-3xl font-semibold tracking-tight md:text-4xl">
                Три шага до результата
              </h2>
            </div>
            <div className="grid gap-8 md:grid-cols-3">
              {steps.map((step) => (
                <div key={step.number} className="flex flex-col gap-3">
                  <span className="font-mono text-sm font-medium text-primary">
                    {step.number}
                  </span>
                  <div className="h-px w-full bg-border" />
                  <h3 className="text-lg font-medium">{step.title}</h3>
                  <p className="text-sm leading-relaxed text-muted-foreground">
                    {step.description}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* CTA */}
        <section className="px-4 pb-20">
          <div className="mx-auto flex w-full max-w-6xl flex-col items-center gap-6 rounded-2xl border border-border bg-primary px-6 py-14 text-center text-primary-foreground">
            <h2 className="text-balance text-3xl font-semibold tracking-tight md:text-4xl">
              Готовы создать свой первый опрос?
            </h2>
            <p className="max-w-md text-pretty text-primary-foreground/80">
              Регистрация занимает меньше минуты. Никаких карт и подписок.
            </p>
            <Button size="lg" variant="secondary" render={<Link href="/register" />}>
              Начать бесплатно
              <ArrowRight className="size-4" />
            </Button>
          </div>
        </section>
      </main>

      <footer className="border-t border-border py-8">
        <div className="mx-auto flex w-full max-w-6xl flex-col items-center justify-between gap-4 px-4 text-sm text-muted-foreground sm:flex-row">
          <Logo />
          <p>© {new Date().getFullYear()} Forma. Пет-проект.</p>
        </div>
      </footer>
    </div>
  )
}
