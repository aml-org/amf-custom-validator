import {BaseRegoRuleGenerator, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {genvar} from "../../VarGen";
import {MaxCountRule} from "../../model/constraints/MaxCountRule";

export class MaxCountRuleGenerator extends BaseRegoRuleGenerator {
    private rule: MaxCountRule;

    constructor(rule: MaxCountRule) {
        super();
        this.rule = rule;
    }
    generateResult(): SimpleRuleResult[] {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name, "maxCount_" + this.rule.valueMD5()).generatePropertyArray();
        const rego: string[] = []

        const inValuesVariable = genvar("propValues");
        rego.push(`${inValuesVariable} = ${pathResult.rule} with data.sourceNode as ${this.rule.variable.name}`)
        if (this.rule.negated) {
            rego.push(`count(${inValuesVariable}) <= ${this.rule.argument}`)
        } else {
            rego.push(`not count(${inValuesVariable}) <= ${this.rule.argument}`)
        }
        return [
            new SimpleRuleResult(
                "maxCount",
                rego,
                [pathResult],
                this.rule.path.source,
                `count(${inValuesVariable})`,
                inValuesVariable,
                `Value not matching maxCount ${this.rule.argument}`
            )
        ];
    }

}