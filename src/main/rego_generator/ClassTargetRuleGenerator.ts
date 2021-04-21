import {BaseRegoRuleGenerator, RegoRuleResult, SimpleRuleResult} from "./BaseRegoRuleGenerator";
import {ClassTarget} from "../model/constraints/ClassTarget";

export class ClassTargetRuleGenerator extends BaseRegoRuleGenerator {
    private rule: ClassTarget;

    constructor(classTargetRule: ClassTarget) {
        super();
        this.rule = classTargetRule;
    }

    generateResult(): RegoRuleResult[] {
        if (this.rule.negated) {
            return [new SimpleRuleResult(
                "classTarget",
                [`target_class_negated[${this.rule.variable.name}] with data.class as "${this.rule.argument}"`],
                "",
                this.rule.variable.name,
                this.rule.variable.name
            )];
        } else {
            return [new SimpleRuleResult(
                "classTarget",
                [`target_class[${this.rule.variable.name}] with data.class as "${this.rule.argument}"`],
                "",
                this.rule.variable.name,
                this.rule.variable.name
            )];
        }
    }

}