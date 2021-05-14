import {AndRule} from "../../model/rules/AndRule";
import {OrRuleGenerator} from "./OrRuleGenerator";
import {BaseRegoRuleGenerator, BranchRuleResult, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RuleDispatcher} from "../RuleDispatcher";


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
                RuleDispatcher.dispatchRule(rule).forEach((result) => {
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

}