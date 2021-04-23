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
        const rego = pathResult.rego;
        if (this.rule.negated) {
            rego.push(`regex.match(${JSON.stringify(this.rule.argument)},${pathResult.variable})`)
        } else {
            rego.push(`not regex.match(${JSON.stringify(this.rule.argument)},${pathResult.variable})`)
        }
        return [
          new SimpleRuleResult(
              "pattern",
              rego,
              this.rule.path[this.rule.path.length-1],
              pathResult.variable,
              pathResult.variable,
              `Value does not match regular expression {${JSON.stringify(this.rule.argument).replace(/"/g,"'")}}`
          )
        ];
    }

}