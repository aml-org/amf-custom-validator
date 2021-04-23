import {BaseRegoRuleGenerator, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {MinCountRule} from "../../model/constraints/MinCountRule";
import {genvar} from "../../VarGen";

export class LessThanPropertyGenerator extends BaseRegoRuleGenerator {
    private rule: MinCountRule;

    constructor(rule: MinCountRule) {
        super();
        this.rule = rule;
    }
    generateResult(): SimpleRuleResult[] {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name, "lessThanProperty_" + this.rule.valueMD5()).generatePropertyArray();
        const rego = pathResult.rego;

        const propAVariable = pathResult.variable;
        const propBVariable = genvar("lessThanPropertyValue");
        rego.push(`${propBVariable} = ${pathResult.pathVariables[pathResult.pathVariables.length-2]}["${this.rule.argumentAsProperty()}"]`);

        if (this.rule.negated) {
            rego.push(`${propAVariable} < ${propBVariable}`);
        } else {
            rego.push(`${propAVariable} >= ${propBVariable}`);
        }
        return [
            new SimpleRuleResult(
                "lessThanProperty",
                rego,
                this.rule.path[this.rule.path.length-1],
                `[${propAVariable},${propBVariable}]`,
                pathResult.variable,
                `Value for property '${path.join("/")}' not less than value for property ${this.rule.argument}`
            )
        ];
    }
}