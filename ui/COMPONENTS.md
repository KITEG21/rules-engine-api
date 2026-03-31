# Rules Engine UI Components

## 📋 Overview

This is a complete Next.js UI implementation for the Rules Engine with the recommended stack:
- **Next.js 16.2.1** (App Router)
- **Tailwind CSS 4** (Styling)
- **shadcn/ui** (Component library)
- **@uiw/react-codemirror** (JSON editing)
- **lucide-react** (Icons)

## 🏗️ Architecture

### Layout Structure

```
┌─────────────────────────────────────────────────────┐
│          Header - Rules Engine                      │
├─────────────────────────────────────────────────────┤
│ 📝 Step 1: Natural Language Input                   │
│ ┌───────────────────────────────────────────────┐  │
│ │ [textarea - describe your rule]               │  │
│ │ [Generate Rule Button]                        │  │
│ └───────────────────────────────────────────────┘  │
├──────────────────────┬──────────────────────────────┤
│ 📋 AST Preview       │ 🎨 Visual Builder           │
│ ┌────────────────┐   │ ┌──────────────────────┐    │
│ │ JSON (editable)│   │ │ Tree visualization   │    │
│ │ CodeMirror     │   │ │ Expandable nodes     │    │
│ └────────────────┘   │ └──────────────────────┘    │
├──────────────────────┼──────────────────────────────┤
│ 🧪 Test & Evaluate   │ 📊 Trace Viewer            │
│ ┌────────────────┐   │ ┌──────────────────────┐    │
│ │ Test data JSON │   │ │ Evaluation steps     │    │
│ │ [Run Button]   │   │ │ Color coded results  │    │
│ │ Result: ✓/✗    │   │ │ Green = True         │    │
│ └────────────────┘   │ │ Red = False          │    │
│                      │ └──────────────────────┘    │
└──────────────────────┴──────────────────────────────┘
```

## 🧩 Components

### 1. **RuleInput** (`./components/RuleInput.tsx`)
Natural language input for rule generation.

**Features:**
- Textarea for natural language prompts
- Character counter
- Loading state with spinner
- Disabled states for validation

**Props:**
```typescript
interface RuleInputProps {
  onSubmit: (input: string) => Promise<void>;
  loading?: boolean;
}
```

### 2. **RulePreview** (`./components/RulePreview.tsx`)
AST visualization with CodeMirror editing.

**Features:**
- JSON display in read-only mode
- Editable mode with CodeMirror
- Format JSON button
- Real-time parsing with error handling

**Props:**
```typescript
interface RulePreviewProps {
  rule: RuleNode | null;
  onUpdate?: (rule: RuleNode) => void;
}
```

### 3. **RuleBuilder** (`./components/RuleBuilder.tsx`)
Visual tree representation of the rule AST.

**Features:**
- Recursive rendering of nested conditions
- Expandable/collapsible nodes
- Color-coded operators and fields
- Supports complex nested logic (AND/OR)

**Props:**
```typescript
interface RuleBuilderProps {
  rule: RuleNode | null;
}
```

**Example Rule Rendered:**
```
┌─ AND (2 conditions) ─┐
│ ├─ age >= 18        │
│ └─ status = active  │
└────────────────────┘
```

### 4. **Evaluator** (`./components/Evaluator.tsx`)
Test data input and rule evaluation runner.

**Features:**
- JSON input for test data
- CodeMirror for editing
- Run evaluation button
- Display result (✓ TRUE or ✗ FALSE)
- Error handling for invalid JSON
- Integration with trace viewer

**Props:**
```typescript
interface EvaluatorProps {
  rule: object | null;
  onEvaluate?: (data: object, trace: TraceStep[]) => void;
}
```

### 5. **TraceViewer** (`./components/TraceViewer.tsx`)
Step-by-step evaluation trace with visual feedback.

**Features:**
- List of evaluation steps
- Color coding (green = true, red = false)
- Operator and value details per step
- Timestamp per evaluation
- Icons for quick visual feedback

**Props:**
```typescript
interface TraceViewerProps {
  steps: TraceStep[];
}
```

## 📦 Types (`./lib/types.ts`)

```typescript
// Rule node structure
interface RuleNode {
  type?: string;
  operator?: string;
  conditions?: RuleNode[];
  field?: string;
  value?: unknown;
}

// Evaluation result
interface EvaluationResult {
  result: boolean;
  trace: TraceStep[];
}

// Single trace step
interface TraceStep {
  path: string;
  operator: string;
  value: unknown;
  result: boolean;
  timestamp: number;
}
```

## 🚀 Getting Started

### Prerequisites
```bash
npm install --legacy-peer-deps
```

### Development
```bash
npm run dev
```
Open [http://localhost:3000](http://localhost:3000)

### Build
```bash
npm run build
npm start
```

### Linting
```bash
npm run lint
```

## 🔗 API Integration Points

The components have placeholder comments for API endpoints. Connect them to your backend:

1. **Generate Rule** (in `page.tsx`)
```javascript
const response = await fetch('/api/generate-rule', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ prompt: input })
});
const { rule } = await response.json();
```

2. **Evaluate Rule** (in `Evaluator.tsx`)
```javascript
const response = await fetch('/api/evaluate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ rule, data: testData })
});
const { result, trace } = await response.json();
```

## 🎨 Styling

All components use Tailwind CSS with a clean slate color palette:
- **Primary**: Blue-600 (buttons, highlights)
- **Success**: Green-500 (true results)
- **Error**: Red-500 (false results)
- **Neutral**: Slate colors (backgrounds, text)

## ✨ Feature Highlights

✅ **Responsive Design** - Mobile-friendly grid layout  
✅ **Real-time Editing** - CodeMirror JSON editing with live parsing  
✅ **Visual Feedback** - Color-coded results and trace steps  
✅ **Type Safety** - Full TypeScript support  
✅ **Loading States** - Spinner animations during async operations  
✅ **Error Handling** - Graceful error messages  
✅ **Expandable Trees** - Collapsible rule conditions  
✅ **Character Counter** - Input validation feedback  

## 📝 Usage Example

```typescript
import { useState } from 'react';
import { RuleInput, RulePreview, RuleBuilder, Evaluator, TraceViewer } from './components';
import { RuleNode, TraceStep } from './lib/types';

export default function Home() {
  const [rule, setRule] = useState<RuleNode | null>(null);
  const [trace, setTrace] = useState<TraceStep[]>([]);

  return (
    <>
      <RuleInput onSubmit={async (input) => {
        // Call your API
        const result = await generateRule(input);
        setRule(result);
      }} />
      
      {rule && (
        <>
          <RulePreview rule={rule} onUpdate={setRule} />
          <RuleBuilder rule={rule} />
          <Evaluator rule={rule} onEvaluate={(data, newTrace) => setTrace(newTrace)} />
          <TraceViewer steps={trace} />
        </>
      )}
    </>
  );
}
```

## 🛠️ Customization

### Change Color Scheme
Update Tailwind classes in components:
- Blue-600 → Yellow-500
- Green-500 → Emerald-500
- Red-500 → Rose-500

### Max Heights & Widths
Adjust `max-h-*` and `max-w-*` classes in components as needed

### CodeMirror Theme
Change `theme="light"` to `theme="dark"` in RulePreview and Evaluator

## 📚 Documentation

Refer to:
- [Next.js App Router](https://nextjs.org/docs/app)
- [Tailwind CSS](https://tailwindcss.com)
- [CodeMirror](https://codemirror.net)
- [Lucide Icons](https://lucide.dev)

## 🐛 Troubleshooting

**CodeMirror not showing?**
- Ensure `@uiw/react-codemirror` is installed
- Check that CodeMirror extensions are properly imported

**Styling issues?**
- Run `npm run build` to ensure Tailwind classes are generated
- Clear `.next` folder if styles don't update

**Type errors?**
- Run `npm run lint` to check TypeScript
- Verify all imports use correct paths from `./lib/types`
