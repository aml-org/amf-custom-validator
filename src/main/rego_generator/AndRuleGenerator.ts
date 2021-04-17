import {AndRule} from "../model/rules/AndRule";
import {AtomicRule, Rule} from "../model/Rule";
import {BaseRegoComplexRuleGenerator, BaseRegoRuleGenerator, RegoRuleResult} from "./BaseRegoRuleGenerator";
import {Expression} from "../model/Expression";
import {InRule} from "../model/constraints/InRule";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {PatternRule} from "../model/constraints/PatternRule";
import {InRuleGenerator} from "./constraints/InRuleGenerator";
import {PatternRuleGenerator} from "./constraints/PatternRuleGenerator";
import {MinCountRuleGenerator} from "./constraints/MinCountRuleGenerator";
import {NestedRule} from "../model/constraints/NestedRule";
import {NestedRuleGenerator} from "./constraints/NestedRuleGenerator";

export class AndRuleGenerator extends BaseRegoComplexRuleGenerator {
    private rule: AndRule;
    private expression: Expression;

    constructor(expression: Expression, rule: AndRule) {
        super();
        this.rule = rule;
        this.expression = expression;
        this.validRule();
    }

    generateResult(): RegoRuleResult[] {
        return this.generateNodeSet();
    }

    private validRule() {
        if (this.rule.body.find((r) => !(r instanceof AtomicRule)) != null) {
            throw new Error("All rules in the AND rule must be atomic, canonical form is required");
        }
    }


    private generateNodeSet(): RegoRuleResult[] {
        const rego: RegoRuleResult[] = [];
        this.rule.body.forEach((rule) => {
            rego.push(this.dispatchRule(rule));
        });
        return rego;
    }

    private dispatchRule(rule: Rule): RegoRuleResult {
        if (rule instanceof InRule) {
            return new InRuleGenerator(rule).generateResult();
        } else if (rule instanceof MinCountRule) {
            return new MinCountRuleGenerator(rule).generateResult();
        } else if (rule instanceof PatternRule) {
            return new PatternRuleGenerator(rule).generateResult();
        } else if (rule instanceof NestedRule) {
            return new NestedRuleGenerator(rule).generateResult();
        } else {
            throw new Error(`Unsupported rule ${rule}`);
        }
    }


}