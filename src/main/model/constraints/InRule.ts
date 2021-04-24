import {AtomicRule, Variable} from "../Rule";
import {PropertyPath} from "../../profile_parser/PathParser";

export class InRule extends AtomicRule {

    constructor(negated: boolean, variable: Variable, path: PropertyPath, argument: any) {
        super(negated, variable, "pattern", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        const vs = this.argument.map((v) => v.toString())
        return `${negation}In(${this.variable.name},'${this.path.source.replace(/\./g,":")}',[${vs}])`
    }

}