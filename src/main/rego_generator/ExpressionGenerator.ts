import {Quantification, Rule} from "../model/Rule";
import {Expression} from "../model/Expression";
import {OrRule} from "../model/rules/OrRule";
import {AndRule} from "../model/rules/AndRule";
import {AndRuleGenerator} from "./AndRuleGenerator";
import {ClassTarget} from "../model/constraints/ClassTarget";
import {RegoRuleResult} from "./BaseRegoRuleGenerator";
import {ClassTargetRuleGenerator} from "./ClassTargetRuleGenerator";

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
                return this.generateAndClause(rule);
            } else {
                throw new Error("Unsupported expression: " + this.expression.toString());
            }
        });

        return regoRules;
    }

    private generateAndClause(rule: AndRule): string {
        const result = this.findClassTargetMapping(rule);
        const classTarget = <ClassTarget>result[0];
        const filteredRule = <AndRule>result[1];
        const regoResults = new AndRuleGenerator(this.expression, filteredRule).generateResult();
        return this.wrapRegoResult(classTarget, regoResults);
    }

    private findClassTargetMapping(rule: AndRule) {
        let classTargetRule: ClassTarget|null = null;
        let remainingRules: Rule[] = [];
        rule.body.forEach((r) => {
            if (r instanceof ClassTarget) {
                classTargetRule = r;
            } else {
                remainingRules.push(r)
            }
        });
        if (classTargetRule == null) {
            throw new Error("Cannot generate AND top-level clause without a classTarget rule");
        }
        const filteredAnd = new AndRule(rule.negated).withBody(remainingRules);
        return [classTargetRule, filteredAnd];
    }

    private wrapRegoResult(classTargetRule: ClassTarget, regoRuleResults: RegoRuleResult[]): string {
        const level = this.expression.level;
        const acc = [];
        const resultBindings: string[] = [];
        let i = 0;

        const classTargetResult = new ClassTargetRuleGenerator(classTargetRule).generateResult();
        const classTargetVariable = classTargetResult.value;

        acc.push(`${level.toLowerCase()}[matches] {`);
        classTargetResult.rego.forEach((l)=> acc.push(" " + l));

        regoRuleResults.forEach((regoResult) => {
            const bindingResult = `_result_${i++}`
            resultBindings.push(bindingResult)

            const matchesLine = `  ${bindingResult} := trace("${regoResult.constraintId}", ${regoResult.value}, "${regoResult.traceMessage || ''}")`

            regoResult.rego.forEach((line) => acc.push("  " + line));
            acc.push(matchesLine);
        });

        acc.push(`  matches := error("${this.expression.name}", ${classTargetVariable}, "${this.expression.message}", [${resultBindings.join(",")}])`);
        acc.push('}');
        return acc.join("\n");
    }

    private canTransform(expression: Expression): boolean {
        return expression.variables[0]!.quantification === Quantification.Exist && expression.negated === true && (expression.rule instanceof OrRule)
    }

    private canonicalDisjunction(): OrRule {
        return <OrRule>this.expression.rule;
    }
}