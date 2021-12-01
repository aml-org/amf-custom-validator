# How to add a test
1. Create a test suite in `./src`. Copy one of the existing ones
2. Add the newly created suite to `./webpack.config.js`. You need to add another entry to `config.entry` like the others
3. Create an HTML file in `./server` that imports the built version of your suite (located in `./suites/out`). Use the others for reference
4. Add your new HTML file in the server in `./server/server.js`
5. Create a Cypress spec in `../cypress/integration`. You can use the assertion that you like, use the others for reference

# How to run tests
1. `npm install`
2. `npm run build:dist` this builds the validator (need to re-run with every change in the validator code)
3. `npm run build:test` this builds the test suites (need to re-run with every change in the test suites code)
4. `npm run serve:test:dev` this starts the server and listens to changes in the server code (no need to re-run)
5. `npm test`

# How does it work?
Cypress runs on different browsers. For that you need to serve an HTML file that imports your test script (as would run in any browser). 
We write suites using ES6 and Node syntax. Therefore, we need to build not only the validator but also the test suites using Webpack.
Once 
