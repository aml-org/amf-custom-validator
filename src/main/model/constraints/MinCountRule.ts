import {AtomicRule, Rule, Variable} from "../Rule";

export class MinCountRule extends AtomicRule {
    constructor(negated: boolean, variable: Variable, path: string[], argument: any) {
        super(negated, variable, "minCount", path, argument);
    }

    toRego(): string {
        return "";
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}MinCount(${this.variable.name},'${this.path.join("/")}',${this.argument})`
    }

    negation(): Rule {
        return new MinCountRule(!this.negated, this.variable, this.path, this.argument);
    }
}