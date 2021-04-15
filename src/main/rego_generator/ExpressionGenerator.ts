import {AtomicRule, Quantification} from "../model/Rule";
import {Expression} from "../model/Expression";
import {OrRule} from "../model/rules/OrRule";
import {AndRule} from "../model/rules/AndRule";
import {AndRuleGenerator} from "./AndRuleGenerator";

export class ExpressionGenerator {
    private expression: Expression;

    constructor(expression: Expression) {
        this.expression = expression;
        if (!this.canTransform(expression)) {
            throw new Error("Expression not supported")
        }
    }

    generate(): string[] {
        const regoRules = this.canonicalDisjunction().body.map((rule) => {
            if (rule instanceof AndRule) {
                return new AndRuleGenerator(this.expression, rule).generate();
            } else if (rule instanceof AtomicRule) {
                throw new Error("Simple conjunction of rules not supported yet")
                //new AtomicRuleGenerator(rule).generate();
            } else {
                throw new Error("Unsupported expression: " + this.expression.toString());
            }
        });

        return regoRules;
    }

    private canTransform(expression: Expression): boolean {
        return expression.variables[0]!.quantification === Quantification.Exist && expression.negated === true && (expression.rule instanceof OrRule)
    }

    private canonicalDisjunction(): OrRule {
        return <OrRule>this.expression.rule;
    }
}