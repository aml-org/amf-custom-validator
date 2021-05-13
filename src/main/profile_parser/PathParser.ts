import * as pegjs from 'pegjs';
import * as fs from "fs";

const parser = pegjs.generate(fs.readFileSync("./src/main/resources/propertyPathGrammar.pegjs").toString());

export interface Property {
    iri: string,
    inverse: boolean,
    transitive: boolean,
    source?: string
}

export interface AndPath {
    and:(AndPath|OrPath|Property)[],
    source?: string
}

export interface OrPath {
    or:(AndPath|OrPath|Property)[],
    source?: string
}

export interface NullPath {
    source?: string
}

export type PropertyPath = Property | AndPath | OrPath | NullPath;

export class PathParser {
    private path: string;
    constructor(path: string) {
        this.path = path;
    }

    parse(): PropertyPath {
        if (this.path == "") {
            return {
                source: ""
            };
        } else {
            const parsed = <PropertyPath>parser.parse(this.path);
            parsed.source = this.path.replace(/\./g,":");
            return parsed;
        }
    }
}