// rollup.config.mjs
import url from "@rollup/plugin-url";
import nodePolyfills from 'rollup-plugin-polyfill-node';

export default {
    input: 'lib/index.js',
    output: {
        file: 'dist/bundle.js',
        format: 'umd',
        name: 'amf-custom-validator',
        inlineDynamicImports: true,
        globals: {
            "pako": "pako",
            "browser-or-node": "browser-or-node"
        }
    },
    external: ['browser-or-node', 'pako'],
    plugins: [
        nodePolyfills(),
        url({

            include: "assets/main.wasm.gz",
            emitFiles: true,
            limit: 1000 * 1000 * 10
        }),
    ]
};