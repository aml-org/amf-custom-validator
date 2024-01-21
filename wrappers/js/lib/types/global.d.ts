export {}

declare global {
    class Go {
        importObject: {
            go: any
        }

        async run(instance, global): void
        exit(code: number): void
    }

    function __AMF__terminateValidator(): void
    function __AMF__generateRego(profile: string): string
    function __AMF__normalizeInput(data: string): string
    function __AMF__validateCustomProfileWithConfiguration(profile: string, data: string, debug: boolean, validationConfig: any | undefined, reportConfig: any): string
    function __AMF__validateCustomProfile(profile: string, data: string, debug: boolean)
}