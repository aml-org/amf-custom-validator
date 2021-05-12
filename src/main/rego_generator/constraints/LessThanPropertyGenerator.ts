import {BaseRegoRuleGenerator, SimpleRuleResult} from "../BaseRegoRuleGenerator";
import {RegoPathGenerator} from "../RegoPathGenerator";
import {AndPath} from "../../profile_parser/PathParser";
import {LessThanPropertyRule} from "../../model/constraints/LessThanPropertyRule";

export class LessThanPropertyGenerator extends BaseRegoRuleGenerator {
    private rule: LessThanPropertyRule;

    constructor(rule: LessThanPropertyRule) {
        super();
        this.rule = rule;
    }
    generateResult(): SimpleRuleResult[] {
        const rego: string[] = [];

        const path = this.rule.path;
        const pathResult = new RegoPathGenerator(path, this.rule.variable.name, "lessThanProperty_" + this.rule.valueMD5()).generatePropertyArray();
        const propAVariable = pathResult.rule + "_propA";


        const pathB = this.alternativePath();
        const pathResultB = new RegoPathGenerator(pathB, this.rule.variable.name, "lessThanPropertyAlt_" + this.rule.valueMD5()).generatePropertyArray();
        const propBVariable = pathResultB.rule + "_propB";

        rego.push(`${propAVariable}s = ${pathResult.rule} with data.sourceNode as ${this.rule.variable.name}`)
        rego.push(`${propBVariable}s = ${pathResultB.rule} with data.sourceNode as ${this.rule.variable.name}`)
        // this will compute all the pairs [a_i,b_j]
        rego.push(`${propAVariable} = ${propAVariable}s[_]`)
        rego.push(`${propBVariable} = ${propBVariable}s[_]`)
        // condition must hold for all pairs, every element in A must be <|>= all the elements in B
        if (this.rule.negated) {
            rego.push(`${propAVariable} < ${propBVariable}`);
        } else {
            rego.push(`${propAVariable} >= ${propBVariable}`);
        }
        return [
            new SimpleRuleResult(
                "lessThanProperty",
                rego,
                [pathResult,pathResultB],
                path.source,
                `[${propAVariable},${propBVariable}]`,
                propAVariable,
                `Value for property '${path.source.replace(/\\./g, ":")}' not less than value for property ${this.rule.argument.replace(/\\./g, ":")}`
            )
        ];
    }

    private alternativePath() {
        const nextProperty = {
            iri: this.rule.argument.replace(".",":"),
            inverse: false,
            transitive: false,
            source: this.rule.argument
        };

        //@ts-ignore
        if (this.rule.path.and) {
            const clonedPath = [].concat((<AndPath>this.rule.path).and);
            clonedPath.push(nextProperty);
            return {
                and: clonedPath,
                source: "(" + this.rule.path.source + ") / " + nextProperty.source
            }
            // @ts-ignore
        } else if (this.rule.path.or) {
            return {
                and: [this.rule.path, nextProperty],
                source: "(" + this.rule.path.source + ") / " + nextProperty.source
            }
        } else {
            return nextProperty
        }

    }
}