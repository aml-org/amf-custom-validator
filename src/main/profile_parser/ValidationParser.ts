import {Rule, Variable} from "../model/Rule";
import {ConstraintParser} from "./ConstraintParser";
import {AndRule} from "../model/rules/AndRule";
import {Expression} from "../model/Expression";
import {OrRule} from "../model/rules/OrRule";

export class ValidationParser {

    private data: any;
    private expression: Expression;
    private variable: Variable;
    private negated: boolean;

    constructor(expression: Expression, variable: Variable, data: any, negated: boolean) {
        this.data = data;
        this.variable = variable;
        this.expression = expression;
        this.negated = negated;
    }


    parse(): Rule {
        switch (this.parseType()) {
            case "implicitAnd":
                return this.parseImplicitAnd();
            case "implicitRego":
                return this.parseImplicitRego();
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
        } else if (this.data.rego) {
            return "implicitRego";
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
        return new AndRule(this.negated).withBody(body);
    }

    private parseImplicitRego() {
        const data = {
            rego: this.data.rego
        };
        const parsed = new ConstraintParser(this.expression, this.variable, "", data).parse();
        return new AndRule(this.negated).withBody(parsed);
    }

    private parseAnd() {
        const nestedRules = this.data.and.map((nestedConstraint) => {
            return new ValidationParser(this.expression, this.variable, nestedConstraint, false).parse();
        });
        return new AndRule(this.negated).withBody(nestedRules);
    }

    private parseOr() {
        const nestedRules = this.data.or.map((nestedConstraint) => {
            return new ValidationParser(this.expression, this.variable, nestedConstraint, false).parse();
        });
        return new OrRule(this.negated).withBody(nestedRules);
    }

    private parseNot() {
        return new ValidationParser(this.expression, this.variable, this.data.not, true).parse();
    }
}