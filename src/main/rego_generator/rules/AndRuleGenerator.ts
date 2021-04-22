import {AndRule} from "../../model/rules/AndRule";
import {Rule} from "../../model/Rule";
import {InRule} from "../../model/constraints/InRule";
import {MinCountRule} from "../../model/constraints/MinCountRule";
import {MinCountRuleGenerator} from "../constraints/MinCountRuleGenerator";
import {PatternRule} from "../../model/constraints/PatternRule";
import {PatternRuleGenerator} from "../constraints/PatternRuleGenerator";
import {Expression} from "../../model/Expression";
import {ExpressionGenerator} from "../ExpressionGenerator";
import {OrRule} from "../../model/rules/OrRule";
import {OrRuleGenerator} from "./OrRuleGenerator";
import {InRuleGenerator} from "../constraints/InRuleGenerator";
import {BaseRegoRuleGenerator, BranchRuleResult, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";


export class AndRuleGenerator extends BaseRegoRuleGenerator {
    private rule: AndRule;

    constructor(rule: AndRule) {
        super();
        this.rule = rule;
    }

    generateResult(): BranchRuleResult[] {
        if (this.rule.negated) {
            const orRule = this.rule.negation()
            return new OrRuleGenerator(orRule).generateResult();
        } else {
            let branches: BranchRuleResult[] = [];
            this.rule.body.forEach((rule) => {
                this.dispatchRule(rule).forEach((result) => {
                    if (result instanceof SimpleRuleResult) {
                        branches.push(new BranchRuleResult(result.constraintId, [result]));
                    } else {
                        branches.push(<BranchRuleResult>result)
                    }
                })
            })
            return branches;
        }
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