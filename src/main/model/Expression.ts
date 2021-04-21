import {Level, Rule, Statement, Variable, VariableCardinality} from "./Rule";
import {VarGenerator} from "../VarGen";

export class Expression extends Rule {

    public variables: Variable[] = [];
    public rule: Rule;
    public readonly message: string
    public readonly name: string;
    public readonly level: string;
    private varGenerator: VarGenerator;

    constructor(negated: boolean, name?: string, message?: string, level?: string, varGenerator: VarGenerator = new VarGenerator()) {
        super(negated);
        this.varGenerator = varGenerator;
        this.name = name;
        this.message = message;
        if (level != null && level != Level.Ignore && level != Level.Info && level != Level.Violation && level != Level.Warning) {
            throw new Error(`Unknown severity level ${level}`)
        }
        this.level = level;
    }


    genVar(quantification: string, cardinality?:VariableCardinality) {
        const variable = this.varGenerator.genExpressionVar(quantification, cardinality);
        this.variables.push(variable);
        return variable;
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        let varsText = [];
        this.variables.forEach((v) => varsText.push(v.toString()));
        if (this.name != null) {
            return `${this.name}[${this.level}] := ${varsText.map((v) => `${negation} ${v}`).join(",")} : ${this.rule.toString()}`
        } else {
            return `${varsText.map((v) => `${negation} ${v}`).join(",")} : ${this.rule.toString()}`
        }

    }

    subExpression(negated: boolean): Expression {
        return new Expression(negated, null, null, null, this.varGenerator);
    }
}