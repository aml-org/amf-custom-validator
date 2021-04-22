import {AtomicRule, Variable} from "../Rule";
import * as md5 from "md5";

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

    valueMD5(): any {
        return md5(`${this.parent.name}_nested_${this.child.name}`)
    }
}