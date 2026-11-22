"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import { Button, ProgressCard } from "@/shared/components";
import { cn } from "@/shared/utils/cn";
import { useNotificationStore } from "@/store/notificationStore";
import ProvidersHeader from "./components/ProvidersHeader";
import ProvidersToolbar from "./components/ProvidersToolbar";
import ProviderCard from "./components/ProviderCard";
import TestResultsModal from "./components/TestResultsModal";
import AddCompatibleModal from "./components/AddCompatibleModal";
import {
  matchesStatusFilter,
  matchesProviderSearch,
  buildProviderEntries,
  computeFleetTelemetry,
  countStatusBuckets,
  countCategoryBuckets,
  partitionConfiguredAndCatalog,
  sortFilteredEntries,
  countActiveConnectionsByMode,
} from "./utils";

const APIKEY_INITIAL_VISIBLE = 20;

const CATEGORY_SECTION_META = {
  custom: {
    title: "Custom Endpoints",
    icon: "terminal",
    testMode: "compatible",
    testLabel: "Test All",
    testTitle: "Test all active custom endpoint connections",
  },
  oauth: {
    title: "OAuth Providers",
    icon: "verified_user",
    testMode: "oauth",
    testLabel: "Test All",
    testTitle: "Test all active OAuth connections",
  },
  apikey: {
    title: "API Key Providers",
    icon: "key",
    testMode: "apikey",
    testLabel: "Test All",
    testTitle: "Test all active API key connections",
  },
  free: {
    title: "Free Tier & Web Session Providers",
    icon: "public",
    testMode: "free",
    testLabel: "Test Free OAuth",
    testTitle: "Test active Free OAuth provider connections",
  },
};

function ProvidersLoadingSkeleton() {
  return (
    <div className="flex min-w-0 flex-col gap-5 animate-pulse">
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="size-9 rounded-xl bg-surface-2" />
          <div className="flex flex-col gap-1.5">
            <div className="h-5 w-28 rounded bg-surface-2" />
            <div className="h-3 w-64 rounded bg-surface-2" />
          </div>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-8 w-32 rounded-lg bg-surface-2" />
          <div className="h-8 w-24 rounded-lg bg-surface-2" />
          <div className="h-8 w-28 rounded-lg bg-surface-2" />
        </div>
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
        {Array.from({ length: 4 }).map((_, i) => (
          <div
            key={i}
            className="h-[66px] rounded-xl border border-border bg-surface p-3.5"
          />
        ))}
      </div>

      <div className="h-24 rounded-2xl border border-border bg-surface p-3.5" />

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {Array.from({ length: 8 }).map((_, i) => (
          <div
            key={i}
            className="h-[118px] rounded-xl border border-border bg-surface p-4"
          />
        ))}
      </div>
    </div>
  );
}

function SectionHeader({
  icon,
  title,
  count,
  subtitle,
  action,
}) {
  return (
    <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between border-b border-border-subtle pb-2">
      <div className="flex items-center gap-2 min-w-0">
        <span className="material-symbols-outlined text-primary text-[18px] shrink-0">
          {icon}
        </span>
        <h2 className="text-sm sm:text-[15px] font-semibold text-text-main truncate">
          {title}
        </h2>
        <span className="text-xs font-mono tabular-nums text-text-muted">
          ({count})
        </span>
        {subtitle && (
          <span className="hidden md:inline text-xs text-text-muted truncate">
            {subtitle}
          </span>
        )}
      </div>
      {action && <div className="flex items-center gap-2 shrink-0">{action}</div>}
    </div>
  );
}

export default function ProvidersPage() {
  const [connections, setConnections] = useState([]);
  const [providerNodes, setProviderNodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState("");
  const [showAllApikey, setShowAllApikey] = useState(false);
  const [addModalVariant, setAddModalVariant] = useState(null);
  const [testingState, setTestingState] = useState(null);
  const [lastTestRequest, setLastTestRequest] = useState(null);
  const [testResults, setTestResults] = useState(null);
  const [togglingProviders, setTogglingProviders] = useState(() => new Set());
  const [statusFilter, setStatusFilter] = useState("all");
  const [localSearch, setLocalSearch] = useState("");
  const [categoryTab, setCategoryTab] = useState("all");
  const notify = useNotificationStore();

  const fetchData = useCallback(async ({ silent = false } = {}) => {
    if (!silent) {
      setLoading(true);
      setFetchError("");
    }
    try {
      const [connectionsRes, nodesRes] = await Promise.all([
        fetch("/api/providers", { cache: "no-store" }),
        fetch("/api/provider-nodes", { cache: "no-store" }),
      ]);
      const [connectionsData, nodesData] = await Promise.all([
        connectionsRes.json().catch(() => ({})),
        nodesRes.json().catch(() => ({})),
      ]);

      if (!connectionsRes.ok || !nodesRes.ok) {
        const msg =
          connectionsData?.error ||
          nodesData?.error ||
          "Failed to load providers data";
        setFetchError(msg);
        return;
      }

      setConnections(connectionsData.connections || []);
      setProviderNodes(nodesData.nodes || []);
      setFetchError("");
    } catch {
      setFetchError("Network error while loading providers");
    } finally {
      if (!silent) setLoading(false);
    }
  }, []);

  const refreshProviders = useCallback(async () => {
    try {
      const res = await fetch("/api/providers", { cache: "no-store" });
      if (!res.ok) return;
      const data = await res.json();
      setConnections(data.connections || []);
    } catch {
      // keep existing state on background refresh failure
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const allEntries = useMemo(
    () => buildProviderEntries(providerNodes, connections),
    [providerNodes, connections],
  );

  const telemetry = useMemo(
    () => computeFleetTelemetry(allEntries, connections),
    [allEntries, connections],
  );

  const activeAllCount = useMemo(
    () => countActiveConnectionsByMode(connections, "all"),
    [connections],
  );

  const matchesEntrySearch = useCallback(
    (entry) =>
      matchesProviderSearch(localSearch, {
        name: entry.info?.name,
        id: entry.id,
        alias: entry.alias,
        prefix: entry.prefix,
        baseUrl: entry.baseUrl,
        typeLabel: entry.typeLabel,
        category: entry.category,
      }),
    [localSearch],
  );

  // Status pill counts reflect current search + category scope
  const statusCounts = useMemo(() => {
    const scoped = allEntries.filter(
      (entry) =>
        (categoryTab === "all" || entry.category === categoryTab) &&
        matchesEntrySearch(entry),
    );
    return countStatusBuckets(scoped);
  }, [allEntries, categoryTab, matchesEntrySearch]);

  // Category tab counts reflect current search + status scope
  const categoryCounts = useMemo(() => {
    const scoped = allEntries.filter(
      (entry) =>
        matchesEntrySearch(entry) &&
        matchesStatusFilter(statusFilter, entry.stats, entry.isNoAuth),
    );
    return countCategoryBuckets(scoped);
  }, [allEntries, statusFilter, matchesEntrySearch]);

  const filteredEntries = useMemo(
    () =>
      allEntries.filter(
        (entry) =>
          (categoryTab === "all" || entry.category === categoryTab) &&
          matchesEntrySearch(entry) &&
          matchesStatusFilter(statusFilter, entry.stats, entry.isNoAuth),
      ),
    [allEntries, categoryTab, statusFilter, matchesEntrySearch],
  );

  const { configured, catalog } = useMemo(
    () => partitionConfiguredAndCatalog(allEntries),
    [allEntries],
  );

  const isSearchingOrStatusFiltered =
    Boolean(localSearch.trim()) || statusFilter !== "all";
  const hasActiveFilters =
    isSearchingOrStatusFiltered || categoryTab !== "all";

  const handleResetFilters = () => {
    setLocalSearch("");
    setStatusFilter("all");
    setCategoryTab("all");
  };

  const handleToggleProvider = async (entry, newActive) => {
    const providerId = entry.id;
    if (togglingProviders.has(providerId)) return;

    const authTypes = Array.isArray(entry.authTypes)
      ? entry.authTypes
      : [entry.authTypes];
    const matches = (c) =>
      c.provider === providerId && authTypes.includes(c.authType);
    const providerConns = connections.filter(matches);
    if (providerConns.length === 0) return;

    const previousStates = new Map(
      providerConns.map((c) => [c.id, c.isActive]),
    );

    setTogglingProviders((prev) => {
      const next = new Set(prev);
      next.add(providerId);
      return next;
    });
    setConnections((prev) =>
      prev.map((c) => (matches(c) ? { ...c, isActive: newActive } : c)),
    );

    try {
      const responses = await Promise.all(
        providerConns.map((c) =>
          fetch(`/api/providers/${c.id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ isActive: newActive }),
          }),
        ),
      );

      if (responses.some((res) => !res.ok)) {
        throw new Error("One or more connection updates failed");
      }

      notify.success(
        `${newActive ? "Enabled" : "Disabled"} ${entry.info?.name || providerId}`,
      );
    } catch {
      setConnections((prev) =>
        prev.map((c) =>
          previousStates.has(c.id)
            ? { ...c, isActive: previousStates.get(c.id) }
            : c,
        ),
      );
      notify.error(
        `Failed to ${newActive ? "enable" : "disable"} ${entry.info?.name || providerId}`,
      );
      await refreshProviders();
    } finally {
      setTogglingProviders((prev) => {
        const next = new Set(prev);
        next.delete(providerId);
        return next;
      });
    }
  };

  const handleBatchTest = async (mode, providerId = null, label = null) => {
    if (testingState) return;
    const requestInfo = { mode, providerId, label };
    setTestingState(requestInfo);
    setLastTestRequest(requestInfo);
    setTestResults(null);

    try {
      const res = await fetch("/api/providers/test-batch", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ mode, providerId }),
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        const errMsg = data?.error || "Provider test request failed";
        setTestResults({ error: errMsg });
        notify.error(errMsg);
        return;
      }

      setTestResults(data);
      if (data.summary) {
        const { passed, failed, total } = data.summary;
        if (total === 0) notify.info("No active connections to test");
        else if (failed === 0) notify.success(`All ${total} tests passed`);
        else notify.warning(`${passed}/${total} passed, ${failed} failed`);
      }
    } catch {
      setTestResults({ error: "Network error while testing connections" });
      notify.error("Provider test failed");
    } finally {
      await refreshProviders();
      setTestingState(null);
    }
  };

  const handleQuickTest = (entry) => {
    handleBatchTest("provider", entry.id, entry.info?.name || entry.id);
  };

  const renderCard = (entry) => (
    <ProviderCard
      key={entry.id}
      entry={entry}
      onToggle={handleToggleProvider}
      onQuickTest={handleQuickTest}
      isToggling={togglingProviders.has(entry.id)}
      isTesting={
        testingState?.mode === "provider" &&
        testingState?.providerId === entry.id
      }
      disableQuickTest={Boolean(testingState)}
    />
  );

  if (loading) {
    return <ProvidersLoadingSkeleton />;
  }

  const testingMessage = testingState
    ? testingState.mode === "provider"
      ? `Testing connections for ${testingState.label || testingState.providerId}...`
      : testingState.mode === "all"
        ? "Testing all active connections across every provider..."
        : `Testing all active ${testingState.label || testingState.mode} connections...`
    : "";

  const renderSectionTestButton = (mode, label, title) => {
    const activeCount = countActiveConnectionsByMode(connections, mode);
    const isRunning = testingState?.mode === mode;
    return (
      <button
        type="button"
        onClick={() => handleBatchTest(mode, null, label)}
        disabled={activeCount === 0 || Boolean(testingState)}
        title={activeCount === 0 ? "No active connections in this group" : title}
        className={cn(
          "inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1 text-xs font-medium transition-colors cursor-pointer",
          isRunning
            ? "bg-primary/15 border-primary/40 text-primary"
            : "bg-surface-2/80 border-border text-text-muted hover:text-text-main hover:border-primary/45",
          (activeCount === 0 || (Boolean(testingState) && !isRunning)) &&
            "opacity-45 cursor-not-allowed",
        )}
      >
        <span
          className={cn(
            "material-symbols-outlined text-[14px]",
            isRunning && "animate-spin",
          )}
        >
          {isRunning ? "progress_activity" : "play_arrow"}
        </span>
        <span>{isRunning ? "Testing..." : label}</span>
      </button>
    );
  };

  const renderApiKeyGridWithDisclosure = (entries) => {
    const visible = showAllApikey
      ? entries
      : entries.slice(0, APIKEY_INITIAL_VISIBLE);
    const hiddenCount = entries.length - APIKEY_INITIAL_VISIBLE;

    return (
      <>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {visible.map(renderCard)}
        </div>
        {hiddenCount > 0 && (
          <div className="flex justify-center pt-1">
            <Button
              variant="secondary"
              size="sm"
              onClick={() => setShowAllApikey((v) => !v)}
            >
              {showAllApikey
                ? "Show fewer providers"
                : `Show all ${entries.length} providers`}
            </Button>
          </div>
        )}
      </>
    );
  };

  const renderDefaultControlStudio = () => {
    const configuredErrorCount = configured.filter(
      (e) => e.stats.error > 0 && !e.stats.allDisabled,
    ).length;

    return (
      <div className="flex flex-col gap-7">
        {/* Zona 1: Your Providers (Configured & Custom Endpoints) */}
        <section className="flex flex-col gap-3">
          <SectionHeader
            icon="tune"
            title="Your Providers"
            count={configured.length}
            subtitle={
              configuredErrorCount > 0
                ? `${configuredErrorCount} needing attention`
                : "Configured providers and custom endpoints"
            }
            action={
              configured.length > 0 &&
              renderSectionTestButton(
                "all",
                "Test All",
                "Test all active provider connections",
              )
            }
          />

          {configured.length === 0 ? (
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between rounded-xl border border-dashed border-border bg-surface-2/50 px-4 py-3.5">
              <div className="flex items-center gap-3 min-w-0">
                <span className="material-symbols-outlined text-primary text-[20px] shrink-0">
                  rocket_launch
                </span>
                <div className="min-w-0">
                  <p className="text-xs font-semibold text-text-main">
                    No providers configured yet
                  </p>
                  <p className="text-xs text-text-muted">
                    Pick any provider from the catalog below to connect OAuth or API keys, or add a custom OpenAI/Anthropic endpoint.
                  </p>
                </div>
              </div>
              <Button
                size="sm"
                variant="secondary"
                icon="add"
                onClick={() => setAddModalVariant("openai")}
                className="shrink-0"
              >
                Add Custom Endpoint
              </Button>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {configured.map(renderCard)}
            </div>
          )}
        </section>

        {/* Zona 2: Explore Catalog (Unconfigured built-in providers only) */}
        <div className="flex flex-col gap-6">
          {catalog.oauth.length > 0 && (
            <section className="flex flex-col gap-3">
              <SectionHeader
                icon="verified_user"
                title="OAuth Providers"
                count={catalog.oauth.length}
                subtitle="Browser or device-code sign-in"
              />
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {catalog.oauth.map(renderCard)}
              </div>
            </section>
          )}

          {catalog.free.length > 0 && (
            <section className="flex flex-col gap-3">
              <SectionHeader
                icon="public"
                title="Free Tier & Web Session"
                count={catalog.free.length}
                subtitle="Public, free-tier, and browser-cookie providers"
              />
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {catalog.free.map(renderCard)}
              </div>
            </section>
          )}

          {catalog.apikey.length > 0 && (
            <section className="flex flex-col gap-3">
              <SectionHeader
                icon="key"
                title="API Key Providers"
                count={catalog.apikey.length}
                subtitle="Direct API key integrations"
              />
              {renderApiKeyGridWithDisclosure(catalog.apikey)}
            </section>
          )}
        </div>
      </div>
    );
  };

  const renderCategoryFocusedView = () => {
    const meta = CATEGORY_SECTION_META[categoryTab];
    const sorted = sortFilteredEntries(filteredEntries);

    return (
      <section className="flex flex-col gap-3">
        <SectionHeader
          icon={meta.icon}
          title={meta.title}
          count={sorted.length}
          action={renderSectionTestButton(
            meta.testMode,
            meta.testLabel,
            meta.testTitle,
          )}
        />

        {sorted.length === 0 && categoryTab === "custom" ? (
          <div className="flex flex-col items-center justify-center gap-3 py-10 px-4 rounded-2xl border border-dashed border-border bg-surface text-center">
            <span className="material-symbols-outlined text-[32px] text-text-muted opacity-50">
              extension
            </span>
            <div>
              <p className="text-sm font-semibold text-text-main">
                No custom endpoints configured yet
              </p>
              <p className="text-xs text-text-muted mt-1">
                Connect any self-hosted, gateway, or third-party OpenAI, Anthropic, or Moonshot-compatible API.
              </p>
            </div>
            <div className="flex flex-wrap items-center justify-center gap-2 pt-1">
              <Button
                size="sm"
                variant="secondary"
                icon="smart_toy"
                onClick={() => setAddModalVariant("openai")}
              >
                OpenAI Compatible
              </Button>
              <Button
                size="sm"
                variant="secondary"
                icon="psychology"
                onClick={() => setAddModalVariant("anthropic")}
              >
                Anthropic Compatible
              </Button>
              <Button
                size="sm"
                variant="secondary"
                icon="dark_mode"
                onClick={() => setAddModalVariant("moonshot")}
              >
                MoonshotAI Compatible
              </Button>
            </div>
          </div>
        ) : categoryTab === "apikey" ? (
          renderApiKeyGridWithDisclosure(sorted)
        ) : (
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {sorted.map(renderCard)}
          </div>
        )}
      </section>
    );
  };

  const renderFilteredResultsView = () => {
    const sorted = sortFilteredEntries(filteredEntries);

    if (sorted.length === 0) {
      return (
        <div className="flex flex-col items-center justify-center gap-2.5 py-12 px-4 rounded-2xl border border-dashed border-border bg-surface text-center">
          <span className="material-symbols-outlined text-[36px] text-text-muted opacity-40">
            search_off
          </span>
          <p className="text-sm font-semibold text-text-main">
            No providers match your current search or filters
          </p>
          <p className="text-xs text-text-muted">
            Try clearing active filters or adding a custom compatible endpoint.
          </p>
          <Button
            size="sm"
            variant="secondary"
            icon="filter_alt_off"
            onClick={handleResetFilters}
            className="mt-1"
          >
            Clear search & filters
          </Button>
        </div>
      );
    }

    return (
      <section className="flex flex-col gap-3">
        <SectionHeader
          icon="filter_list"
          title="Matching Providers"
          count={sorted.length}
          subtitle={`of ${allEntries.length} total`}
          action={
            <Button
              size="sm"
              variant="ghost"
              icon="filter_alt_off"
              onClick={handleResetFilters}
            >
              Clear filters
            </Button>
          }
        />
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {sorted.map(renderCard)}
        </div>
      </section>
    );
  };

  return (
    <div className="flex min-w-0 flex-col gap-5">
      <ProvidersHeader
        telemetry={telemetry}
        statusFilter={statusFilter}
        onSelectStatusFilter={setStatusFilter}
        activeAllCount={activeAllCount}
        testingState={testingState}
        onTestAll={() => handleBatchTest("all", null, "All Active")}
        onOpenAddModal={(variant) => setAddModalVariant(variant)}
        onCooldownCleared={refreshProviders}
      />

      {fetchError && (
        <div
          role="alert"
          className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between rounded-xl border border-danger/35 bg-danger/8 px-4 py-3 text-xs text-danger"
        >
          <div className="flex items-center gap-2 min-w-0">
            <span className="material-symbols-outlined text-[18px] shrink-0">
              error
            </span>
            <span className="font-medium break-words">{fetchError}</span>
          </div>
          <Button
            size="sm"
            variant="secondary"
            icon="refresh"
            onClick={() => fetchData()}
            className="shrink-0"
          >
            Retry
          </Button>
        </div>
      )}

      <ProvidersToolbar
        search={localSearch}
        onSearchChange={setLocalSearch}
        statusFilter={statusFilter}
        onStatusFilterChange={setStatusFilter}
        statusCounts={statusCounts}
        categoryTab={categoryTab}
        onCategoryTabChange={setCategoryTab}
        categoryCounts={categoryCounts}
        hasActiveFilters={hasActiveFilters}
        onResetFilters={handleResetFilters}
      />

      {testingState && (
        <ProgressCard
          title="Testing connections"
          message={testingMessage}
          section="Card statuses will refresh automatically when testing completes"
        />
      )}

      {isSearchingOrStatusFiltered
        ? renderFilteredResultsView()
        : categoryTab === "all"
          ? renderDefaultControlStudio()
          : renderCategoryFocusedView()}

      <AddCompatibleModal
        variant={addModalVariant || "openai"}
        isOpen={Boolean(addModalVariant)}
        onClose={() => setAddModalVariant(null)}
        onCreated={(node) => {
          setProviderNodes((prev) => [...prev, node]);
          setAddModalVariant(null);
        }}
      />

      <TestResultsModal
        results={testResults}
        onClose={() => setTestResults(null)}
        onRerun={
          lastTestRequest
            ? () =>
                handleBatchTest(
                  lastTestRequest.mode,
                  lastTestRequest.providerId,
                  lastTestRequest.label,
                )
            : null
        }
        rerunning={Boolean(testingState)}
      />
    </div>
  );
}
