/// <reference types="cypress" />

describe('generate', () => {
  it('should render generated Rego', () => {
    cy.visit('cypress/html/generate.html')
    cy.get('#generate').should('include.text', "package profile_kiali")
  })
})
