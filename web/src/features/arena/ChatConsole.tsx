import {
  Check,
  Clock,
  Copy,
  Cpu,
  MessageSquare,
  Plus,
  Search,
  Send,
  SlidersHorizontal,
  Sparkles,
  Trash2,
  User,
  X,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  useArenaModels,
  useChatSessionDetail,
  useChatSessions,
  useCreateChatSessionMutation,
  useDeleteChatSessionMutation,
  useSendChatMessageMutation,
  useUpdateChatSessionMutation,
} from "./hooks";
import {
  computeSessionTelemetry,
  filterChatSessions,
  formatLatencyMs,
  formatTokenCount,
} from "./playgroundUtils";
import {
  chatMessageFormSchema,
  createChatSessionFormSchema,
  updateChatSessionFormSchema,
} from "./schemas";

const STARTER_PROMPTS = {
  en: [
    "Summarize the key trade-offs between latency and context length.",
    "Write a concise curl command to test an OpenAI-compatible /v1/chat/completions endpoint.",
    "Explain how circuit-breaker failover works in an AI model gateway.",
  ],
  id: [
    "Ringkas trade-off utama antara latensi dan panjang konteks model AI.",
    "Buat perintah curl singkat untuk menguji endpoint /v1/chat/completions yang kompatibel dengan OpenAI.",
    "Jelaskan cara kerja circuit-breaker failover pada AI model gateway.",
  ],
} as const;

export function ChatConsole() {
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";

  const modelsQuery = useArenaModels();
  const sessionsQuery = useChatSessions();

  const createSession = useCreateChatSessionMutation();
  const updateSession = useUpdateChatSessionMutation();
  const deleteSession = useDeleteChatSessionMutation();
  const sendMessage = useSendChatMessageMutation();

  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null);
  const [sessionSearch, setSessionSearch] = useState("");
  const [createPanelOpen, setCreatePanelOpen] = useState(false);
  const [paramsOpen, setParamsOpen] = useState(false);

  // Create session form state
  const [newTitle, setNewTitle] = useState("");
  const [newModel, setNewModel] = useState("");
  const [newSystemPrompt, setNewSystemPrompt] = useState("");
  const [newTemperature, setNewTemperature] = useState("");
  const [newMaxTokens, setNewMaxTokens] = useState("");

  // Edit active session parameters state
  const [editTitle, setEditTitle] = useState("");
  const [editModel, setEditModel] = useState("");
  const [editSystemPrompt, setEditSystemPrompt] = useState("");
  const [editTemperature, setEditTemperature] = useState("");
  const [editMaxTokens, setEditMaxTokens] = useState("");
  const [paramsSaved, setParamsSaved] = useState(false);

  // Composer & feedback state
  const [draft, setDraft] = useState("");
  const [validationError, setValidationError] = useState<string | null>(null);
  const [copiedMessageId, setCopiedMessageId] = useState<string | null>(null);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const composerRef = useRef<HTMLTextAreaElement>(null);
  const copyTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const savedTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const models = modelsQuery.data?.models ?? [];
  const sessions = sessionsQuery.data?.sessions ?? [];
  const activeSessionId = selectedSessionId ?? sessions[0]?.id ?? null;

  const detailQuery = useChatSessionDetail(activeSessionId);
  const activeSession = detailQuery.data?.session;
  const messages = detailQuery.data?.messages ?? [];

  const visibleSessions = useMemo(
    () => filterChatSessions(sessions, sessionSearch),
    [sessions, sessionSearch],
  );

  const telemetry = useMemo(() => computeSessionTelemetry(messages), [messages]);

  const showCreatePanel = createPanelOpen || (!sessionsQuery.isLoading && sessions.length === 0);

  // Clean up timers on unmount
  useEffect(() => {
    return () => {
      if (copyTimerRef.current) {
        clearTimeout(copyTimerRef.current);
      }
      if (savedTimerRef.current) {
        clearTimeout(savedTimerRef.current);
      }
    };
  }, []);

  // Sync edit form whenever active session changes
  useEffect(() => {
    if (!activeSession) {
      setParamsOpen(false);
      return;
    }
    setEditTitle(activeSession.title);
    setEditModel(activeSession.model);
    setEditSystemPrompt(activeSession.system_prompt);
    setEditTemperature(
      activeSession.temperature !== undefined ? String(activeSession.temperature) : "",
    );
    setEditMaxTokens(activeSession.max_tokens > 0 ? String(activeSession.max_tokens) : "");
  }, [activeSession]);

  // Auto-scroll message stream when new messages arrive or generation starts
  const messageCount = messages.length;
  useEffect(() => {
    if (messageCount === 0 && !sendMessage.isPending) {
      return;
    }
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth", block: "end" });
  }, [messageCount, sendMessage.isPending]);

  const clearMutationErrors = () => {
    setValidationError(null);
    createSession.reset();
    updateSession.reset();
    deleteSession.reset();
    sendMessage.reset();
  };

  const handleCopyMessage = (messageId: string, content: string) => {
    if (typeof navigator === "undefined" || !navigator.clipboard) {
      return;
    }
    void navigator.clipboard.writeText(content).then(() => {
      setCopiedMessageId(messageId);
      if (copyTimerRef.current) {
        clearTimeout(copyTimerRef.current);
      }
      copyTimerRef.current = setTimeout(() => {
        setCopiedMessageId(null);
      }, 1800);
    });
  };

  const handleCreateSession = (e: React.FormEvent) => {
    e.preventDefault();
    clearMutationErrors();
    const chosenModel = newModel.trim() || models[0]?.id || "";
    const parsed = createChatSessionFormSchema.safeParse({
      title: newTitle,
      model: chosenModel,
      system_prompt: newSystemPrompt,
      temperature: newTemperature.trim() === "" ? undefined : Number(newTemperature),
      max_tokens: newMaxTokens.trim() === "" ? undefined : Number(newMaxTokens),
    });
    if (!parsed.success) {
      setValidationError(parsed.error.issues[0]?.message ?? "Model is required");
      return;
    }
    createSession.mutate(
      {
        title: parsed.data.title || undefined,
        model: parsed.data.model,
        system_prompt: parsed.data.system_prompt || undefined,
        temperature: parsed.data.temperature,
        max_tokens: parsed.data.max_tokens,
      },
      {
        onSuccess: (created) => {
          setSelectedSessionId(created.id);
          setNewTitle("");
          setNewSystemPrompt("");
          setNewTemperature("");
          setNewMaxTokens("");
          setCreatePanelOpen(false);
          setParamsOpen(false);
        },
      },
    );
  };

  const handleQuickStartSession = () => {
    const defaultModel = models[0]?.id;
    if (!defaultModel) {
      setCreatePanelOpen(true);
      return;
    }
    clearMutationErrors();
    createSession.mutate(
      { model: defaultModel },
      {
        onSuccess: (created) => {
          setSelectedSessionId(created.id);
          setCreatePanelOpen(false);
        },
      },
    );
  };

  const handleUpdateSession = (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeSession) {
      return;
    }
    clearMutationErrors();
    const parsed = updateChatSessionFormSchema.safeParse({
      title: editTitle,
      model: editModel,
      system_prompt: editSystemPrompt,
      temperature: editTemperature.trim() === "" ? undefined : Number(editTemperature),
      max_tokens: editMaxTokens.trim() === "" ? undefined : Number(editMaxTokens),
    });
    if (!parsed.success) {
      setValidationError(parsed.error.issues[0]?.message ?? "Invalid session parameters");
      return;
    }
    updateSession.mutate(
      {
        id: activeSession.id,
        body: {
          title: parsed.data.title,
          model: parsed.data.model,
          system_prompt: parsed.data.system_prompt,
          temperature: parsed.data.temperature,
          max_tokens: parsed.data.max_tokens,
        },
      },
      {
        onSuccess: () => {
          setParamsSaved(true);
          if (savedTimerRef.current) {
            clearTimeout(savedTimerRef.current);
          }
          savedTimerRef.current = setTimeout(() => {
            setParamsSaved(false);
          }, 2000);
        },
      },
    );
  };

  const submitMessage = () => {
    if (!activeSessionId || sendMessage.isPending) {
      return;
    }
    clearMutationErrors();
    const parsed = chatMessageFormSchema.safeParse({ content: draft });
    if (!parsed.success) {
      setValidationError(parsed.error.issues[0]?.message ?? "Message cannot be empty");
      return;
    }
    sendMessage.mutate(
      { id: activeSessionId, body: { content: parsed.data.content } },
      {
        onSuccess: () => {
          setDraft("");
        },
      },
    );
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    submitMessage();
  };

  const handleComposerKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      submitMessage();
    }
  };

  const handleStarterPrompt = (promptText: string) => {
    setDraft(promptText);
    composerRef.current?.focus();
  };

  const activeError =
    modelsQuery.error ??
    sessionsQuery.error ??
    detailQuery.error ??
    createSession.error ??
    updateSession.error ??
    deleteSession.error ??
    sendMessage.error;

  const starterPrompts = isId ? STARTER_PROMPTS.id : STARTER_PROMPTS.en;
  const quickModels = models.slice(0, 3);

  return (
    <div className="space-y-5">
      {/* Compact, proportional page header */}
      <section className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div className="min-w-0">
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            {isId ? "Studio Interaktif" : "Interactive Studio"}
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Playground
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            {isId
              ? "Uji model, atur system prompt, dan pantau latensi serta penggunaan token per giliran langsung melalui gateway."
              : "Test models, tune system prompts, and inspect per-turn latency and token telemetry directly through the gateway."}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="secondary" className="tabular-nums">
            {sessions.length} {isId ? "SESI" : sessions.length === 1 ? "SESSION" : "SESSIONS"}
          </Badge>
          <Badge variant="secondary" className="tabular-nums">
            {models.length} {isId ? "MODEL" : models.length === 1 ? "MODEL" : "MODELS"}
          </Badge>
        </div>
      </section>

      <EmberRule />

      {activeError || validationError ? (
        <Alert variant="destructive" className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <AlertTitle>{isId ? "Kesalahan Playground" : "Playground error"}</AlertTitle>
            <AlertDescription>{validationError ?? activeError?.message}</AlertDescription>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            aria-label={isId ? "Tutup pesan kesalahan" : "Dismiss error"}
            onClick={clearMutationErrors}
          >
            <X className="h-3.5 w-3.5" aria-hidden="true" />
          </Button>
        </Alert>
      ) : null}

      {/* Shared model datalist for create & edit forms */}
      <datalist id="playground-model-options">
        {models.map((m) => (
          <option key={m.id} value={m.id} />
        ))}
      </datalist>

      {/* Unified split-pane Playground workspace */}
      <div className="overflow-hidden rounded-xl border border-[var(--border-subtle)] bg-[var(--bg-surface)]">
        <div className="grid lg:h-[calc(100vh-14.5rem)] lg:min-h-[560px] lg:max-h-[840px] lg:grid-cols-[292px_minmax(0,1fr)] xl:grid-cols-[312px_minmax(0,1fr)]">
          {/* Left rail: Sessions & New Session drawer */}
          <aside className="flex min-h-0 flex-col border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)]/55 lg:border-r lg:border-b-0">
            <div className="flex items-center justify-between gap-2 border-b border-[var(--border-subtle)] px-4 py-3">
              <div className="flex items-center gap-2">
                <MessageSquare
                  className="h-4 w-4 text-[var(--accent-primary)]"
                  aria-hidden="true"
                />
                <h3 className="text-sm font-semibold text-[var(--text-primary)]">
                  {isId ? "Sesi" : "Sessions"}
                </h3>
                <span className="rounded bg-[var(--bg-surface-raised)] px-1.5 py-0.5 font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                  {sessions.length}
                </span>
              </div>

              <Button
                type="button"
                size="xs"
                variant={createPanelOpen ? "secondary" : "default"}
                onClick={() => {
                  clearMutationErrors();
                  setCreatePanelOpen((prev) => !prev);
                }}
              >
                {createPanelOpen ? (
                  <>
                    <X className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{isId ? "Tutup" : "Close"}</span>
                  </>
                ) : (
                  <>
                    <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{isId ? "Baru" : "New"}</span>
                  </>
                )}
              </Button>
            </div>

            {/* Collapsible New Session composer */}
            {showCreatePanel ? (
              <form
                onSubmit={handleCreateSession}
                className="space-y-3 border-b border-[var(--border-subtle)] bg-[var(--bg-surface)] p-3.5"
              >
                <div className="space-y-1">
                  <Label htmlFor="playground-new-model" className="text-xs">
                    Model
                  </Label>
                  <Input
                    id="playground-new-model"
                    list="playground-model-options"
                    placeholder={models[0]?.id ?? "openai/gpt-4o"}
                    value={newModel}
                    onChange={(e) => setNewModel(e.target.value)}
                    className="h-9 font-mono text-xs"
                  />
                  {quickModels.length > 0 ? (
                    <div className="flex flex-wrap gap-1 pt-1">
                      {quickModels.map((m) => (
                        <button
                          key={m.id}
                          type="button"
                          onClick={() => setNewModel(m.id)}
                          className="max-w-full truncate rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--text-secondary)] transition-colors hover:border-[var(--accent-primary)] hover:text-[var(--text-primary)]"
                        >
                          {m.id}
                        </button>
                      ))}
                    </div>
                  ) : null}
                </div>

                <div className="space-y-1">
                  <Label htmlFor="playground-new-title" className="text-xs">
                    {isId ? "Judul (opsional)" : "Title (optional)"}
                  </Label>
                  <Input
                    id="playground-new-title"
                    placeholder={isId ? "Eksperimen prompt..." : "Prompt experiment..."}
                    value={newTitle}
                    onChange={(e) => setNewTitle(e.target.value)}
                    className="h-9 text-xs"
                  />
                </div>

                <div className="space-y-1">
                  <Label htmlFor="playground-new-sys" className="text-xs">
                    {isId ? "System Prompt (opsional)" : "System Prompt (optional)"}
                  </Label>
                  <textarea
                    id="playground-new-sys"
                    rows={2}
                    placeholder={
                      isId ? "Anda adalah asisten yang ringkas." : "You are a concise assistant."
                    }
                    value={newSystemPrompt}
                    onChange={(e) => setNewSystemPrompt(e.target.value)}
                    className="w-full resize-none rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 text-xs text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)]"
                  />
                </div>

                <div className="grid grid-cols-2 gap-2">
                  <div className="space-y-1">
                    <Label htmlFor="playground-new-temp" className="text-xs">
                      Temp
                    </Label>
                    <Input
                      id="playground-new-temp"
                      type="number"
                      step="0.1"
                      min="0"
                      max="2"
                      placeholder="0.7"
                      value={newTemperature}
                      onChange={(e) => setNewTemperature(e.target.value)}
                      className="h-8 font-mono text-xs tabular-nums"
                    />
                  </div>
                  <div className="space-y-1">
                    <Label htmlFor="playground-new-max-tokens" className="text-xs">
                      Max Tokens
                    </Label>
                    <Input
                      id="playground-new-max-tokens"
                      type="number"
                      min="0"
                      max="128000"
                      placeholder="auto"
                      value={newMaxTokens}
                      onChange={(e) => setNewMaxTokens(e.target.value)}
                      className="h-8 font-mono text-xs tabular-nums"
                    />
                  </div>
                </div>

                <Button
                  type="submit"
                  size="sm"
                  className="w-full"
                  disabled={createSession.isPending}
                >
                  <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {createSession.isPending
                      ? isId
                        ? "Membuat..."
                        : "Creating..."
                      : isId
                        ? "Buat Sesi Playground"
                        : "Create Session"}
                  </span>
                </Button>
              </form>
            ) : null}

            {/* Session search filter when list grows */}
            {sessions.length >= 4 ? (
              <div className="border-b border-[var(--border-subtle)] px-3 py-2.5">
                <div className="relative">
                  <Search
                    className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
                    aria-hidden="true"
                  />
                  <Input
                    value={sessionSearch}
                    onChange={(e) => setSessionSearch(e.target.value)}
                    placeholder={isId ? "Cari sesi atau model..." : "Filter sessions or model..."}
                    aria-label={isId ? "Cari sesi" : "Filter sessions"}
                    className="h-8 pl-8 text-xs"
                  />
                </div>
              </div>
            ) : null}

            {/* Scrollable session list */}
            <div className="max-h-[260px] flex-1 space-y-1.5 overflow-y-auto p-2.5 lg:max-h-none">
              {sessionsQuery.isLoading ? (
                <div className="space-y-2 p-1">
                  {[0, 1, 2].map((item) => (
                    <Skeleton
                      key={item}
                      className="h-14 rounded-lg border border-[var(--border-subtle)]"
                    />
                  ))}
                </div>
              ) : visibleSessions.length > 0 ? (
                visibleSessions.map((session) => {
                  const isActive = activeSessionId === session.id;
                  return (
                    <div
                      key={session.id}
                      className={`group relative flex items-center justify-between gap-2 rounded-lg border px-3 py-2.5 transition-colors ${
                        isActive
                          ? "border-[var(--accent-primary)]/60 bg-[var(--bg-surface)]"
                          : "border-transparent bg-transparent hover:border-[var(--border-subtle)] hover:bg-[var(--bg-surface)]/70"
                      }`}
                    >
                      {isActive ? (
                        <span
                          className="absolute top-2 bottom-2 left-0 w-0.5 rounded-r bg-[var(--accent-primary)]"
                          aria-hidden="true"
                        />
                      ) : null}
                      <button
                        type="button"
                        className="min-w-0 flex-1 text-left"
                        onClick={() => {
                          clearMutationErrors();
                          setSelectedSessionId(session.id);
                        }}
                      >
                        <p className="truncate text-sm font-medium text-[var(--text-primary)]">
                          {session.title}
                        </p>
                        <div className="mt-0.5 flex items-center gap-1.5 font-mono text-[11px] text-[var(--text-muted)]">
                          <span className="truncate">{session.model}</span>
                          <span aria-hidden="true">·</span>
                          <span className="shrink-0 tabular-nums">
                            {session.message_count} {isId ? "pesan" : "msgs"}
                          </span>
                        </div>
                      </button>
                      <Button
                        type="button"
                        size="icon-sm"
                        variant="ghost"
                        className="shrink-0 opacity-70 transition-opacity hover:text-[var(--status-error)] group-hover:opacity-100 focus-visible:opacity-100"
                        aria-label={`${isId ? "Hapus" : "Delete"} ${session.title}`}
                        disabled={deleteSession.isPending}
                        onClick={() => {
                          deleteSession.mutate(session.id, {
                            onSuccess: () => {
                              if (selectedSessionId === session.id) {
                                setSelectedSessionId(null);
                              }
                            },
                          });
                        }}
                      >
                        <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                      </Button>
                    </div>
                  );
                })
              ) : (
                <div className="rounded-lg border border-dashed border-[var(--border-strong)] px-3 py-6 text-center">
                  <p className="text-xs font-medium text-[var(--text-secondary)]">
                    {sessionSearch
                      ? isId
                        ? "Tidak ada sesi yang cocok."
                        : "No matching sessions."
                      : isId
                        ? "Belum ada sesi Playground."
                        : "No Playground sessions yet."}
                  </p>
                  <p className="mt-1 text-[11px] text-[var(--text-muted)]">
                    {sessionSearch
                      ? isId
                        ? "Coba kata kunci lain."
                        : "Try a different keyword."
                      : isId
                        ? "Buat sesi baru di atas untuk mulai."
                        : "Create a session above to get started."}
                  </p>
                </div>
              )}
            </div>
          </aside>

          {/* Right main stage: Active session conversation & parameters */}
          <section className="flex min-h-[480px] min-w-0 flex-col bg-[var(--bg-surface)] lg:min-h-0">
            {/* Stage header */}
            <header className="flex flex-wrap items-center justify-between gap-3 border-b border-[var(--border-subtle)] px-4 py-3 sm:px-5">
              <div className="min-w-0 flex-1">
                <div className="flex flex-wrap items-center gap-2">
                  <h3 className="truncate text-sm font-semibold text-[var(--text-primary)] sm:text-base">
                    {activeSession
                      ? activeSession.title
                      : isId
                        ? "Pilih atau buat sesi"
                        : "Select or create a session"}
                  </h3>
                  {activeSession ? (
                    <Badge variant="secondary" className="max-w-[220px] truncate font-mono">
                      {activeSession.model}
                    </Badge>
                  ) : null}
                </div>

                {activeSession ? (
                  <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                    <span className="inline-flex items-center gap-1">
                      <MessageSquare className="h-3 w-3" aria-hidden="true" />
                      {messages.length} {isId ? "pesan" : "msgs"}
                    </span>
                    {telemetry.assistantTurns > 0 ? (
                      <>
                        <span className="inline-flex items-center gap-1">
                          <Clock className="h-3 w-3" aria-hidden="true" />
                          {isId ? "rata-rata" : "avg"} {formatLatencyMs(telemetry.avgLatencyMs)}
                        </span>
                        <span className="inline-flex items-center gap-1">
                          <Cpu className="h-3 w-3" aria-hidden="true" />
                          {formatTokenCount(telemetry.totalTokens)} tok
                        </span>
                      </>
                    ) : null}
                    {activeSession.temperature !== undefined ? (
                      <span>temp {activeSession.temperature}</span>
                    ) : null}
                    {activeSession.max_tokens > 0 ? (
                      <span>max {activeSession.max_tokens}</span>
                    ) : null}
                  </div>
                ) : (
                  <p className="mt-0.5 text-xs text-[var(--text-muted)]">
                    {isId
                      ? "Mulai sesi baru dari panel kiri untuk menguji model."
                      : "Start a session from the left rail to begin testing prompts."}
                  </p>
                )}
              </div>

              {activeSession ? (
                <Button
                  type="button"
                  size="xs"
                  variant={paramsOpen ? "default" : "secondary"}
                  onClick={() => {
                    clearMutationErrors();
                    setParamsOpen((prev) => !prev);
                  }}
                >
                  <SlidersHorizontal className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Parameter" : "Parameters"}</span>
                </Button>
              ) : null}
            </header>

            {/* Collapsible Session Parameters editor */}
            {paramsOpen && activeSession ? (
              <form
                onSubmit={handleUpdateSession}
                className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)]/70 px-4 py-3.5 sm:px-5"
              >
                <div className="mx-auto max-w-3xl space-y-3">
                  <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
                    <div className="space-y-1">
                      <Label htmlFor="playground-edit-title" className="text-xs">
                        {isId ? "Judul Sesi" : "Session Title"}
                      </Label>
                      <Input
                        id="playground-edit-title"
                        value={editTitle}
                        onChange={(e) => setEditTitle(e.target.value)}
                        className="h-8 text-xs"
                      />
                    </div>
                    <div className="space-y-1">
                      <Label htmlFor="playground-edit-model" className="text-xs">
                        Model
                      </Label>
                      <Input
                        id="playground-edit-model"
                        list="playground-model-options"
                        value={editModel}
                        onChange={(e) => setEditModel(e.target.value)}
                        className="h-8 font-mono text-xs"
                      />
                    </div>
                    <div className="space-y-1">
                      <Label htmlFor="playground-edit-temp" className="text-xs">
                        Temperature
                      </Label>
                      <Input
                        id="playground-edit-temp"
                        type="number"
                        step="0.1"
                        min="0"
                        max="2"
                        placeholder="0.7"
                        value={editTemperature}
                        onChange={(e) => setEditTemperature(e.target.value)}
                        className="h-8 font-mono text-xs tabular-nums"
                      />
                    </div>
                    <div className="space-y-1">
                      <Label htmlFor="playground-edit-max-tokens" className="text-xs">
                        Max Tokens
                      </Label>
                      <Input
                        id="playground-edit-max-tokens"
                        type="number"
                        min="0"
                        max="128000"
                        placeholder="auto"
                        value={editMaxTokens}
                        onChange={(e) => setEditMaxTokens(e.target.value)}
                        className="h-8 font-mono text-xs tabular-nums"
                      />
                    </div>
                  </div>

                  <div className="space-y-1">
                    <Label htmlFor="playground-edit-sys" className="text-xs">
                      System Prompt
                    </Label>
                    <textarea
                      id="playground-edit-sys"
                      rows={2}
                      value={editSystemPrompt}
                      onChange={(e) => setEditSystemPrompt(e.target.value)}
                      placeholder={
                        isId
                          ? "Instruksi sistem untuk memandu respons model..."
                          : "System instructions to guide model behavior..."
                      }
                      className="w-full resize-none rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 text-xs text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)]"
                    />
                  </div>

                  <div className="flex items-center justify-end gap-2">
                    {paramsSaved ? (
                      <span className="inline-flex items-center gap-1 font-mono text-xs text-[var(--status-ok)]">
                        <Check className="h-3.5 w-3.5" aria-hidden="true" />
                        {isId ? "Tersimpan" : "Saved"}
                      </span>
                    ) : null}
                    <Button
                      type="button"
                      size="xs"
                      variant="ghost"
                      onClick={() => setParamsOpen(false)}
                    >
                      {isId ? "Tutup" : "Close"}
                    </Button>
                    <Button type="submit" size="xs" disabled={updateSession.isPending}>
                      {updateSession.isPending
                        ? isId
                          ? "Menyimpan..."
                          : "Saving..."
                        : isId
                          ? "Simpan Parameter"
                          : "Save Parameters"}
                    </Button>
                  </div>
                </div>
              </form>
            ) : null}

            {/* Scrollable message stream */}
            <div className="min-h-0 flex-1 overflow-y-auto px-4 py-5 sm:px-6">
              <div className="mx-auto flex w-full max-w-3xl flex-col gap-4">
                {/* Active system prompt callout */}
                {activeSession?.system_prompt ? (
                  <div className="rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3.5 py-2.5">
                    <p className="font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
                      System Prompt
                    </p>
                    <p className="mt-1 whitespace-pre-wrap text-xs leading-relaxed text-[var(--text-secondary)]">
                      {activeSession.system_prompt}
                    </p>
                  </div>
                ) : null}

                {detailQuery.isLoading && activeSessionId ? (
                  <div className="space-y-3 py-4" role="status" aria-busy="true">
                    <Skeleton className="ml-auto h-16 w-2/3 rounded-xl" />
                    <Skeleton className="h-24 w-4/5 rounded-xl" />
                  </div>
                ) : messages.length > 0 ? (
                  <>
                    {messages.map((msg) => {
                      const isUser = msg.role === "user";
                      const isCopied = copiedMessageId === msg.id;
                      return (
                        <article
                          key={msg.id}
                          className={`group rounded-xl border p-4 transition-colors ${
                            isUser
                              ? "border-[var(--border-strong)] bg-[var(--bg-surface-raised)]"
                              : "border-[var(--border-subtle)] bg-[var(--bg-canvas)]"
                          }`}
                        >
                          <div className="mb-2 flex items-center justify-between gap-2">
                            <div className="flex items-center gap-2">
                              <span
                                className={`flex h-6 w-6 items-center justify-center rounded-md border ${
                                  isUser
                                    ? "border-[var(--border-strong)] bg-[var(--bg-surface)] text-[var(--text-primary)]"
                                    : "border-[var(--accent-primary)]/30 bg-[var(--bg-surface)] text-[var(--accent-text)]"
                                }`}
                              >
                                {isUser ? (
                                  <User className="h-3.5 w-3.5" aria-hidden="true" />
                                ) : (
                                  <Sparkles className="h-3.5 w-3.5" aria-hidden="true" />
                                )}
                              </span>
                              <span className="font-mono text-xs font-medium uppercase tracking-wide text-[var(--text-secondary)]">
                                {isUser ? (isId ? "Anda" : "You") : isId ? "Asisten" : "Assistant"}
                              </span>
                            </div>

                            <Button
                              type="button"
                              variant="ghost"
                              size="icon-sm"
                              aria-label={isId ? "Salin pesan" : "Copy message"}
                              onClick={() => handleCopyMessage(msg.id, msg.content)}
                              className="opacity-70 transition-opacity group-hover:opacity-100 focus-visible:opacity-100"
                            >
                              {isCopied ? (
                                <Check
                                  className="h-3.5 w-3.5 text-[var(--status-ok)]"
                                  aria-hidden="true"
                                />
                              ) : (
                                <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                              )}
                            </Button>
                          </div>

                          <div className="whitespace-pre-wrap break-words font-sans text-sm leading-relaxed text-[var(--text-primary)]">
                            {msg.content}
                          </div>

                          {!isUser ? (
                            <div className="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1 border-t border-[var(--border-subtle)] pt-2.5 font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                              <span>{msg.model}</span>
                              <span aria-hidden="true">·</span>
                              <span>{formatLatencyMs(msg.latency_ms)}</span>
                              <span aria-hidden="true">·</span>
                              <span>
                                {msg.prompt_tokens} in / {msg.completion_tokens} out (
                                {msg.total_tokens} tok)
                              </span>
                            </div>
                          ) : null}
                        </article>
                      );
                    })}

                    {sendMessage.isPending ? (
                      <div className="rounded-xl border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
                        <div className="flex items-center gap-2 font-mono text-xs text-[var(--accent-text)]">
                          <Sparkles className="h-3.5 w-3.5 animate-pulse" aria-hidden="true" />
                          <span>{isId ? "Menghasilkan respons..." : "Generating response..."}</span>
                        </div>
                      </div>
                    ) : null}
                  </>
                ) : activeSession ? (
                  <div className="flex flex-col items-center justify-center py-10 text-center">
                    <div className="flex h-11 w-11 items-center justify-center rounded-xl border border-[var(--border-strong)] bg-[var(--bg-canvas)] text-[var(--accent-text)]">
                      <Sparkles className="h-5 w-5" aria-hidden="true" />
                    </div>
                    <h4 className="mt-3.5 text-base font-semibold text-[var(--text-primary)]">
                      {isId ? "Siap menguji prompt Anda" : "Ready to test your prompt"}
                    </h4>
                    <p className="mt-1 max-w-md text-xs leading-relaxed text-[var(--text-muted)]">
                      {isId
                        ? `Sesi terhubung ke ${activeSession.model}. Ketik pesan di bawah atau pilih contoh prompt untuk mulai.`
                        : `Connected to ${activeSession.model}. Type a message below or pick a starter prompt to begin.`}
                    </p>
                    <div className="mt-5 grid w-full max-w-xl gap-2 sm:grid-cols-3">
                      {starterPrompts.map((promptText) => (
                        <button
                          key={promptText}
                          type="button"
                          onClick={() => handleStarterPrompt(promptText)}
                          className="rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 text-left text-xs leading-snug text-[var(--text-secondary)] transition-colors hover:border-[var(--accent-primary)]/60 hover:text-[var(--text-primary)]"
                        >
                          {promptText}
                        </button>
                      ))}
                    </div>
                  </div>
                ) : (
                  <div className="flex flex-col items-center justify-center py-12 text-center">
                    <div className="flex h-11 w-11 items-center justify-center rounded-xl border border-[var(--border-strong)] bg-[var(--bg-canvas)] text-[var(--accent-text)]">
                      <MessageSquare className="h-5 w-5" aria-hidden="true" />
                    </div>
                    <h4 className="mt-3.5 text-base font-semibold text-[var(--text-primary)]">
                      {isId ? "Selamat datang di Playground" : "Welcome to Playground"}
                    </h4>
                    <p className="mt-1 max-w-md text-xs leading-relaxed text-[var(--text-muted)]">
                      {isId
                        ? "Buat sesi untuk menguji routing model, system prompt, serta memantau latensi dan token secara langsung."
                        : "Create a session to test model routing, tune system prompts, and inspect live latency and token usage."}
                    </p>
                    <Button
                      type="button"
                      size="sm"
                      className="mt-4"
                      disabled={createSession.isPending}
                      onClick={handleQuickStartSession}
                    >
                      <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>
                        {createSession.isPending
                          ? isId
                            ? "Membuat sesi..."
                            : "Creating session..."
                          : models[0]
                            ? isId
                              ? `Mulai cepat dengan ${models[0].id}`
                              : `Quick start with ${models[0].id}`
                            : isId
                              ? "Buat Sesi Pertama"
                              : "Create First Session"}
                      </span>
                    </Button>
                  </div>
                )}

                <div ref={messagesEndRef} />
              </div>
            </div>

            {/* Unified modern prompt composer */}
            <footer className="border-t border-[var(--border-subtle)] bg-[var(--bg-surface)] px-4 py-3.5 sm:px-6">
              <form onSubmit={handleSendMessage} className="mx-auto max-w-3xl">
                <div className="rounded-xl border border-[var(--border-strong)] bg-[var(--bg-input)] p-2.5 transition-colors focus-within:border-[var(--accent-primary)]">
                  <textarea
                    ref={composerRef}
                    rows={2}
                    aria-label={isId ? "Prompt pesan Playground" : "Playground message prompt"}
                    placeholder={
                      activeSession
                        ? isId
                          ? "Tulis prompt Anda di sini..."
                          : "Write your prompt here..."
                        : isId
                          ? "Buat atau pilih sesi terlebih dahulu..."
                          : "Create or select a session first..."
                    }
                    disabled={!activeSession || sendMessage.isPending}
                    value={draft}
                    onChange={(e) => setDraft(e.target.value)}
                    onKeyDown={handleComposerKeyDown}
                    className="max-h-40 min-h-[52px] w-full resize-y bg-transparent px-1.5 py-1 text-sm leading-relaxed text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:outline-none disabled:cursor-not-allowed disabled:opacity-60"
                  />
                  <div className="mt-2 flex flex-wrap items-center justify-between gap-2 border-t border-[var(--border-subtle)] pt-2">
                    <p className="font-mono text-[11px] text-[var(--text-muted)]">
                      {isId
                        ? "Enter kirim · Shift+Enter baris baru"
                        : "Enter to send · Shift+Enter for newline"}
                    </p>
                    <Button
                      type="submit"
                      size="sm"
                      disabled={!activeSession || !draft.trim() || sendMessage.isPending}
                    >
                      <Send className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>
                        {sendMessage.isPending
                          ? isId
                            ? "Mengirim..."
                            : "Sending..."
                          : isId
                            ? "Kirim"
                            : "Send"}
                      </span>
                    </Button>
                  </div>
                </div>
              </form>
            </footer>
          </section>
        </div>
      </div>
    </div>
  );
}
