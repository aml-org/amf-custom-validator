import {BaseRegoRuleGenerator, RegoRuleResult, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoRule} from "../../model/constraints/RegoRule";
import {RegoPathGenerator} from "../RegoPathGenerator";

export class RegoRuleGenerator extends BaseRegoRuleGenerator {

    private rule: RegoRule;

    constructor(rule: RegoRule) {
        super();
        this.rule = rule;
    }

    generateResult(): RegoRuleResult[] {
        const path = this.rule.path;
        const rego = [];

        // by default we use the expression variable, this is the right one if the rego rule is top-level
        let checkVariable = this.rule.variable.name;

        // let's try generate the path rule for the constraint
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name, "rego_" + this.rule.valueMD5()).generatePropertyValues();

        // if this is not a top-level rego rule (the path generates code), we use the bind the check variable for the path computation result
        if (pathResult.rego.length > 0) {
            checkVariable =  `${pathResult.rule}_node`
            rego.push(`${checkVariable}_array = ${pathResult.rule} with data.sourceNode as ${this.rule.variable.name}`)
            rego.push(`${checkVariable} = ${checkVariable}_array[_]`)
        }

        // this is the value where we will store the result of the custom rego code
        const resultVariable = `rego_result_${this.rule.valueMD5()}`

        // let's add all custom rego code to the code to be generated
        let text = this.rule.argument.code;
        // we first need to replace the variables in the rego templat fro the right check and result variables
        text.replace(/\$node/g,checkVariable).replace(/\$result/g,resultVariable).split("\n").forEach((line) => rego.push(line));
        // now we can negate or not the resultVariable, we are checking that the result is erroneous
        if (this.rule.negated) {
            rego.push(`${resultVariable} == true`)
        } else {
            rego.push(`${resultVariable} != true`)
        }

        return [
            new SimpleRuleResult(
                "rego",
                rego,
                [pathResult], // this can be an empty path result
                this.rule.path.source,
                checkVariable,
                checkVariable,
                this.rule.argument.message
            )
        ];
    }
}