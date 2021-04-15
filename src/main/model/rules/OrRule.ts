import {AtomicRule, ComplexRule, Rule} from "../Rule";
import {AndRule} from "./AndRule";
import {Canonical} from "../CanonicalCheck";

export class OrRule extends ComplexRule {
    constructor(negated: boolean) {
        super(negated);
    }

    negation(): Rule {
        const body = <Rule[]>this.body.map((rule) => rule.negation());
        const newRule = new AndRule(this.negated);
        newRule.withBody(body);
        return newRule;
    }

    toRego(): string {
        return "";
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        const body = this.body.map((rule) => rule.toString()).join(" ∨ ")
        return `${negation}( ${body} )`
    }

    toCanonicalForm(): Rule {
        if (this.negated) {
            return <Rule>this.negation().toCanonicalForm();
        } else {
            if (this.body.length === 1) {
                return <Rule>this.body[0];
            } else if (Canonical.check(this)) {
                return this;
            } else {
                const tmp = this.body.map((r) => <Rule>r.toCanonicalForm());
                let acc: Rule = <Rule>tmp.shift();
                tmp.forEach((e) => {
                    acc = this.distributeOr(acc, e);
                });
                return acc;
            }
        }
    }

    distributeOr(exp1: Rule, exp2: Rule): Rule {
        if (exp1 instanceof AndRule && exp2 instanceof AndRule) {
            return new OrRule(false).withBody([exp1, exp2]);
        } else if (exp2 instanceof AndRule) {
            const orBodyWithAnds = exp2.body.map((orBody) => {
                return new AndRule(false).withBody([orBody, exp1]).toCanonicalForm();
            });
            return <Rule>(new OrRule(false).withBody(<Rule[]>orBodyWithAnds)).toCanonicalForm();
        } else if (exp1 instanceof AndRule) {
            const orBodyWithAnds = exp1.body.map((orBody) => {
                return new AndRule(false).withBody([orBody, exp2]).toCanonicalForm();
            });
            return <Rule>(new OrRule(false).withBody(<Rule[]>orBodyWithAnds)).toCanonicalForm();
        } else if (exp1 instanceof OrRule && exp2 instanceof OrRule) {
            return <Rule>new OrRule(false).withBody(exp1.body.concat(exp2.body)).toCanonicalForm();
        } else if (exp1 instanceof OrRule) {
            return <Rule>new OrRule(false).withBody(exp1.body.concat(exp2)).toCanonicalForm();
        } else if (exp2 instanceof OrRule) {
            return <Rule>new OrRule(false).withBody(exp2.body.concat(exp1)).toCanonicalForm();
        } else {
            return <Rule>(new OrRule(false).withBody([exp1, exp2]));
        }
    }
}
