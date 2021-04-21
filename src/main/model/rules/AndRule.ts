import {ComplexRule} from "../Rule";

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
}