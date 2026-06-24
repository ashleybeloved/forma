import { config } from "./config";
import type {
  NewPollRequest,
  NewVoteRequest,
  Poll,
  Stats,
  UpdatePollRequest,
  UserResponse,
} from "./types";

/**
 * Ошибка API с HTTP-статусом и сообщением от сервера.
 */
export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  // Не выбрасывать ошибку на 4xx — вернуть тело как есть (используется редко).
  query?: Record<string, string | number | undefined>;
}

async function request<T>(
  path: string,
  options: RequestOptions = {},
): Promise<T> {
  const { method = "GET", body, query } = options;

  let url = `${config.apiUrl}${path}`;
  if (query) {
    const params = new URLSearchParams();
    for (const [key, value] of Object.entries(query)) {
      if (value !== undefined) params.set(key, String(value));
    }
    const qs = params.toString();
    if (qs) url += `?${qs}`;
  }

  let res: Response;
  try {
    res = await fetch(url, {
      method,
      // Передаём httpOnly cookie forma_token / гостевой токен.
      credentials: "include",
      headers: body ? { "Content-Type": "application/json" } : undefined,
      body: body ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError(
      0,
      "Не удалось связаться с сервером. Проверьте, что бэкенд запущен и доступен через фронтенд-прокси.",
    );
  }

  // Пустой ответ (204 No Content и т.п.)
  if (res.status === 204) {
    return undefined as T;
  }

  const text = await res.text();
  let data: unknown = undefined;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = text;
    }
  }

  if (!res.ok) {
    const message =
      (data && typeof data === "object" && "error" in data
        ? String((data as { error: unknown }).error)
        : undefined) ?? `Ошибка запроса (${res.status})`;
    throw new ApiError(res.status, message);
  }

  return data as T;
}

/* ----------------------------- Аутентификация ----------------------------- */

export const api = {
  async register(username: string, password: string): Promise<void> {
    await request<{ message: string }>("/register", {
      method: "POST",
      body: { username, password },
    });
  },

  async login(username: string, password: string): Promise<void> {
    await request<{ message: string }>("/login", {
      method: "POST",
      body: { username, password },
    });
  },

  async logout(): Promise<void> {
    await request<{ message: string }>("/logout", { method: "POST" });
  },

  me(): Promise<UserResponse> {
    return request<UserResponse>("/me");
  },

  /* -------------------------------- Опросы -------------------------------- */

  getMyPolls(limit: number, offset: number): Promise<Poll[]> {
    return request<Poll[]>("/poll", { query: { limit, offset } });
  },

  createPoll(payload: NewPollRequest): Promise<Poll> {
    return request<Poll>("/poll", { method: "POST", body: payload });
  },

  updatePoll(shortId: string, payload: UpdatePollRequest): Promise<void> {
    return request<void>(`/poll/${shortId}`, {
      method: "PATCH",
      body: payload,
    });
  },

  deletePoll(shortId: string): Promise<void> {
    return request<void>(`/poll/${shortId}`, { method: "DELETE" });
  },

  getPollStats(shortId: string): Promise<Stats> {
    return request<Stats>(`/poll/${shortId}/stats`);
  },

  /* ------------------------- Публичные (голосование) ----------------------- */

  getPoll(shortId: string): Promise<Poll> {
    return request<Poll>(`/poll/${shortId}`);
  },

  async vote(shortId: string, payload: NewVoteRequest): Promise<void> {
    await request<{ message: string }>(`/poll/${shortId}/vote`, {
      method: "POST",
      body: payload,
    });
  },

  checkVote(shortId: string): Promise<{ voted: boolean }> {
    return request<{ voted: boolean }>(`/poll/${shortId}/check`, {
      method: "POST",
    });
  },
};
