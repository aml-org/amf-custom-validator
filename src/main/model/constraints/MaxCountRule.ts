import {AtomicRule, Variable} from "../Rule";
import {PropertyPath} from "../../profile_parser/PathParser";

export class MaxCountRule extends AtomicRule {
    constructor(negated: boolean, variable: Variable, path: PropertyPath, argument: any) {
        super(negated, variable, "maxCount", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}MaxCount(${this.variable.name},'${this.path.source.replace(/\./g, ":")}',${this.argument})`
    }

}