import {PathParser} from "./PathParser";
import {CardinalityOperation, Quantification, Rule, Variable, VariableCardinality} from "../model/Rule";
import {MinCountRule} from "../model/constraints/MinCountRule";
import {PatternRule} from "../model/constraints/PatternRule";
import {InRule} from "../model/constraints/InRule";
import {Expression} from "../model/Expression";
import {ValidationParser} from "./ValidationParser";
import {Implication} from "../model/Implication";
import {NestedRule} from "../model/mappers/NestedRule";


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
                case "atLeast":
                    return this.parseQualifiedNested(this.constraints[constraint], CardinalityOperation.GTEQ);
                case "atMost":
                    return this.parseQualifiedNested(this.constraints[constraint], CardinalityOperation.LTEQ);
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
        const nestedExpression = this.expression.subExpression(false);
        const nextVar = nestedExpression.genVar(Quantification.ForAll);

        const nested = new ValidationParser(nestedExpression, nextVar, constraint).parse()
        const nestedRule = new NestedRule(false, this.variable, nextVar, this.path);
        nestedExpression.rule = new Implication(false, this.variable, nestedRule,nested);
        return nestedExpression;
    }

    private parseQualifiedNested(constraint: any, cardinality: string) {
        const count = <number>constraint.count;
        let variableCardinality: VariableCardinality;
        switch (cardinality) {
            case CardinalityOperation.GTEQ:
                variableCardinality = VariableCardinality.greaterThanOrEqual(count);
                break;
            case CardinalityOperation.GT:
                variableCardinality = VariableCardinality.greaterThan(count);
                break;
            case CardinalityOperation.EQ:
                variableCardinality = VariableCardinality.equal(count);
                break;
            case CardinalityOperation.NEQ:
                variableCardinality = VariableCardinality.notEqual(count);
                break;
            case CardinalityOperation.LT:
                variableCardinality = VariableCardinality.lessThan(count);
                break;
            case CardinalityOperation.LTEQ:
                variableCardinality = VariableCardinality.lessThanOrEqual(count);
                break;
            default:
                throw new Error("Unknown cardinality "+ cardinality);
        }

        const nestedExpression = this.expression.subExpression(false);
        const nextVar = nestedExpression.genVar(Quantification.Exist, variableCardinality);
        const nested = new ValidationParser(nestedExpression, nextVar, constraint.validation).parse()
        const nestedRule = new NestedRule(false, this.variable, nextVar, this.path);
        nestedExpression.rule = new Implication(false, this.variable, nestedRule,nested);
        return nestedExpression;
    }
}