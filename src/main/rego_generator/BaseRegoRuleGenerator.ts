export interface RegoRuleResult {
    constraintId: string
    rego: string[]
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