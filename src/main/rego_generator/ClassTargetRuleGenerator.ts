import {BaseRegoAtomicRuleGenerator, BaseRegoRuleGenerator, RegoRuleResult} from "./BaseRegoRuleGenerator";
import {ClassTarget} from "../model/constraints/ClassTarget";

export class ClassTargetRuleGenerator extends BaseRegoAtomicRuleGenerator {
    private rule: ClassTarget;

    constructor(classTargetRule: ClassTarget) {
        super();
        this.rule = classTargetRule;
    }

    generateResult(): RegoRuleResult {
        if (this.rule.negated) {
            return {
                constraintId: "classTarget",
                path: "",
                rego: [`target_class_negated[${this.rule.variable.name}] with data.class as "${this.rule.argument}"`],
                value: this.rule.variable.name,
            }
        } else {
            return {
                constraintId: "classTarget",
                path: "",
                rego: [`target_class[${this.rule.variable.name}] with data.class as "${this.rule.argument}"`],
                value: this.rule.variable.name
            }
        }
    }

}