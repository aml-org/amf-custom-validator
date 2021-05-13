import {AtomicRule, Variable} from "../Rule";
import {PropertyPath} from "../../profile_parser/PathParser";

export type RegoRuleArgument = {
    message: string,
    code: string
}

export class RegoRule extends AtomicRule {
    constructor(negated: boolean, variable:Variable, path: PropertyPath, argument: RegoRuleArgument) {
        super(negated, variable, "rego", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}Rego(${this.variable.name},'${this.path.source.replace(/\./g, ":")}', '${this.argument.message.substring(0,5)}...')`
    }
}