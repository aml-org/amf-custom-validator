import {ComplexRule, Rule, Variable} from "./Rule";

export class Implication extends ComplexRule {
    private variable: Variable;
    public readonly head: Rule;

    constructor(negation: boolean, variable: Variable, head: Rule, body: Rule) {
        super(negation);
        this.variable = variable;
        this.body = [body]
        this.head = head;
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}( ${this.head.toString()} → ${this.body.toString()} )`
    }
}