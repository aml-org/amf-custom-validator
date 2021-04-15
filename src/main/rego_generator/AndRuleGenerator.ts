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

export class AndRuleGenerator extends BaseRegoComplexRuleGenerator {
    private rule: AndRule;
    private ruleGroups: {[variable: string]: Rule[]} = {};
    private expression: Expression;

    constructor(expression: Expression, rule: AndRule) {
        super();
        this.rule = rule;
        this.expression = expression;
        this.validRule();
    }

    generateResult(): RegoRuleResult[] {
        this.groupRules();
        return this.generateFromGroups();
    }

    private validRule() {
        if (this.rule.body.find((r) => !(r instanceof AtomicRule)) != null) {
            throw new Error("All rules in the AND rule must be atomic, canonical form is required");
        }
    }

    private groupRules() {
        this.rule.body.forEach((r) => {
            const atomicRule = <AtomicRule> r;
            const group = this.ruleGroups[atomicRule.variable.name] || [];
            group.push(atomicRule);
            this.ruleGroups[atomicRule.variable.name] = group;
        });
    }

    private generateFromGroups() {
        if (Object.keys(this.ruleGroups).length > 1) {
            throw new Error("More than one rule group per AND condition not supported yet");
        }
        const results: RegoRuleResult[] = [];
        Object.keys(this.ruleGroups).map((variable) => {
            const group = this.ruleGroups[variable];
            this.generateNodeSet(group).forEach((result) => results.push(result));
        });
        return results;
    }

    private generateNodeSet(group: Rule[]): RegoRuleResult[] {
        const rego: RegoRuleResult[] = [];
        group.forEach((rule) => {
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
        } else {
            throw new Error(`Unsupported rule ${rule}`);
        }
    }


}