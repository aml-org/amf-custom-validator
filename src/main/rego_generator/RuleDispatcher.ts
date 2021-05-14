import {Rule} from "../model/Rule";
import {RegoRuleResult} from "./BaseRegoRuleGenerator";
import {InRuleGenerator} from "./constraints/InRuleGenerator";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {MinCountRuleGenerator} from "./constraints/MinCountRuleGenerator";
import {PatternRule} from "../model/constraints/PatternRule";
import {PatternRuleGenerator} from "./constraints/PatternRuleGenerator";
import {RegoRule} from "../model/constraints/RegoRule";
import {RegoRuleGenerator} from "./constraints/RegoRuleGenerator";
import {LessThanPropertyRule} from "../model/constraints/LessThanPropertyRule";
import {LessThanPropertyGenerator} from "./constraints/LessThanPropertyGenerator";
import {Expression} from "../model/Expression";
import {ExpressionGenerator} from "./ExpressionGenerator";
import {AndRule} from "../model/rules/AndRule";
import {AndRuleGenerator} from "./rules/AndRuleGenerator";
import {OrRule} from "../model/rules/OrRule";
import {OrRuleGenerator} from "./rules/OrRuleGenerator";
import {InRule} from "../model/constraints/InRule";
import {MaxCountRule} from "../model/constraints/MaxCountRule";
import {MaxCountRuleGenerator} from "./constraints/MaxCountRuleGenerator";

export class RuleDispatcher {

    static dispatchRule(rule: Rule): RegoRuleResult[] {
        if (rule instanceof InRule) {
            return new InRuleGenerator(rule).generateResult();
        } else if (rule instanceof MinCountRule) {
            return new MinCountRuleGenerator(rule).generateResult();
        } else if (rule instanceof MaxCountRule) {
            return new MaxCountRuleGenerator(rule).generateResult();
        } else if (rule instanceof PatternRule) {
            return new PatternRuleGenerator(rule).generateResult();
        } else if (rule instanceof RegoRule) {
            return new RegoRuleGenerator(rule).generateResult();
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