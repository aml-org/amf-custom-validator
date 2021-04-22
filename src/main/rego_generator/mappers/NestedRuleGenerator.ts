import {BaseRegoRuleGenerator, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {NestedRule} from "../../model/mappers/NestedRule";

export class NestedRuleGenerator extends BaseRegoRuleGenerator {
    private rule: NestedRule;

    constructor(rule: NestedRule) {
        super();
        this.rule = rule;
    }

    generateResult() {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.parent.name, "nested_" + this.rule.valueMD5()).generateNodeArray();
        const rego = pathResult.rego;
        const pluralName = `${this.rule.child.name}s`
        rego.push(`${pluralName} = ${pathResult.variable}`)
        return [new SimpleRuleResult(
            "nested",
            rego,
            this.rule.path[this.rule.path.length-1],
            pluralName,
            pluralName,
            ""
        )];
    }

}