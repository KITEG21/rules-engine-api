export interface RuleNode {
  type?: string;
  operator?: string;
  conditions?: RuleNode[];
  field?: string;
  value?: unknown;
  [key: string]: unknown;
}

export interface Rule {
  id: number;
  name: string;
  description: string;
  definition: RuleNode;
  created_at: string;
  updated_at: string;
}

export interface EvaluationResult {
  result: boolean;
  trace: TraceStep[];
}

export interface EvaluationResultWithRuleId {
  ruleId: number;
  ruleName: string;
  result: boolean;
  trace: TraceStep[];
}

export interface TraceStep {
  path: string;
  operator: string;
  value: unknown;
  result: boolean;
  timestamp: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}
