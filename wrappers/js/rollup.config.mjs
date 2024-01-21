// rollup.config.mjs
import url from "@rollup/plugin-url";
import nodePolyfills from 'rollup-plugin-polyfill-node';
import typescript from "@rollup/plugin-typescript";

export default {
    input: 'lib/index.ts',
    output: {
        file: 'dist/bundle.js',
        format: 'umd',
        name: 'amf-custom-validator',
        inlineDynamicImports: true,
        globals: {
            "pako": "pako",
            "browser-or-node": "browser-or-node"
        },
        sourcemap: true
    },
    external: ['browser-or-node', 'pako'],
    plugins: [
        typescript(),
        nodePolyfills(),
        url({

            include: "assets/main.wasm.gz",
            emitFiles: true,
            limit: 1000 * 1000 * 10
        }),
    ]
};