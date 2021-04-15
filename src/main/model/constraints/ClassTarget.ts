import {AtomicRule, Statement, Variable} from "../Rule";

export class ClassTarget extends AtomicRule {

    constructor(negated: boolean, variable: Variable, argument: any) {
        super(negated, variable, "classTarget", [], argument.replace(".", ":"));
    }


    negation(): Statement {
        return new ClassTarget(!this.negated, this.variable, this.argument);
    }

    toRego(): string {
        return "";
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}Class(${this.variable.name},'${this.argument}')`
    }

}