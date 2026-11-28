import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryClient } from "@tanstack/react-query";
import {
  Activity,
  AlertCircle,
  AlertTriangle,
  ArrowDownAZ,
  ArrowDownWideNarrow,
  ArrowLeft,
  ArrowUpNarrowWide,
  ArrowUpZA,
  Bot,
  Brain,
  Check,
  CheckCircle2,
  ChevronDown,
  ChevronUp,
  Copy,
  Download,
  ExternalLink,
  Eye,
  EyeOff,
  Info,
  KeyRound,
  Lock,
  Plus,
  RefreshCw,
  Save,
  ScanEye,
  Search,
  Server,
  Settings2,
  Trash2,
  X,
} from "lucide-react";
import { lazy, Suspense, useDeferredValue, useEffect, useMemo, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate, useSearchParams } from "react-router";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  type ModelIdentityResult,
  runModelIdentityProbe,
  runSingleModelTest,
} from "../models/hooks";
import type { CatalogTestResult } from "../models/types";
import {
  buildProviderEntries,
  canAddProviderConnection,
  canProbeProviderModels,
  type ProviderEntry,
} from "./catalogUtils";
import { removeModelFromList } from "./customModelUtils";
import {
  type FetchConnectionModelsResponse,
  invalidateModelRoutingQueries,
  providersQueryKeys,
  runWithConcurrency,
  testProviderConnection,
  updateProviderConnection,
  useDeleteProviderConnectionMutation,
  useDeleteProviderNodeMutation,
  useFetchConnectionModelsMutation,
  useProviderConnection,
  useProviderConnections,
  useProvidersCatalog,
  useTestProviderConnectionMutation,
  useUpdateProviderConnectionMutation,
} from "./hooks";
import {
  getModelIntelligenceProfile,
  MODEL_TIER_BADGE,
  type ModelSortDirection,
  type ModelSortField,
  sortProviderModels,
} from "./modelIntelligence";
import { ProviderLogo } from "./ProviderLogo";
import { ProviderStrategyCard } from "./ProviderStrategyCard";
import {
  type CatalogAdapterDescriptor,
  connectionToEditFormValues,
  DEFAULT_UPDATE_CONNECTION_VALUES,
  type ProviderConnectionSummary,
  type ProviderNode,
  type TestProviderConnectionResponse,
  toUpdateConnectionRequest,
  type UpdateConnectionFormInput,
  type UpdateConnectionFormOutput,
  updateConnectionSchema,
} from "./schemas";

const AddApiKeyModal = lazy(() =>
  import("./ProvidersModals").then((m) => ({ default: m.AddApiKeyModal })),
);
const AddCustomModelModal = lazy(() =>
  import("./AddCustomModelModal").then((m) => ({ default: m.AddCustomModelModal })),
);
const EditConnectionModal = lazy(() =>
  import("./ProvidersModals").then((m) => ({ default: m.EditConnectionModal })),
);
const OAuthWizardModal = lazy(() =>
  import("./ProvidersModals").then((m) => ({ default: m.OAuthWizardModal })),
);

const EMPTY_ADAPTERS: CatalogAdapterDescriptor[] = [];
const EMPTY_CUSTOM_NODES: ProviderNode[] = [];
const EMPTY_CONNECTIONS: ProviderConnectionSummary[] = [];

const TEXTAREA_CLASS_NAME =
  "flex min-h-[92px] w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-xs text-[var(--text-primary)] placeholder:font-sans placeholder:text-[var(--text-muted)] transition-colors focus:border-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)]";

type ConnectionStatusFilter = "all" | "issues" | "active" | "disabled";

const KNOWN_FREE_OPENCODE_MODELS = new Set(["big-pickle", "union-alpha"]);

function isModelFreeByRule(id: string, providerId: string, explicit?: boolean): boolean {
  if (explicit) {
    return true;
  }
  const lower = id.trim().toLowerCase();
  if (!lower) {
    return false;
  }
  if (
    lower.includes(":free") ||
    lower.includes("-free") ||
    lower.includes("/free") ||
    lower.endsWith(" free") ||
    lower.startsWith("cline-free/")
  ) {
    return true;
  }
  if (
    (providerId === "opencode" || providerId === "opencode-zen") &&
    KNOWN_FREE_OPENCODE_MODELS.has(lower)
  ) {
    return true;
  }
  return false;
}

function isConnectionInError(
  conn: ProviderConnectionSummary,
  live: TestProviderConnectionResponse | undefined,
): boolean {
  if (live) {
    return !live.ok;
  }
  return conn.last_test_status === "error" || Boolean(conn.last_test_error?.trim());
}

function formatTimestamp(iso: string | null): string {
  if (!iso) {
    return "Never";
  }
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return iso;
  }
  return date.toLocaleString([], {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function ProviderStudioDetail({ providerId }: { providerId: string }) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";

  const catalogQuery = useProvidersCatalog();
  const connectionsQuery = useProviderConnections();
  const updateConnMutation = useUpdateProviderConnectionMutation();
  const deleteConnMutation = useDeleteProviderConnectionMutation();
  const testConnMutation = useTestProviderConnectionMutation();
  const fetchModelsMutation = useFetchConnectionModelsMutation();
  const deleteNodeMutation = useDeleteProviderNodeMutation();

  const [showAddApiKeyModal, setShowAddApiKeyModal] = useState(false);
  const [showOAuthModal, setShowOAuthModal] = useState(false);
  const [showAddModelModal, setShowAddModelModal] = useState(false);
  const [removingModelId, setRemovingModelId] = useState<string | null>(null);
  const [editingConnection, setEditingConnection] = useState<ProviderConnectionSummary | null>(
    null,
  );
  const [modalsActivated, setModalsActivated] = useState(false);
  const lastEditingConnectionRef = useRef<ProviderConnectionSummary | null>(null);
  if (editingConnection) {
    lastEditingConnectionRef.current = editingConnection;
  }
  useEffect(() => {
    if (
      !modalsActivated &&
      (showAddApiKeyModal || showOAuthModal || showAddModelModal || editingConnection !== null)
    ) {
      setModalsActivated(true);
    }
  }, [editingConnection, modalsActivated, showAddApiKeyModal, showAddModelModal, showOAuthModal]);
  const [feedbackBanner, setFeedbackBanner] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [busyConnectionId, setBusyConnectionId] = useState<string | null>(null);
  const [testingAll, setTestingAll] = useState(false);
  const [testAllProgress, setTestAllProgress] = useState<{
    total: number;
    completed: number;
    passed: number;
    failed: number;
    currentConnectionId: string | null;
    currentConnectionName: string | null;
  } | null>(null);
  const [liveTestResults, setLiveTestResults] = useState<
    Record<string, TestProviderConnectionResponse>
  >({});
  const [loadedModelsResult, setLoadedModelsResult] =
    useState<FetchConnectionModelsResponse | null>(null);
  const autoLoadedConnRef = useRef<string | null>(null);
  const fetchModelsMutate = fetchModelsMutation.mutate;
  const [confirmDeleteConnId, setConfirmDeleteConnId] = useState<string | null>(null);
  const [confirmDeleteNode, setConfirmDeleteNode] = useState(false);
  const [connectionFilter, setConnectionFilter] = useState<ConnectionStatusFilter>("all");
  const [connectionSearch, setConnectionSearch] = useState("");
  const [modelSearch, setModelSearch] = useState("");
  const [freeModelsOnly, setFreeModelsOnly] = useState(false);
  const [modelSortField, setModelSortField] = useState<ModelSortField>("name");
  const [modelSortDirection, setModelSortDirection] = useState<ModelSortDirection>("asc");
  const [copiedModelId, setCopiedModelId] = useState<string | null>(null);
  const copyTimerRef = useRef<number | null>(null);
  const [testingModelIds, setTestingModelIds] = useState<Record<string, boolean>>({});
  const [modelTestResults, setModelTestResults] = useState<Record<string, CatalogTestResult>>({});
  const modelTestControllersRef = useRef<Map<string, AbortController>>(new Map());
  const [identifyingModelIds, setIdentifyingModelIds] = useState<Record<string, boolean>>({});
  const [modelIdentityResults, setModelIdentityResults] = useState<
    Record<string, ModelIdentityResult>
  >({});
  const [openIdentityModelIds, setOpenIdentityModelIds] = useState<Record<string, boolean>>({});
  const modelIdentityControllersRef = useRef<Map<string, AbortController>>(new Map());

  const adapters = catalogQuery.data?.adapters ?? EMPTY_ADAPTERS;
  const customNodes = catalogQuery.data?.custom_nodes ?? EMPTY_CUSTOM_NODES;
  const allConnections = connectionsQuery.data ?? EMPTY_CONNECTIONS;
  const deferredConnectionSearch = useDeferredValue(connectionSearch);
  const deferredModelSearch = useDeferredValue(modelSearch);

  const allEntries = useMemo(
    () => buildProviderEntries(adapters, customNodes, allConnections),
    [adapters, customNodes, allConnections],
  );

  const entry: ProviderEntry | null = useMemo(
    () => allEntries.find((item) => item.id === providerId) ?? null,
    [allEntries, providerId],
  );

  const sortedConnections = useMemo(() => {
    if (!entry) return [];
    return [...entry.connections].sort((a, b) => a.priority - b.priority);
  }, [entry]);

  const connectionCounts = useMemo(() => {
    let issues = 0;
    let active = 0;
    let disabled = 0;
    for (const conn of sortedConnections) {
      const hasError = isConnectionInError(conn, liveTestResults[conn.id]);
      if (hasError) {
        issues += 1;
      }
      if (!conn.enabled) {
        disabled += 1;
      } else if (!hasError) {
        active += 1;
      }
    }
    return {
      all: sortedConnections.length,
      issues,
      active,
      disabled,
    };
  }, [sortedConnections, liveTestResults]);

  const isConnectionFiltered = connectionFilter !== "all" || connectionSearch.trim().length > 0;

  const filteredConnections = useMemo(() => {
    const q = deferredConnectionSearch.trim().toLowerCase();
    return sortedConnections.filter((conn) => {
      const hasError = isConnectionInError(conn, liveTestResults[conn.id]);
      if (connectionFilter === "issues" && !hasError) {
        return false;
      }
      if (connectionFilter === "active" && (!conn.enabled || hasError)) {
        return false;
      }
      if (connectionFilter === "disabled" && conn.enabled) {
        return false;
      }
      if (!q) {
        return true;
      }
      return (
        conn.name.toLowerCase().includes(q) ||
        conn.auth_type.toLowerCase().includes(q) ||
        (conn.last_test_error?.toLowerCase().includes(q) ?? false)
      );
    });
  }, [sortedConnections, liveTestResults, connectionFilter, deferredConnectionSearch]);

  const activeConnectionForLoad = sortedConnections.find((c) => c.enabled) ?? sortedConnections[0];

  const allModels = useMemo(() => {
    if (!entry) return [];
    const seen = new Set<string>();
    const catalogIds = new Set([
      ...entry.models.map((m) => m.id.trim()),
      ...entry.defaultModels.map((m) => m.trim()),
    ]);
    const liveIds = new Set(
      loadedModelsResult && loadedModelsResult.source === "upstream"
        ? loadedModelsResult.models.map((m) => m.id.trim())
        : [],
    );
    const savedConnIds = new Set<string>();
    for (const conn of entry.connections) {
      for (const m of conn.models) {
        const trimmed = m.trim();
        if (trimmed) {
          savedConnIds.add(trimmed);
        }
      }
    }

    const list: Array<{
      id: string;
      name?: string;
      isFree?: boolean;
      isLive?: boolean;
      isCustom?: boolean;
      isSaved?: boolean;
    }> = [];

    const hasLiveUpstream =
      loadedModelsResult !== null &&
      loadedModelsResult.source === "upstream" &&
      loadedModelsResult.models.length > 0;
    const catalogById = new Map(entry.models.map((m) => [m.id.trim(), m]));

    if (hasLiveUpstream) {
      for (const m of loadedModelsResult.models) {
        const trimmed = m.id.trim();
        if (trimmed && !seen.has(trimmed)) {
          seen.add(trimmed);
          const catalogEntry = catalogById.get(trimmed);
          list.push({
            id: trimmed,
            name: m.name || catalogEntry?.name,
            isFree:
              typeof m.is_free === "boolean"
                ? m.is_free
                : isModelFreeByRule(trimmed, entry.id, catalogEntry?.isFree),
            isLive: !catalogIds.has(trimmed),
            isSaved: savedConnIds.has(trimmed),
          });
        }
      }
    }
    for (const m of entry.models) {
      const trimmed = m.id.trim();
      if (trimmed && !seen.has(trimmed)) {
        seen.add(trimmed);
        list.push({
          id: trimmed,
          name: m.name,
          isFree: hasLiveUpstream ? false : isModelFreeByRule(trimmed, entry.id, m.isFree),
          isSaved: savedConnIds.has(trimmed),
        });
      }
    }
    for (const m of entry.defaultModels) {
      const trimmed = m.trim();
      if (trimmed && !seen.has(trimmed)) {
        seen.add(trimmed);
        list.push({
          id: trimmed,
          isFree: hasLiveUpstream ? false : isModelFreeByRule(trimmed, entry.id, false),
          isSaved: savedConnIds.has(trimmed),
        });
      }
    }
    for (const conn of entry.connections) {
      for (const m of conn.models) {
        const trimmed = m.trim();
        if (trimmed && !seen.has(trimmed)) {
          seen.add(trimmed);
          const canClassifyCustom =
            !fetchModelsMutation.isPending &&
            (loadedModelsResult !== null || !activeConnectionForLoad?.enabled);
          list.push({
            id: trimmed,
            isCustom: canClassifyCustom && !catalogIds.has(trimmed) && !liveIds.has(trimmed),
            isSaved: true,
          });
        }
      }
    }
    return list;
  }, [entry, loadedModelsResult, fetchModelsMutation.isPending, activeConnectionForLoad?.enabled]);

  const freeModelsCount = useMemo(
    () => allModels.filter((m) => Boolean(m.isFree)).length,
    [allModels],
  );

  const filteredModels = useMemo(() => {
    const q = deferredModelSearch.trim().toLowerCase();
    const matched = allModels.filter((m) => {
      if (freeModelsOnly && !m.isFree) {
        return false;
      }
      if (!q) {
        return true;
      }
      return m.id.toLowerCase().includes(q) || (m.name?.toLowerCase().includes(q) ?? false);
    });
    return sortProviderModels(matched, modelSortField, modelSortDirection);
  }, [allModels, freeModelsOnly, deferredModelSearch, modelSortField, modelSortDirection]);

  // Auto-detect & load live models from upstream whenever any provider with an active connection is opened.
  useEffect(() => {
    if (!activeConnectionForLoad?.enabled) return;
    const key = `${providerId}:${activeConnectionForLoad.id}`;
    if (autoLoadedConnRef.current === key) return;
    autoLoadedConnRef.current = key;
    fetchModelsMutate(activeConnectionForLoad.id, {
      onSuccess: (res) => {
        setLoadedModelsResult(res);
      },
    });
  }, [
    providerId,
    activeConnectionForLoad?.id,
    activeConnectionForLoad?.enabled,
    fetchModelsMutate,
  ]);

  useEffect(() => {
    return () => {
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
        copyTimerRef.current = null;
      }
      for (const controller of modelTestControllersRef.current.values()) {
        controller.abort();
      }
      modelTestControllersRef.current.clear();
      for (const controller of modelIdentityControllersRef.current.values()) {
        controller.abort();
      }
      modelIdentityControllersRef.current.clear();
    };
  }, []);

  if (catalogQuery.isPending || connectionsQuery.isPending) {
    return (
      <div
        role="status"
        aria-busy="true"
        aria-label={isId ? "Memuat detail provider" : "Loading provider details"}
        className="min-w-0 space-y-5"
      >
        <Skeleton className="h-5 w-40" />
        <div className="flex items-center gap-3">
          <Skeleton className="size-12 rounded-xl" />
          <div className="space-y-2">
            <Skeleton className="h-6 w-48" />
            <Skeleton className="h-3.5 w-28" />
          </div>
        </div>
        <Skeleton className="h-44 rounded-xl border border-[var(--color-border)]" />
        <Skeleton className="h-44 rounded-xl border border-[var(--color-border)]" />
      </div>
    );
  }

  if (!entry) {
    return (
      <Card>
        <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
          <CardTitle>{isId ? "Provider Tidak Ditemukan" : "Provider Not Found"}</CardTitle>
          <Button asChild variant="outline" size="xs" className="w-full sm:w-auto">
            <Link to="/dashboard/providers/">
              <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{isId ? "Kembali ke Providers" : "Back to Providers"}</span>
            </Link>
          </Button>
        </CardHeader>
        <CardContent>
          <Alert variant="warning">
            <AlertCircle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
              aria-hidden="true"
            />
            <div>
              <AlertTitle>{isId ? "Provider tidak terdaftar" : "Unknown provider ID"}</AlertTitle>
              <AlertDescription className="break-all font-mono">{providerId}</AlertDescription>
            </div>
          </Alert>
        </CardContent>
      </Card>
    );
  }

  const isOAuth = entry.isNativeOAuth || entry.category === "oauth";
  const isCustom = entry.category === "custom";
  const hasDualAuthModes = !isCustom && entry.isNativeOAuth;
  const externalLinkUrl = entry.apiKeyUrl ?? entry.signupUrl ?? entry.website;
  const modelPrefix = entry.alias || entry.id;
  const enabledCount = sortedConnections.filter((c) => c.enabled).length;
  const canTestModels = canProbeProviderModels(entry);

  const oauthButtonLabel =
    entry.id === "aws-quick"
      ? isId
        ? "Login Amazon Quick"
        : "Sign In Amazon Quick"
      : entry.adapter?.oauth?.flow_type === "device_code" ||
          entry.id === "github-copilot" ||
          entry.id === "kiro" ||
          entry.id === "freepi" ||
          entry.id === "codebuddy-cn" ||
          entry.id === "codebuddy-intl"
        ? isId
          ? "Device Login"
          : "Device Login"
        : "OAuth";

  const triggerPrimaryAdd = () => {
    setActionError(null);
    if (entry.isNativeOAuth) {
      setShowOAuthModal(true);
    } else {
      setShowAddApiKeyModal(true);
    }
  };

  const handleSwapPriority = async (indexA: number, indexB: number) => {
    const connA = sortedConnections[indexA];
    const connB = sortedConnections[indexB];
    if (!connA || !connB || busyConnectionId !== null) return;
    setActionError(null);
    setBusyConnectionId(connA.id);
    const nextPriorityA = connB.priority === connA.priority ? indexB + 1 : connB.priority;
    const nextPriorityB = connA.priority === connB.priority ? indexA + 1 : connA.priority;
    try {
      await Promise.all([
        updateProviderConnection({ id: connA.id, body: { priority: nextPriorityA } }),
        updateProviderConnection({ id: connB.id, body: { priority: nextPriorityB } }),
      ]);
      await Promise.all([
        invalidateModelRoutingQueries(queryClient, connA.id),
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(connB.id) }),
      ]);
    } catch (err) {
      setActionError(err instanceof Error ? err.message : "Failed to reorder priority");
    } finally {
      setBusyConnectionId(null);
    }
  };

  const handleTestSingle = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null || testingAll) return;
    setActionError(null);
    setBusyConnectionId(conn.id);
    testConnMutation.mutate(conn.id, {
      onSuccess: (res) => {
        setLiveTestResults((prev) => ({ ...prev, [conn.id]: res }));
      },
      onError: (err) => setActionError(err.message),
      onSettled: () => setBusyConnectionId(null),
    });
  };

  const handleTestAll = async () => {
    const enabled = sortedConnections.filter((c) => c.enabled);
    if (enabled.length === 0 || testingAll) return;
    setActionError(null);
    setTestingAll(true);
    let completed = 0;
    let passed = 0;
    let failed = 0;
    setTestAllProgress({
      total: enabled.length,
      completed: 0,
      passed: 0,
      failed: 0,
      currentConnectionId: enabled[0]?.id ?? null,
      currentConnectionName: enabled[0]?.name ?? null,
    });
    try {
      await runWithConcurrency(enabled, 4, async (conn) => {
        setTestAllProgress((prev) =>
          prev
            ? {
                ...prev,
                currentConnectionId: conn.id,
                currentConnectionName: conn.name,
              }
            : prev,
        );
        try {
          const res = await testProviderConnection(conn.id);
          setLiveTestResults((prev) => ({ ...prev, [conn.id]: res }));
          if (res.ok) {
            passed += 1;
          } else {
            failed += 1;
          }
        } catch (err) {
          failed += 1;
          setActionError(err instanceof Error ? err.message : `Failed testing ${conn.name}`);
        }
        completed += 1;
        setTestAllProgress({
          total: enabled.length,
          completed,
          passed,
          failed,
          currentConnectionId: completed < enabled.length ? conn.id : null,
          currentConnectionName: completed < enabled.length ? conn.name : null,
        });
      });
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.connections }),
        ...enabled.map((conn) =>
          queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(conn.id) }),
        ),
      ]);
    } finally {
      setTestingAll(false);
    }
  };

  const handleToggleConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) return;
    setActionError(null);
    setBusyConnectionId(conn.id);
    updateConnMutation.mutate(
      { id: conn.id, body: { enabled: !conn.enabled } },
      {
        onError: (err) => setActionError(err.message),
        onSettled: () => setBusyConnectionId(null),
      },
    );
  };

  const handleDeleteConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) return;
    setActionError(null);
    setBusyConnectionId(conn.id);
    deleteConnMutation.mutate(conn.id, {
      onSuccess: () => {
        setConfirmDeleteConnId(null);
        setFeedbackBanner(
          isId ? `Koneksi "${conn.name}" dihapus.` : `Deleted connection "${conn.name}".`,
        );
      },
      onError: (err) => setActionError(err.message),
      onSettled: () => setBusyConnectionId(null),
    });
  };

  const handleDeleteNode = () => {
    if (!entry.node || deleteNodeMutation.isPending) return;
    setActionError(null);
    deleteNodeMutation.mutate(entry.node.id, {
      onSuccess: () => {
        void navigate("/dashboard/providers/");
      },
      onError: (err) => setActionError(err.message),
    });
  };

  const handleCopyModel = async (modelId: string) => {
    const qualified = modelId.includes("/") ? modelId : `${modelPrefix}/${modelId}`;
    try {
      await navigator.clipboard.writeText(qualified);
      setCopiedModelId(modelId);
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
      copyTimerRef.current = window.setTimeout(() => {
        copyTimerRef.current = null;
        setCopiedModelId((cur) => (cur === modelId ? null : cur));
      }, 1800);
    } catch {
      // Ignore clipboard failures
    }
  };

  const handleTestModel = async (modelId: string) => {
    if (!canTestModels || testingModelIds[modelId]) {
      return;
    }
    const previous = modelTestControllersRef.current.get(modelId);
    if (previous) {
      previous.abort();
    }
    const controller = new AbortController();
    modelTestControllersRef.current.set(modelId, controller);
    setTestingModelIds((prev) => ({ ...prev, [modelId]: true }));
    try {
      const result = await runSingleModelTest(entry.id, modelId, controller.signal);
      if (controller.signal.aborted) {
        return;
      }
      setModelTestResults((prev) => ({ ...prev, [modelId]: result }));
    } finally {
      if (modelTestControllersRef.current.get(modelId) === controller) {
        modelTestControllersRef.current.delete(modelId);
      }
      if (!controller.signal.aborted) {
        setTestingModelIds((prev) => ({ ...prev, [modelId]: false }));
      }
    }
  };

  const handleIdentifyModel = async (modelId: string) => {
    if (!canTestModels || identifyingModelIds[modelId]) {
      return;
    }
    const previous = modelIdentityControllersRef.current.get(modelId);
    if (previous) {
      previous.abort();
    }
    const controller = new AbortController();
    modelIdentityControllersRef.current.set(modelId, controller);
    setIdentifyingModelIds((prev) => ({ ...prev, [modelId]: true }));
    try {
      const result = await runModelIdentityProbe(entry.id, modelId, controller.signal);
      if (controller.signal.aborted) {
        return;
      }
      setModelIdentityResults((prev) => ({ ...prev, [modelId]: result }));
      window.requestAnimationFrame(() => {
        if (!controller.signal.aborted) {
          setOpenIdentityModelIds((prev) => ({ ...prev, [modelId]: true }));
        }
      });
    } finally {
      if (modelIdentityControllersRef.current.get(modelId) === controller) {
        modelIdentityControllersRef.current.delete(modelId);
      }
      if (!controller.signal.aborted) {
        setIdentifyingModelIds((prev) => ({ ...prev, [modelId]: false }));
      }
    }
  };

  const handleLoadModels = () => {
    if (!activeConnectionForLoad || fetchModelsMutation.isPending) return;
    setActionError(null);
    fetchModelsMutate(activeConnectionForLoad.id, {
      onSuccess: (res) => {
        setLoadedModelsResult(res);
      },
      onError: (err) => {
        setActionError(err.message);
      },
    });
  };

  const handleRemoveCustomModel = async (modelId: string) => {
    if (removingModelId) return;
    const targets = sortedConnections.filter((conn) =>
      conn.models.some((m) => m.trim() === modelId),
    );
    if (targets.length === 0) return;

    setActionError(null);
    setRemovingModelId(modelId);
    try {
      await runWithConcurrency(targets, 4, async (conn) => {
        const nextModels = removeModelFromList(conn.models, modelId);
        return updateProviderConnection({
          id: conn.id,
          body: { models: nextModels },
        });
      });
      await Promise.all([
        invalidateModelRoutingQueries(queryClient),
        ...targets.map((conn) =>
          queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(conn.id) }),
        ),
      ]);
      setFeedbackBanner(
        isId
          ? `Model manual "${modelId}" dihapus dari koneksi ${entry.name}.`
          : `Removed custom model "${modelId}" from ${entry.name}.`,
      );
    } catch (err) {
      setActionError(err instanceof Error ? err.message : "Failed to remove custom model.");
    } finally {
      setRemovingModelId(null);
    }
  };

  return (
    <div className="min-w-0 space-y-5">
      {/* OmniRoute Breadcrumb + Provider Header */}
      <div className="space-y-4">
        <div className="flex flex-wrap items-center gap-2 text-xs text-[var(--color-text-muted)]">
          <Link
            to="/dashboard/providers/"
            className="inline-flex items-center gap-1 font-medium transition-colors hover:text-[var(--color-primary)]"
          >
            <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{isId ? "Providers" : "Providers"}</span>
          </Link>
          <span>/</span>
          <span className="font-medium text-[var(--color-text-main)]">{entry.name}</span>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex min-w-0 items-center gap-3.5">
            <div className="flex size-12 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/80 p-2">
              <ProviderLogo
                providerId={entry.id}
                name={entry.name}
                protocol={entry.protocol}
                isCustom={entry.isCustom}
                size={30}
              />
            </div>
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="truncate text-xl font-bold tracking-tight text-[var(--color-text-main)] sm:text-2xl">
                  {entry.name}
                </h1>
                {externalLinkUrl ? (
                  <a
                    href={externalLinkUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-xs font-medium text-[var(--color-primary)] hover:underline"
                  >
                    <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>
                      {entry.apiKeyUrl
                        ? isId
                          ? "Ambil API Key"
                          : "Get API Key"
                        : isId
                          ? "Daftar / Pelajari"
                          : "Sign up / Learn more"}
                    </span>
                  </a>
                ) : null}
              </div>
              <p className="mt-0.5 text-xs text-[var(--color-text-muted)]">
                {sortedConnections.length}{" "}
                {isId ? "koneksi" : `connection${sortedConnections.length === 1 ? "" : "s"}`}
                {" · "}
                <span className="font-mono uppercase">{entry.protocol}</span>
              </p>
            </div>
          </div>

          <Button
            type="button"
            variant="secondary"
            size="sm"
            className="w-full shrink-0 sm:w-auto"
            onClick={() => void Promise.all([catalogQuery.refetch(), connectionsQuery.refetch()])}
            disabled={catalogQuery.isFetching || connectionsQuery.isFetching}
          >
            <RefreshCw
              className={cn(
                "h-3.5 w-3.5",
                (catalogQuery.isFetching || connectionsQuery.isFetching) && "animate-spin",
              )}
              aria-hidden="true"
            />
            <span>{isId ? "Segarkan" : "Refresh"}</span>
          </Button>
        </div>
      </div>

      {/* OmniRoute Deprecation / Info Notice Banner */}
      {entry.deprecated ? (
        <div className="flex items-start gap-2.5 rounded-xl border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-3.5 py-2.5">
          <AlertTriangle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
            aria-hidden="true"
          />
          <p className="text-xs leading-relaxed text-[var(--color-text-main)]">
            {entry.noticeText ||
              (isId ? "Provider ini sudah ditandai deprecated." : "This provider is deprecated.")}
          </p>
        </div>
      ) : entry.noticeText ? (
        <div className="flex flex-col gap-2 rounded-xl border border-[var(--color-primary)]/25 bg-[var(--color-primary)]/8 px-3.5 py-2.5 sm:flex-row sm:items-center">
          <div className="flex min-w-0 flex-1 items-start gap-2">
            <Info
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--color-primary)]"
              aria-hidden="true"
            />
            <p className="min-w-0 flex-1 text-xs leading-relaxed text-[var(--color-text-main)]">
              {entry.noticeText}
            </p>
          </div>
          {entry.apiKeyUrl ? (
            <a
              href={entry.apiKeyUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex shrink-0 items-center justify-center rounded-lg bg-[var(--color-primary)] px-2.5 py-1.5 text-xs font-medium text-[var(--color-on-primary)] transition-opacity hover:opacity-90 sm:py-1"
            >
              {isId ? "Ambil API Key →" : "Get API Key →"}
            </a>
          ) : null}
        </div>
      ) : null}

      {/* Feedback / Error Banners */}
      {actionError ? (
        <Alert variant="destructive">
          <div className="flex items-center justify-between gap-2">
            <AlertDescription className="break-words">{actionError}</AlertDescription>
            <button
              type="button"
              onClick={() => setActionError(null)}
              aria-label="Dismiss error"
              className="cursor-pointer rounded p-1 hover:bg-black/10"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        </Alert>
      ) : null}

      {feedbackBanner ? (
        <Alert>
          <div className="flex items-center justify-between gap-2">
            <AlertDescription>{feedbackBanner}</AlertDescription>
            <button
              type="button"
              onClick={() => setFeedbackBanner(null)}
              aria-label="Dismiss notice"
              className="cursor-pointer rounded p-1 hover:bg-black/10"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        </Alert>
      ) : null}

      {/* Custom Compatible Node Details Card (OmniRoute parity) */}
      {isCustom && entry.node ? (
        <Card>
          <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="min-w-0">
              <CardTitle>
                {entry.protocol === "anthropic"
                  ? "Anthropic Compatible Details"
                  : "OpenAI Compatible Details"}
              </CardTitle>
              <CardDescription className="mt-1 break-all font-mono">
                {entry.protocol.toUpperCase()} · {entry.baseUrl}
              </CardDescription>
            </div>
            <div className="grid w-full grid-cols-1 gap-2 sm:flex sm:w-auto sm:items-center">
              <Button
                type="button"
                size="sm"
                className="w-full sm:w-auto"
                onClick={() => setShowAddApiKeyModal(true)}
              >
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Tambah API Key" : "Add API Key"}</span>
              </Button>
              {confirmDeleteNode ? (
                <div className="grid grid-cols-2 gap-2 sm:flex sm:items-center">
                  <Button
                    type="button"
                    size="sm"
                    variant="destructive"
                    className="w-full sm:w-auto"
                    disabled={deleteNodeMutation.isPending}
                    onClick={handleDeleteNode}
                  >
                    {isId ? "Ya, Hapus" : "Confirm"}
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="secondary"
                    className="w-full sm:w-auto"
                    onClick={() => setConfirmDeleteNode(false)}
                  >
                    {isId ? "Batal" : "Cancel"}
                  </Button>
                </div>
              ) : (
                <Button
                  type="button"
                  size="sm"
                  variant="secondary"
                  className="w-full sm:w-auto"
                  onClick={() => setConfirmDeleteNode(true)}
                >
                  <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Hapus" : "Delete"}</span>
                </Button>
              )}
            </div>
          </CardHeader>
        </Card>
      ) : null}

      {/* OmniRoute Connections Card */}
      <Card>
        <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2">
            <CardTitle className="text-base">{isId ? "Koneksi" : "Connections"}</CardTitle>
            <span className="rounded-full bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
              {isConnectionFiltered
                ? `${filteredConnections.length} / ${sortedConnections.length}`
                : sortedConnections.length}
            </span>
            {connectionCounts.issues > 0 ? (
              <Badge variant="destructive" className="tabular-nums">
                {connectionCounts.issues} {isId ? "bermasalah" : "issues"}
              </Badge>
            ) : null}
          </div>
          {enabledCount > 0 ? (
            <Button
              type="button"
              size="sm"
              variant="secondary"
              className="w-full sm:w-auto"
              disabled={testingAll || busyConnectionId !== null}
              onClick={() => void handleTestAll()}
            >
              <RefreshCw
                className={cn("h-3.5 w-3.5", testingAll && "animate-spin")}
                aria-hidden="true"
              />
              <span>
                {testingAll && testAllProgress
                  ? isId
                    ? `Menguji ${testAllProgress.completed}/${testAllProgress.total}...`
                    : `Testing ${testAllProgress.completed}/${testAllProgress.total}...`
                  : isId
                    ? "Tes Semua Koneksi"
                    : "Test Connections"}
              </span>
            </Button>
          ) : null}
        </CardHeader>

        <CardContent>
          {sortedConnections.length === 0 ? (
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex items-center gap-3">
                <div
                  className={cn(
                    "inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-full",
                    entry.noAuth
                      ? "bg-[var(--status-ok-bg)] text-[var(--status-ok)]"
                      : "bg-[var(--color-primary)]/10 text-[var(--color-primary)]",
                  )}
                >
                  {entry.noAuth ? (
                    <CheckCircle2 className="h-4 w-4" aria-hidden="true" />
                  ) : isOAuth ? (
                    <Lock className="h-4 w-4" aria-hidden="true" />
                  ) : (
                    <KeyRound className="h-4 w-4" aria-hidden="true" />
                  )}
                </div>
                <div className="min-w-0">
                  <p className="text-sm font-medium text-[var(--color-text-main)]">
                    {entry.noAuth
                      ? isId
                        ? "Siap digunakan (Tanpa Auth)"
                        : "Ready to use (No Auth Required)"
                      : isId
                        ? "Belum ada koneksi"
                        : "No connections yet"}
                  </p>
                  <p className="text-xs text-[var(--color-text-muted)]">
                    {entry.noAuth
                      ? isId
                        ? `${entry.name} berjalan tanpa API key atau login — langsung gunakan model di bawah.`
                        : `${entry.name} works without an API key or login — use the models below directly.`
                      : isId
                        ? `Tambahkan akun ${entry.name} untuk mulai merutekan request.`
                        : `Add a ${entry.name} account to start routing requests.`}
                  </p>
                </div>
              </div>

              {canAddProviderConnection(entry) ? (
                <div className="grid grid-cols-1 gap-2 sm:flex sm:items-center">
                  {hasDualAuthModes ? (
                    <>
                      <Button
                        type="button"
                        size="sm"
                        className="w-full sm:w-auto"
                        onClick={() => setShowOAuthModal(true)}
                      >
                        <Lock className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>{oauthButtonLabel}</span>
                      </Button>
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        className="w-full sm:w-auto"
                        onClick={() => setShowAddApiKeyModal(true)}
                      >
                        <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>API Key</span>
                      </Button>
                    </>
                  ) : (
                    <Button
                      type="button"
                      size="sm"
                      className="w-full sm:w-auto"
                      onClick={triggerPrimaryAdd}
                    >
                      <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>{isId ? "Tambah" : "Add"}</span>
                    </Button>
                  )}
                </div>
              ) : null}
            </div>
          ) : (
            <div className="space-y-3">
              {testAllProgress ? (
                <div
                  className={cn(
                    "space-y-2 rounded-lg border px-3 py-2.5 text-xs",
                    testingAll
                      ? "border-[var(--color-primary)]/35 bg-[var(--color-primary)]/8"
                      : testAllProgress.failed > 0
                        ? "border-[var(--status-error-border)] bg-[var(--status-error-bg)]"
                        : "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)]",
                  )}
                >
                  <div className="flex items-center justify-between gap-2">
                    <div className="flex min-w-0 items-center gap-2">
                      {testingAll ? (
                        <RefreshCw
                          className="h-3.5 w-3.5 shrink-0 animate-spin text-[var(--color-primary)]"
                          aria-hidden="true"
                        />
                      ) : testAllProgress.failed > 0 ? (
                        <AlertTriangle
                          className="h-3.5 w-3.5 shrink-0 text-[var(--status-error)]"
                          aria-hidden="true"
                        />
                      ) : (
                        <CheckCircle2
                          className="h-3.5 w-3.5 shrink-0 text-[var(--status-ok)]"
                          aria-hidden="true"
                        />
                      )}
                      <span className="truncate font-medium text-[var(--color-text-main)]">
                        {testingAll
                          ? isId
                            ? `Memeriksa ${testAllProgress.currentConnectionName ?? "koneksi"} (${testAllProgress.completed}/${testAllProgress.total})...`
                            : `Checking ${testAllProgress.currentConnectionName ?? "connection"} (${testAllProgress.completed}/${testAllProgress.total})...`
                          : isId
                            ? `Selesai memeriksa ${testAllProgress.completed}/${testAllProgress.total} koneksi`
                            : `Finished checking ${testAllProgress.completed}/${testAllProgress.total} connections`}
                      </span>
                    </div>
                    <div className="flex shrink-0 items-center gap-2 font-mono text-[11px] tabular-nums">
                      <span className="text-[var(--status-ok)]">{testAllProgress.passed} OK</span>
                      {testAllProgress.failed > 0 ? (
                        <span className="text-[var(--status-error)]">
                          {testAllProgress.failed} {isId ? "gagal" : "failed"}
                        </span>
                      ) : null}
                      {!testingAll ? (
                        <button
                          type="button"
                          onClick={() => setTestAllProgress(null)}
                          aria-label={isId ? "Tutup ringkasan tes" : "Dismiss test summary"}
                          className="rounded p-0.5 text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]"
                        >
                          <X className="h-3.5 w-3.5" aria-hidden="true" />
                        </button>
                      ) : null}
                    </div>
                  </div>
                  <div className="h-1.5 w-full overflow-hidden rounded-full bg-[var(--color-surface-2)]">
                    <div
                      className={cn(
                        "h-full rounded-full transition-all duration-200",
                        !testingAll && testAllProgress.failed > 0
                          ? "bg-[var(--status-error)]"
                          : !testingAll
                            ? "bg-[var(--status-ok)]"
                            : "bg-[var(--color-primary)]",
                      )}
                      style={{
                        width: `${
                          testAllProgress.total > 0
                            ? Math.round((testAllProgress.completed / testAllProgress.total) * 100)
                            : 0
                        }%`,
                      }}
                    />
                  </div>
                </div>
              ) : null}

              <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
                <div className="flex flex-wrap items-center gap-1.5">
                  {(
                    [
                      {
                        id: "all",
                        label: isId ? "Semua" : "All",
                        count: connectionCounts.all,
                      },
                      {
                        id: "issues",
                        label: isId ? "Bermasalah" : "Issues",
                        count: connectionCounts.issues,
                      },
                      {
                        id: "active",
                        label: isId ? "Aktif" : "Active",
                        count: connectionCounts.active,
                      },
                      {
                        id: "disabled",
                        label: isId ? "Nonaktif" : "Disabled",
                        count: connectionCounts.disabled,
                      },
                    ] as const
                  ).map((tab) => {
                    const isSelected = connectionFilter === tab.id;
                    const isIssueHighlight = tab.id === "issues" && tab.count > 0 && !isSelected;
                    return (
                      <Button
                        key={tab.id}
                        type="button"
                        size="xs"
                        variant={
                          isSelected
                            ? tab.id === "issues"
                              ? "destructive"
                              : "default"
                            : "secondary"
                        }
                        aria-pressed={isSelected}
                        className={cn(
                          isIssueHighlight &&
                            "border-[var(--status-error-border)] text-[var(--status-error)]",
                        )}
                        onClick={() => setConnectionFilter(tab.id)}
                      >
                        <span>{tab.label}</span>
                        <span className="rounded-full bg-[var(--color-surface)]/25 px-1.5 py-0.25 font-mono text-[10px] tabular-nums">
                          {tab.count}
                        </span>
                      </Button>
                    );
                  })}
                </div>
                {sortedConnections.length > 5 ? (
                  <div className="relative w-full sm:w-56">
                    <Search
                      className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--color-text-muted)]"
                      aria-hidden="true"
                    />
                    <Input
                      value={connectionSearch}
                      onChange={(e) => setConnectionSearch(e.target.value)}
                      placeholder={
                        isId ? "Cari koneksi atau error..." : "Filter connections or errors..."
                      }
                      className="h-8 pl-8 text-xs"
                    />
                  </div>
                ) : null}
              </div>

              {filteredConnections.length === 0 ? (
                <p className="py-4 text-center text-xs text-[var(--color-text-muted)]">
                  {connectionFilter === "issues" && !connectionSearch.trim()
                    ? isId
                      ? "Tidak ada key/OAuth yang bermasalah."
                      : "No keys or OAuth connections with issues."
                    : isId
                      ? "Tidak ada koneksi yang cocok dengan filter saat ini."
                      : "No connections match the current filter."}
                </p>
              ) : (
                <div className="max-h-[360px] divide-y divide-[var(--color-border-subtle)] overflow-y-auto pr-1">
                  {filteredConnections.map((conn, idx) => {
                    const isFirst = isConnectionFiltered || idx === 0;
                    const isLast = isConnectionFiltered || idx === filteredConnections.length - 1;
                    const live = liveTestResults[conn.id];
                    const isCheckingNow =
                      testingAll && testAllProgress?.currentConnectionId === conn.id;
                    const isBusy = busyConnectionId === conn.id || isCheckingNow;
                    const isConfirming = confirmDeleteConnId === conn.id;

                    return (
                      <div
                        key={conn.id}
                        className={cn(
                          "group flex min-w-0 flex-col gap-3 py-3 first:pt-0 last:pb-0 sm:flex-row sm:items-center sm:justify-between",
                          !conn.enabled && "opacity-65",
                          isCheckingNow && "rounded-lg bg-[var(--color-primary)]/6 px-2",
                        )}
                      >
                        <div className="flex min-w-0 flex-1 items-start gap-2.5 sm:items-center">
                          {/* OmniRoute Priority Up/Down Arrows */}
                          <div className="flex shrink-0 flex-col">
                            <button
                              type="button"
                              disabled={isFirst || busyConnectionId !== null}
                              onClick={() => void handleSwapPriority(idx, idx - 1)}
                              aria-label={`Move ${conn.name} up`}
                              className={cn(
                                "rounded p-0.5 transition-colors",
                                isFirst
                                  ? "cursor-not-allowed text-[var(--color-text-muted)]/30"
                                  : "cursor-pointer text-[var(--color-text-muted)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-primary)]",
                              )}
                            >
                              <ChevronUp className="h-3.5 w-3.5" aria-hidden="true" />
                            </button>
                            <button
                              type="button"
                              disabled={isLast || busyConnectionId !== null}
                              onClick={() => void handleSwapPriority(idx, idx + 1)}
                              aria-label={`Move ${conn.name} down`}
                              className={cn(
                                "rounded p-0.5 transition-colors",
                                isLast
                                  ? "cursor-not-allowed text-[var(--color-text-muted)]/30"
                                  : "cursor-pointer text-[var(--color-text-muted)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-primary)]",
                              )}
                            >
                              <ChevronDown className="h-3.5 w-3.5" aria-hidden="true" />
                            </button>
                          </div>

                          <span className="mt-0.5 shrink-0 text-[var(--color-text-muted)] sm:mt-0">
                            {conn.auth_type === "oauth" ? (
                              <Lock className="h-4 w-4" aria-hidden="true" />
                            ) : (
                              <KeyRound className="h-4 w-4" aria-hidden="true" />
                            )}
                          </span>

                          <div className="min-w-0 flex-1">
                            <p className="truncate text-sm font-medium text-[var(--color-text-main)]">
                              {conn.name}
                            </p>
                            <div className="mt-1 flex min-w-0 flex-wrap items-center gap-1.5">
                              <Badge
                                variant={
                                  !conn.enabled
                                    ? "outline"
                                    : conn.last_test_status === "error"
                                      ? "destructive"
                                      : "success"
                                }
                              >
                                ●{" "}
                                {!conn.enabled
                                  ? isId
                                    ? "nonaktif"
                                    : "disabled"
                                  : conn.last_test_status === "error"
                                    ? "error"
                                    : isId
                                      ? "aktif"
                                      : "active"}
                              </Badge>
                              <Badge variant="secondary">
                                {conn.auth_type === "oauth" ? "OAuth" : "API Key"}
                              </Badge>
                              <span className="font-mono text-xs text-[var(--color-text-muted)]">
                                #{conn.priority}
                              </span>
                              {isCheckingNow ? (
                                <Badge variant="secondary" className="gap-1">
                                  <RefreshCw
                                    className="h-3 w-3 animate-spin text-[var(--color-primary)]"
                                    aria-hidden="true"
                                  />
                                  <span>{isId ? "Memeriksa..." : "Checking..."}</span>
                                </Badge>
                              ) : live ? (
                                <Badge
                                  variant={live.ok ? "success" : "destructive"}
                                  className="tabular-nums"
                                >
                                  {live.ok
                                    ? `OK ${live.status_code} · ${live.latency_ms}ms`
                                    : `FAIL ${live.status_code || "ERR"} · ${live.latency_ms}ms`}
                                </Badge>
                              ) : null}
                              {conn.last_test_error && conn.enabled ? (
                                <span
                                  className="max-w-full truncate font-mono text-xs text-[var(--status-error)] sm:max-w-[280px]"
                                  title={conn.last_test_error}
                                >
                                  {conn.last_test_error}
                                </span>
                              ) : null}
                            </div>
                          </div>
                        </div>

                        {/* Connection Row Actions */}
                        <div className="flex w-full flex-wrap items-center justify-between gap-2 sm:w-auto sm:justify-end">
                          <div className="flex flex-1 flex-wrap items-center gap-1.5 sm:flex-none">
                            <Button
                              type="button"
                              size="xs"
                              variant="secondary"
                              className="flex-1 sm:flex-none"
                              disabled={busyConnectionId !== null || testingAll}
                              onClick={() => handleTestSingle(conn)}
                            >
                              {isBusy ? (
                                <RefreshCw
                                  className="h-3.5 w-3.5 animate-spin"
                                  aria-hidden="true"
                                />
                              ) : (
                                <Activity className="h-3.5 w-3.5" aria-hidden="true" />
                              )}
                              <span>{isId ? "Tes" : "Test"}</span>
                            </Button>

                            <Button
                              type="button"
                              size="xs"
                              variant="secondary"
                              className="flex-1 sm:flex-none"
                              onClick={() => setEditingConnection(conn)}
                            >
                              <Settings2 className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>{isId ? "Edit" : "Edit"}</span>
                            </Button>
                          </div>

                          <div className="flex items-center gap-2">
                            <button
                              type="button"
                              role="switch"
                              aria-checked={conn.enabled}
                              disabled={busyConnectionId !== null}
                              onClick={() => handleToggleConnection(conn)}
                              aria-label={`${conn.enabled ? "Disable" : "Enable"} ${conn.name}`}
                              className={cn(
                                "relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border transition-colors disabled:opacity-50",
                                conn.enabled
                                  ? "border-[var(--color-primary)] bg-[var(--color-primary)]"
                                  : "border-[var(--color-border)] bg-[var(--color-surface-3)]",
                              )}
                            >
                              <span
                                className={cn(
                                  "inline-block h-3.5 w-3.5 rounded-full bg-[var(--color-on-primary)] shadow-xs transition-transform",
                                  conn.enabled ? "translate-x-4" : "translate-x-0.5",
                                )}
                              />
                            </button>

                            {isConfirming ? (
                              <div className="flex items-center gap-1">
                                <Button
                                  type="button"
                                  size="xs"
                                  variant="destructive"
                                  disabled={busyConnectionId !== null}
                                  onClick={() => handleDeleteConnection(conn)}
                                >
                                  {isId ? "Hapus" : "Confirm"}
                                </Button>
                                <Button
                                  type="button"
                                  size="xs"
                                  variant="secondary"
                                  onClick={() => setConfirmDeleteConnId(null)}
                                >
                                  {isId ? "Batal" : "Cancel"}
                                </Button>
                              </div>
                            ) : (
                              <button
                                type="button"
                                disabled={busyConnectionId !== null}
                                onClick={() => setConfirmDeleteConnId(conn.id)}
                                aria-label={`Delete ${conn.name}`}
                                className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--status-error-bg)] hover:text-[var(--status-error)] disabled:opacity-40"
                              >
                                <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                              </button>
                            )}
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}

              {!isCustom && canAddProviderConnection(entry) ? (
                <div className="grid grid-cols-1 gap-2 border-t border-[var(--color-border-subtle)] pt-3 sm:flex">
                  {hasDualAuthModes ? (
                    <>
                      <Button
                        type="button"
                        size="sm"
                        className="w-full sm:w-auto"
                        onClick={() => setShowOAuthModal(true)}
                      >
                        <Lock className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>{oauthButtonLabel}</span>
                      </Button>
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        className="w-full sm:w-auto"
                        onClick={() => setShowAddApiKeyModal(true)}
                      >
                        <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>API Key</span>
                      </Button>
                    </>
                  ) : (
                    <Button
                      type="button"
                      size="sm"
                      className="w-full sm:w-auto"
                      onClick={triggerPrimaryAdd}
                    >
                      <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>{isId ? "Tambah" : "Add"}</span>
                    </Button>
                  )}
                </div>
              ) : null}
            </div>
          )}
        </CardContent>
      </Card>

      {/* OmniRoute Account Routing Strategy Card */}
      <ProviderStrategyCard
        providerId={entry.id}
        enabledConnectionCount={enabledCount}
        isId={isId}
      />

      {/* OmniRoute Available Models Card (ModelRow compact chip styling + Load Models) */}
      <Card>
        <CardHeader className="flex-col items-start gap-2.5 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2">
            <CardTitle className="text-base">
              {isId ? "Available Models" : "Available Models"}
            </CardTitle>
            <span className="rounded-full bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
              {freeModelsOnly || modelSearch.trim()
                ? `${filteredModels.length} / ${allModels.length}`
                : allModels.length}
            </span>
          </div>
          <div className="flex w-full flex-col gap-2 sm:w-auto sm:flex-row sm:flex-wrap sm:items-center">
            {allModels.length > 10 ? (
              <div className="relative w-full sm:w-56">
                <Search
                  className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--color-text-muted)]"
                  aria-hidden="true"
                />
                <Input
                  value={modelSearch}
                  onChange={(e) => setModelSearch(e.target.value)}
                  placeholder={isId ? "Cari model..." : "Filter models..."}
                  className="h-8 pl-8 text-xs"
                />
              </div>
            ) : null}
            {freeModelsCount > 0 || freeModelsOnly ? (
              <Button
                type="button"
                size="xs"
                variant={freeModelsOnly ? "default" : "secondary"}
                className="w-full shrink-0 sm:w-auto"
                aria-pressed={freeModelsOnly}
                onClick={() => setFreeModelsOnly((prev) => !prev)}
                title={
                  freeModelsOnly
                    ? isId
                      ? "Menampilkan hanya model gratis (klik untuk tampilkan semua)"
                      : "Showing free models only (click to show all)"
                    : isId
                      ? "Filter hanya model gratis"
                      : "Filter free models only"
                }
              >
                <span>{isId ? "Gratis" : "Free"}</span>
                <span className="rounded-full bg-[var(--color-surface)]/25 px-1.5 py-0.25 font-mono text-[10px] tabular-nums">
                  {freeModelsCount}
                </span>
              </Button>
            ) : null}
            {allModels.length > 1 ? (
              <div className="flex w-full items-center gap-1.5 sm:w-auto">
                <Button
                  type="button"
                  size="xs"
                  variant={modelSortField === "intelligence" ? "default" : "secondary"}
                  className="flex-1 shrink-0 sm:flex-none"
                  onClick={() => {
                    if (modelSortField === "intelligence") {
                      setModelSortField("name");
                      setModelSortDirection("asc");
                    } else {
                      setModelSortField("intelligence");
                      setModelSortDirection("desc");
                    }
                  }}
                  title={
                    modelSortField === "intelligence"
                      ? isId
                        ? "Diurutkan berdasarkan kecerdasan (klik untuk urutkan nama)"
                        : "Sorted by intelligence (click to sort by name)"
                      : isId
                        ? "Urutkan berdasarkan tingkat kecerdasan"
                        : "Sort by intelligence level"
                  }
                >
                  <Brain className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {modelSortField === "intelligence"
                      ? isId
                        ? "Kecerdasan"
                        : "Intelligence"
                      : isId
                        ? "Nama"
                        : "Name"}
                  </span>
                </Button>
                <Button
                  type="button"
                  size="xs"
                  variant="secondary"
                  className="flex-1 shrink-0 sm:flex-none"
                  onClick={() => setModelSortDirection((prev) => (prev === "asc" ? "desc" : "asc"))}
                  aria-label={
                    modelSortField === "intelligence"
                      ? modelSortDirection === "desc"
                        ? isId
                          ? "Urutkan kecerdasan dari terendah"
                          : "Sort intelligence lowest first"
                        : isId
                          ? "Urutkan kecerdasan dari tertinggi"
                          : "Sort intelligence highest first"
                      : modelSortDirection === "asc"
                        ? isId
                          ? "Urutkan model Z ke A"
                          : "Sort models Z to A"
                        : isId
                          ? "Urutkan model A ke Z"
                          : "Sort models A to Z"
                  }
                  title={
                    modelSortField === "intelligence"
                      ? modelSortDirection === "desc"
                        ? isId
                          ? "Kecerdasan tertinggi dulu (klik untuk terendah)"
                          : "Highest intelligence first (click for lowest)"
                        : isId
                          ? "Kecerdasan terendah dulu (klik untuk tertinggi)"
                          : "Lowest intelligence first (click for highest)"
                      : modelSortDirection === "asc"
                        ? isId
                          ? "Urutan A–Z (klik untuk Z–A)"
                          : "Sorted A–Z (click for Z–A)"
                        : isId
                          ? "Urutan Z–A (klik untuk A–Z)"
                          : "Sorted Z–A (click for A–Z)"
                  }
                >
                  {modelSortField === "intelligence" ? (
                    modelSortDirection === "desc" ? (
                      <ArrowDownWideNarrow className="h-3.5 w-3.5" aria-hidden="true" />
                    ) : (
                      <ArrowUpNarrowWide className="h-3.5 w-3.5" aria-hidden="true" />
                    )
                  ) : modelSortDirection === "asc" ? (
                    <ArrowDownAZ className="h-3.5 w-3.5" aria-hidden="true" />
                  ) : (
                    <ArrowUpZA className="h-3.5 w-3.5" aria-hidden="true" />
                  )}
                  <span>
                    {modelSortField === "intelligence"
                      ? modelSortDirection === "desc"
                        ? isId
                          ? "Tertinggi"
                          : "Highest"
                        : isId
                          ? "Terendah"
                          : "Lowest"
                      : modelSortDirection === "asc"
                        ? "A–Z"
                        : "Z–A"}
                  </span>
                </Button>
              </div>
            ) : null}
            <div className="flex w-full items-center gap-1.5 sm:w-auto">
              <Button
                type="button"
                size="xs"
                variant="secondary"
                className="flex-1 shrink-0 sm:flex-none"
                disabled={sortedConnections.length === 0}
                title={
                  sortedConnections.length > 0
                    ? isId
                      ? "Tambahkan model secara manual ke koneksi provider ini"
                      : "Manually add a custom model ID to this provider"
                    : isId
                      ? "Tambahkan koneksi terlebih dahulu untuk menambah model manual"
                      : "Add a connection first to add custom models"
                }
                onClick={() => setShowAddModelModal(true)}
              >
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Tambah Model" : "Add Model"}</span>
              </Button>
              <Button
                type="button"
                size="xs"
                variant="secondary"
                className="flex-1 shrink-0 sm:flex-none"
                disabled={!activeConnectionForLoad || fetchModelsMutation.isPending}
                title={
                  activeConnectionForLoad
                    ? isId
                      ? `Muat daftar model terbaru dari ${activeConnectionForLoad.name}`
                      : `Fetch live models from ${activeConnectionForLoad.name}`
                    : isId
                      ? "Tambahkan koneksi aktif terlebih dahulu untuk memuat model"
                      : "Add an active connection first to load live models"
                }
                onClick={handleLoadModels}
              >
                {fetchModelsMutation.isPending ? (
                  <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                ) : (
                  <Download className="h-3.5 w-3.5" aria-hidden="true" />
                )}
                <span>
                  {fetchModelsMutation.isPending
                    ? isId
                      ? "Memuat Model..."
                      : "Loading Models..."
                    : isId
                      ? "Muat Model"
                      : "Load Models"}
                </span>
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-3">
          {loadedModelsResult ? (
            <div
              className={cn(
                "flex flex-col gap-2 rounded-lg border px-3 py-2 text-xs sm:flex-row sm:items-center sm:justify-between",
                loadedModelsResult.source === "upstream"
                  ? "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] text-[var(--color-text-main)]"
                  : "border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] text-[var(--color-text-main)]",
              )}
            >
              <div className="min-w-0">
                {loadedModelsResult.source === "upstream" ? (
                  <p className="font-medium">
                    {isId
                      ? `${loadedModelsResult.models.length} model berhasil dimuat dari upstream (${loadedModelsResult.latency_ms}ms) & otomatis siap digunakan.`
                      : `Loaded ${loadedModelsResult.models.length} live models from upstream (${loadedModelsResult.latency_ms}ms) & auto-synced for routing.`}
                  </p>
                ) : (
                  <p className="break-words">
                    {isId ? "Fallback ke katalog statis:" : "Fell back to static catalog:"}{" "}
                    <span className="font-mono text-[11px]">{loadedModelsResult.warning}</span>
                  </p>
                )}
              </div>
            </div>
          ) : null}

          {filteredModels.length === 0 ? (
            <p className="text-xs text-[var(--color-text-muted)]">
              {allModels.length > 0 && (freeModelsOnly || modelSearch.trim())
                ? isId
                  ? "Tidak ada model yang cocok dengan filter saat ini."
                  : "No models match the current filter."
                : isId
                  ? "Belum ada daftar model bawaan untuk provider ini."
                  : "No built-in models listed for this provider."}
            </p>
          ) : (
            <div className="grid grid-cols-1 items-start gap-2.5 sm:grid-cols-2 lg:grid-cols-3">
              {filteredModels.map((model) => {
                const fullModel = model.id.includes("/") ? model.id : `${modelPrefix}/${model.id}`;
                const isCopied = copiedModelId === model.id;
                const isTestingModel = Boolean(testingModelIds[model.id]);
                const modelTestResult = modelTestResults[model.id];
                const isIdentifyingModel = Boolean(identifyingModelIds[model.id]);
                const identityResult = modelIdentityResults[model.id];
                const isIdentityOpen = Boolean(identityResult && openIdentityModelIds[model.id]);
                const hasSubtitle = Boolean(model.name && model.name.trim() !== model.id);
                const intelligence = getModelIntelligenceProfile(model.id, model.name);
                const tierBadge = MODEL_TIER_BADGE[intelligence.tier];

                return (
                  <div
                    key={model.id}
                    className="group flex min-h-[96px] min-w-0 flex-col justify-between gap-2.5 rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-2.5 transition-colors hover:border-[var(--border-strong)] hover:bg-[var(--color-surface-2)]/50"
                  >
                    <div className="min-w-0 space-y-1">
                      <div className="flex min-w-0 items-center gap-1.5">
                        <Bot
                          className="h-3.5 w-3.5 shrink-0 text-[var(--color-text-muted)]"
                          aria-hidden="true"
                        />
                        <code
                          className="min-w-0 flex-1 truncate rounded bg-[var(--color-surface-2)] px-1.5 py-0.5 font-mono text-[11px] text-[var(--color-text-main)]"
                          title={fullModel}
                        >
                          {fullModel}
                        </code>
                        <span
                          className={cn(
                            "shrink-0 rounded px-1.5 py-0.5 text-[9px] font-bold leading-none",
                            tierBadge.className,
                          )}
                          title={
                            isId
                              ? `Tier ${tierBadge.label} · estimasi skor ~${intelligence.score} (${intelligence.source === "curated" ? "kurasi" : "heuristik"})`
                              : `${tierBadge.label} tier · est. score ~${intelligence.score} (${intelligence.source})`
                          }
                        >
                          {tierBadge.label}
                        </span>
                        {model.isLive ? (
                          <span className="shrink-0 rounded bg-[var(--color-primary)]/15 px-1.5 py-0.5 text-[9px] font-bold leading-none text-[var(--color-primary)]">
                            LIVE
                          </span>
                        ) : null}
                        {model.isCustom ? (
                          <span className="shrink-0 rounded bg-[var(--status-warn-bg)] px-1.5 py-0.5 text-[9px] font-bold leading-none text-[var(--status-warn)]">
                            MANUAL
                          </span>
                        ) : null}
                        {model.isFree ? (
                          <span className="shrink-0 rounded bg-[var(--status-ok-bg)] px-1.5 py-0.5 text-[9px] font-bold leading-none text-[var(--status-ok)]">
                            FREE
                          </span>
                        ) : null}
                      </div>
                      {hasSubtitle ? (
                        <p
                          className="truncate pl-5 text-[11px] text-[var(--color-text-muted)]"
                          title={model.name}
                        >
                          {model.name}
                        </p>
                      ) : null}
                    </div>

                    <div className="border-t border-[var(--color-border-subtle)] pt-2">
                      <div className="flex items-center justify-between gap-2">
                        <div className="min-w-0 flex-1">
                          {isTestingModel ? (
                            <span className="inline-flex items-center gap-1 font-mono text-[11px] text-[var(--color-text-muted)]">
                              <RefreshCw className="h-3 w-3 animate-spin" aria-hidden="true" />
                              <span>{isId ? "Menguji..." : "Testing..."}</span>
                            </span>
                          ) : modelTestResult ? (
                            <Badge
                              variant={
                                modelTestResult.status === "ok"
                                  ? "success"
                                  : modelTestResult.status === "slow"
                                    ? "warning"
                                    : "destructive"
                              }
                              className="tabular-nums"
                            >
                              {modelTestResult.status === "ok"
                                ? `OK · ${modelTestResult.latencyMs ?? 0}ms`
                                : modelTestResult.status === "slow"
                                  ? `SLOW · ${modelTestResult.latencyMs ?? 0}ms`
                                  : `FAIL${modelTestResult.statusCode ? ` ${modelTestResult.statusCode}` : ""}`}
                            </Badge>
                          ) : (
                            <span className="font-mono text-[10.5px] text-[var(--color-text-subtle)]">
                              {canTestModels
                                ? isId
                                  ? "Siap diuji"
                                  : "Ready"
                                : isId
                                  ? "Butuh koneksi aktif"
                                  : "Needs active connection"}
                            </span>
                          )}
                        </div>

                        <div className="flex shrink-0 items-center gap-1">
                          <Button
                            type="button"
                            size="xs"
                            variant="secondary"
                            disabled={!canTestModels || isTestingModel}
                            onClick={() => void handleTestModel(model.id)}
                            aria-label={`Test ${fullModel}`}
                            title={
                              canTestModels
                                ? isId
                                  ? `Tes model ${fullModel}`
                                  : `Test model ${fullModel}`
                                : isId
                                  ? "Tambahkan atau aktifkan koneksi terlebih dahulu untuk mengetes model"
                                  : "Add or enable a connection first to test this model"
                            }
                            className="h-7 px-2 text-[11px]"
                          >
                            {isTestingModel ? (
                              <RefreshCw className="h-3 w-3 animate-spin" aria-hidden="true" />
                            ) : (
                              <Activity className="h-3 w-3" aria-hidden="true" />
                            )}
                            <span>{isId ? "Tes" : "Test"}</span>
                          </Button>
                          <Button
                            type="button"
                            size="xs"
                            variant="secondary"
                            disabled={!canTestModels || isIdentifyingModel}
                            onClick={() => void handleIdentifyModel(model.id)}
                            aria-label={`Unmask ${fullModel}`}
                            title={
                              canTestModels
                                ? isId
                                  ? `Cek identitas asli model ${fullModel} (deteksi masking)`
                                  : `Probe true identity of ${fullModel} (detect masking)`
                                : isId
                                  ? "Tambahkan atau aktifkan koneksi terlebih dahulu untuk cek identitas model"
                                  : "Add or enable a connection first to inspect model identity"
                            }
                            className="h-7 px-2 text-[11px]"
                          >
                            {isIdentifyingModel ? (
                              <RefreshCw className="h-3 w-3 animate-spin" aria-hidden="true" />
                            ) : (
                              <ScanEye className="h-3 w-3" aria-hidden="true" />
                            )}
                            <span>{isId ? "Cek Asli" : "Unmask"}</span>
                          </Button>
                          <button
                            type="button"
                            onClick={() => void handleCopyModel(model.id)}
                            title={isCopied ? "Copied!" : `Copy ${fullModel}`}
                            aria-label={`Copy ${fullModel}`}
                            className="inline-flex h-7 w-7 shrink-0 cursor-pointer items-center justify-center rounded border border-[var(--color-border)] text-[var(--color-text-muted)] transition-colors hover:border-[var(--border-strong)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-primary)]"
                          >
                            {isCopied ? (
                              <Check
                                className="h-3.5 w-3.5 text-[var(--status-ok)]"
                                aria-hidden="true"
                              />
                            ) : (
                              <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                            )}
                          </button>
                          {model.isCustom && model.isSaved ? (
                            <button
                              type="button"
                              disabled={removingModelId !== null}
                              onClick={() => void handleRemoveCustomModel(model.id)}
                              title={
                                isId
                                  ? `Hapus model manual ${model.id}`
                                  : `Remove custom model ${model.id}`
                              }
                              aria-label={`Remove ${fullModel}`}
                              className="inline-flex h-7 w-7 shrink-0 cursor-pointer items-center justify-center rounded border border-[var(--color-border)] text-[var(--color-text-muted)] transition-colors hover:border-[var(--status-error)]/40 hover:bg-[var(--status-error-bg)] hover:text-[var(--status-error)] disabled:opacity-40"
                            >
                              {removingModelId === model.id ? (
                                <RefreshCw className="h-3 w-3 animate-spin" aria-hidden="true" />
                              ) : (
                                <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                              )}
                            </button>
                          ) : null}
                        </div>
                      </div>

                      <div
                        className={cn(
                          "grid transition-[grid-template-rows,opacity] duration-200 ease-out motion-reduce:transition-none",
                          modelTestResult?.status === "error" && modelTestResult.error
                            ? "grid-rows-[1fr] opacity-100"
                            : "grid-rows-[0fr] opacity-0",
                        )}
                      >
                        <div className="overflow-hidden">
                          {modelTestResult?.status === "error" && modelTestResult.error ? (
                            <p
                              className="pt-1.5 truncate font-mono text-[10.5px] text-[var(--status-error)]"
                              title={modelTestResult.error}
                            >
                              {modelTestResult.error}
                            </p>
                          ) : null}
                        </div>
                      </div>

                      <div
                        className={cn(
                          "grid transition-[grid-template-rows,opacity] duration-200 ease-out motion-reduce:transition-none",
                          isIdentityOpen
                            ? "grid-rows-[1fr] opacity-100"
                            : "grid-rows-[0fr] opacity-0",
                        )}
                        aria-hidden={!isIdentityOpen}
                      >
                        <div className="overflow-hidden">
                          {identityResult ? (
                            <div className="mt-1.5 space-y-1.5 rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/60 p-2 text-[10.5px]">
                              <div className="flex items-center justify-between gap-1.5">
                                <div className="flex min-w-0 flex-wrap items-center gap-1.5">
                                  <Badge
                                    variant={
                                      identityResult.verdict === "verified"
                                        ? "success"
                                        : identityResult.verdict === "suspected_masked"
                                          ? "destructive"
                                          : identityResult.verdict === "error"
                                            ? "destructive"
                                            : "warning"
                                    }
                                  >
                                    {identityResult.verdict === "verified"
                                      ? isId
                                        ? `Asli · ${identityResult.confidence}`
                                        : `Verified · ${identityResult.confidence}`
                                      : identityResult.verdict === "suspected_masked"
                                        ? `Masked → ${identityResult.detectedModel || identityResult.detectedFamily || "?"} · ${identityResult.confidence}`
                                        : identityResult.verdict === "error"
                                          ? isId
                                            ? "Cek Gagal"
                                            : "Probe Failed"
                                          : isId
                                            ? "Tidak Pasti"
                                            : "Inconclusive"}
                                  </Badge>
                                  {identityResult.reportedModel ? (
                                    <span
                                      className="truncate font-mono text-[10px] text-[var(--color-text-muted)]"
                                      title={`response.model: ${identityResult.reportedModel}`}
                                    >
                                      model: {identityResult.reportedModel}
                                    </span>
                                  ) : null}
                                </div>
                                <button
                                  type="button"
                                  tabIndex={isIdentityOpen ? 0 : -1}
                                  onClick={() =>
                                    setOpenIdentityModelIds((prev) => ({
                                      ...prev,
                                      [model.id]: false,
                                    }))
                                  }
                                  aria-label={
                                    isId ? "Tutup hasil cek asli" : "Dismiss identity result"
                                  }
                                  className="shrink-0 cursor-pointer rounded p-0.5 text-[var(--color-text-muted)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]"
                                >
                                  <X className="h-3 w-3" aria-hidden="true" />
                                </button>
                              </div>

                              {identityResult.error ? (
                                <p
                                  className="truncate font-mono text-[10px] text-[var(--status-error)]"
                                  title={identityResult.error}
                                >
                                  {identityResult.error}
                                </p>
                              ) : null}

                              {identityResult.signals.length > 0 ? (
                                <ul className="space-y-0.5 text-[10px] text-[var(--color-text-secondary)]">
                                  {identityResult.signals.slice(0, 3).map((sig) => (
                                    <li key={sig} className="truncate" title={sig}>
                                      • {sig}
                                    </li>
                                  ))}
                                </ul>
                              ) : null}

                              {identityResult.selfReport ? (
                                <p
                                  className="line-clamp-2 font-mono text-[10px] text-[var(--color-text-muted)]"
                                  title={identityResult.selfReport}
                                >
                                  &ldquo;{identityResult.selfReport}&rdquo;
                                </p>
                              ) : null}

                              <div
                                className="flex items-start gap-1 border-t border-[var(--color-border-subtle)] pt-1 text-[9.5px] leading-snug text-[var(--color-text-muted)]"
                                title={identityResult.limitationNote}
                              >
                                <Info className="mt-0.5 h-3 w-3 shrink-0" aria-hidden="true" />
                                <span>
                                  {isId
                                    ? "Batasan: Proxy yang menulis ulang seluruh metadata/header sekaligus menyuntikkan system prompt penyamaran tetap dapat menyamarkan model asli; gunakan bukti di atas sebagai indikasi."
                                    : identityResult.limitationNote}
                                </span>
                              </div>
                            </div>
                          ) : null}
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Scoped Add API Key, OAuth & Edit Modals (lazy-loaded on first open, kept mounted for smooth close animation) */}
      {modalsActivated ? (
        <Suspense fallback={null}>
          <AddApiKeyModal
            open={showAddApiKeyModal}
            onClose={() => setShowAddApiKeyModal(false)}
            initialProviderId={entry.id}
            allEntries={allEntries}
            isId={isId}
            lockProvider
            onCreated={(msg) => setFeedbackBanner(msg)}
          />

          <OAuthWizardModal
            open={showOAuthModal}
            onClose={() => setShowOAuthModal(false)}
            initialProviderId={entry.id}
            allEntries={allEntries}
            isId={isId}
            lockProvider
            onConnected={(msg) => setFeedbackBanner(msg)}
            onSwitchToApiKey={() => {
              setShowOAuthModal(false);
              setShowAddApiKeyModal(true);
            }}
          />

          <EditConnectionModal
            open={editingConnection !== null}
            connection={editingConnection ?? lastEditingConnectionRef.current}
            isId={isId}
            onClose={() => setEditingConnection(null)}
            onUpdated={(msg) => setFeedbackBanner(msg)}
          />

          <AddCustomModelModal
            open={showAddModelModal}
            onClose={() => setShowAddModelModal(false)}
            providerId={entry.id}
            providerName={entry.name}
            providerAlias={entry.alias}
            connections={sortedConnections}
            isId={isId}
            onAdded={(msg) => setFeedbackBanner(msg)}
          />
        </Suspense>
      ) : null}
    </div>
  );
}

export function ProviderDetailView() {
  const [searchParams] = useSearchParams();
  const providerId = (searchParams.get("provider") ?? "").trim();
  const connectionId = (searchParams.get("id") ?? "").trim();

  if (providerId && !connectionId) {
    return <ProviderStudioDetail key={providerId} providerId={providerId} />;
  }

  return <ConnectionDetailEditor connectionId={connectionId} />;
}

function ConnectionDetailEditor({ connectionId }: { connectionId: string }) {
  const navigate = useNavigate();

  const connectionQuery = useProviderConnection(connectionId);
  const updateMutation = useUpdateProviderConnectionMutation();
  const deleteMutation = useDeleteProviderConnectionMutation();
  const testMutation = useTestProviderConnectionMutation();

  const [saveFeedback, setSaveFeedback] = useState<string | null>(null);
  const [liveTestResult, setLiveTestResult] = useState<TestProviderConnectionResponse | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);
  const [showAccessToken, setShowAccessToken] = useState(false);
  const [showRefreshToken, setShowRefreshToken] = useState(false);

  const editForm = useForm<UpdateConnectionFormInput, unknown, UpdateConnectionFormOutput>({
    resolver: zodResolver(updateConnectionSchema),
    defaultValues: DEFAULT_UPDATE_CONNECTION_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const {
    reset,
    formState: { errors, isDirty },
  } = editForm;

  useEffect(() => {
    if (connectionQuery.data && !isDirty) {
      reset(connectionToEditFormValues(connectionQuery.data));
    }
  }, [connectionQuery.data, isDirty, reset]);

  if (!connectionId) {
    return (
      <Card>
        <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
          <CardTitle>Provider Connection Not Specified</CardTitle>
          <Button asChild variant="outline" size="xs" className="w-full sm:w-auto">
            <Link to="/dashboard/providers/">
              <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Back to Providers</span>
            </Link>
          </Button>
        </CardHeader>
        <CardContent>
          <Alert variant="warning">
            <AlertCircle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
              aria-hidden="true"
            />
            <div>
              <AlertTitle>Missing connection query parameter</AlertTitle>
              <AlertDescription>
                Open a connection from the Providers Control Studio to inspect or rotate its
                configuration.
              </AlertDescription>
            </div>
          </Alert>
        </CardContent>
      </Card>
    );
  }

  const connection = connectionQuery.data;

  const onSaveSubmit = (values: UpdateConnectionFormOutput) => {
    if (!connection) {
      return;
    }
    setSaveFeedback(null);
    updateMutation.reset();

    updateMutation.mutate(
      {
        id: connection.id,
        body: toUpdateConnectionRequest(values, connection.auth_type),
      },
      {
        onSuccess: (updated) => {
          reset(connectionToEditFormValues(updated));
          setSaveFeedback(`Saved changes to "${updated.name}".`);
        },
      },
    );
  };

  const handleRunTest = () => {
    if (!connection || testMutation.isPending) {
      return;
    }
    testMutation.reset();
    testMutation.mutate(connection.id, {
      onSuccess: (result) => {
        setLiveTestResult(result);
      },
    });
  };

  const handleConfirmDelete = () => {
    if (!connection || deleteMutation.isPending) {
      return;
    }
    deleteMutation.mutate(connection.id, {
      onSuccess: () => {
        void navigate("/dashboard/providers/");
      },
    });
  };

  return (
    <div className="min-w-0 space-y-6">
      <section className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between sm:gap-4">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <Button asChild variant="ghost" size="xs" className="-ml-2 shrink-0">
              <Link
                to={
                  connection
                    ? `/dashboard/providers/detail/?provider=${encodeURIComponent(connection.provider_id)}`
                    : "/dashboard/providers/"
                }
              >
                <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{connection ? connection.provider_id : "Providers Studio"}</span>
              </Link>
            </Button>
            <span className="break-all font-mono text-xs text-[var(--text-muted)]">
              / {connectionId}
            </span>
          </div>
          <h2 className="mt-1 break-words text-xl font-semibold tracking-tight text-[var(--text-primary)] sm:text-2xl">
            {connection ? connection.name : "Provider Connection Detail"}
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Configure routing priority, allowed model list, upstream headers, and rotate vaulted
            credentials.
          </p>
        </div>

        <Button
          type="button"
          variant="outline"
          size="sm"
          className="w-full shrink-0 sm:w-auto"
          onClick={() => void connectionQuery.refetch()}
          disabled={connectionQuery.isFetching}
        >
          <RefreshCw className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
          <span>{connectionQuery.isFetching ? "Refreshing" : "Refresh"}</span>
        </Button>
      </section>

      {connectionQuery.isError ? (
        <Alert variant="destructive">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div>
            <AlertTitle>Unable to load provider connection</AlertTitle>
            <AlertDescription className="font-mono">
              {connectionQuery.error.message}
            </AlertDescription>
          </div>
        </Alert>
      ) : connectionQuery.isPending || !connection ? (
        <Card>
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading connection details...
          </CardContent>
        </Card>
      ) : (
        <>
          <Card>
            <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex min-w-0 items-center gap-2">
                <Server
                  className="h-4 w-4 shrink-0 text-[var(--accent-primary)]"
                  aria-hidden="true"
                />
                <CardTitle>Connection Telemetry &amp; Connectivity</CardTitle>
              </div>
              <Button
                type="button"
                size="xs"
                className="w-full shrink-0 sm:w-auto"
                onClick={handleRunTest}
                disabled={testMutation.isPending}
              >
                <Activity className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{testMutation.isPending ? "Testing..." : "Run Connectivity Test"}</span>
              </Button>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Status &amp; Protocol
                  </p>
                  <div className="mt-2 flex flex-wrap gap-1.5">
                    <Badge variant={connection.enabled ? "success" : "warning"}>
                      {connection.enabled ? "Enabled" : "Disabled"}
                    </Badge>
                    <Badge variant="outline">{connection.protocol}</Badge>
                  </div>
                </div>

                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Adapter &amp; Auth
                  </p>
                  <div className="mt-2 flex flex-wrap gap-1.5">
                    <Badge variant="secondary">{connection.provider_id}</Badge>
                    <Badge variant={connection.auth_type === "oauth" ? "default" : "secondary"}>
                      {connection.auth_type === "oauth" ? "OAuth" : "API Key"}
                    </Badge>
                    <Badge variant={connection.has_credentials ? "success" : "destructive"}>
                      {connection.has_credentials ? "Vaulted" : "Missing"}
                    </Badge>
                  </div>
                </div>

                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Last Test
                  </p>
                  <div className="mt-2 flex items-center gap-2">
                    <Badge
                      variant={
                        connection.last_test_status === "ok"
                          ? "success"
                          : connection.last_test_status === "error"
                            ? "destructive"
                            : "outline"
                      }
                    >
                      {connection.last_test_status || "Untested"}
                    </Badge>
                  </div>
                  <p className="mt-1.5 font-mono text-xs text-[var(--text-muted)]">
                    {formatTimestamp(connection.last_tested_at)}
                  </p>
                </div>

                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Timestamps
                  </p>
                  <p className="mt-1.5 font-mono text-xs text-[var(--text-muted)]">
                    Updated: {formatTimestamp(connection.updated_at)}
                  </p>
                  {connection.token_expires_at ? (
                    <p className="mt-0.5 font-mono text-xs text-[var(--text-secondary)]">
                      Token Exp: {formatTimestamp(connection.token_expires_at)}
                    </p>
                  ) : null}
                </div>
              </div>

              {testMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Connectivity test request failed</AlertTitle>
                    <AlertDescription className="font-mono">
                      {testMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {liveTestResult ? (
                <Alert variant={liveTestResult.ok ? "default" : "destructive"}>
                  {liveTestResult.ok ? (
                    <CheckCircle2
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                      aria-hidden="true"
                    />
                  ) : (
                    <AlertCircle
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                      aria-hidden="true"
                    />
                  )}
                  <div>
                    <AlertTitle>
                      {liveTestResult.ok
                        ? `Live Test Succeeded (HTTP ${liveTestResult.status_code} · ${liveTestResult.latency_ms}ms)`
                        : `Live Test Failed (${liveTestResult.status_code ? `HTTP ${liveTestResult.status_code}` : "Network Error"} · ${liveTestResult.latency_ms}ms)`}
                    </AlertTitle>
                    {liveTestResult.error ? (
                      <AlertDescription className="font-mono">
                        {liveTestResult.error}
                      </AlertDescription>
                    ) : null}
                  </div>
                </Alert>
              ) : connection.last_test_error ? (
                <Alert variant="warning">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Last recorded test error</AlertTitle>
                    <AlertDescription className="font-mono">
                      {connection.last_test_error}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}
            </CardContent>
          </Card>

          <EmberRule />

          <Card>
            <CardHeader className="flex-col items-start gap-2 sm:flex-row sm:items-center sm:justify-between">
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  <KeyRound
                    className="h-4 w-4 shrink-0 text-[var(--accent-primary)]"
                    aria-hidden="true"
                  />
                  <CardTitle>Configuration &amp; Credential Rotation</CardTitle>
                </div>
                <CardDescription className="mt-0.5">
                  Leave credential fields empty to keep existing encrypted secrets untouched.
                </CardDescription>
              </div>
              <span className="break-all font-mono text-xs text-[var(--text-muted)]">
                PATCH /api/providers/connections/{connection.id}
              </span>
            </CardHeader>
            <CardContent>
              <form onSubmit={editForm.handleSubmit(onSaveSubmit)} noValidate className="space-y-5">
                {updateMutation.error ? (
                  <Alert variant="destructive">
                    <AlertCircle
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                      aria-hidden="true"
                    />
                    <div>
                      <AlertTitle>Failed to update connection</AlertTitle>
                      <AlertDescription className="font-mono">
                        {updateMutation.error.message}
                      </AlertDescription>
                    </div>
                  </Alert>
                ) : null}

                {saveFeedback ? (
                  <Alert>
                    <CheckCircle2
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                      aria-hidden="true"
                    />
                    <div>
                      <AlertTitle>Connection Updated</AlertTitle>
                      <AlertDescription>{saveFeedback}</AlertDescription>
                    </div>
                  </Alert>
                ) : null}

                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <div>
                    <Label htmlFor="detail-name">Connection Name</Label>
                    <Input
                      id="detail-name"
                      type="text"
                      aria-invalid={Boolean(errors.name)}
                      className="mt-1.5"
                      {...editForm.register("name")}
                    />
                    {errors.name ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.name.message}
                      </p>
                    ) : null}
                  </div>

                  <div>
                    <Label htmlFor="detail-priority">Routing Priority (lower runs first)</Label>
                    <Input
                      id="detail-priority"
                      type="number"
                      min={0}
                      max={100000}
                      step={1}
                      aria-invalid={Boolean(errors.priority)}
                      className="mt-1.5 tabular-nums"
                      {...editForm.register("priority")}
                    />
                    {errors.priority ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.priority.message}
                      </p>
                    ) : null}
                  </div>

                  <div className="md:col-span-2">
                    <Label htmlFor="detail-base-url">Upstream Base URL</Label>
                    <Input
                      id="detail-base-url"
                      type="url"
                      aria-invalid={Boolean(errors.base_url)}
                      className="mt-1.5"
                      {...editForm.register("base_url")}
                    />
                    {errors.base_url ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.base_url.message}
                      </p>
                    ) : null}
                  </div>

                  <div className="md:col-span-2">
                    <Label htmlFor="detail-models">
                      Allowed / Advertised Models (comma-separated)
                    </Label>
                    <Input
                      id="detail-models"
                      type="text"
                      aria-invalid={Boolean(errors.models)}
                      className="mt-1.5"
                      {...editForm.register("models")}
                    />
                    {errors.models ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.models.message}
                      </p>
                    ) : null}
                  </div>

                  <div className="md:col-span-2">
                    <Label htmlFor="detail-headers">
                      Custom Upstream Headers (one Key: Value per line)
                    </Label>
                    <textarea
                      id="detail-headers"
                      aria-invalid={Boolean(errors.headers)}
                      className={`${TEXTAREA_CLASS_NAME} mt-1.5`}
                      {...editForm.register("headers")}
                    />
                    {errors.headers ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.headers.message}
                      </p>
                    ) : null}
                  </div>
                </div>

                {connection.auth_type === "api_key" ? (
                  <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
                    <Label htmlFor="detail-api-key">Rotate API Key (optional)</Label>
                    <div className="relative mt-1.5">
                      <Input
                        id="detail-api-key"
                        type={showApiKey ? "text" : "password"}
                        autoComplete="off"
                        placeholder="Leave blank to keep existing vaulted API key"
                        className="pr-10"
                        {...editForm.register("api_key")}
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon-sm"
                        onClick={() => setShowApiKey((prev) => !prev)}
                        aria-label={showApiKey ? "Hide API key" : "Show API key"}
                        className="absolute top-1.5 right-1.5"
                      >
                        {showApiKey ? (
                          <EyeOff className="h-4 w-4" aria-hidden="true" />
                        ) : (
                          <Eye className="h-4 w-4" aria-hidden="true" />
                        )}
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="space-y-3 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <span className="text-sm font-medium text-[var(--text-primary)]">
                        OAuth Token Rotation (optional)
                      </span>
                      <Link
                        to="/dashboard/providers/"
                        className="font-mono text-xs text-[var(--accent-text)] hover:underline"
                      >
                        Re-authorize via OAuth Wizard
                      </Link>
                    </div>
                    <div className="grid grid-cols-1 gap-3.5 md:grid-cols-2">
                      <div>
                        <Label htmlFor="detail-access-token">New Access Token</Label>
                        <div className="relative mt-1.5">
                          <Input
                            id="detail-access-token"
                            type={showAccessToken ? "text" : "password"}
                            autoComplete="off"
                            placeholder="Leave blank to keep existing"
                            className="pr-10"
                            {...editForm.register("access_token")}
                          />
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => setShowAccessToken((prev) => !prev)}
                            aria-label={showAccessToken ? "Hide access token" : "Show access token"}
                            className="absolute top-1.5 right-1.5"
                          >
                            {showAccessToken ? (
                              <EyeOff className="h-4 w-4" aria-hidden="true" />
                            ) : (
                              <Eye className="h-4 w-4" aria-hidden="true" />
                            )}
                          </Button>
                        </div>
                      </div>

                      <div>
                        <Label htmlFor="detail-refresh-token">New Refresh Token</Label>
                        <div className="relative mt-1.5">
                          <Input
                            id="detail-refresh-token"
                            type={showRefreshToken ? "text" : "password"}
                            autoComplete="off"
                            placeholder="Leave blank to keep existing"
                            className="pr-10"
                            {...editForm.register("refresh_token")}
                          />
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => setShowRefreshToken((prev) => !prev)}
                            aria-label={
                              showRefreshToken ? "Hide refresh token" : "Show refresh token"
                            }
                            className="absolute top-1.5 right-1.5"
                          >
                            {showRefreshToken ? (
                              <EyeOff className="h-4 w-4" aria-hidden="true" />
                            ) : (
                              <Eye className="h-4 w-4" aria-hidden="true" />
                            )}
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                <label
                  htmlFor="detail-enabled"
                  className="flex cursor-pointer items-center justify-between rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-3"
                >
                  <span className="text-sm text-[var(--text-primary)]">
                    Connection enabled for gateway routing
                  </span>
                  <input
                    id="detail-enabled"
                    type="checkbox"
                    className="h-4 w-4 accent-[var(--accent-primary)]"
                    {...editForm.register("enabled")}
                  />
                </label>

                <div className="flex flex-col-reverse gap-2 pt-1 sm:flex-row sm:items-center sm:justify-end sm:gap-3">
                  <Button
                    type="button"
                    variant="secondary"
                    size="sm"
                    className="w-full sm:w-auto"
                    onClick={() => {
                      updateMutation.reset();
                      setSaveFeedback(null);
                      reset(connectionToEditFormValues(connection));
                    }}
                  >
                    Discard Changes
                  </Button>
                  <Button
                    type="submit"
                    size="sm"
                    className="w-full sm:w-auto"
                    disabled={updateMutation.isPending}
                  >
                    <Save className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{updateMutation.isPending ? "Saving..." : "Save Changes"}</span>
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          <Card className="border-[var(--status-error-border)]">
            <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="min-w-0">
                <CardTitle>Delete Connection</CardTitle>
                <CardDescription className="mt-0.5">
                  Permanently remove this connection and purge its vaulted credentials.
                </CardDescription>
              </div>
              {confirmDelete ? (
                <div className="grid w-full grid-cols-2 gap-2 sm:flex sm:w-auto sm:items-center">
                  <Button
                    type="button"
                    variant="destructive"
                    size="xs"
                    className="w-full sm:w-auto"
                    onClick={handleConfirmDelete}
                    disabled={deleteMutation.isPending}
                  >
                    {deleteMutation.isPending ? "Deleting..." : "Confirm Delete"}
                  </Button>
                  <Button
                    type="button"
                    variant="secondary"
                    size="xs"
                    className="w-full sm:w-auto"
                    onClick={() => setConfirmDelete(false)}
                    disabled={deleteMutation.isPending}
                  >
                    Cancel
                  </Button>
                </div>
              ) : (
                <Button
                  type="button"
                  variant="destructive"
                  size="xs"
                  className="w-full shrink-0 sm:w-auto"
                  onClick={() => setConfirmDelete(true)}
                >
                  <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Delete Connection</span>
                </Button>
              )}
            </CardHeader>
            {deleteMutation.error ? (
              <CardContent>
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to delete connection</AlertTitle>
                    <AlertDescription className="font-mono">
                      {deleteMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              </CardContent>
            ) : null}
          </Card>
        </>
      )}
    </div>
  );
}
