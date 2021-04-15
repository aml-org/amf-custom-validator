export class PathParser {
    private path: string;
    constructor(path: string) {
        this.path = path;
    }

    parse(): string[] {
        this.path = this.path.trim().replace(/ /g, "")
        return this.path.split("/").map((component) => component.replace(".", ":"));
    }
}