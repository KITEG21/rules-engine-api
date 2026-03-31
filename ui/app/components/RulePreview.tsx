'use client';

import { useState } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import { json } from '@codemirror/lang-json';
import { RuleNode } from '@/app/lib/types';

interface RulePreviewProps {
  rule: RuleNode | null;
  onUpdate?: (rule: RuleNode) => void;
}

export function RulePreview({ rule, onUpdate }: RulePreviewProps) {
  const [editMode, setEditMode] = useState(false);
  const [jsonContent, setJsonContent] = useState(JSON.stringify(rule, null, 2));

  const handleJsonChange = (value: string) => {
    setJsonContent(value);
    try {
      const parsed = JSON.parse(value);
      onUpdate?.(parsed);
    } catch {
      // Invalid JSON, don't update
    }
  };

  const handleFormatJson = () => {
    try {
      const parsed = JSON.parse(jsonContent);
      setJsonContent(JSON.stringify(parsed, null, 2));
    } catch {
      alert('Invalid JSON format');
    }
  };

  if (!rule) {
    return (
      <div className="w-full h-full border border-[var(--input-border)] rounded-lg bg-[var(--bg-tertiary)] flex items-center justify-center">
        <p className="text-[var(--text-secondary)]">No rule generated yet</p>
      </div>
    );
  }

  return (
    <div className="w-full border border-[var(--input-border)] rounded-lg overflow-hidden bg-[var(--bg-secondary)]">
      <div className="bg-[var(--bg-tertiary)] border-b border-[var(--input-border)] px-4 py-3 flex justify-between items-center">
        <h3 className="font-semibold text-[var(--text-primary)]">AST (JSON)</h3>
        <div className="flex gap-2">
          <button
            onClick={() => setEditMode(!editMode)}
            className="px-3 py-1 text-sm bg-[var(--bg-secondary)] border border-[var(--input-border)] rounded hover:bg-[var(--bg-tertiary)] text-[var(--text-primary)] transition-colors"
          >
            {editMode ? 'Done' : 'Edit'}
          </button>
          {editMode && (
            <button
              onClick={handleFormatJson}
              className="px-3 py-1 text-sm bg-[var(--bg-secondary)] border border-[var(--input-border)] rounded hover:bg-[var(--bg-tertiary)] text-[var(--text-primary)] transition-colors"
            >
              Format
            </button>
          )}
        </div>
      </div>

      <div className="overflow-auto max-h-96 bg-[var(--bg-secondary)]">
        {editMode ? (
          <div className="bg-[var(--bg-secondary)]">
            <CodeMirror
              value={jsonContent}
              onChange={handleJsonChange}
              extensions={[json()]}
              basicSetup={{
                lineNumbers: true,
                highlightActiveLineGutter: true,
                foldGutter: true,
                dropCursor: true,
                allowMultipleSelections: true,
                indentOnInput: true,
                bracketMatching: true,
                closeBrackets: true,
                autocompletion: true,
                rectangularSelection: true,
                highlightSelectionMatches: true,
                searchKeymap: true,
              }}
              className="text-sm [&_.cm-editor]:bg-[var(--bg-tertiary)] [&_.cm-editor]:text-[var(--text-primary)] [&_.cm-gutters]:bg-[var(--bg-secondary)] [&_.cm-gutters]:border-r-[var(--input-border)] [&_.cm-activeLineGutter]:bg-[var(--bg-tertiary)] [&_.cm-cursor]:border-l-[var(--accent)]"
            />
          </div>
        ) : (
          <pre className="p-4 text-sm overflow-auto bg-[var(--bg-tertiary)] text-[var(--text-primary)] border-t border-[var(--input-border)] font-mono">
            {JSON.stringify(rule, null, 2)}
          </pre>
        )}
      </div>
    </div>
  );
}
