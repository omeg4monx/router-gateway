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
