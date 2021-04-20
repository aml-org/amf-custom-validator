import {AtomicRule, Rule, Variable} from "../Rule";

export class NestedRule extends AtomicRule {
    public readonly parent: Variable;
    public readonly child: Variable;

    constructor(negated, parent: Variable, child: Variable, path: string[]) {
        super(negated, parent, "nested", path, child);
        this.parent = parent;
        this.child = child;
    }
    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }

        if (this.child.cardinality != null) {
            return `${negation}Nested(${this.parent.name},${this.child.name},'${this.path.join("/")}')`;
        } else {
            return `${negation}Nested(${this.parent.name},${this.child.name},'${this.path.join("/")}')`;
        }
    }

    negation(): Rule {
        if (this.child.cardinality != null) {
            return new NestedRule(!this.negated, this.parent, this.child.negation(), this.path)
        } else {
            return new NestedRule(!this.negated, this.parent, this.child, this.path)
        }
    }
}