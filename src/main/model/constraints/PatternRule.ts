import {AtomicRule, Variable} from "../Rule";
import {PropertyPath} from "../../profile_parser/PathParser";


export class PatternRule extends AtomicRule {
    constructor(negated: boolean, variable: Variable, path: PropertyPath, argument: any) {
        super(negated, variable, "pattern", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}Pattern(${this.variable.name},'${this.path.source.replace(/\./g, ":")}','${this.argument}')`
    }
}