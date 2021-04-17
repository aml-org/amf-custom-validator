import {Level, Rule, Statement, Variable} from "./Rule";

export class Expression extends Statement {

    public variables: Variable[] = [];
    public rule: Rule;
    public readonly message: string
    public readonly name: string;
    public readonly level: string;

    constructor(negated: boolean, name: string, message: string, level: string) {
        super(negated);
        this.name = name;
        this.message = message;
        if (level != Level.Ignore && level != Level.Info && level != Level.Violation && level != Level.Warning) {
            throw new Error(`Unknown severity level ${level}`)
        }
        this.level = level;
    }

    private varCounter = 0
    private vars = ["x", "y", "z", "p", "q", "r", "s", "t", "u", "v", "w", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"]

    genVar(quantification: string): Variable {
        let varName;
        if (this.varCounter < this.vars.length) {
            varName = this.vars[this.varCounter];

        } else {
            varName = `X${this.varCounter}`
        }
        const v = new Variable(varName, quantification);
        this.varCounter++;
        this.variables.push(v)
        return v;
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
        return `${this.name}[${this.level}] := ${varsText.map((v) => `${negation} ${v}`).join(",")} : ${this.rule.toString()}`
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

        return canonical;
    }
}