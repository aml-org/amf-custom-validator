const { defineConfig } = require('cypress')

module.exports = defineConfig({
  fixturesFolder: false,
  downloadsFolder: 'cypress/downloads',
  video: false,
  screenshotOnRunFailure: false,
  e2e: {
    setupNodeEvents(on, config) {},
    supportFile: false,
  },
})
