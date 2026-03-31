'use client';

import { RuleNode } from '@/app/lib/types';

interface RuleBuilderProps {
  rule: RuleNode | null;
}

function Condition({ node }: { node: RuleNode }) {
  const field = node.field || 'N/A';
  const operator = node.operator || '?';
  const value = node.value !== undefined ? String(node.value) : 'N/A';

  return (
    <div className="inline-block px-4 py-2 rounded-full bg-[var(--bg-tertiary)] border border-[var(--input-border)] text-[var(--text-primary)] text-sm font-medium">
      <span className="text-[var(--accent-2)]">{field}</span>
      <span className="mx-2 text-[var(--text-secondary)]">{operator}</span>
      <span className="text-[var(--text-primary)]">{value}</span>
    </div>
  );
}

function LogicBadge({ logic }: { logic: string }) {
  const bgColor =
    logic.toUpperCase() === 'AND' ? 'var(--and-color)' : 'var(--or-color)';

  return (
    <span
      style={{ backgroundColor: bgColor }}
      className="inline-block px-3 py-1 mx-2 text-white text-sm font-bold rounded-full"
    >
      {logic.toUpperCase()}
    </span>
  );
}

interface NodeProps {
  node: RuleNode;
}

function Node({ node }: NodeProps) {
  if (!node.conditions || node.conditions.length === 0) {
    return <Condition node={node} />;
  }

  return (
    <div className="flex flex-wrap items-center gap-2">
      {node.conditions.map((child, idx) => (
        <div key={idx} className="flex items-center gap-2">
          {idx > 0 && <LogicBadge logic={node.operator || 'AND'} />}
          <Node node={child} />
        </div>
      ))}
    </div>
  );
}

export function RuleBuilder({ rule }: RuleBuilderProps) {
  if (!rule) {
    return (
      <div className="w-full h-full border border-[var(--input-border)] rounded-lg bg-[var(--bg-tertiary)] flex items-center justify-center">
        <p className="text-[var(--text-secondary)]">No rule to visualize</p>
      </div>
    );
  }

  return (
    <div className="w-full border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
      <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-6 py-4">
        <h3 className="font-semibold text-[var(--text-primary)]">Visual Rule</h3>
      </div>
      <div className="p-6">
        <div className="flex flex-wrap items-center gap-3">
          <Node node={rule} />
        </div>
      </div>
    </div>
  );
}
