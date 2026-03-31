'use client';

import { useState, useEffect, useCallback } from 'react';
import { ChevronDown, Trash2, Copy } from 'lucide-react';
import { Rule } from '@/app/lib/types';

interface RulesListProps {
  apiUrl: string;
  selectedRuleIds: number[];
  onSelectionChange: (ids: number[]) => void;
  onRuleDeleted: () => void;
}

export function RulesList({ apiUrl, selectedRuleIds, onSelectionChange, onRuleDeleted }: RulesListProps) {
  const [rules, setRules] = useState<Rule[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedId, setExpandedId] = useState<number | null>(null);
  const [deleting, setDeleting] = useState<number | null>(null);

  const fetchRules = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${apiUrl}/api/v1/rules`);
      if (!response.ok) throw new Error('Failed to fetch rules');
      const data = await response.json();
      setRules(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch rules');
      setRules([]);
    } finally {
      setLoading(false);
    }
  }, [apiUrl]);

  useEffect(() => {
    fetchRules();
  }, [fetchRules]);

  const handleSelect = (ruleId: number) => {
    const isSelected = selectedRuleIds.includes(ruleId);
    if (isSelected) {
      onSelectionChange(selectedRuleIds.filter(id => id !== ruleId));
    } else {
      onSelectionChange([...selectedRuleIds, ruleId]);
    }
  };

  const handleDelete = async (ruleId: number) => {
    if (!confirm('Are you sure you want to delete this rule?')) return;

    setDeleting(ruleId);
    try {
      const response = await fetch(`${apiUrl}/api/v1/rules/${ruleId}`, {
        method: 'DELETE',
      });
      if (!response.ok) throw new Error('Failed to delete rule');
      setRules(rules.filter(r => r.id !== ruleId));
      onSelectionChange(selectedRuleIds.filter(id => id !== ruleId));
      onRuleDeleted();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete rule');
    } finally {
      setDeleting(null);
    }
  };

  const handleDuplicate = async (rule: Rule) => {
    try {
      const response = await fetch(`${apiUrl}/api/v1/rules`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: `${rule.name} (Copy)`,
          description: rule.description,
          definition: rule.definition,
        }),
      });
      if (!response.ok) throw new Error('Failed to duplicate rule');
      const newRule = await response.json();
      setRules([...rules, newRule]);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to duplicate rule');
    }
  };

  if (loading) {
    return (
      <div className="w-full border border-[var(--input-border)] rounded-lg bg-[var(--bg-secondary)] p-8 flex items-center justify-center">
        <p className="text-[var(--text-secondary)]">Loading rules...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full border border-[var(--error)] rounded-lg bg-[rgba(239,68,68,0.05)] p-8">
        <p className="text-[var(--error)] font-semibold">Error: {error}</p>
        <button
          onClick={() => fetchRules()}
          className="mt-4 px-4 py-2 bg-[var(--error)] text-white rounded hover:bg-opacity-90 transition-all hover:shadow-lg hover:shadow-[rgba(220,38,38,0.3)] hover:-translate-y-0.5"
        >
          Retry
        </button>
      </div>
    );
  }

  if (rules.length === 0) {
    return (
      <div className="w-full border border-[var(--input-border)] rounded-lg bg-[var(--bg-tertiary)] p-8 flex items-center justify-center">
        <p className="text-[var(--text-secondary)]">No rules created yet</p>
      </div>
    );
  }

  return (
    <div className="w-full border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
      <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-6 py-4">
        <h3 className="font-semibold text-[var(--text-primary)]">
          Existing Rules ({selectedRuleIds.length} selected)
        </h3>
      </div>

      <div className="divide-y divide-[var(--input-border)]">
        {rules.map(rule => (
          <div
            key={rule.id}
            className="transition-colors hover:bg-[var(--bg-tertiary)]"
          >
            <div className="p-4 flex items-center gap-4">
              <input
                type="checkbox"
                checked={selectedRuleIds.includes(rule.id)}
                onChange={() => handleSelect(rule.id)}
                className="w-5 h-5 rounded border-[var(--input-border)] text-[var(--accent)] focus:ring-2 focus:ring-[var(--accent)] cursor-pointer"
              />

              <button
                onClick={() => setExpandedId(expandedId === rule.id ? null : rule.id)}
                className="flex-1 flex items-center justify-between text-left hover:opacity-80 transition-opacity"
                type="button"
              >
                <div>
                  <p className="font-semibold text-[var(--text-primary)]">{rule.name}</p>
                  {rule.description && (
                    <p className="text-sm text-[var(--text-secondary)] mt-1">{rule.description}</p>
                  )}
                  <p className="text-xs text-[var(--text-muted)] mt-2">
                    Updated: {new Date(rule.updated_at).toLocaleDateString()}
                  </p>
                </div>
                <ChevronDown
                  className={`w-5 h-5 text-[var(--text-secondary)] transition-transform ${
                    expandedId === rule.id ? 'rotate-180' : ''
                  }`}
                />
              </button>

              <div className="flex items-center gap-2">
                <button
                  onClick={() => handleDuplicate(rule)}
                  className="p-2 hover:bg-[var(--bg-secondary)] rounded transition-colors text-[var(--text-secondary)] hover:text-[var(--accent)]"
                  title="Duplicate rule"
                  type="button"
                >
                  <Copy className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(rule.id)}
                  disabled={deleting === rule.id}
                  className="p-2 hover:bg-[var(--bg-secondary)] rounded transition-colors text-[var(--text-secondary)] hover:text-[var(--error)]"
                  title="Delete rule"
                  type="button"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>

            {expandedId === rule.id && (
              <div className="bg-[var(--bg-tertiary)] border-t border-[var(--input-border)] px-6 py-4">
                <pre className="text-xs text-[var(--text-secondary)] font-mono overflow-auto max-h-64 bg-[var(--bg-secondary)] p-3 rounded border border-[var(--input-border)]">
                  {JSON.stringify(rule.definition, null, 2)}
                </pre>
              </div>
            )}
          </div>
        ))}
      </div>

      {selectedRuleIds.length > 0 && (
        <div className="bg-[var(--bg-tertiary)] border-t border-[var(--input-border)] px-6 py-4">
          <p className="text-sm text-[var(--text-secondary)]">
            {selectedRuleIds.length} rule{selectedRuleIds.length !== 1 ? 's' : ''} selected for evaluation
          </p>
        </div>
      )}
    </div>
  );
}
