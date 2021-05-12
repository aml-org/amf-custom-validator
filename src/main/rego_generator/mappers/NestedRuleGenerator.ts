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
        const rego = [];
        const pluralName = `${this.rule.child.name}s`
        rego.push(`${pluralName} = ${pathResult.rule} with data.sourceNode as ${this.rule.variable.name}`)
        return [new SimpleRuleResult(
            "nested",
            rego,
            [pathResult],
            this.rule.path.source,
            pluralName,
            pluralName,
            ""
        )];
    }

}