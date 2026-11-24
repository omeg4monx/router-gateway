import { BookOpen, Download, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { Button } from "@/shared/ui/button";
import {
  useCreateSkillMutation,
  useDeleteSkillMutation,
  useSkillsList,
  useUpdateSkillMutation,
} from "./hooks";

export function SkillsStudio() {
  const skillsQuery = useSkillsList();
  const createMutation = useCreateSkillMutation();
  const updateMutation = useUpdateSkillMutation();
  const deleteMutation = useDeleteSkillMutation();

  const [slug, setSlug] = useState("");
  const [name, setName] = useState("");
  const [category, setCategory] = useState("operations");
  const [description, setDescription] = useState("");
  const [content, setContent] = useState("# Custom Skill\n\nDescribe the workflow steps here.");
  const [notice, setNotice] = useState<string | null>(null);

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    setNotice(null);
    createMutation.mutate(
      {
        slug: slug.trim().toLowerCase(),
        name: name.trim(),
        category: category.trim() || "general",
        description: description.trim(),
        content,
        enabled: true,
      },
      {
        onSuccess: (created) => {
          setSlug("");
          setName("");
          setDescription("");
          setContent("# Custom Skill\n\nDescribe the workflow steps here.");
          setNotice(`Created skill ${created.slug}.`);
        },
      },
    );
  };

  return (
    <div className="space-y-6">
      <div className="border-b border-border pb-4">
        <div className="flex items-center gap-2 text-primary">
          <BookOpen className="h-4 w-4" />
          <span className="font-mono text-xs uppercase tracking-widest">Agent & IDE Catalog</span>
        </div>
        <h1 className="mt-1 text-2xl font-semibold tracking-tight text-foreground">
          Skills Catalog
        </h1>
        <p className="text-sm text-muted-foreground">
          Manage built-in and custom SKILL.md blueprints for coding agents and local workflows.
        </p>
      </div>

      {notice && (
        <div className="rounded-md border border-primary/40 bg-card px-4 py-2.5 text-sm text-foreground">
          {notice}
        </div>
      )}

      {(createMutation.error || deleteMutation.error || updateMutation.error) && (
        <div className="rounded-md border border-destructive/50 bg-card px-4 py-2.5 text-sm text-destructive">
          {createMutation.error?.message ??
            deleteMutation.error?.message ??
            updateMutation.error?.message}
        </div>
      )}

      <div className="grid gap-6 lg:grid-cols-3">
        <form
          onSubmit={handleCreate}
          className="space-y-3 rounded-lg border border-border bg-card p-5 lg:col-span-1"
        >
          <h2 className="text-base font-semibold text-foreground">Create Custom Skill</h2>
          <div className="space-y-1.5">
            <label htmlFor="skill-slug" className="block text-xs font-medium text-muted-foreground">
              Slug
            </label>
            <input
              id="skill-slug"
              type="text"
              required
              value={slug}
              onChange={(e) => setSlug(e.target.value)}
              placeholder="router-failover-audit"
              className="w-full rounded-md border border-border bg-background px-3 py-1.5 font-mono text-xs text-foreground"
            />
          </div>
          <div className="space-y-1.5">
            <label htmlFor="skill-name" className="block text-xs font-medium text-muted-foreground">
              Display Name
            </label>
            <input
              id="skill-name"
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Router Failover Audit"
              className="w-full rounded-md border border-border bg-background px-3 py-1.5 text-xs text-foreground"
            />
          </div>
          <div className="space-y-1.5">
            <label
              htmlFor="skill-category"
              className="block text-xs font-medium text-muted-foreground"
            >
              Category
            </label>
            <input
              id="skill-category"
              type="text"
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="w-full rounded-md border border-border bg-background px-3 py-1.5 text-xs text-foreground"
            />
          </div>
          <div className="space-y-1.5">
            <label
              htmlFor="skill-description"
              className="block text-xs font-medium text-muted-foreground"
            >
              Description
            </label>
            <input
              id="skill-description"
              type="text"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Short summary for SKILL.md frontmatter"
              className="w-full rounded-md border border-border bg-background px-3 py-1.5 text-xs text-foreground"
            />
          </div>
          <div className="space-y-1.5">
            <label
              htmlFor="skill-content"
              className="block text-xs font-medium text-muted-foreground"
            >
              Markdown Instructions
            </label>
            <textarea
              id="skill-content"
              rows={7}
              required
              value={content}
              onChange={(e) => setContent(e.target.value)}
              className="w-full rounded-md border border-border bg-background p-3 font-mono text-xs text-foreground"
            />
          </div>
          <Button type="submit" size="sm" disabled={createMutation.isPending} className="w-full">
            <Plus className="mr-1.5 h-3.5 w-3.5" />
            Add Skill
          </Button>
        </form>

        <div className="space-y-3 lg:col-span-2">
          {skillsQuery.isLoading ? (
            <div className="rounded-lg border border-border bg-card p-6 text-sm text-muted-foreground">
              Loading skills...
            </div>
          ) : (
            (skillsQuery.data?.skills ?? []).map((skill) => (
              <div key={skill.id} className="rounded-lg border border-border bg-card p-5 space-y-3">
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <div className="flex items-center gap-2">
                      <h3 className="text-base font-semibold text-foreground">{skill.name}</h3>
                      <span className="rounded border border-border bg-background px-2 py-0.5 font-mono text-[11px] text-muted-foreground">
                        {skill.slug}
                      </span>
                      {skill.builtin && (
                        <span className="rounded border border-primary/40 bg-background px-2 py-0.5 font-mono text-[11px] text-primary">
                          built-in
                        </span>
                      )}
                    </div>
                    <p className="mt-1 text-xs text-muted-foreground">{skill.description}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        updateMutation.mutate({
                          id: skill.id,
                          body: { enabled: !skill.enabled },
                        })
                      }
                    >
                      {skill.enabled ? "Disable" : "Enable"}
                    </Button>
                    <a
                      href={`/api/skills/${encodeURIComponent(skill.id)}/export`}
                      download={`${skill.slug}-SKILL.md`}
                      className="inline-flex h-8 items-center rounded-md border border-border bg-background px-3 text-xs font-medium text-foreground hover:border-primary"
                    >
                      <Download className="mr-1.5 h-3.5 w-3.5" />
                      SKILL.md
                    </a>
                    {!skill.builtin && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => deleteMutation.mutate(skill.id)}
                      >
                        <Trash2 className="h-3.5 w-3.5 text-destructive" />
                      </Button>
                    )}
                  </div>
                </div>
                <pre className="max-h-40 overflow-auto rounded border border-border bg-background p-3 font-mono text-xs text-foreground">
                  {skill.content}
                </pre>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
