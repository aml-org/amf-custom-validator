export const Quantification = {
    ForAll: "ForAll",
    Exist: "Exist"
}

export const Level = {
    Violation: "VIOLATION",
    Warning: "WARNING",
    Info: "INFO",
    Ignore: "IGNORE"
}

export abstract class Statement {
    constructor(negated: boolean) {
        this.negated = negated;
    }

    public readonly negated: boolean
    abstract negation(): Statement;
    abstract toCanonicalForm(): Statement;
}

export class Variable {
    public quantification: string;
    public name: string;

    constructor(name: string, quantification: string) {
        this.name = name;
        this.quantification = quantification;
        if (this.quantification != Quantification.ForAll && this.quantification != Quantification.Exist) {
            throw new Error(`Quantification values must be ${Quantification.ForAll} or ${Quantification.Exist}`)
        }
    }

    protected quantifierNegation() {
        if (this.quantification == Quantification.ForAll) {
            return Quantification.Exist
        } else {
            return Quantification.ForAll
        }
    }

    negation(): Variable {
        return new Variable(this.name, this.quantifierNegation())
    }

    toString() {
        if (this.quantification === Quantification.ForAll) {
            return `∀${this.name}`
        } else {
            return `∃${this.name}`
        }
    }
}

export abstract class Rule extends Statement {

    public name: string|null;

    constructor(negated: boolean) {
        super(negated);
    }

    abstract toString(): string;
}


export abstract class ComplexRule extends Rule {
    public body: Rule[] = [];

    withBody(rules: Rule[]): ComplexRule {
        this.body = rules;
        return this;
    }
}

export abstract class AtomicRule extends Rule {
    public readonly argument: any;
    public readonly path: string[];
    private constraint: string;
    public readonly variable: Variable;

    constructor(negated: boolean, variable: Variable, constraint: string, path: string[], argument: any) {
        super(negated);
        this.variable = variable;
        this.path = path;
        this.argument = argument;
        this.constraint = constraint
    }

    toCanonicalForm(): Statement {
        return this;
    }
}

