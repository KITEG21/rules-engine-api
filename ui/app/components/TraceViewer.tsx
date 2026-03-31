'use client';

import { TraceStep } from '@/app/lib/types';

interface TraceViewerProps {
  steps: TraceStep[];
}

export function TraceViewer({ steps }: TraceViewerProps) {
  if (steps.length === 0) {
    return (
      <div className="w-full border border-[var(--input-border)] rounded-lg bg-[var(--bg-tertiary)] flex items-center justify-center h-48">
        <p className="text-[var(--text-secondary)]">No trace data</p>
      </div>
    );
  }

  return (
    <div className="w-full border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
      <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-6 py-4">
        <h3 className="font-semibold text-[var(--text-primary)]">Evaluation Trace</h3>
      </div>
      <div className="overflow-auto max-h-96">
        <div className="divide-y divide-[var(--card-border)]">
          {steps.map((step, index) => (
            <div
              key={index}
              className={`p-4 flex gap-4 items-start border-l-4 ${
                step.result
                  ? 'border-l-[var(--success)] bg-[rgba(34,197,94,0.05)]'
                  : 'border-l-[var(--error)] bg-[rgba(239,68,68,0.05)]'
              }`}
            >
              <div className="flex-shrink-0 text-sm font-semibold pt-0.5 w-16 text-[var(--text-secondary)]">
                {step.result ? 'PASS' : 'FAIL'}
              </div>

              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-2">
                  <p className="text-sm font-mono text-[var(--accent-2)]">
                    {step.path}
                  </p>
                  <span className="inline-block px-2 py-1 bg-[var(--bg-tertiary)] rounded text-xs text-[var(--text-secondary)] font-mono">
                    {step.operator}
                  </span>
                </div>

                <p className="text-sm text-[var(--text-secondary)] font-mono">
                  {String(step.value)}
                </p>

                <p
                  className={`text-xs font-bold mt-2 ${
                    step.result
                      ? 'text-[var(--success)]'
                      : 'text-[var(--error)]'
                  }`}
                >
                  → {step.result ? 'TRUE' : 'FALSE'}
                </p>
              </div>

              <div className="text-xs text-[var(--text-muted)] whitespace-nowrap">
                {new Date(step.timestamp).toLocaleTimeString()}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
