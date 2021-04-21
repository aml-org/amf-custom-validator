import * as md5 from "md5";

export interface RegoPathResult {
    rego: string[],
    variable: string
}

export class RegoPathGenerator {

    private path: string[];
    private variable: string;
    private id: any;
    private hint: string;

    constructor(path: string[], variable: string, hint: string) {
        this.id = md5(path.join("/"))
        this.path = path;
        this.variable = variable;
        this.hint = hint;
    }

    public generatePropertyValues(): RegoPathResult {
        const lines = [];
        let previous_binding = this.variable;
        let binding = this.variable + "_0_" + this.id + "_" + this.hint;
        for (let i=0; i<this.path.length; i++) {
            const nextPath = this.path[i];
            if (i == this.path.length-1) {
                lines.push(`${binding} = ${previous_binding}["${nextPath}"]`)
            } else {
                lines.push(`nested[${binding}] with data.nodes as ${previous_binding}[${nextPath}]`)
            }
            previous_binding = binding;
            binding = this.variable + "_" + (i+1) + "_" + this.id + "_" + this.hint
        }

        return {
            rego: lines,
            variable: previous_binding
        }
    }

    public generatePropertyArray(): RegoPathResult {
        const lines = [];
        let previous_binding = this.variable;
        let binding = this.variable + "_0_" + this.id;
        for (let i=0; i<this.path.length; i++) {
            const nextPath = this.path[i];
            if (i == this.path.length-1) {
                lines.push(`${binding} = object.get(${previous_binding},"${nextPath}",[])`)
            } else {
                lines.push(`nested[${binding}] with data.nodes as ${previous_binding}[${nextPath}]`)
            }
            previous_binding = binding;
            binding = this.variable + "_" + (i+1) + "_" + this.id
        }

        return {
            rego: lines,
            variable: previous_binding
        }
    }
}