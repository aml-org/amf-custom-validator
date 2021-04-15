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
            default:
                throw new Error(`Unknown validation type ${JSON.stringify(this.data)}`);
        }
    }

    private parseType() {
        if (this.data.propertyConstraints != null) {
            return "implicitAnd";
        } else if (this.data.or != null) {
            return "or"
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

    private parseOr() {
        const variables:{[name: string]: Variable} = {};
        const nestedRules = this.data.or.map((nestedConstraint) => {
            return new ValidationParser(this.expression, this.variable, nestedConstraint).parse();
        });
        return new OrRule(false).withBody(nestedRules);
    }
}