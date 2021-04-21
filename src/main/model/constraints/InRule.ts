import {AtomicRule, Variable} from "../Rule";

export class InRule extends AtomicRule {

    constructor(negated: boolean, variable: Variable, path: string[], argument: any) {
        super(negated, variable, "pattern", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        const vs = this.argument.map((v) => v.toString())
        return `${negation}In(${this.variable.name},'${this.path.join("/")}',[${vs}])`
    }

}