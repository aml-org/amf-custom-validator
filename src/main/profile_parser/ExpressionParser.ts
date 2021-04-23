import {Expression} from "../model/Expression";
import {ValidationParser} from "./ValidationParser";
import {Quantification} from "../model/Rule";
import {Implication} from "../model/Implication";
import {ClassTarget} from "../model/mappers/ClassTarget";

export class ExpressionParser {

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

        const expression = new Expression(false, this.name, this.message, this.level);
        const v = expression.genVar(Quantification.ForAll);

        const validation = new ValidationParser(expression, v, this.data, false).parse();
        const headTarget = new ClassTarget(false, v, this.targetClass);
        const implication = new Implication(false, v, headTarget, validation)

        expression.rule = implication;

        return expression;

    }
}