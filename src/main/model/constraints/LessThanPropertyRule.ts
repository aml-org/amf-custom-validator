import {AtomicRule, Variable} from "../Rule";

export class LessThanPropertyRule extends AtomicRule {

    constructor(negated: boolean, variable: Variable, path: string[], argument: any) {
        super(negated, variable, "lessThanProperty", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}(Property(${this.variable.name},'${this.path.join("/")}') < Property(${this.variable.name},'${this.argument}'))`
    }

}