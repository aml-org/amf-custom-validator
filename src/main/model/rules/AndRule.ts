import {ComplexRule} from "../Rule";
import {OrRule} from "./OrRule";

export class AndRule extends ComplexRule {
    constructor(negated: boolean) {
        super(negated);
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        const body = this.body.map((rule) => rule.toString()).join(" ∧ ")
        return `${negation}( ${body} )`
    }

    negation() {
        const orRule = new OrRule(!this.negated)
        orRule.withBody(this.body.map((r) => r.negation()))
        return orRule;
    }
}