'use client';

import { useState } from 'react';
import { Loader2 } from 'lucide-react';

interface RuleInputProps {
  onSubmit: (params: { name: string; description: string; definition: string }) => Promise<void>;
  loading?: boolean;
}

export function RuleInput({ onSubmit, loading = false }: RuleInputProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [input, setInput] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !description.trim() || !input.trim()) {
      return;
    }

    await onSubmit({
      name: name.trim(),
      description: description.trim(),
      definition: input.trim(),
    });

    setName('');
    setDescription('');
    setInput('');
  };

  return (
    <form onSubmit={handleSubmit} className="w-full flex flex-col items-center justify-center gap-6">
      <div className="w-full max-w-2xl space-y-3">
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Rule name (e.g., 'Premium User Check')"
          className="w-full px-4 py-2 border border-[var(--input-border)] bg-[var(--bg-tertiary)] text-[var(--text-primary)] text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:border-transparent"
          disabled={loading}
        />

        <input
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="What does this rule do?"
          className="w-full px-4 py-2 border border-[var(--input-border)] bg-[var(--bg-tertiary)] text-[var(--text-secondary)] text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:border-transparent"
          disabled={loading}
        />

        <div className="relative">
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder='Ask what you want: e.g. "Users older than 18 and VIP members"'
            className="w-full h-32 px-5 py-4 border border-[var(--input-border)] bg-[var(--bg-tertiary)] text-[var(--text-primary)] text-base rounded-xl focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:border-transparent resize-none"
            disabled={loading}
          />
          {input.length > 0 && (
            <span className="absolute bottom-3 right-4 text-xs text-[var(--text-muted)]">
              {input.length} characters
            </span>
          )}
        </div>

        <button
          type="submit"
          disabled={
            loading || !name.trim() || !description.trim() || !input.trim()
          }
          className="w-full px-6 py-3 bg-[var(--accent)] text-white font-semibold rounded-lg hover:bg-[var(--accent-hover)] disabled:bg-[var(--text-muted)] disabled:cursor-not-allowed transition-all items-center justify-center gap-2 flex text-lg hover:shadow-lg hover:shadow-[rgba(79,70,229,0.3)] hover:-translate-y-0.5"
        >
          {loading && <Loader2 className="w-5 h-5 animate-spin" />}
          {loading ? 'Generating Rule...' : 'Generate Rule'}
        </button>
      </div>
    </form>
  );
}
