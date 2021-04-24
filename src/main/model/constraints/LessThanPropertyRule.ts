import {AtomicRule, Variable} from "../Rule";
import {PropertyPath} from "../../profile_parser/PathParser";

export class LessThanPropertyRule extends AtomicRule {

    constructor(negated: boolean, variable: Variable, path: PropertyPath, argument: any) {
        super(negated, variable, "lessThanProperty", path, argument);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}(Property(${this.variable.name},'${this.path.source.replace(/\./g, ":")}') < Property(${this.variable.name},'${this.argument}'))`
    }

}