import {BaseRegoAtomicRuleGenerator, RegoRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {MinCountRule} from "../../model/constraints/MinCountRule";
import {genvar} from "../../VarGen";

export class MinCountRuleGenerator extends BaseRegoAtomicRuleGenerator {
    private rule: MinCountRule;

    constructor(rule: MinCountRule) {
        super();
        this.rule = rule;
    }
    generateResult(): RegoRuleResult {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name).generatePropertyArray();
        const rego = pathResult.rego;

        const inValuesVariable = genvar("propValues");
        rego.push(`${inValuesVariable} = nodes_array with data.nodes as ${pathResult.variable}`)
        if (this.rule.negated) {
            rego.push(`not count(${inValuesVariable}) >= ${this.rule.argument}`)
        } else {
            rego.push(`count(${inValuesVariable}) >= ${this.rule.argument}`)
        }
        return {
            constraintId: "minCount",
            rego: rego,
            path: this.rule.path[this.rule.path.length-1],
            value: `count(${inValuesVariable})`,
            traceMessage: `Value not matching minCount ${this.rule.argument}`
        }
    }

}