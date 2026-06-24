import type { NextRequest } from "next/server";

const RAW_BACKEND_URL = (
  process.env.FORMA_BACKEND_URL ??
  process.env.API_URL ??
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080"
).replace(/\/$/, "");

const BACKEND_URL = RAW_BACKEND_URL.endsWith("/api")
  ? RAW_BACKEND_URL
  : `${RAW_BACKEND_URL}/api`;

const HOP_BY_HOP_HEADERS = new Set([
  "connection",
  "content-length",
  "host",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade",
]);

export const dynamic = "force-dynamic";

type RouteContext = {
  params: Promise<{ path: string[] }>;
};

function buildTargetUrl(path: string[], request: NextRequest) {
  const pathname = path.map(encodeURIComponent).join("/");
  const url = new URL(`${BACKEND_URL}/${pathname}`);
  url.search = request.nextUrl.search;
  return url;
}

function buildRequestHeaders(request: NextRequest) {
  const headers = new Headers();

  for (const [key, value] of request.headers.entries()) {
    const lowerKey = key.toLowerCase();
    if (HOP_BY_HOP_HEADERS.has(lowerKey)) continue;
    if (lowerKey === "origin" || lowerKey === "referer") continue;
    headers.set(key, value);
  }

  return headers;
}

function buildResponseHeaders(upstreamHeaders: Headers) {
  const headers = new Headers();

  for (const [key, value] of upstreamHeaders.entries()) {
    const lowerKey = key.toLowerCase();
    if (HOP_BY_HOP_HEADERS.has(lowerKey)) continue;
    if (lowerKey.startsWith("access-control-")) continue;
    headers.append(key, value);
  }

  return headers;
}

async function proxy(request: NextRequest, { params }: RouteContext) {
  const { path } = await params;
  const url = buildTargetUrl(path, request);
  const method = request.method;

  let body: BodyInit | undefined;
  if (method !== "GET" && method !== "HEAD") {
    body = await request.arrayBuffer();
  }

  let upstream: Response;
  try {
    upstream = await fetch(url, {
      method,
      headers: buildRequestHeaders(request),
      body,
      redirect: "manual",
      cache: "no-store",
    });
  } catch {
    return Response.json({ error: "failed to reach backend" }, { status: 502 });
  }

  return new Response(upstream.body, {
    status: upstream.status,
    statusText: upstream.statusText,
    headers: buildResponseHeaders(upstream.headers),
  });
}

export async function GET(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function POST(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function PATCH(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function DELETE(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function PUT(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}

export async function OPTIONS(request: NextRequest, context: RouteContext) {
  return proxy(request, context);
}
