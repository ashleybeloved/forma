/**
 * Конфигурация приложения Forma.
 *
 * Во фронтенде по умолчанию используем same-origin прокси `/api`, чтобы браузер
 * не упирался в CORS при обращении к Go-бэкенду с другого порта.
 *
 * Клиент всегда обращается к `/api`, а реальный адрес бэкенда настраивается
 * на стороне Next.js route handler через `FORMA_BACKEND_URL` или `API_URL`.
 */
export const config = {
  apiUrl: "/api",
  appName: "Forma",
  // Количество опросов, подгружаемых за раз в дашборде.
  pollsPageSize: 12,
} as const;
