import {BaseRegoRuleGenerator, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {NestedRule} from "../../model/constraints/NestedRule";
import {RegoPathGenerator} from "../RegoPathGenerator";

export class NestedRuleGenerator extends BaseRegoRuleGenerator {
    private rule: NestedRule;

    constructor(rule: NestedRule) {
        super();
        this.rule = rule;
    }

    generateResult(): SimpleRuleResult[] {
        if (this.rule.child.cardinality != null) {
            return this.generateQuantifiedNestedResult()
        } else {
            return this.generatedNestedResult();
        }

    }

    generateQuantifiedNestedResult() {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.parent.name, "nested_" + this.rule.valueMD5()).generatePropertyValues();
        const rego = pathResult.rego;

        if (this.rule.negated) {
            // rego.push(`${this.rule.child.name} = find with data.link as${pathResult.variable}`);
            throw new Error("Not supported yet")
        } else {
            rego.push(`${this.rule.child.name} = find with data.link as ${pathResult.variable}`);
        }
        rego.push(this.rule.child.toQuantifiedRego());
        return [
          new SimpleRuleResult(
              this.rule.child.cardinality!.operator,
              rego,
              this.rule.path[this.rule.path.length-1],
              `count(${this.rule.child.name})`,
              this.rule.child.name,
              `violated quantified constraint ${this.rule.child.cardinality.toString()} `
          )
        ];
    }

    generatedNestedResult() {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.parent.name, "nested_" + this.rule.valueMD5()).generatePropertyValues();
        const rego = pathResult.rego;

        if (this.rule.negated) {
            // rego.push(`${this.rule.child.name} = find with data.link as${pathResult.variable}`);
            throw new Error("Not supported yet")
        } else {
            rego.push(`${this.rule.child.name} = find with data.link as ${pathResult.variable}`);
        }
        return [new SimpleRuleResult(
            "nested",
            rego,
            this.rule.path[this.rule.path.length-1],
            this.rule.child.name,
            this.rule.child.name,
            `Not nested matching constraints for parent ${this.rule.parent} and child ${this.rule.child} under ${this.rule.path.join("/")}`
        )];
    }

}