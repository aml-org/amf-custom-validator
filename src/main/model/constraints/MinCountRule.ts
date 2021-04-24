import {AtomicRule, Variable} from "../Rule";
import {PropertyPath} from "../../profile_parser/PathParser";

export class MinCountRule extends AtomicRule {
    constructor(negated: boolean, variable: Variable, path: PropertyPath, argument: any) {
        super(negated, variable, "minCount", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}MinCount(${this.variable.name},'${this.path.source.replace(/\./g, ":")}',${this.argument})`
    }

}