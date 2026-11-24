export type CatalogTab = "models" | "combos";

export type ProviderHealthState = "healthy" | "degraded" | "down" | "untested";

export type TestResultStatus = "ok" | "slow" | "error";

export type TestErrorClass = "rate-limited" | "quota" | "timeout" | "other";

export interface CatalogModelRow {
  providerId: string;
  provider: string;
  id: string;
  qualifiedId: string;
  name: string;
  type: string;
  subtype?: string;
  custom?: boolean;
  free?: boolean;
  capabilities?: Record<string, boolean>;
  context_length?: number;
  max_output_tokens?: number;
  input_modalities?: string[];
  output_modalities?: string[];
  supported_endpoints?: string[];
  hidden?: boolean;
}

export interface CatalogComboMember {
  model: string;
  kind?: "combo-ref" | "provider-wildcard" | "model";
  label?: string;
}

export interface CatalogComboRow {
  id: string;
  name: string;
  displayName: string;
  strategy: string;
  description?: string;
  models: CatalogComboMember[];
  memberCount: number;
  status: "active" | "paused";
  enabled: boolean;
}

export interface CatalogTestResult {
  id: string;
  targetType: "model" | "combo";
  providerId?: string;
  modelId?: string;
  comboName?: string;
  status: TestResultStatus;
  latencyMs?: number;
  statusCode?: number;
  error?: string;
  errorClass?: TestErrorClass;
  testedAt: number;
  resolvedBy?: string;
}

export interface ModelCatalogFilters {
  query: string;
  providerId: string;
  type: string;
  subtype: string;
  capability: string;
  pricing: "all" | "free" | "paid";
  providerHealth: "all" | ProviderHealthState;
  testResult: "all" | TestResultStatus | "untested";
  visibility: "all" | "visible" | "hidden";
  minContextLength?: number;
  minMaxOutputTokens?: number;
}

export interface ComboCatalogFilters {
  query: string;
  strategy: string;
  status: "all" | "active" | "paused";
  testResult: "all" | TestResultStatus | "untested";
  minMembers?: number;
  maxMembers?: number;
}

export type ModelSortField =
  | "id"
  | "provider"
  | "type"
  | "context_length"
  | "max_output_tokens"
  | "testResult";

export type ComboSortField = "name" | "strategy" | "status" | "memberCount" | "testResult";

export type SortDirection = "asc" | "desc";

export interface BulkTestProgress {
  kind: "models" | "combos";
  completed: number;
  total: number;
  cancelled: boolean;
}

export interface ApiKeyAccessItem {
  id: string;
  name: string;
  prefix: string;
  disabled: boolean;
  allowedModels: string[];
  modelAccessMode: "all" | "restricted";
}
