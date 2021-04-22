import * as md5 from "md5";

export const Quantification = {
    ForAll: "ForAll",
    Exist: "Exist"
}

export const CardinalityOperation = {
    LTEQ: "lessThanOrEqual",
    LT: "lessThan",
    EQ: "equal",
    NEQ: "notEqual",
    GT: "greaterThan",
    GTEQ: "greaterThanOrEqual"
}

export class VariableCardinality {
    public readonly operator: string;
    public readonly value: number;

    constructor(operator: string, value: number) {
        this.operator = operator;
        this.value = value;
        if (this.operator !== CardinalityOperation.LTEQ && this.operator !== CardinalityOperation.LT && this.operator !== CardinalityOperation.EQ
            && this.operator !== CardinalityOperation.NEQ && this.operator !== CardinalityOperation.GT && this.operator !== CardinalityOperation.GTEQ) {
            throw new Error(`Unknown cardinality ${operator} ${value}`);
        }
    }

    toRego(total: string, error_count: string, negated: boolean) {
        if (negated) {
            switch(this.operator) {
                case CardinalityOperation.GTEQ:
                    return `count(${total}) - count(${error_count}) < ${this.value}`
                case CardinalityOperation.GT:
                    return `count(${total}) - count(${error_count}) <= ${this.value}`
                case CardinalityOperation.EQ:
                    return `count(${total}) - count(${error_count}) != ${this.value}`
                case CardinalityOperation.NEQ:
                    return `count(${total}) - count(${error_count}) == ${this.value}`
                case CardinalityOperation.LT:
                    return `count(${total}) - count(${error_count}) >= ${this.value}`
                case CardinalityOperation.LTEQ:
                    return `count(${total}) - count(${error_count}) > ${this.value}`
                default:
                    throw new Error("Cannot negate unknown cardinallity: " + this.operator); // should never happen;
            }
        } else {
            switch(this.operator) {
                case CardinalityOperation.GTEQ:
                    return `count(${total}) - count(${error_count}) >= ${this.value}`
                case CardinalityOperation.GT:
                    return `count(${total}) - count(${error_count}) > ${this.value}`
                case CardinalityOperation.EQ:
                    return `count(${total}) - count(${error_count}) == ${this.value}`
                case CardinalityOperation.NEQ:
                    return `count(${total}) - count(${error_count}) != ${this.value}`
                case CardinalityOperation.LT:
                    return `count(${total}) - count(${error_count}) < ${this.value}`
                case CardinalityOperation.LTEQ:
                    return `count(${total}) - count(${error_count}) <= ${this.value}`
                default:
                    throw new Error("Cannot negate unknown cardinality: " + this.operator); // should never happen;
            }
        }
    }

    toString() {
        switch(this.operator) {
            case CardinalityOperation.GTEQ:
                return `>= ${this.value}`
            case CardinalityOperation.GT:
                return `> ${this.value}`
            case CardinalityOperation.EQ:
                return `= ${this.value}`
            case CardinalityOperation.NEQ:
                return `<> ${this.value}`
            case CardinalityOperation.LT:
                return `< ${this.value}`
            case CardinalityOperation.LTEQ:
                return `<= ${this.value}`
            default:
                throw new Error("Cannot negate unknown cardinality: " + this.operator); // should never happen;
        }
    }

    public negation(): VariableCardinality {
        switch(this.operator) {
            case CardinalityOperation.GTEQ:
                return new VariableCardinality(CardinalityOperation.LT, -this.value);
            case CardinalityOperation.GT:
                return new VariableCardinality(CardinalityOperation.LTEQ, -this.value);
            case CardinalityOperation.EQ:
                return new VariableCardinality(CardinalityOperation.NEQ, -this.value);
            case CardinalityOperation.NEQ:
                return new VariableCardinality(CardinalityOperation.EQ, -this.value);
            case CardinalityOperation.LT:
                return new VariableCardinality(CardinalityOperation.GTEQ, -this.value);
            case CardinalityOperation.LTEQ:
                return new VariableCardinality(CardinalityOperation.GT, -this.value);
            default:
                throw new Error("Cannot negate unknown cardinallity: " + this.operator); // should never happen;
        }
    }

    static lessThan(n: number) {
        return new VariableCardinality(CardinalityOperation.LT, n)
    }
    static lessThanOrEqual(n: number) {
        return new VariableCardinality(CardinalityOperation.LTEQ, n)
    }

    static equal(n: number) {
        return new VariableCardinality(CardinalityOperation.EQ, n)
    }

    static notEqual(n: number) {
        return new VariableCardinality(CardinalityOperation.NEQ, n)
    }

    static greaterThan(n: number) {
        return new VariableCardinality(CardinalityOperation.GT, n)
    }

    static greaterThanOrEqual(n: number) {
        return new VariableCardinality(CardinalityOperation.GTEQ, n)
    }
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

    public negated: boolean

    negation(): Statement {
        this.negated = !this.negated;
        return this;
    }

}

export class Variable {
    public readonly quantification: string;
    public readonly name: string;
    public readonly cardinality: VariableCardinality|undefined;

    constructor(name: string, quantification: string, cardinality?: VariableCardinality) {
        this.name = name;
        this.quantification = quantification;
        if (this.quantification != Quantification.ForAll && this.quantification != Quantification.Exist) {
            throw new Error(`Quantification values must be ${Quantification.ForAll} or ${Quantification.Exist}`)
        }
        this.cardinality = cardinality;
        if (this.quantification == Quantification.ForAll && this.cardinality != null) {
            throw new Error("ForAll variables cannot have an associated cardinality");
        }
    }

    protected quantifierNegation() {
        if (this.quantification == Quantification.ForAll) {
            return Quantification.Exist;
        } else {
            if (this.cardinality == null) {
                return Quantification.ForAll;
            } else {
                return Quantification.Exist;
            }

        }
    }

    negation(): Variable {
        if (this.cardinality == null) {
            return new Variable(this.name, this.quantifierNegation());
        } else {
            return new Variable(this.name, this.quantifierNegation(), this.cardinality.negation());
        }

    }

    toString() {
        if (this.quantification === Quantification.ForAll) {
            return `∀${this.name}`
        } else if (this.cardinality) {
            return `∃${this.name}:${this.cardinality.toString()}`;
        } else{
            return `∃${this.name}`;
        }
    }
}

export abstract class Rule extends Statement {

    public name: string|null;

    constructor(negated: boolean) {
        super(negated);
    }

    abstract toString(): string;

    negation(): Rule {
        return <Rule>super.negation();
    }
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

    valueMD5() {
        return md5(JSON.stringify(this.argument)).toString();
    }

}

