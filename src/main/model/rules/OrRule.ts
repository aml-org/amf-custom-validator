import {ComplexRule} from "../Rule";
import {AndRule} from "./AndRule";

export class OrRule extends ComplexRule {
    constructor(negated: boolean) {
        super(negated);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        const body = this.body.map((rule) => rule.toString()).join(" ∨ ")
        return `${negation}( ${body} )`
    }

    negation() {
        const andRule = new AndRule(!this.negated)
        andRule.withBody(this.body.map((r) => r.negation()))
        return andRule;
    }
}
