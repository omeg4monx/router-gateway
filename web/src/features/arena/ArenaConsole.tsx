import { Award, CheckCircle2, Flame, Play, Trash2, Trophy } from "lucide-react";
import { useState } from "react";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useArenaBattles,
  useArenaLeaderboard,
  useArenaModels,
  useDeleteArenaBattleMutation,
  useRunArenaBattleMutation,
  useVoteArenaBattleMutation,
} from "./hooks";
import { type ArenaBattle, arenaBattleFormSchema } from "./schemas";

export function ArenaConsole() {
  const modelsQuery = useArenaModels();
  const battlesQuery = useArenaBattles();
  const leaderboardQuery = useArenaLeaderboard();

  const runBattle = useRunArenaBattleMutation();
  const voteBattle = useVoteArenaBattleMutation();
  const deleteBattle = useDeleteArenaBattleMutation();

  const [selectedModels, setSelectedModels] = useState<string[]>([]);
  const [customModel, setCustomModel] = useState("");
  const [prompt, setPrompt] = useState("");
  const [systemPrompt, setSystemPrompt] = useState("");
  const [temperature, setTemperature] = useState("0.7");
  const [maxTokens, setMaxTokens] = useState("1024");
  const [validationError, setValidationError] = useState<string | null>(null);
  const [activeBattleId, setActiveBattleId] = useState<string | null>(null);

  const availableModels = modelsQuery.data?.models ?? [];
  const battles = battlesQuery.data?.battles ?? [];
  const leaderboard = leaderboardQuery.data?.entries ?? [];

  const activeBattle: ArenaBattle | undefined =
    battles.find((b) => b.id === activeBattleId) ?? battles[0];

  const toggleModel = (modelId: string) => {
    setValidationError(null);
    setSelectedModels((prev) => {
      if (prev.includes(modelId)) {
        return prev.filter((m) => m !== modelId);
      }
      if (prev.length >= 4) {
        return prev;
      }
      return [...prev, modelId];
    });
  };

  const addCustomModel = () => {
    const trimmed = customModel.trim();
    if (!trimmed) {
      return;
    }
    if (!selectedModels.includes(trimmed) && selectedModels.length < 4) {
      setSelectedModels((prev) => [...prev, trimmed]);
    }
    setCustomModel("");
  };

  const clearMutationErrors = () => {
    setValidationError(null);
    runBattle.reset();
    voteBattle.reset();
    deleteBattle.reset();
  };

  const handleRunBattle = (e: React.FormEvent) => {
    e.preventDefault();
    clearMutationErrors();

    const parsed = arenaBattleFormSchema.safeParse({
      prompt,
      system_prompt: systemPrompt,
      models: selectedModels,
      temperature: temperature.trim() === "" ? undefined : Number(temperature),
      max_tokens: maxTokens.trim() === "" ? undefined : Number(maxTokens),
    });
    if (!parsed.success) {
      setValidationError(parsed.error.issues[0]?.message ?? "Invalid battle configuration");
      return;
    }

    runBattle.mutate(
      {
        prompt: parsed.data.prompt,
        system_prompt: parsed.data.system_prompt,
        models: parsed.data.models,
        temperature: parsed.data.temperature,
        max_tokens: parsed.data.max_tokens,
      },
      {
        onSuccess: (created) => {
          setActiveBattleId(created.id);
        },
      },
    );
  };

  const activeError =
    modelsQuery.error ??
    battlesQuery.error ??
    leaderboardQuery.error ??
    runBattle.error ??
    voteBattle.error ??
    deleteBattle.error;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Evaluation Workbench
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Model Arena
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Execute one prompt concurrently across 2 to 4 gateway models, compare latency and token
            usage side-by-side, and score winners onto the persistent leaderboard.
          </p>
        </div>

        <Badge variant="secondary">{battles.length} BATTLES RECORDED</Badge>
      </section>

      <EmberRule />

      {activeError || validationError ? (
        <Alert variant="destructive">
          <AlertTitle>Arena error</AlertTitle>
          <AlertDescription>{validationError ?? activeError?.message}</AlertDescription>
        </Alert>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Flame className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            Configure Contenders &amp; Prompt
          </CardTitle>
          <CardDescription>
            Select 2 to 4 models ({selectedModels.length}/4 selected) from your configured gateway
            inventory or enter explicit model IDs.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleRunBattle} className="space-y-4">
            <div className="space-y-2">
              <Label>Contender Models</Label>
              <div className="flex flex-wrap gap-2">
                {availableModels.map((m) => {
                  const active = selectedModels.includes(m.id);
                  return (
                    <Button
                      key={m.id}
                      type="button"
                      size="sm"
                      variant={active ? "default" : "outline"}
                      onClick={() => toggleModel(m.id)}
                    >
                      <span className="font-mono text-xs">{m.id}</span>
                    </Button>
                  );
                })}
                {availableModels.length === 0 ? (
                  <p className="text-xs text-[var(--text-muted)]">
                    No active provider models discovered yet — add custom model IDs below.
                  </p>
                ) : null}
              </div>

              <div className="flex max-w-md gap-2 pt-1">
                <Input
                  placeholder="Custom model ID (e.g. openai/gpt-4o)"
                  value={customModel}
                  onChange={(e) => setCustomModel(e.target.value)}
                  className="font-mono text-xs"
                />
                <Button type="button" size="sm" variant="secondary" onClick={addCustomModel}>
                  Add
                </Button>
              </div>

              {selectedModels.length > 0 ? (
                <div className="flex flex-wrap gap-1.5 pt-1">
                  {selectedModels.map((id) => (
                    <Badge key={id} variant="default" className="gap-1.5 font-mono text-xs">
                      <span>{id}</span>
                      <button
                        type="button"
                        onClick={() => toggleModel(id)}
                        className="ml-1 text-xs underline"
                        aria-label={`Remove ${id}`}
                      >
                        x
                      </button>
                    </Badge>
                  ))}
                </div>
              ) : null}
            </div>

            <div className="grid gap-4 sm:grid-cols-3">
              <div className="space-y-1.5 sm:col-span-1">
                <Label htmlFor="arena-sys">System Prompt (optional)</Label>
                <Input
                  id="arena-sys"
                  placeholder="Be concise and factual."
                  value={systemPrompt}
                  onChange={(e) => setSystemPrompt(e.target.value)}
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="arena-temp">Temperature</Label>
                <Input
                  id="arena-temp"
                  type="number"
                  step="0.1"
                  min="0"
                  max="2"
                  value={temperature}
                  onChange={(e) => setTemperature(e.target.value)}
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="arena-max-tokens">Max Tokens</Label>
                <Input
                  id="arena-max-tokens"
                  type="number"
                  min="0"
                  max="128000"
                  value={maxTokens}
                  onChange={(e) => setMaxTokens(e.target.value)}
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="arena-prompt">Battle Prompt</Label>
              <textarea
                id="arena-prompt"
                rows={3}
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="Enter a prompt to run concurrently across all selected contenders..."
                className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-2 text-sm text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:border-[var(--accent-primary)] focus:outline-none"
              />
            </div>

            <Button type="submit" disabled={runBattle.isPending}>
              <Play className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{runBattle.isPending ? "Running Battle..." : "Run Battle"}</span>
            </Button>
          </form>
        </CardContent>
      </Card>

      {activeBattle ? (
        <Card>
          <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <CardTitle className="text-base">Battle Comparison</CardTitle>
              <CardDescription className="line-clamp-2">{activeBattle.prompt}</CardDescription>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              {activeBattle.winner_model ? (
                <Badge variant="default">
                  <Award className="mr-1 h-3 w-3" aria-hidden="true" />
                  Winner: {activeBattle.winner_model}
                </Badge>
              ) : null}
              <Button
                type="button"
                size="sm"
                variant="outline"
                disabled={voteBattle.isPending}
                onClick={() => voteBattle.mutate({ id: activeBattle.id, winner_model: "tie" })}
              >
                Vote Tie
              </Button>
              <Button
                type="button"
                size="sm"
                variant="secondary"
                disabled={deleteBattle.isPending}
                onClick={() => deleteBattle.mutate(activeBattle.id)}
              >
                <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                <span>Delete</span>
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2">
              {activeBattle.contenders.map((c) => {
                const isWinner = activeBattle.winner_model === c.model;
                return (
                  <div
                    key={c.id}
                    className={`flex flex-col justify-between rounded-md border p-4 ${
                      isWinner
                        ? "border-[var(--accent-primary)] bg-[var(--bg-surface-raised)]"
                        : "border-[var(--border-subtle)] bg-[var(--bg-canvas)]"
                    }`}
                  >
                    <div className="space-y-3">
                      <div className="flex items-center justify-between gap-2">
                        <span className="font-mono text-sm font-semibold text-[var(--text-primary)]">
                          {c.model}
                        </span>
                        <div className="flex items-center gap-1.5">
                          <Badge variant={c.status === "ok" ? "secondary" : "destructive"}>
                            {c.status}
                          </Badge>
                          <span className="font-mono text-xs tabular-nums text-[var(--text-muted)]">
                            {c.latency_ms}ms
                          </span>
                        </div>
                      </div>

                      {c.status === "ok" ? (
                        <pre className="max-h-72 overflow-y-auto whitespace-pre-wrap font-sans text-sm text-[var(--text-primary)]">
                          {c.content || "(empty response)"}
                        </pre>
                      ) : (
                        <p className="font-mono text-xs text-[var(--status-error)]">{c.error}</p>
                      )}
                    </div>

                    <div className="mt-4 flex items-center justify-between border-t border-[var(--border-subtle)] pt-3">
                      <span className="font-mono text-xs tabular-nums text-[var(--text-muted)]">
                        {c.total_tokens} tok ({c.prompt_tokens} in / {c.completion_tokens} out)
                      </span>
                      <Button
                        type="button"
                        size="sm"
                        variant={isWinner ? "default" : "outline"}
                        disabled={voteBattle.isPending}
                        onClick={() => {
                          clearMutationErrors();
                          voteBattle.mutate({ id: activeBattle.id, winner_model: c.model });
                        }}
                      >
                        <CheckCircle2 className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>{isWinner ? "Winner" : "Vote Winner"}</span>
                      </Button>
                    </div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>
      ) : null}

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Trophy className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              Model Leaderboard
            </CardTitle>
            <CardDescription>
              Aggregated win rate and average latency across battles.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {leaderboard.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-[var(--border-subtle)] font-mono text-xs text-[var(--text-muted)]">
                    <tr>
                      <th className="py-2 pr-3">Model</th>
                      <th className="py-2 px-2 text-right">Win %</th>
                      <th className="py-2 px-2 text-right">W/T/L</th>
                      <th className="py-2 pl-2 text-right">Avg ms</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-[var(--border-subtle)]">
                    {leaderboard.map((row) => (
                      <tr key={row.model}>
                        <td className="py-2 pr-3 font-mono text-xs text-[var(--text-primary)]">
                          {row.model}
                        </td>
                        <td className="py-2 px-2 text-right font-mono text-xs tabular-nums">
                          {row.win_rate.toFixed(1)}%
                        </td>
                        <td className="py-2 px-2 text-right font-mono text-xs tabular-nums text-[var(--text-secondary)]">
                          {row.wins}/{row.ties}/{row.losses}
                        </td>
                        <td className="py-2 pl-2 text-right font-mono text-xs tabular-nums text-[var(--text-muted)]">
                          {row.avg_latency_ms}ms
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <p className="text-sm text-[var(--text-muted)]">
                Run and vote on a battle to populate the leaderboard.
              </p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Recent Battles</CardTitle>
            <CardDescription>Select a past comparison to inspect or re-score.</CardDescription>
          </CardHeader>
          <CardContent>
            {battles.length > 0 ? (
              <div className="space-y-2">
                {battles.map((b) => (
                  <button
                    key={b.id}
                    type="button"
                    onClick={() => setActiveBattleId(b.id)}
                    className={`w-full rounded-md border px-3 py-2.5 text-left transition-colors ${
                      activeBattle?.id === b.id
                        ? "border-[var(--accent-primary)] bg-[var(--bg-surface-raised)]"
                        : "border-[var(--border-subtle)] bg-[var(--bg-canvas)] hover:bg-[var(--bg-surface-raised)]"
                    }`}
                  >
                    <p className="line-clamp-1 text-sm font-medium text-[var(--text-primary)]">
                      {b.prompt}
                    </p>
                    <div className="mt-1 flex items-center justify-between font-mono text-xs text-[var(--text-muted)]">
                      <span>{b.contenders.map((c) => c.model).join(" vs ")}</span>
                      <span>{b.winner_model ? `winner: ${b.winner_model}` : "unvoted"}</span>
                    </div>
                  </button>
                ))}
              </div>
            ) : (
              <p className="text-sm text-[var(--text-muted)]">No battles yet.</p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
