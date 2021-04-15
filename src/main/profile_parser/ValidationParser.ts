import {Quantification, Rule} from "../model/Rule";
import {ConstraintParser} from "./ConstraintParser";
import {AndRule} from "../model/rules/AndRule";
import {Expression} from "../model/Expression";
import {Implication} from "../model/Implication";
import {ClassTarget} from "../model/constraints/ClassTarget";

export class ValidationParser {

    private data: any;
    private message: string;
    private targetClass: string;
    private name: string|null;
    private level: string;

    constructor(name: string|null, data: any, level: string) {
        this.name = name;
        this.data = data;
        this.level = level;
    }


    parse(): Expression {
        this.targetClass = this.data.targetClass;
        this.message = this.data.message || "Validation error"
        switch (this.parseType()) {
            case "implicitAnd":
                return this.parseImplicitAnd();
            default:
                throw new Error("Unknown validation type");
        }
    }

    private parseType() {
        if (this.data.propertyConstraints != null) {
            return "implicitAnd";
        } else {
            return null;
        }
    }

    private parseImplicitAnd() {
        const expression = new Expression(false, this.name, this.message, this.level);
        const v = expression.genVar(Quantification.ForAll);

        const body: Rule[] = [];
        Object.keys(this.data.propertyConstraints).map((path) => {
            const value = this.data.propertyConstraints[path];
            new ConstraintParser(expression, v, path, value).parse().forEach(r => body.push(r));
        });
        const implicitAnd = new AndRule(false).withBody(body);

        const headTarget = new ClassTarget(false, v, this.targetClass);

        const implication = new Implication(false, v, headTarget, implicitAnd)

        expression.rule = implication;

        return expression;
    }
}