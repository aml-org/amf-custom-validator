import {BaseRegoAtomicRuleGenerator, RegoRuleResult} from "../BaseRegoRuleGenerator";
import {NestedRule} from "../../model/constraints/NestedRule";
import {RegoPathGenerator} from "../RegoPathGenerator";

export class NestedRuleGenerator extends BaseRegoAtomicRuleGenerator {
    private rule: NestedRule;

    constructor(rule: NestedRule) {
        super();
        this.rule = rule;
    }

    generateResult(): RegoRuleResult {
        if (this.rule.child.cardinality != null) {
            return this.generateQuantifiedNestedResult()
        } else {
            return this.generatedNestedResult();
        }

    }

    generateQuantifiedNestedResult() {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.parent.name).generatePropertyValues();
        const rego = pathResult.rego;

        if (this.rule.negated) {
            rego.push(`not ${this.rule.child.name} = find with data.link as${pathResult.variable}`);
        } else {
            rego.push(`${this.rule.child.name} = find with data.link as ${pathResult.variable}`);
        }
        rego.push(this.rule.child.toQuantifiedRego());
        return {
            constraintId: this.rule.child.cardinality!.operator,
            rego: rego,
            path: this.rule.path[this.rule.path.length-1],
            value: `count(${this.rule.child.name})`,
            traceMessage: `violated quantified constraint ${this.rule.child.cardinality.toString()} `
        }
    }

    generatedNestedResult() {
        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.parent.name).generatePropertyValues();
        const rego = pathResult.rego;

        if (this.rule.negated) {
            rego.push(`not ${this.rule.child.name} = find with data.link as${pathResult.variable}`);
        } else {
            rego.push(`${this.rule.child.name} = find with data.link as ${pathResult.variable}`);
        }
        return {
            constraintId: "nested",
            rego: rego,
            path: this.rule.path[this.rule.path.length-1],
            value: this.rule.child.name,
            traceMessage: `Not nested matching constraints for parent ${this.rule.parent} and child ${this.rule.child} under ${this.rule.path.join("/")}`
        }
    }

}