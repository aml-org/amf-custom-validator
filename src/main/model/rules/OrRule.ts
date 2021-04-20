import {ComplexRule, Rule} from "../Rule";
import {AndRule} from "./AndRule";
import {Canonical} from "../CanonicalCheck";
import {Expression} from "../Expression";

export class OrRule extends ComplexRule {
    constructor(negated: boolean) {
        super(negated);
    }

    negation(): Rule {
        return new OrRule(!this.negated).withBody(this.body)
    }

    deMorgan(): Rule {
        return new AndRule(!this.negated).withBody(<Rule[]>this.body.map((e) => e.negation()));
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
            return <Rule>this.deMorgan().toCanonicalForm();
        } else {
            // flatten nested ORs
            let flattenedBody = []
            this.body.map((e) => {
                if (e instanceof OrRule && !e.negated) {
                    flattenedBody = flattenedBody.concat(e.body);
                } else {
                    flattenedBody.push(e);
                }
            });
            this.body = flattenedBody;

            if (this.body.length === 1) {
                return <Rule>this.body[0].toCanonicalForm();
            } else if (Canonical.check(this)) {
                return this;
            } else {
                const tmp = this.body.map((r) => <Rule>r.toCanonicalForm());
                const tmpOr = new OrRule(this.negated).withBody(tmp);
                if (Canonical.check(tmpOr)) {
                    return tmpOr
                } else {
                    let acc: Rule = <Rule>tmp.shift();
                    tmp.forEach((e) => {
                        acc = this.distributeOr(acc, e);
                    });
                    return acc;
                }
            }
        }
    }

    distributeOr(exp1: Rule, exp2: Rule): Rule {
        // accumulate expressions
        if (exp1 instanceof Expression && exp2 instanceof Expression) {
            exp2.variables.forEach((v) => exp1.variables.push(v));
            const body = this.distributeOr(exp1.rule, exp2.rule);
            exp1.rule = body;
            return exp1;
        } else if (exp1 instanceof Expression) {
            const body = this.distributeOr(exp1.rule, exp2);
            exp1.rule = body;
            return exp1;
        } else if (exp2 instanceof Expression) {
            const body = this.distributeOr(exp1, exp2.rule);
            exp2.rule = body;
            return exp2;
        }

        // distribute
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
