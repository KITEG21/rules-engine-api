'use client';

import { useState } from 'react';
import Link from 'next/link';
import { RuleInput } from './components/RuleInput';
import { RulePreview } from './components/RulePreview';
import { RuleBuilder } from './components/RuleBuilder';
import { Evaluator } from './components/Evaluator';
import { TraceViewer } from './components/TraceViewer';
import { RuleNode, TraceStep } from './lib/types';

const rawApiUrl = process.env.NEXT_PUBLIC_API_URL?.trim() ?? '';
const normalizedApiUrl = rawApiUrl
  ? rawApiUrl.replace(/\/+$/, '')
  : '';

function resolveApiUrl(): string {
  if (!normalizedApiUrl) return 'http://localhost:8080';
  if (/^https?:\/\//i.test(normalizedApiUrl)) {
    return normalizedApiUrl;
  }
  return `http://${normalizedApiUrl}`;
}

export default function Home() {
  const [rule, setRule] = useState<RuleNode | null>(null);
  const [ruleId, setRuleId] = useState<number | null>(null);
  const [trace, setTrace] = useState<TraceStep[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGenerateRule = async (params: { name: string; description: string; definition: string }) => {
    setError(null);
    setLoading(true);
    try {
      const targetHost = resolveApiUrl();
      const response = await fetch(`${targetHost}/api/v1/rules`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: params.name,
          description: params.description,
          definition: params.definition,
        }),
      });

      if (!response.ok) {
        const errBody = await response.json().catch(() => null);
        const message = errBody?.error || 'Failed to create rule';

        if (typeof message === 'string' && message.toLowerCase().includes('ai client')) {
          const fallbackRule: RuleNode = {
            operator: 'AND',
            conditions: [
              { field: 'age', operator: 'gte', value: 18 },
              { field: 'status', operator: 'equals', value: 'active' },
            ],
          };
          setRule(fallbackRule);
          setRuleId(null);
          setError('AI client not configured; using fallback locally generated AST.');
          return;
        }

        setError(message);
        return;
      }

      const data = await response.json();
      setRule(data.definition as RuleNode);
      setRuleId(data.id);
    } catch (error) {
      console.error('Error generating rule:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)]">
      {/* Header */}
      <div className="bg-[var(--bg-secondary)] border-b border-[var(--card-border)]">
        <div className="max-w-7xl mx-auto px-6 py-6 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-[var(--text-primary)]">Rules Engine</h1>
            <p className="text-[var(--text-secondary)] mt-2">Generate, preview, and test business rules using natural language</p>
          </div>
          <Link
            href="/rules"
            className="px-6 py-2 bg-[var(--accent)] !text-[var(--text-primary)] rounded-lg hover:bg-[var(--accent-hover)] transition-all font-semibold whitespace-nowrap hover:shadow-lg hover:shadow-[rgba(79,70,229,0.3)] hover:-translate-y-0.5"
          >
            View All Rules
          </Link>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto p-6 space-y-6">
        {/* Section 1: Natural Language Input */}
        <div className="card rounded-lg shadow-sm p-6">
          <h2 className="text-lg font-semibold text-[var(--text-primary)] mb-4">Step 1: Generate Rule</h2>
          <RuleInput onSubmit={handleGenerateRule} loading={loading} />
          {error && (
            <div className="mt-4 p-3 bg-amber-50 border border-amber-200 rounded text-sm text-amber-700">
              {error}
            </div>
          )}
        </div>

        {rule && (
          <>
            {/* Section 2: Two Column Layout - AST and Visual Builder */}
            <div className="grid grid-cols-2 gap-6">
              <div className="card rounded-lg shadow-sm p-6">
                <h2 className="text-lg font-semibold text-[var(--text-primary)] mb-4">Step 2a: Preview AST</h2>
                <RulePreview rule={rule} onUpdate={setRule} />
              </div>

              <div className="card rounded-lg shadow-sm p-6">
                <h2 className="text-lg font-semibold text-[var(--text-primary)] mb-4">Step 2b: Rule Builder</h2>
                <RuleBuilder rule={rule} />
              </div>
            </div>

            {/* Section 3: Two Column Layout - Test Data and Result */}
            <div className="grid grid-cols-2 gap-6">
              <div className="card rounded-lg shadow-sm p-6">
                <h2 className="text-lg font-semibold text-[var(--text-primary)] mb-4">Step 3a: Test & Evaluate</h2>
                <Evaluator
                  ruleId={ruleId}
                  onEvaluate={(data, newTrace) => {
                    setTrace(newTrace);
                  }}
                />
              </div>

              <div className="card rounded-lg shadow-sm p-6">
                <h2 className="text-lg font-semibold text-[var(--text-primary)] mb-4">Step 3b: Trace Viewer</h2>
                <TraceViewer steps={trace} />
              </div>
            </div>
          </>
        )}

        {/* Empty State */}
        {!rule && (
          <div className="bg-[var(--bg-secondary)] rounded-xl border-2  border-[var(--accent)] p-12 flex justify-center">
            <div className="text-center">
              <p className="text-[var(--text-primary)] text-lg font-semibold">👆 Start by generating a rule above</p>
              <p className="text-[var(--text-secondary)] text-sm mt-2">Your generated rule will appear here</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
