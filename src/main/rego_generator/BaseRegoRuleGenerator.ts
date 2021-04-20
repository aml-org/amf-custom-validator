export interface RegoRuleResult {
    constraintId: string
    rego: string[]
    path: string
    value: string
    traceMessage?: string
}

export abstract class BaseRegoRuleGenerator {
    abstract generate(): string
}

export abstract class BaseRegoAtomicRuleGenerator extends  BaseRegoRuleGenerator {
    generate(): string {
        return this.generateResult().rego.join("\n");
    }
    abstract generateResult(): RegoRuleResult
}

export abstract class BaseRegoComplexRuleGenerator extends  BaseRegoRuleGenerator {
    generate(): string {
        return this.generateResult().map((r) => r.rego).join("\n");
    }
    abstract generateResult(): RegoRuleResult[]
}