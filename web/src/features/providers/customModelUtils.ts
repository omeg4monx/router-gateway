export function stripProviderModelPrefix(raw: string, providerPrefixes: string[]): string {
  let clean = raw.trim();
  if (!clean) {
    return "";
  }
  const normalizedPrefixes = providerPrefixes.map((p) => p.trim()).filter((p) => p.length > 0);
  for (const prefix of normalizedPrefixes) {
    const withSlash = `${prefix}/`;
    if (clean.startsWith(withSlash)) {
      clean = clean.slice(withSlash.length).trim();
      break;
    }
  }
  return clean;
}

export function normalizeManualModelInput(raw: string, providerPrefixes: string[]): string[] {
  const seen = new Set<string>();
  const result: string[] = [];

  for (const part of raw.split(/[,\r\n]+/)) {
    const clean = stripProviderModelPrefix(part, providerPrefixes);
    if (!clean || seen.has(clean)) {
      continue;
    }
    seen.add(clean);
    result.push(clean);
  }

  return result;
}

export interface MergeManualModelsResult {
  merged: string[];
  added: string[];
  duplicates: string[];
}

export function mergeManualModels(
  existing: string[],
  additions: string[],
): MergeManualModelsResult {
  const existingSet = new Set(existing.map((id) => id.trim()).filter(Boolean));
  const merged = Array.from(existingSet);
  const added: string[] = [];
  const duplicates: string[] = [];
  const seenAddition = new Set<string>();

  for (const raw of additions) {
    const id = raw.trim();
    if (!id || seenAddition.has(id)) {
      continue;
    }
    seenAddition.add(id);
    if (existingSet.has(id)) {
      duplicates.push(id);
      continue;
    }
    existingSet.add(id);
    merged.push(id);
    added.push(id);
  }

  return { merged, added, duplicates };
}

export function removeModelFromList(existing: string[], modelId: string): string[] {
  const target = modelId.trim();
  return existing.filter((id) => id.trim() !== target);
}
