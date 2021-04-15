import {AndRule} from "../model/rules/AndRule";
import {AtomicRule, Rule} from "../model/Rule";
import {ClassTarget} from "../model/constraints/ClassTarget";
import {BaseRegoRuleGenerator, RegoRuleResult} from "./BaseRegoRuleGenerator";
import {Expression} from "../model/Expression";
import {ClassTargetRuleGenerator} from "./ClassTargetRuleGenerator";
import {InRule} from "../model/constraints/InRule";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {PatternRule} from "../model/constraints/PatternRule";
import {InRuleGenerator} from "./constraints/InRuleGenerator";
import {PatternRuleGenerator} from "./constraints/PatternRuleGenerator";
import {MinCountRuleGenerator} from "./constraints/MinCountRuleGenerator";

export class AndRuleGenerator extends BaseRegoRuleGenerator {
    private rule: AndRule;
    private ruleGroups: {[variable: string]: Rule[]} = {};
    private classTargetVariable: string|undefined;
    private regoText: string[] = [];
    private expression: Expression;

    constructor(expression: Expression, rule: AndRule) {
        super();
        this.rule = rule;
        this.expression = expression;
        this.validRule();
    }

    generate(): string {
        this.groupRules();
        this.generateFromGroups();
        return this.regoText.join("\n\n");
    }

    private validRule() {
        if (this.rule.body.find((r) => !(r instanceof AtomicRule)) != null) {
            throw new Error("All rules in the AND rule must be atomic, canonical form is required");
        }
    }

    private groupRules() {
        this.rule.body.forEach((r) => {
            const atomicRule = <AtomicRule> r;
            if (atomicRule instanceof ClassTarget) {
                this.classTargetVariable = atomicRule.variable.name
            }
            const group = this.ruleGroups[atomicRule.variable.name] || [];
            group.push(atomicRule);
            this.ruleGroups[atomicRule.variable.name] = group;
        });
    }

    private generateFromGroups() {
        this.generateTargetNodeSet().forEach((line) => this.regoText.push(line));
        if (Object.keys(this.ruleGroups).length > 0) {
            throw new Error("More than one rule group per AND condition not supported yet");
        }
    }

    private generateTargetNodeSet(): string[] {
        if (this.classTargetVariable) {
            const rego = [];
            const classTargetGroup = this.ruleGroups[this.classTargetVariable];
            delete this.ruleGroups[this.classTargetVariable];
            const classTargetRule = <ClassTarget>classTargetGroup.find((r) => r instanceof ClassTarget)!;
            const remainingRules = classTargetGroup.filter((r) => !(r instanceof ClassTarget));
            const classTargetLine = new ClassTargetRuleGenerator(classTargetRule).generate();
            remainingRules.forEach((rule) => {
                rego.push(this.wrapRegoResult(classTargetLine, this.dispatchRule(rule)));
            });

            return rego;
        } else {
            return [];
        }
    }

    private dispatchRule(rule: Rule) {
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

    private wrapRegoResult(classTargetLine: string, regoRuleResult: RegoRuleResult) {
        const level = this.expression.level;
        const firstLine = `${level.toLowerCase()}[matches] {`
        const targetLine = "  " + classTargetLine;
        const matchesLine = `  matches := error("${this.expression.name}","${regoRuleResult.constraintId}",${this.classTargetVariable}, ${regoRuleResult.value}, "${regoRuleResult.traceMessage || ''}", "${this.expression.message}")`
        const lastLine = `}`

        const acc = [];
        acc.push(firstLine);
        acc.push(targetLine);
        regoRuleResult.rego.forEach((l) => acc.push("  " + l));
        acc.push(matchesLine);
        acc.push(lastLine);
        return acc.join("\n");
    }
}