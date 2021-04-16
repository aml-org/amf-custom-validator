import {Rule, Variable} from "../model/Rule";
import {ConstraintParser} from "./ConstraintParser";
import {AndRule} from "../model/rules/AndRule";
import {Expression} from "../model/Expression";
import {OrRule} from "../model/rules/OrRule";

export class ValidationParser {

    private data: any;
    private expression: Expression;
    private variable: Variable;

    constructor(expression: Expression, variable: Variable, data: any) {
        this.data = data;
        this.variable = variable;
        this.expression = expression;
    }


    parse(): Rule {
        switch (this.parseType()) {
            case "implicitAnd":
                return this.parseImplicitAnd();
            case "or":
                return this.parseOr();
            case "not":
                return this.parseNot();
            case "and":
                return this.parseAnd();
            default:
                throw new Error(`Unknown validation type ${JSON.stringify(this.data)}`);
        }
    }

    private parseType() {
        if (this.data.propertyConstraints != null) {
            return "implicitAnd";
        } else if (this.data.or != null) {
            return "or";
        } else if (this.data.not != null) {
            return "not";
        } else if (this.data.and != null) {
            return "and";
        } else {
            return null;
        }
    }

    private parseImplicitAnd() {
        const body: Rule[] = [];
        Object.keys(this.data.propertyConstraints).map((path) => {
            const value = this.data.propertyConstraints[path];
            new ConstraintParser(this.expression, this.variable, path, value).parse().forEach(r => body.push(r));
        });
        return new AndRule(false).withBody(body);
    }

    private parseAnd() {
        const nestedRules = this.data.and.map((nestedConstraint) => {
            return new ValidationParser(this.expression, this.variable, nestedConstraint).parse();
        });
        return new AndRule(false).withBody(nestedRules);
    }

    private parseOr() {
        const nestedRules = this.data.or.map((nestedConstraint) => {
            return new ValidationParser(this.expression, this.variable, nestedConstraint).parse();
        });
        return new OrRule(false).withBody(nestedRules);
    }

    private parseNot() {
        const nested = new ValidationParser(this.expression, this.variable, this.data.not).parse();
        return <Rule>nested.negation();
    }
}