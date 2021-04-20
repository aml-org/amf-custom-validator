import {Level, Rule, Statement, Variable, VariableCardinality} from "./Rule";
import {VarGenerator} from "../VarGen";

export class Expression extends Statement {

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

    negation(): Expression {
        const acc = this.variables.map((variable) => variable.negation())
        const negatedExpression = new Expression(!this.negated, this.name, this.message, this.level);
        negatedExpression.variables = acc;
        negatedExpression.rule = <Rule>this.rule.negation()
        return negatedExpression;
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

    toCanonicalForm(): Expression {
        let canonical: Expression;
        if (this.negated == false) {
            canonical = this.negation();
        } else {
            canonical = new Expression(this.negated, this.name, this.message, this.level);
            canonical.variables = this.variables;
            canonical.rule = this.rule;
        }

        canonical.rule = <Rule>canonical.rule.toCanonicalForm();

        if (canonical.rule instanceof Expression) {
            canonical.rule.variables.forEach((v) => canonical.variables.push(v));
            canonical.rule = canonical.rule.rule;
        }
        return canonical;
    }

    subExpression(negated: boolean): Expression {
        return new Expression(negated, null, null, null, this.varGenerator);
    }
}