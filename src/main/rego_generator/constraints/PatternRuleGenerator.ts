import {BaseRegoAtomicRuleGenerator, RegoRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {PatternRule} from "../../model/constraints/PatternRule";
import {genvar} from "../../VarGen";

export class PatternRuleGenerator extends  BaseRegoAtomicRuleGenerator {
    private rule: PatternRule;

    constructor(rule: PatternRule) {
        super();
        this.rule = rule;
    }

    generateResult(): RegoRuleResult {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name).generatePropertyValues();
        const rego = pathResult.rego;
        if (this.rule.negated) {
            rego.push(`not regex.match(${pathResult.variable}, ${JSON.stringify(this.rule.argument)})`)
        } else {
            rego.push(`regex.match(${pathResult.variable}, ${JSON.stringify(this.rule.argument)})`)
        }
        return {
            constraintId: "match",
            rego: rego,
            path: this.rule.path[this.rule.path.length-1],
            value: pathResult.variable,
            traceMessage: `Value does not match regular expression {${JSON.stringify(this.rule.argument).replace(/"/g,"'")}}`
        }
    }

}