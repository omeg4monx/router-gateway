import { MessageSquare, Plus, Send, Trash2 } from "lucide-react";
import { useState } from "react";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useArenaModels,
  useChatSessionDetail,
  useChatSessions,
  useCreateChatSessionMutation,
  useDeleteChatSessionMutation,
  useSendChatMessageMutation,
} from "./hooks";
import { chatMessageFormSchema, createChatSessionFormSchema } from "./schemas";

export function ChatConsole() {
  const modelsQuery = useArenaModels();
  const sessionsQuery = useChatSessions();

  const createSession = useCreateChatSessionMutation();
  const deleteSession = useDeleteChatSessionMutation();
  const sendMessage = useSendChatMessageMutation();

  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null);
  const [newModel, setNewModel] = useState("");
  const [newSystemPrompt, setNewSystemPrompt] = useState("");
  const [draft, setDraft] = useState("");
  const [validationError, setValidationError] = useState<string | null>(null);

  const models = modelsQuery.data?.models ?? [];
  const sessions = sessionsQuery.data?.sessions ?? [];
  const activeSessionId = selectedSessionId ?? sessions[0]?.id ?? null;

  const detailQuery = useChatSessionDetail(activeSessionId);
  const activeSession = detailQuery.data?.session;
  const messages = detailQuery.data?.messages ?? [];

  const clearMutationErrors = () => {
    setValidationError(null);
    createSession.reset();
    deleteSession.reset();
    sendMessage.reset();
  };

  const handleCreateSession = (e: React.FormEvent) => {
    e.preventDefault();
    clearMutationErrors();
    const chosenModel = newModel.trim() || models[0]?.id || "";
    const parsed = createChatSessionFormSchema.safeParse({
      model: chosenModel,
      system_prompt: newSystemPrompt,
    });
    if (!parsed.success) {
      setValidationError(parsed.error.issues[0]?.message ?? "Model is required");
      return;
    }
    createSession.mutate(
      {
        model: parsed.data.model,
        system_prompt: parsed.data.system_prompt,
      },
      {
        onSuccess: (created) => {
          setSelectedSessionId(created.id);
          setNewSystemPrompt("");
        },
      },
    );
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeSessionId) {
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

  const activeError =
    modelsQuery.error ??
    sessionsQuery.error ??
    detailQuery.error ??
    createSession.error ??
    deleteSession.error ??
    sendMessage.error;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Interactive Console
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Basic Chat
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Persistent multi-session chat backed directly by the internal gateway pipeline with
            per-turn latency and token telemetry.
          </p>
        </div>

        <Badge variant="secondary">{sessions.length} SESSIONS</Badge>
      </section>

      <EmberRule />

      {activeError || validationError ? (
        <Alert variant="destructive">
          <AlertTitle>Chat error</AlertTitle>
          <AlertDescription>{validationError ?? activeError?.message}</AlertDescription>
        </Alert>
      ) : null}

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <MessageSquare className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              Sessions
            </CardTitle>
            <CardDescription>Create or switch persistent chat threads.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <form onSubmit={handleCreateSession} className="space-y-2.5">
              <div className="space-y-1">
                <Label htmlFor="chat-model">Model</Label>
                <Input
                  id="chat-model"
                  list="chat-model-options"
                  placeholder={models[0]?.id ?? "openai/gpt-4o"}
                  value={newModel}
                  onChange={(e) => setNewModel(e.target.value)}
                  className="font-mono text-xs"
                />
                <datalist id="chat-model-options">
                  {models.map((m) => (
                    <option key={m.id} value={m.id} />
                  ))}
                </datalist>
              </div>
              <div className="space-y-1">
                <Label htmlFor="chat-sys">System Prompt (optional)</Label>
                <Input
                  id="chat-sys"
                  placeholder="You are a concise assistant."
                  value={newSystemPrompt}
                  onChange={(e) => setNewSystemPrompt(e.target.value)}
                />
              </div>
              <Button type="submit" size="sm" className="w-full" disabled={createSession.isPending}>
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{createSession.isPending ? "Creating..." : "New Chat Session"}</span>
              </Button>
            </form>

            <div className="space-y-1.5 border-t border-[var(--border-subtle)] pt-3">
              {sessions.length > 0 ? (
                sessions.map((s) => (
                  <div
                    key={s.id}
                    className={`flex items-center justify-between gap-2 rounded-md border px-2.5 py-2 ${
                      activeSessionId === s.id
                        ? "border-[var(--accent-primary)] bg-[var(--bg-surface-raised)]"
                        : "border-[var(--border-subtle)] bg-[var(--bg-canvas)]"
                    }`}
                  >
                    <button
                      type="button"
                      className="min-w-0 flex-1 text-left"
                      onClick={() => {
                        clearMutationErrors();
                        setSelectedSessionId(s.id);
                      }}
                    >
                      <p className="truncate text-sm font-medium text-[var(--text-primary)]">
                        {s.title}
                      </p>
                      <p className="truncate font-mono text-[11px] text-[var(--text-muted)]">
                        {s.model} · {s.message_count} msgs
                      </p>
                    </button>
                    <Button
                      type="button"
                      size="icon"
                      variant="secondary"
                      className="h-7 w-7 shrink-0"
                      aria-label={`Delete ${s.title}`}
                      onClick={() => {
                        deleteSession.mutate(s.id, {
                          onSuccess: () => {
                            if (selectedSessionId === s.id) {
                              setSelectedSessionId(null);
                            }
                          },
                        });
                      }}
                    >
                      <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                    </Button>
                  </div>
                ))
              ) : (
                <p className="text-xs text-[var(--text-muted)]">No chat sessions yet.</p>
              )}
            </div>
          </CardContent>
        </Card>

        <Card className="flex flex-col lg:col-span-2">
          <CardHeader className="border-b border-[var(--border-subtle)]">
            <div className="flex items-center justify-between gap-2">
              <div>
                <CardTitle className="text-base">
                  {activeSession ? activeSession.title : "Select or create a session"}
                </CardTitle>
                <CardDescription>
                  {activeSession
                    ? `Model: ${activeSession.model}`
                    : "Start a new session from the left panel to begin chatting."}
                </CardDescription>
              </div>
              {activeSession ? <Badge variant="secondary">{activeSession.model}</Badge> : null}
            </div>
          </CardHeader>

          <CardContent className="flex flex-1 flex-col justify-between gap-4 pt-4">
            <div className="max-h-[460px] space-y-3 overflow-y-auto pr-1">
              {messages.length > 0 ? (
                messages.map((msg) => (
                  <div
                    key={msg.id}
                    className={`rounded-md border p-3 ${
                      msg.role === "user"
                        ? "border-[var(--accent-primary)] bg-[var(--bg-surface-raised)]"
                        : "border-[var(--border-subtle)] bg-[var(--bg-canvas)]"
                    }`}
                  >
                    <div className="mb-1 flex items-center justify-between font-mono text-[11px] text-[var(--text-muted)]">
                      <span className="uppercase">{msg.role}</span>
                      {msg.role === "assistant" ? (
                        <span>
                          {msg.latency_ms}ms · {msg.total_tokens} tok
                        </span>
                      ) : null}
                    </div>
                    <pre className="whitespace-pre-wrap font-sans text-sm text-[var(--text-primary)]">
                      {msg.content}
                    </pre>
                  </div>
                ))
              ) : (
                <p className="py-8 text-center text-sm text-[var(--text-muted)]">
                  {activeSession
                    ? "Send a message below to start this conversation."
                    : "No active session selected."}
                </p>
              )}
            </div>

            <form
              onSubmit={handleSendMessage}
              className="flex gap-2 border-t border-[var(--border-subtle)] pt-3"
            >
              <Input
                placeholder={activeSession ? "Type your message..." : "Create a session first..."}
                disabled={!activeSession || sendMessage.isPending}
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
              />
              <Button type="submit" disabled={!activeSession || sendMessage.isPending}>
                <Send className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{sendMessage.isPending ? "Sending..." : "Send"}</span>
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
