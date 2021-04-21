import {InRule} from "../../model/constraints/InRule";
import {BaseRegoRuleGenerator, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";
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
        const rego = pathResult.rego;

        const inValuesVariable = genvar("invalues");
        rego.push(`${inValuesVariable} = {${this.argumentValue()}}`)
        if (this.rule.negated) {
            rego.push(`${inValuesVariable}[${pathResult.variable}]`)
        } else {
            rego.push(`not ${inValuesVariable}[${pathResult.variable}]`)
        }
        return [
            new SimpleRuleResult(
                "in",
                rego,
                this.rule.path[this.rule.path.length-1],
                pathResult.variable,pathResult.variable,
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