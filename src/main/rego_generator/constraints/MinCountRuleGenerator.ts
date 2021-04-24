import {BaseRegoRuleGenerator, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {MinCountRule} from "../../model/constraints/MinCountRule";
import {genvar} from "../../VarGen";

export class MinCountRuleGenerator extends BaseRegoRuleGenerator {
    private rule: MinCountRule;

    constructor(rule: MinCountRule) {
        super();
        this.rule = rule;
    }
    generateResult(): SimpleRuleResult[] {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name, "minCount_" + this.rule.valueMD5()).generatePropertyArray();
        const rego = pathResult.rego;

        const inValuesVariable = genvar("propValues");
        rego.push(`${inValuesVariable} = nodes_array with data.nodes as ${pathResult.variable}`)
        if (this.rule.negated) {
            rego.push(`count(${inValuesVariable}) >= ${this.rule.argument}`)
        } else {
            rego.push(`not count(${inValuesVariable}) >= ${this.rule.argument}`)
        }
        return [
          new SimpleRuleResult(
              "minCount",
              rego,
              this.rule.path.source,
              `count(${inValuesVariable})`,
              pathResult.variable,
              `Value not matching minCount ${this.rule.argument}`
          )
        ];
    }

}