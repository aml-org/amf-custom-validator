import {AtomicRule, Variable} from "../Rule";


export class PatternRule extends AtomicRule {
    constructor(negated: boolean, variable: Variable, path: string[], argument: any) {
        super(negated, variable, "pattern", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}Pattern(${this.variable.name},'${this.path.join("/")}','${this.argument}')`
    }
}