import {InRule} from "../../model/constraints/InRule";
import {BaseRegoRuleGenerator, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {genvar} from "../../VarGen";

export class InRuleGenerator extends BaseRegoRuleGenerator {
    private rule: InRule;

    constructor(rule: InRule) {
        super();
        this.rule = rule;
    }

    generateResult(): SimpleRuleResult[] {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name, "in_" + this.rule.valueMD5()).generatePropertyValues();
        const rego: string[] = []
        const inValuesVariable = genvar("invalues");
        const inValuesTestVariable = `${this.rule.variable.name}_check`;
        rego.push(`${inValuesTestVariable}_array = ${pathResult.rule} with data.sourceNode as ${this.rule.variable.name}`)
        rego.push(`${inValuesTestVariable} = ${inValuesTestVariable}_array[_]`)
        rego.push(`${inValuesVariable} = {${this.argumentValue()}}`)
        if (this.rule.negated) {
            rego.push(`${inValuesVariable}[${inValuesTestVariable}]`)
        } else {
            rego.push(`not ${inValuesVariable}[${inValuesTestVariable}]`)
        }
        return [
            new SimpleRuleResult(
                "in",
                rego,
                [pathResult],
                path.source,
                inValuesTestVariable,
                inValuesTestVariable,
                `Value no in set {${this.argumentValue().replace(/"/g, "'")}}`
            )
        ];
    }

    protected argumentValue(): string {
        if (this.rule.argument instanceof Array) {
            return this.rule.argument.map((v) => JSON.stringify(v)).join(",")
        } else {
            return JSON.stringify(this.rule.argument);
        }
    }

}