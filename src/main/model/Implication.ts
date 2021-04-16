import {ComplexRule, Rule, Statement, Variable} from "./Rule";
import {AndRule} from "./rules/AndRule";
import {OrRule} from "./rules/OrRule";

export class Implication extends ComplexRule {
    private variable: Variable;
    private head: Rule;

    constructor(negation: boolean, variable: Variable, head: Rule, body: Rule) {
        super(negation);
        this.variable = variable;
        this.body = [body]
        this.head = head;
    }

    negation(): Rule {
        if (this.negated) {
            return new Implication(false, this.variable, this.head, this.body[0])
        } else {
            const andRule = new AndRule(false);
            const negatedBody = <Rule>this.body[0].negation();
            andRule.withBody([this.head, negatedBody]);
            return andRule;
        }
    }

    toRego(): string {
        return "";
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}( ${this.head.toString()} → ${this.body.toString()} )`
    }

    toCanonicalForm(): Statement {
        if (this.negated) {
            const and = this.negation();
            return and.toCanonicalForm();
        } else {
            const or = new OrRule(false);
            const negatedHead = <Rule>this.head.negation();
            or.withBody([negatedHead, this.body[0]]);
            return or.toCanonicalForm();
        }


    }
}