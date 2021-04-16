import {AtomicRule, Rule, Variable} from "../Rule";

export class NestedRule extends AtomicRule {
    private parent: Variable;
    private child: Variable;

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

        return `${negation}Nested(${this.parent.name},${this.child.name},'${this.path.join("/")}')`;
    }

    negation(): Rule {
        return new NestedRule(!this.negated, this.parent, this.child, this.path)
    }
}