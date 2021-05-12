import {RegoPathResult} from "./RegoPathGenerator";

export abstract class RegoRuleResult {
    public readonly constraintId: string

    constructor(constraintId: string) {
        this.constraintId = constraintId;
    }

}

export class SimpleRuleResult extends RegoRuleResult {
    public readonly rego: string[];
    public readonly path: string;
    public readonly value: string;
    public readonly variable: string;
    public readonly traceMessage: string;
    public readonly pathRules: RegoPathResult[];

    constructor(constraintId: string, rego: string[], pathRules: RegoPathResult[], path: string, value: string, variable: string,traceMessage?: string|any) {
        super(constraintId);
        this.rego = rego;
        this.path = path;
        this.value = value;
        this.variable = variable;
        this.traceMessage = traceMessage
        this.pathRules = pathRules;
    }
}

export class BranchRuleResult extends RegoRuleResult {
    public readonly branch: SimpleRuleResult[];

    constructor(constraintId: string, branch: SimpleRuleResult[]) {
        super(constraintId);
        this.branch = branch;
    }
}



export abstract class BaseRegoRuleGenerator {
    abstract generateResult(): RegoRuleResult[]
}