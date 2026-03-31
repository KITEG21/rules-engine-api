'use client';

import { useState } from 'react';
import { RulesList } from '@/app/components/RulesList';
import { TraceViewer } from '@/app/components/TraceViewer';
import { EvaluationResultWithRuleId } from '@/app/lib/types';
import Link from 'next/link';

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

export default function RulesPage() {
  const apiUrl = resolveApiUrl();
  const [selectedRuleIds, setSelectedRuleIds] = useState<number[]>([]);
  const [testData, setTestData] = useState('{}');
  const [evaluationResults, setEvaluationResults] = useState<EvaluationResultWithRuleId[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [rulesRefresh, setRulesRefresh] = useState(0);
  const [testDataError, setTestDataError] = useState<string | null>(null);

  const handleEvaluate = async () => {
    setError(null);
    setTestDataError(null);

    if (selectedRuleIds.length === 0) {
      setError('Please select at least one rule to evaluate');
      return;
    }

    let parsedData: Record<string, unknown>;
    try {
      parsedData = JSON.parse(testData);
    } catch (err) {
      setTestDataError('Invalid JSON: ' + (err instanceof Error ? err.message : 'Unknown error'));
      return;
    }

    setLoading(true);
    try {
      const response = await fetch(`${apiUrl}/api/v1/rules/evaluate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          data: parsedData,
          ruleIds: selectedRuleIds,
        }),
      });

      if (!response.ok) {
        const errBody = await response.json().catch(() => null);
        const message = errBody?.error || 'Evaluation failed';
        throw new Error(message);
      }

      const data = await response.json();
      setEvaluationResults(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Evaluation failed');
      setEvaluationResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleTestDataChange = (value: string) => {
    setTestData(value);
    setTestDataError(null);
  };

  const handleFormatTestData = () => {
    try {
      const parsed = JSON.parse(testData);
      setTestData(JSON.stringify(parsed, null, 2));
      setTestDataError(null);
    } catch (err) {
      setTestDataError('Invalid JSON: ' + (err instanceof Error ? err.message : 'Unknown error'));
    }
  };

  return (
    <main className="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)] p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header with Navigation */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold mb-2">Rules List & Evaluator</h1>
            <p className="text-[var(--text-secondary)]">Select rules and test them with data</p>
          </div>
          <Link
            href="/"
            className="px-6 py-2 bg-[var(--accent)] !text-white rounded-lg hover:bg-[var(--accent-hover)] transition-all font-semibold hover:shadow-lg hover:shadow-[rgba(79,70,229,0.3)] hover:-translate-y-0.5"
          >
            ← New Rule
          </Link>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-[rgba(239,68,68,0.1)] border border-[var(--error)] rounded-lg">
            <p className="text-[var(--error)] font-semibold">Error: {error}</p>
          </div>
        )}

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Rules List */}
          <div>
            <RulesList
              key={rulesRefresh}
              apiUrl={apiUrl}
              selectedRuleIds={selectedRuleIds}
              onSelectionChange={setSelectedRuleIds}
              onRuleDeleted={() => setRulesRefresh(prev => prev + 1)}
            />
          </div>

          {/* Test Data Input */}
          <div className="flex flex-col gap-4">
            <div className="border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
              <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-6 py-4 flex justify-between items-center">
                <h3 className="font-semibold text-[var(--text-primary)]">Test Data</h3>
                <button
                  onClick={handleFormatTestData}
                  className="px-3 py-1 text-sm bg-[var(--bg-secondary)] border border-[var(--input-border)] rounded hover:bg-[var(--bg-tertiary)] text-[var(--text-primary)] transition-colors"
                >
                  Format
                </button>
              </div>

              <textarea
                value={testData}
                onChange={e => handleTestDataChange(e.target.value)}
                className={`w-full p-4 bg-[var(--bg-secondary)] text-[var(--text-primary)] font-mono text-sm border-none outline-none resize-none h-72 ${testDataError ? 'ring-2 ring-[var(--error)]' : 'hover:bg-[var(--bg-tertiary)]'}`}
                placeholder="Enter test data as JSON..."
              />

              {testDataError && (
                <div className="bg-[rgba(239,68,68,0.05)] border-t border-[var(--error)] px-4 py-3">
                  <p className="text-xs text-[var(--error)]">{testDataError}</p>
                </div>
              )}
            </div>

            {/* Evaluate Button */}
            <button
              onClick={handleEvaluate}
              disabled={loading || selectedRuleIds.length === 0}
              className="w-full px-6 py-4 bg-[var(--accent)] text-white rounded-lg hover:bg-[var(--accent-hover)] disabled:bg-[var(--text-muted)] disabled:cursor-not-allowed transition-all font-semibold text-lg hover:shadow-lg hover:shadow-[rgba(79,70,229,0.3)] hover:-translate-y-0.5"
            >
              {loading ? 'Evaluating...' : `Evaluate ${selectedRuleIds.length} Rule${selectedRuleIds.length !== 1 ? 's' : ''}`}
            </button>
          </div>
        </div>

        {/* Evaluation Results */}
        {evaluationResults.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Results Summary */}
            <div className="border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
              <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-6 py-4">
                <h3 className="font-semibold text-[var(--text-primary)]">Results Summary</h3>
              </div>

              <div className="divide-y divide-[var(--input-border)]">
                {evaluationResults.map((result, idx) => (
                  <div
                    key={idx}
                    className={`p-4 flex items-center justify-between ${
                      result.result
                        ? 'bg-[rgba(34,197,94,0.05)]'
                        : 'bg-[rgba(239,68,68,0.05)]'
                    } border-l-4 ${
                      result.result
                        ? 'border-l-[var(--success)]'
                        : 'border-l-[var(--error)]'
                    }`}
                  >
                    <div>
                      <p className="font-semibold text-[var(--text-primary)]">
                        {result.ruleName}
                      </p>
                      <p className="text-xs text-[var(--text-secondary)] mt-1">
                        Rule ID: {result.ruleId}
                      </p>
                    </div>
                    <div className={`text-3xl font-bold ${
                      result.result
                        ? 'text-[var(--success)]'
                        : 'text-[var(--error)]'
                    }`}>
                      {result.result ? 'PASS' : 'FAIL'}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Detailed Traces */}
            <div className="flex flex-col gap-4">
              {evaluationResults.map((result, idx) => (
                <div key={idx} className="border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
                  <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-6 py-3">
                    <h4 className="font-semibold text-[var(--text-primary)] text-sm">
                      {result.ruleName} - Trace
                    </h4>
                  </div>
                  <TraceViewer steps={result.trace} />
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Empty State */}
        {evaluationResults.length === 0 && selectedRuleIds.length > 0 && !loading && (
          <div className="border border-[var(--input-border)] rounded-lg bg-[var(--bg-tertiary)] p-8 text-center">
            <p className="text-[var(--text-secondary)]">
              Click &quot;Evaluate&quot; button to see results
            </p>
          </div>
        )}
      </div>
    </main>
  );
}
