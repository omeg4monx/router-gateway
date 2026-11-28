import { Check, ChevronsUpDown, FileText, Search, Sparkles } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";
import type { UILocale } from "@/shared/store/uiStore";
import type { BuiltinPromptTemplate } from "./builtinPromptTemplates";
import { estimatePromptTokens, MODEL_EDITOR_COPY } from "./parameterGuide";

export function filterPromptTemplates(
  templates: readonly BuiltinPromptTemplate[],
  query: string,
): BuiltinPromptTemplate[] {
  const normalized = query.trim().toLowerCase();
  if (!normalized) return [...templates];
  return templates.filter(
    (tpl) =>
      tpl.name.toLowerCase().includes(normalized) ||
      tpl.content.slice(0, 240).toLowerCase().includes(normalized),
  );
}

function templatePreviewSnippet(content: string): string {
  const compact = content.replace(/\s+/g, " ").trim();
  return compact.length > 72 ? `${compact.slice(0, 72)}…` : compact;
}

interface PromptTemplateComboboxProps {
  selectedId: string;
  templates?: readonly BuiltinPromptTemplate[];
  locale: UILocale;
  onSelectCustom: () => void;
  onSelectTemplate: (template: BuiltinPromptTemplate) => void;
}

export function PromptTemplateCombobox({
  selectedId,
  templates,
  locale,
  onSelectCustom,
  onSelectTemplate,
}: PromptTemplateComboboxProps) {
  const copy = MODEL_EDITOR_COPY[locale];
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [activeIndex, setActiveIndex] = useState(0);
  const [loadedTemplates, setLoadedTemplates] = useState<readonly BuiltinPromptTemplate[] | null>(
    null,
  );
  const [isLoadingTemplates, setIsLoadingTemplates] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const activeTemplates = templates ?? loadedTemplates ?? [];

  useEffect(() => {
    if (templates || loadedTemplates || (!open && selectedId === "custom")) {
      return;
    }
    let cancelled = false;
    setIsLoadingTemplates(true);
    void import("./builtinPromptTemplates")
      .then((mod) => {
        if (!cancelled) {
          setLoadedTemplates(mod.BUILTIN_PROMPT_TEMPLATES);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setIsLoadingTemplates(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [open, selectedId, templates, loadedTemplates]);

  const selectedTemplate = useMemo(
    () => activeTemplates.find((tpl) => tpl.id === selectedId) ?? null,
    [activeTemplates, selectedId],
  );

  const filtered = useMemo(
    () => filterPromptTemplates(activeTemplates, search),
    [activeTemplates, search],
  );

  // Index 0 is "Free Text" ("custom"), indices 1..filtered.length are templates.
  const totalOptions = 1 + filtered.length;

  useEffect(() => {
    if (!open) return;
    const frame = window.requestAnimationFrame(() => {
      searchInputRef.current?.focus();
    });
    const handlePointerDown = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", handlePointerDown);
    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.cancelAnimationFrame(frame);
      document.removeEventListener("mousedown", handlePointerDown);
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [open]);

  const handleSearchKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "ArrowDown") {
      event.preventDefault();
      setActiveIndex((prev) => (prev + 1) % totalOptions);
    } else if (event.key === "ArrowUp") {
      event.preventDefault();
      setActiveIndex((prev) => (prev - 1 + totalOptions) % totalOptions);
    } else if (event.key === "Enter") {
      event.preventDefault();
      if (activeIndex === 0) {
        onSelectCustom();
      } else {
        const chosen = filtered[activeIndex - 1];
        if (chosen) onSelectTemplate(chosen);
      }
      setSearch("");
      setOpen(false);
    }
  };

  const activeTemplate = activeIndex > 0 ? filtered[activeIndex - 1] : undefined;
  const activeOptionId =
    open && activeIndex === 0
      ? "prompt-tpl-opt-custom"
      : open && activeTemplate
        ? `prompt-tpl-opt-${activeTemplate.id}`
        : undefined;

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        role="combobox"
        aria-expanded={open}
        aria-haspopup="listbox"
        aria-controls="prompt-template-listbox"
        aria-label={copy.promptTemplatePickerLabel}
        onClick={() => {
          setActiveIndex(0);
          setOpen((prev) => !prev);
        }}
        className="flex h-8 w-full cursor-pointer items-center justify-between gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-2.5 text-left text-xs transition-colors hover:border-[var(--border-strong)] focus:border-[var(--accent-primary)] focus:outline-none"
      >
        <div className="flex min-w-0 flex-1 items-center gap-1.5">
          {selectedTemplate ? (
            <Sparkles
              className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
          ) : (
            <FileText
              className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]"
              aria-hidden="true"
            />
          )}
          <span className="truncate font-medium text-[var(--text-primary)]">
            {selectedTemplate ? selectedTemplate.name : copy.promptFreeTextLabel}
          </span>
          <span className="shrink-0 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-1.5 py-0.2 font-mono text-[10px] text-[var(--text-muted)]">
            {activeTemplates.length > 0
              ? `${activeTemplates.length} ${copy.promptTemplatesUnit}`
              : copy.promptTemplatesUnit}
          </span>
        </div>
        <ChevronsUpDown
          className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]"
          aria-hidden="true"
        />
      </button>

      {open ? (
        <div className="absolute left-0 right-0 top-full z-40 mt-1 flex max-h-64 flex-col overflow-hidden rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-xl">
          <div className="flex shrink-0 items-center gap-2 border-b border-[var(--border-subtle)] px-2.5 py-2">
            <Search className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]" aria-hidden="true" />
            <input
              ref={searchInputRef}
              type="text"
              value={search}
              onChange={(event) => {
                setSearch(event.target.value);
                setActiveIndex(0);
              }}
              onKeyDown={handleSearchKeyDown}
              placeholder={copy.promptTemplateSearchPlaceholder}
              aria-label={copy.promptTemplateSearchPlaceholder}
              aria-activedescendant={activeOptionId}
              className="w-full bg-transparent text-xs text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:outline-none"
            />
          </div>

          <div
            id="prompt-template-listbox"
            role="listbox"
            className="min-h-0 flex-1 overflow-y-auto p-1"
          >
            {/* Free Text option */}
            <button
              id="prompt-tpl-opt-custom"
              type="button"
              role="option"
              aria-selected={selectedId === "custom"}
              onMouseEnter={() => setActiveIndex(0)}
              onClick={() => {
                onSelectCustom();
                setSearch("");
                setOpen(false);
              }}
              className={cn(
                "mb-1 flex w-full cursor-pointer items-center justify-between gap-2 rounded border-b border-[var(--border-subtle)] px-2 py-1.5 text-left text-xs transition-colors",
                selectedId === "custom"
                  ? "bg-[var(--accent-primary)]/15 text-[var(--accent-text)]"
                  : activeIndex === 0
                    ? "bg-[var(--bg-surface-raised)] text-[var(--text-primary)]"
                    : "text-[var(--text-primary)] hover:bg-[var(--bg-surface-raised)]",
              )}
            >
              <div className="min-w-0 flex-1">
                <p className="font-medium">{copy.promptFreeTextLabel}</p>
                <p className="truncate text-[10px] text-[var(--text-muted)]">
                  {copy.promptFreeTextHint}
                </p>
              </div>
              {selectedId === "custom" ? (
                <Check
                  className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
                  aria-hidden="true"
                />
              ) : null}
            </button>

            {isLoadingTemplates ? (
              <p className="px-2.5 py-3 text-center text-xs text-[var(--text-muted)]">
                {copy.promptTemplateLoading}
              </p>
            ) : filtered.length === 0 ? (
              <p className="px-2.5 py-3 text-center text-xs text-[var(--text-muted)]">
                {copy.promptTemplateEmpty}
              </p>
            ) : (
              filtered.map((tpl, idx) => {
                const optionIdx = idx + 1;
                const isSelected = tpl.id === selectedId;
                const isActive = activeIndex === optionIdx;
                const tokens = estimatePromptTokens(tpl.content);
                return (
                  <button
                    key={tpl.id}
                    id={`prompt-tpl-opt-${tpl.id}`}
                    type="button"
                    role="option"
                    aria-selected={isSelected}
                    onMouseEnter={() => setActiveIndex(optionIdx)}
                    onClick={() => {
                      onSelectTemplate(tpl);
                      setSearch("");
                      setOpen(false);
                    }}
                    className={cn(
                      "flex w-full cursor-pointer items-center justify-between gap-2 rounded px-2 py-1.5 text-left text-xs transition-colors",
                      isSelected
                        ? "bg-[var(--accent-primary)]/15 text-[var(--accent-text)]"
                        : isActive
                          ? "bg-[var(--bg-surface-raised)] text-[var(--text-primary)]"
                          : "text-[var(--text-primary)] hover:bg-[var(--bg-surface-raised)]",
                    )}
                  >
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-1.5">
                        <span className="truncate font-medium">{tpl.name}</span>
                        <span className="shrink-0 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-1 font-mono text-[9px] text-[var(--text-muted)]">
                          {tpl.position === "suffix"
                            ? copy.promptSuffixLabel
                            : copy.promptPrefixLabel}
                        </span>
                        <span className="shrink-0 font-mono text-[10px] text-[var(--text-muted)]">
                          ~{tokens} tok
                        </span>
                      </div>
                      <p className="mt-0.5 truncate text-[10px] text-[var(--text-muted)]">
                        {templatePreviewSnippet(tpl.content)}
                      </p>
                    </div>
                    {isSelected ? (
                      <Check
                        className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
                        aria-hidden="true"
                      />
                    ) : null}
                  </button>
                );
              })
            )}
          </div>
        </div>
      ) : null}
    </div>
  );
}
