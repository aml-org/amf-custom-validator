import {PathParser} from "./PathParser";
import {Quantification, Rule, Variable} from "../model/Rule";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {PatternRule} from "../model/constraints/PatternRule";
import {InRule} from "../model/constraints/InRule";
import {Expression} from "../model/Expression";
import {ValidationParser} from "./ValidationParser";
import {AndRule} from "../model/rules/AndRule";
import {NestedRule} from "../model/rules/NestedRule";


export class ConstraintParser {
    private path: string[];
    private constraints: any;
    private expression: Expression;
    private variable: Variable;
    constructor(expression: Expression, variable: Variable, path: string, constraints: any) {
        this.expression = expression;
        this.variable = variable;
        this.path = new PathParser(path).parse();
        this.constraints = constraints;
    }

    parse(): Rule[] {
        return Object.keys(this.constraints).map((constraint) => {
            switch (constraint) {
                case "minCount":
                    return this.parseMinCount(this.constraints[constraint]);
                case "pattern":
                    return this.parsePattern(this.constraints[constraint]);
                case "in":
                    return this.parseIn(this.constraints[constraint]);
                case "nested":
                    return this.parseNested(this.constraints[constraint]);
                default:
                    throw new Error(`Constraint ${constraint} not supported yet`);
            }
        })

    }


    private parseMinCount(constraint: any): MinCountRule {
        return new MinCountRule(false, this.variable, this.path, constraint);
    }

    private parsePattern(constraint: any): Rule {
        return new PatternRule(false, this.variable, this.path, constraint);
    }

    private parseIn(constraint: any): Rule {
        return new InRule(false, this.variable, this.path, constraint);
    }

    private parseNested(constraint: any) {
        const nextVar = this.expression.genVar(Quantification.ForAll);
        const nested = new ValidationParser(this.expression, nextVar, constraint).parse()
        const nestedRule = new NestedRule(false, this.variable, nextVar, this.path);
        return new AndRule(false).withBody([nestedRule,nested]);
    }

}