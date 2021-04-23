import {Quantification, Rule, VariableCardinality} from "../model/Rule";
import {Expression} from "../model/Expression";
import {
    BaseRegoRuleGenerator,
    BranchRuleResult,
    RegoRuleResult,
    SimpleRuleResult
} from "./BaseRegoRuleGenerator";
import {Implication} from "../model/Implication";
import {InRuleGenerator} from "./constraints/InRuleGenerator";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {MinCountRuleGenerator} from "./constraints/MinCountRuleGenerator";
import {PatternRule} from "../model/constraints/PatternRule";
import {PatternRuleGenerator} from "./constraints/PatternRuleGenerator";
import {AndRule} from "../model/rules/AndRule";
import {OrRule} from "../model/rules/OrRule";
import {InRule} from "../model/constraints/InRule";
import {ClassTarget} from "../model/mappers/ClassTarget";
import {NestedRule} from "../model/mappers/NestedRule";
import {NestedRuleGenerator} from "./mappers/NestedRuleGenerator";
import {ClassTargetRuleGenerator} from "./mappers/ClassTargetRuleGenerator";
import {AndRuleGenerator} from "./rules/AndRuleGenerator";
import {OrRuleGenerator} from "./rules/OrRuleGenerator";
import {LessThanPropertyRule} from "../model/constraints/LessThanPropertyRule";
import {LessThanPropertyGenerator} from "./constraints/LessThanPropertyGenerator";

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
            return this.generateNestedQualified()
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
        let bodyResult;
        if (this.expression.negated) {
            bodyResult = this.generatedBodyNegated();
        } else {
            bodyResult = this.generateBody();
        }
        return this.wrapRegoResult(<ClassTarget>this.implication().head, bodyResult)
    }

    private generatedNested() {
        const bodyResult = this.generateBody();
        const nestedRule = <NestedRule>this.implication().head
        const headResult: SimpleRuleResult = new NestedRuleGenerator(nestedRule).generateResult()[0];
        return bodyResult.map((result) => {
            if (result instanceof SimpleRuleResult) {
                return this.wrapNestedRegoResult(nestedRule.child.name, headResult, new BranchRuleResult("nested exp", [result]));
            } else {
                return this.wrapNestedRegoResult(nestedRule.child.name, headResult, <BranchRuleResult>result);

            }
        })
    }

    private generateNestedQualified() {
        const bodyResult = this.generateBody();
        const nestedRule = <NestedRule>this.implication().head
        const headResult: SimpleRuleResult = new NestedRuleGenerator(nestedRule).generateResult()[0];
        return bodyResult.map((result) => {
            if (result instanceof SimpleRuleResult) {
                return this.wrapNestedQualifiedRegoResult(nestedRule.child.name, nestedRule.child.cardinality, headResult, new BranchRuleResult("nested exp", [result]));
            } else {
                return this.wrapNestedQualifiedRegoResult(nestedRule.child.name, nestedRule.child.cardinality, headResult, <BranchRuleResult>result);

            }
        })
    }

    private wrapNestedRegoResult(nestedVariable: string, headResult: SimpleRuleResult, bodyResult: BranchRuleResult) {
        const rego = [].concat(headResult.rego);
        const variable = headResult.variable
        rego.push(`${variable}_errors = [ ${variable}_error |`) // we generated a comprehension to look for errors in the collection
        rego.push(`  ${nestedVariable} = ${variable}[_]`) // the underlying rules expect the quantified variable that was passed on profile parsing
        this.wrapBranch(bodyResult, `${variable}_error`, nestedVariable, rego)
        rego.push(']');
        if (this.expression.negated) {
            rego.push(`count(${variable}_errors) == 0`)
        } else {
            rego.push(`count(${variable}_errors) > 0`)

        }
        const nestedSimpleResult = new SimpleRuleResult(
            "nested",
            rego,
            headResult.path,
            `{"failed": count(${variable}_errors), "success":(count(${variable}) - count(${variable}_errors))}`,
            `${variable}_errors`,
            {"code" :`[e | e := ${variable}_errors[_].trace]`}
        );

        return new BranchRuleResult("nested", [nestedSimpleResult]);
    }

    private wrapNestedQualifiedRegoResult(nestedVariable: string, cardinality: VariableCardinality, headResult: SimpleRuleResult, bodyResult: BranchRuleResult) {
        const rego = [].concat(headResult.rego);
        const variable = headResult.variable
        rego.push(`${variable}_errors = [ ${variable}_error |`) // we generated a comprehension to look for errors in the collection
        rego.push(`  ${nestedVariable} = ${variable}[_]`) // the underlying rules expect the quantified variable that was passed on profile parsing
        this.wrapBranch(bodyResult, `${variable}_error`, nestedVariable, rego)
        rego.push(']');
        if (this.expression.negated) {
            rego.push(cardinality.toRego(variable, `${variable}_errors`, true));
        } else {
            rego.push(cardinality.toRego(variable, `${variable}_errors`, false));

        }

        const nestedSimpleResult = new SimpleRuleResult(
            "nested",
            rego,
            headResult.path,
            `{"failed": count(${variable}_errors), "success":(count(${variable}) - count(${variable}_errors))}`,
            `${variable}_errors`,
            {"code" :`[e | e := ${variable}_errors[_].trace]`}
        );

        return new BranchRuleResult("nested", [nestedSimpleResult]);
    }

    private generateBody() {
        const body = this.implication().body;
        const acc: RegoRuleResult[] = [];
        body.forEach((r) => this.dispatchRule(r).forEach((result) => acc.push(result)))
        return acc;
    }

    private generatedBodyNegated() {
        const body = this.implication().body;
        let negatedRule: Rule;
        if (body.length == 1) {
            negatedRule = body[0].negation();
        } else {
            negatedRule = new AndRule(false).withBody(body).negation();
        }
        const acc: RegoRuleResult[] = [];
        this.dispatchRule(negatedRule).forEach((result) => acc.push(result))
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
            this.wrapBranch(branch, 'matches', classTargetVariable, acc);
            acc.push('}');
            return acc.join("\n");
        });

        return validations.join("\n\n");
    }

    wrapBranch(branch: BranchRuleResult, matchesVariable: string, mappingVariable: string, acc: string[]) {
        let i = 0;
        const resultBindings: string[] = [];
        branch.branch.forEach((regoResult) => {
            const bindingResult = `_result_${i++}`
            resultBindings.push(bindingResult)

            let traceMessage = regoResult.traceMessage || '';
            if (traceMessage["code"]) {
                traceMessage = traceMessage["code"]
            } else {
                traceMessage = `"${traceMessage}"`
            }

            const matchesLine = `  ${bindingResult} := trace("${regoResult.constraintId}", "${regoResult.path}", ${regoResult.value}, ${traceMessage})`

            regoResult.rego.forEach((line) => acc.push("  " + line));
            acc.push(matchesLine);
        });
        acc.push(`  ${matchesVariable} := error("${this.expression.name}", ${mappingVariable}, "${this.expression.message}", [${resultBindings.join(",")}])`);
    }

    dispatchRule(rule: Rule): RegoRuleResult[] {
        if (rule instanceof InRule) {
            return new InRuleGenerator(rule).generateResult();
        } else if (rule instanceof MinCountRule) {
            return new MinCountRuleGenerator(rule).generateResult();
        } else if (rule instanceof PatternRule) {
            return new PatternRuleGenerator(rule).generateResult();
        } else if (rule instanceof LessThanPropertyRule) {
            return new LessThanPropertyGenerator(rule).generateResult();
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