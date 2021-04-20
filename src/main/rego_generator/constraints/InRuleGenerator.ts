import {InRule} from "../../model/constraints/InRule";
import {RegoRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {genvar} from "../../VarGen";
import {BaseRegoAtomicRuleGenerator} from "../BaseRegoRuleGenerator";

export class InRuleGenerator extends BaseRegoAtomicRuleGenerator {
    private rule: InRule;

    constructor(rule: InRule) {
        super();
        this.rule = rule;
    }

    generateResult(): RegoRuleResult {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name).generatePropertyValues();
        const rego = pathResult.rego;

        const inValuesVariable = genvar("invalues");
        rego.push(`${inValuesVariable} = {${this.argumentValue()}}`)
        if (this.rule.negated) {
            rego.push(`not ${inValuesVariable}[${pathResult.variable}]`)
        } else {
            rego.push(`${inValuesVariable}[${pathResult.variable}]`)
        }
        return {
            constraintId: "in",
            rego: rego,
            path: this.rule.path[this.rule.path.length-1],
            value: pathResult.variable,
            traceMessage: `Value no in set {${this.argumentValue().replace(/"/g, "'")}}`
        }
    }

    protected argumentValue(): string {
        if (this.rule.argument instanceof Array) {
            return this.rule.argument.map((v) => JSON.stringify(v)).join(",")
        } else {
            return JSON.stringify(this.rule.argument);
        }
    }

}