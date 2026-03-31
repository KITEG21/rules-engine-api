'use client';

import { useState } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import { json } from '@codemirror/lang-json';
import { Loader2 } from 'lucide-react';
import { TraceStep } from '@/app/lib/types';

interface EvaluatorProps {
  ruleId: number | null;
  onEvaluate?: (data: object, trace: TraceStep[]) => void;
}

const rawApiUrl = process.env.NEXT_PUBLIC_API_URL?.trim() ?? '';
const normalizedApiUrl = rawApiUrl ? rawApiUrl.replace(/\/+$/, '') : '';

function resolveApiUrl(): string {
  if (!normalizedApiUrl) return 'http://localhost:8080';
  if (/^https?:\/\//i.test(normalizedApiUrl)) {
    return normalizedApiUrl;
  }
  return `http://${normalizedApiUrl}`;
}

export function Evaluator({ ruleId, onEvaluate }: EvaluatorProps) {
  const [testData, setTestData] = useState('{\n  \n}');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<boolean | null>(null);
  const [trace, setTrace] = useState<TraceStep[]>([]);
  const [error, setError] = useState<string | null>(null);

  const handleRun = async () => {
    try {
      setError(null);
      const data = JSON.parse(testData);
      setLoading(true);

      if (!ruleId) {
        throw new Error('Rule ID is not available. Save a rule first.');
      }

      const targetHost = resolveApiUrl();
      const response = await fetch(`${targetHost}/api/v1/rules/evaluate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ data, ruleIds: [ruleId] }),
      });

      if (!response.ok) {
        const errorBody = await response.json().catch(() => null);
        throw new Error(errorBody?.error || 'Evaluation API error');
      }

      const evalResult = await response.json();
      const resultObj = evalResult.results?.[0];

      const isMatched = !!resultObj?.matched;
      const step: TraceStep = {
        path: `rule.${ruleId}`,
        operator: resultObj?.error ? 'error' : 'evaluate',
        value: resultObj?.value ?? null,
        result: isMatched,
        timestamp: Date.now(),
      };

      setResult(isMatched);
      setTrace([step]);
      onEvaluate?.(data, [step]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Invalid JSON');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full space-y-3">
      <div className="border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
        <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-4 py-3">
          <h3 className="font-semibold text-[var(--text-primary)]">Test Data Input</h3>
        </div>
        <div className="max-h-48 overflow-auto">
          <CodeMirror
            value={testData}
            onChange={setTestData}
            extensions={[json()]}
            basicSetup={{
              lineNumbers: true,
              highlightActiveLineGutter: true,
              foldGutter: true,
            }}
            theme="dark"
            className="text-sm text-[var(--text-primary)] bg-[var(--code-bg)]"
          />
        </div>
      </div>

      <button
        onClick={handleRun}
        disabled={loading || !ruleId || !testData.trim()}
        className="w-full px-4 py-3 bg-[var(--accent)] text-[var(--text-primary)] font-medium rounded-lg hover:bg-[var(--accent-hover)] disabled:bg-[var(--text-muted)] disabled:cursor-not-allowed transition-all flex items-center justify-center gap-2 hover:shadow-lg hover:shadow-[rgba(79,70,229,0.3)] hover:-translate-y-0.5"
      >
        {loading && <Loader2 className="w-4 h-4 animate-spin" />}
        {loading ? 'Evaluating...' : 'Run Evaluation'}
      </button>

      {error && (
        <div className="p-3 bg-[rgba(220,38,38,0.1)] border border-[var(--error)] rounded-lg text-sm text-[var(--error)]">
          {error}
        </div>
      )}

      {result !== null && (
        <div className="space-y-4">
          <div
            className={`p-8 rounded-2xl border-2 flex flex-col items-center justify-center ${
              result
                ? 'bg-[rgba(34,197,94,0.1)] border-[var(--success)]'
                : 'bg-[rgba(239,68,68,0.1)] border-[var(--error)]'
            }`}
          >
            <div className="text-3xl font-bold mb-3 text-[var(--text-secondary)]">
              {result ? 'PASS' : 'FAIL'}
            </div>
            <p
              className={`text-2xl font-semibold ${
                result
                  ? 'text-[var(--success)]'
                  : 'text-[var(--error)]'
              }`}
            >
              {result ? 'True' : 'False'}
            </p>
            <p className="text-[var(--text-secondary)] mt-3">
              {trace.length} evaluation step{trace.length !== 1 ? 's' : ''}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
