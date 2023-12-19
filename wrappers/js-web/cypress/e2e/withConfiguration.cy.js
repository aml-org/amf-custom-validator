/// <reference types="cypress" />

describe('withConfiguration', () => {
  it('should render report', () => {
    cy.visit('cypress/html/withConfiguration.html')

    let dateText = "does not have date"
    let reportSchema = "http://a.ml/report#/declarations/"
    let lexicalSchema = "http://a.ml/lexical#/declarations/"
    let expectedText = `${dateText}, ${reportSchema}, ${lexicalSchema}`

    cy.get('#report').should('have.text', expectedText)
  })
})
