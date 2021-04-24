import {AtomicRule, Variable} from "../Rule";

export class ClassTarget extends AtomicRule {

    constructor(negated: boolean, variable: Variable, argument: any) {
        super(negated, variable, "classTarget", {and:[]}, argument.replace(".", ":"));
    }


    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}Class(${this.variable.name},'${this.argument}')`
    }

}