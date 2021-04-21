import {OrRule} from "../model/rules/OrRule";
import {
    BaseRegoRuleGenerator,
    BranchRuleResult,
    RegoRuleResult,
    SimpleRuleResult
} from "./BaseRegoRuleGenerator";
import {Rule} from "../model/Rule";
import {InRule} from "../model/constraints/InRule";
import {InRuleGenerator} from "./constraints/InRuleGenerator";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {MinCountRuleGenerator} from "./constraints/MinCountRuleGenerator";
import {PatternRule} from "../model/constraints/PatternRule";
import {PatternRuleGenerator} from "./constraints/PatternRuleGenerator";
import {Expression} from "../model/Expression";
import {ExpressionGenerator} from "./ExpressionGenerator";
import {AndRule} from "../model/rules/AndRule";
import {AndRuleGenerator} from "./AndRuleGenerator";

export class OrRuleGenerator extends BaseRegoRuleGenerator {
    private rule: OrRule;

    constructor(rule: OrRule) {
        super();
        this.rule = rule;
    }

    generateResult(): BranchRuleResult[] {
        const rego: RegoRuleResult[][] = [];
        this.rule.body.forEach((rule) => {
            rego.push(this.dispatchRule(rule));
        });

        const regoBranches: RegoRuleResult[][] = [];
        const regoResults: SimpleRuleResult[] = [];
        rego.forEach((result) => {
            if (result[0] instanceof SimpleRuleResult) {
                regoResults.push(<SimpleRuleResult>result[0])
            } else {
                regoBranches.push(result)
            }
        });

        let regoExpandedBranches: RegoRuleResult[][] = [regoResults];
        regoBranches.forEach((branches) => {
            const acc: RegoRuleResult[][] = [];
            branches.forEach((branch) => {
                regoExpandedBranches.forEach((expanded) => {
                    acc.push(expanded.concat([branch]))
                });
            });
            regoExpandedBranches = acc;
        })

        return regoExpandedBranches.map((branch) => {
            const variables = [];
            const simpleResults: SimpleRuleResult[] = [];
            branch.forEach((result) => {
                if (result instanceof SimpleRuleResult) {
                    simpleResults.push(result);
                    variables.push(result.variable);
                } else if (result instanceof BranchRuleResult) {
                    result.branch.forEach((r) => simpleResults.push(r));
                    variables.push(result.branch[result.branch.length-1].variable);
                }
            });
            let unificationString: string[] = []
            for (let i=0; i<variables.length-1; i++) {
                unificationString.push(`${variables[i]} == ${variables[i+1]}`);
            }
            simpleResults.push({
                constraintId: "or",
                rego: unificationString,
                path: "",
                variable: variables[variables.length-1],
                value: variables[variables.length-1],
                traceMessage: "Failed or constraint"
            });
            return new BranchRuleResult("or", simpleResults);
        });
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