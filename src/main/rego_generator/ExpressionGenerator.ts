import {Quantification, Rule} from "../model/Rule";
import {Expression} from "../model/Expression";
import {ClassTarget} from "../model/constraints/ClassTarget";
import {
    BaseRegoRuleGenerator,
    BranchRuleResult,
    RegoRuleResult,
    SimpleRuleResult
} from "./BaseRegoRuleGenerator";
import {ClassTargetRuleGenerator} from "./ClassTargetRuleGenerator";
import {Implication} from "../model/Implication";
import {NestedRule} from "../model/constraints/NestedRule";
import {NestedRuleGenerator} from "./constraints/NestedRuleGenerator";
import {InRuleGenerator} from "./constraints/InRuleGenerator";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {MinCountRuleGenerator} from "./constraints/MinCountRuleGenerator";
import {PatternRule} from "../model/constraints/PatternRule";
import {PatternRuleGenerator} from "./constraints/PatternRuleGenerator";
import {AndRule} from "../model/rules/AndRule";
import {OrRule} from "../model/rules/OrRule";
import {OrRuleGenerator} from "./OrRuleGenerator";
import {AndRuleGenerator} from "./AndRuleGenerator";
import {InRule} from "../model/constraints/InRule";

export class ExpressionGenerator extends BaseRegoRuleGenerator {

    private expression: Expression;

    constructor(expression: Expression) {
        super()
        this.expression = expression;
        this.checkValidExpression()
    }

    generate(): string {
        const variable = this.expression.variables[0];
        if (variable.quantification == Quantification.ForAll) {
            if (this.topLevel()) {
                return this.generateTopLevel()
            } else if(this.nested()) {
                throw new Error("Only universal classTargets can be generated as stand-alone validations")
            }
        } else {
            throw new Error("Existential not supported yet");
        }
    }

    generateResult(): BranchRuleResult[] {
        const variable = this.expression.variables[0];
        if (variable.quantification == Quantification.ForAll) {
            if (this.topLevel()) {
                throw new Error("Universally quantified classTargets cannot be generated as intermediate Rego results")
            } else if(this.nested()) {
                return this.generatedNested();
            }
        } else {
            throw new Error("Existential not supported yet");
        }
    }

    private topLevel() {
        return this.implication().head instanceof ClassTarget
    }

    private nested() {
        return this.implication().head instanceof NestedRule
    }

    private implication(): Implication {
        return <Implication>this.expression.rule;
    }

    private checkValidExpression() {
        if (!(this.expression.rule instanceof Implication)) {
            throw new Error("Only implications are supported");
        }
    }


    private generateTopLevel() {
        const bodyResult = this.generateBody();
        return this.wrapRegoResult(<ClassTarget>this.implication().head, bodyResult)
    }

    private generatedNested() {
        const bodyResult = this.generateBody();
        const headResult: SimpleRuleResult = new NestedRuleGenerator(<NestedRule>this.implication().head).generatedNestedResult()[0];
        return bodyResult.map((result) => {
            if (result instanceof SimpleRuleResult) {
                return new BranchRuleResult("nested exp", [headResult, result]);
            } else {
                return new BranchRuleResult(result.constraintId, [headResult].concat((<BranchRuleResult>result).branch))
            }
        })
    }

    private generateBody() {
        const body = this.implication().body;
        const acc: RegoRuleResult[] = [];
        body.forEach((r) => this.dispatchRule(r).forEach((result) => acc.push(result)))
        return acc;
    }

    private wrapRegoResult(classTargetRule: ClassTarget, regoRuleResults: RegoRuleResult[]): string {
        const level = this.expression.level;

        const classTargetResult = <SimpleRuleResult>new ClassTargetRuleGenerator(classTargetRule).generateResult()[0];
        const classTargetVariable = classTargetResult.value;

        const branches: BranchRuleResult[] = regoRuleResults.map((result) => {
            if (result instanceof SimpleRuleResult) {
                return new BranchRuleResult("exp", [result]);
            } else {
                return <BranchRuleResult>result
            }
        });

        const validations = branches.map((branch) => {
            const acc = [];
            acc.push(`${level.toLowerCase()}[matches] {`);
            classTargetResult.rego.forEach((l)=> acc.push(" " + l));
            let i = 0;
            const resultBindings: string[] = [];

            branch.branch.forEach((regoResult) => {
                const bindingResult = `_result_${i++}`
                resultBindings.push(bindingResult)

                const matchesLine = `  ${bindingResult} := trace("${regoResult.constraintId}", "${regoResult.path}", ${regoResult.value}, "${regoResult.traceMessage || ''}")`

                regoResult.rego.forEach((line) => acc.push("  " + line));
                acc.push(matchesLine);
            });

            acc.push(`  matches := error("${this.expression.name}", ${classTargetVariable}, "${this.expression.message}", [${resultBindings.join(",")}])`);
            acc.push('}');
            return acc.join("\n");
        });

        return validations.join("\n\n");
    }

    dispatchRule(rule: Rule): RegoRuleResult[] {
        if (rule instanceof InRule) {
            return new InRuleGenerator(rule).generateResult();
        } else if (rule instanceof MinCountRule) {
            return new MinCountRuleGenerator(rule).generateResult();
        } else if (rule instanceof PatternRule) {
            return new PatternRuleGenerator(rule).generateResult();
        } else if (rule instanceof Expression) {
            return new ExpressionGenerator(rule).generateResult();
        } else if (rule instanceof AndRule) {
            return new AndRuleGenerator(rule).generateResult();
        } else if (rule instanceof OrRule) {
            return new OrRuleGenerator(rule).generateResult();
        } else {
            throw new Error(`Unsupported rule ${rule}`);
        }
    }
}