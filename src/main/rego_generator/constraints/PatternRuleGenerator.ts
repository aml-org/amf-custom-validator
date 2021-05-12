import {BaseRegoRuleGenerator, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {PatternRule} from "../../model/constraints/PatternRule";

export class PatternRuleGenerator extends  BaseRegoRuleGenerator {
    private rule: PatternRule;

    constructor(rule: PatternRule) {
        super();
        this.rule = rule;
    }

    generateResult(): SimpleRuleResult[] {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name, "pattern_" + this.rule.valueMD5()).generatePropertyValues();
        const rego = [];
        const checkVariable = `${pathResult.rule}_node`
        rego.push(`${checkVariable}_array = ${pathResult.rule} with data.sourceNode as ${this.rule.variable.name}`)
        rego.push(`${checkVariable} = ${checkVariable}_array[_]`)
        if (this.rule.negated) {
            rego.push(`regex.match(${JSON.stringify(this.rule.argument)},${checkVariable})`)
        } else {
            rego.push(`not regex.match(${JSON.stringify(this.rule.argument)},${checkVariable})`)
        }
        return [
          new SimpleRuleResult(
              "pattern",
              rego,
              [pathResult],
              this.rule.path.source,
              checkVariable,
              checkVariable,
              `Value does not match regular expression {${JSON.stringify(this.rule.argument).replace(/"/g,"'")}}`
          )
        ];
    }

}