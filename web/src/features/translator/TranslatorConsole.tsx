"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRightLeft, Copy, Pause, Play, Radio, TerminalSquare, Trash2 } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Label } from "@/shared/ui/label";
import { useTranslatePayloadMutation } from "./hooks";
import {
  type ConsoleLogEntry,
  parseJsonObject,
  SAMPLE_PAYLOADS,
  TRANSLATOR_MODES,
  TRANSLATOR_PROTOCOLS,
  type TranslateFormValues,
  type TranslatePayloadResponse,
  type TranslatorMode,
  type TranslatorProtocol,
  translateFormSchema,
} from "./schemas";

type LogLevelFilter = "ALL" | "DEBUG" | "INFO" | "WARN" | "ERROR";

interface KeyedLogEntry {
  seq: number;
  entry: ConsoleLogEntry;
}

const MAX_LIVE_LOGS = 250;

function levelBadgeVariant(level: string): "destructive" | "warning" | "secondary" | "default" {
  const upper = level.toUpperCase();
  if (upper === "ERROR") {
    return "destructive";
  }
  if (upper === "WARN" || upper === "WARNING") {
    return "warning";
  }
  if (upper === "DEBUG") {
    return "secondary";
  }
  return "default";
}

export function TranslatorConsole() {
  const translateMutation = useTranslatePayloadMutation();
  const [result, setResult] = useState<TranslatePayloadResponse | null>(null);
  const [translateError, setTranslateError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const [levelFilter, setLevelFilter] = useState<LogLevelFilter>("ALL");
  const [paused, setPaused] = useState(false);
  const [streamConnected, setStreamConnected] = useState(false);
  const [logs, setLogs] = useState<KeyedLogEntry[]>([]);
  const seqRef = useRef(0);
  const pausedRef = useRef(paused);
  pausedRef.current = paused;

  const form = useForm<TranslateFormValues>({
    resolver: zodResolver(translateFormSchema),
    defaultValues: {
      mode: "request",
      from_protocol: "openai-chat",
      to_protocol: "anthropic",
      payloadText: SAMPLE_PAYLOADS.request["openai-chat"],
    },
  });

  const watchedMode = form.watch("mode");
  const watchedFrom = form.watch("from_protocol");

  useEffect(() => {
    const params = new URLSearchParams({ tail: "50" });
    if (levelFilter !== "ALL") {
      params.set("level", levelFilter);
    }
    const source = new EventSource(`/api/translator/logs/stream?${params.toString()}`, {
      withCredentials: true,
    });

    source.onopen = () => setStreamConnected(true);
    source.onerror = () => setStreamConnected(false);

    const onLog = (event: MessageEvent<string>) => {
      if (pausedRef.current) {
        return;
      }
      try {
        const entry = JSON.parse(event.data) as ConsoleLogEntry;
        const fp = `${entry.time}|${entry.level}|${entry.message}|${JSON.stringify(entry.attrs ?? null)}`;
        setLogs((prev) => {
          if (
            prev.some(
              (item) =>
                `${item.entry.time}|${item.entry.level}|${item.entry.message}|${JSON.stringify(item.entry.attrs ?? null)}` ===
                fp,
            )
          ) {
            return prev;
          }
          seqRef.current += 1;
          const combined = [...prev, { seq: seqRef.current, entry }];
          return combined.length > MAX_LIVE_LOGS ? combined.slice(-MAX_LIVE_LOGS) : combined;
        });
      } catch {
        // Ignore malformed SSE frame.
      }
    };

    source.addEventListener("log", onLog as EventListener);
    return () => {
      source.removeEventListener("log", onLog as EventListener);
      source.close();
      setStreamConnected(false);
    };
  }, [levelFilter]);

  const handleLoadSample = (mode: TranslatorMode, proto: TranslatorProtocol) => {
    form.setValue("payloadText", SAMPLE_PAYLOADS[mode][proto], { shouldValidate: true });
  };

  const onSubmit = form.handleSubmit((values) => {
    setTranslateError(null);
    setCopied(false);
    const parsed = parseJsonObject(values.payloadText);
    if (!parsed) {
      setTranslateError("Payload must be a valid JSON object.");
      return;
    }
    translateMutation.mutate(
      {
        mode: values.mode,
        from_protocol: values.from_protocol,
        to_protocol: values.to_protocol,
        payload: parsed,
      },
      {
        onSuccess: (res) => setResult(res),
        onError: (err) => {
          setResult(null);
          setTranslateError(err.message);
        },
      },
    );
  });

  const formattedOutput = result ? JSON.stringify(result.output, null, 2) : "";

  const handleCopyOutput = async () => {
    if (!formattedOutput) {
      return;
    }
    try {
      await navigator.clipboard.writeText(formattedOutput);
      setCopied(true);
      setTimeout(() => setCopied(false), 1800);
    } catch {
      setCopied(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-xl font-semibold tracking-tight text-[var(--text-primary)]">
            Protocol Translator &amp; Live Console
          </h2>
          <p className="text-sm text-[var(--text-secondary)]">
            Test Canonical IR request/response translations across OpenAI Chat, Anthropic, OpenAI
            Responses, and Gemini, and inspect live structured server logs.
          </p>
        </div>
        <Badge variant={streamConnected ? "success" : "secondary"}>
          <Radio className="mr-1.5 h-3 w-3" aria-hidden="true" />
          {streamConnected ? "SSE Connected" : "SSE Reconnecting"}
        </Badge>
      </div>

      {translateError ? (
        <Alert variant="destructive">
          <AlertTitle>Translation Error</AlertTitle>
          <AlertDescription>{translateError}</AlertDescription>
        </Alert>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <ArrowRightLeft className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            Interactive Protocol Workbench
          </CardTitle>
          <CardDescription>
            Executes the exact production gateway Canonical IR and fast-path translation pipeline.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid gap-3 sm:grid-cols-4">
              <div className="space-y-1">
                <Label htmlFor="tr-mode">Mode</Label>
                <select
                  id="tr-mode"
                  {...form.register("mode")}
                  onChange={(e) => {
                    const nextMode = e.target.value as TranslatorMode;
                    form.setValue("mode", nextMode);
                    handleLoadSample(nextMode, watchedFrom);
                  }}
                  className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 py-2 text-sm text-[var(--text-primary)] focus:border-[#E56A4A] focus:outline-none"
                >
                  {TRANSLATOR_MODES.map((m) => (
                    <option key={m} value={m}>
                      {m}
                    </option>
                  ))}
                </select>
              </div>

              <div className="space-y-1">
                <Label htmlFor="tr-from">Source Protocol</Label>
                <select
                  id="tr-from"
                  {...form.register("from_protocol")}
                  onChange={(e) => {
                    const nextFrom = e.target.value as TranslatorProtocol;
                    form.setValue("from_protocol", nextFrom);
                    handleLoadSample(watchedMode, nextFrom);
                  }}
                  className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 py-2 font-mono text-xs text-[var(--text-primary)] focus:border-[#E56A4A] focus:outline-none"
                >
                  {TRANSLATOR_PROTOCOLS.map((p) => (
                    <option key={p} value={p}>
                      {p}
                    </option>
                  ))}
                </select>
              </div>

              <div className="space-y-1">
                <Label htmlFor="tr-to">Target Protocol</Label>
                <select
                  id="tr-to"
                  {...form.register("to_protocol")}
                  className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 py-2 font-mono text-xs text-[var(--text-primary)] focus:border-[#E56A4A] focus:outline-none"
                >
                  {TRANSLATOR_PROTOCOLS.map((p) => (
                    <option key={p} value={p}>
                      {p}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex items-end gap-2">
                <Button type="submit" className="w-full" disabled={translateMutation.isPending}>
                  {translateMutation.isPending ? "Translating..." : "Translate"}
                </Button>
              </div>
            </div>

            <div className="grid gap-4 lg:grid-cols-2">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="tr-input">Input JSON</Label>
                  <button
                    type="button"
                    onClick={() => handleLoadSample(watchedMode, watchedFrom)}
                    className="font-mono text-xs text-[var(--accent-text)] hover:underline"
                  >
                    Reset sample
                  </button>
                </div>
                <textarea
                  id="tr-input"
                  rows={12}
                  spellCheck={false}
                  className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 font-mono text-xs text-[var(--text-primary)] focus:border-[#E56A4A] focus:outline-none"
                  {...form.register("payloadText")}
                />
                {form.formState.errors.payloadText ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {form.formState.errors.payloadText.message}
                  </p>
                ) : null}
              </div>

              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Label htmlFor="tr-output">Translated Output</Label>
                    {result ? (
                      <>
                        <Badge variant={result.fast_path ? "success" : "secondary"}>
                          {result.fast_path ? "fast-path" : "canonical-ir"}
                        </Badge>
                        <span className="font-mono text-xs text-[var(--text-muted)]">
                          {result.duration_us} us
                        </span>
                      </>
                    ) : null}
                  </div>
                  <Button
                    type="button"
                    size="sm"
                    variant="secondary"
                    disabled={!formattedOutput}
                    onClick={handleCopyOutput}
                  >
                    <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{copied ? "Copied" : "Copy"}</span>
                  </Button>
                </div>
                <textarea
                  id="tr-output"
                  readOnly
                  rows={12}
                  value={formattedOutput}
                  placeholder="Click Translate to inspect wire output..."
                  className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 font-mono text-xs text-[var(--text-primary)] focus:outline-none"
                />
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <CardTitle className="flex items-center gap-2 text-base">
                <TerminalSquare className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                Live Server Console Stream
              </CardTitle>
              <CardDescription>
                Real-time structured slog stream via /api/translator/logs/stream (ring buffer tail +
                SSE).
              </CardDescription>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <select
                aria-label="Filter log level"
                value={levelFilter}
                onChange={(e) => {
                  setLogs([]);
                  setLevelFilter(e.target.value as LogLevelFilter);
                }}
                className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 py-1.5 font-mono text-xs text-[var(--text-primary)]"
              >
                <option value="ALL">ALL LEVELS</option>
                <option value="DEBUG">DEBUG+</option>
                <option value="INFO">INFO+</option>
                <option value="WARN">WARN+</option>
                <option value="ERROR">ERROR</option>
              </select>
              <Button
                type="button"
                size="sm"
                variant="secondary"
                onClick={() => setPaused((p) => !p)}
              >
                {paused ? (
                  <>
                    <Play className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>Resume</span>
                  </>
                ) : (
                  <>
                    <Pause className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>Pause</span>
                  </>
                )}
              </Button>
              <Button type="button" size="sm" variant="secondary" onClick={() => setLogs([])}>
                <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                <span>Clear</span>
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="max-h-96 overflow-y-auto rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 font-mono text-xs">
            {logs.length === 0 ? (
              <p className="py-8 text-center text-[var(--text-muted)]">
                Waiting for structured log events...
              </p>
            ) : (
              <div className="space-y-1.5">
                {logs.map(({ seq, entry }) => (
                  <div
                    key={seq}
                    className="flex flex-wrap items-baseline gap-2 border-b border-[var(--border-subtle)] pb-1.5 last:border-b-0"
                  >
                    <span className="shrink-0 text-[var(--text-muted)]">
                      {entry.time.slice(11, 23)}
                    </span>
                    <Badge variant={levelBadgeVariant(entry.level)}>{entry.level}</Badge>
                    <span className="text-[var(--text-primary)]">{entry.message}</span>
                    {entry.attrs && Object.keys(entry.attrs).length > 0 ? (
                      <span className="text-[var(--text-secondary)]">
                        {JSON.stringify(entry.attrs)}
                      </span>
                    ) : null}
                  </div>
                ))}
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
