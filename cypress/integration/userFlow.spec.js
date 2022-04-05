/// <reference types="cypress"/>

describe("UserFlow Test", () => {
    before(() => {
        cy.testRequest("delete", false)
        cy.clearLocalStorage()
    })

    it("Login flow works fully", () => {
        cy.visit("/login")
        // shows login page
        cy.get("img[alt='meiki-logo']").should("be.visible")
        cy.get("#Username").should("be.visible")
        cy.get("#Password").should("be.visible")
        cy.get("Button").should("include.text", "Login").and("be.visible")
        cy.get("a[href='/create']").and("be.visible").click()

        // user creates an account
        cy.get("Button")
            .should("include.text", "Create Meiki account")
            .and("be.visible")
        cy.get("#Username").type("shnoo")
        cy.get("#Password").type("thisisveryunsafe")
        cy.get("Button").click()

        // goes to account creation success page
        cy.contains("Your account has successfully been created").should(
            "be.visible"
        )
        cy.get("a[href='/login']").should("be.visible").click()

        // user logs in
        cy.get("#Username").type("shnoo")
        cy.get("#Password").type("thisisveryunsafe")
        cy.get("Button").click()

        // assert it goes to the app
        cy.get("nav").should("be.visible")
        cy.get("[data-cy='profile']").should("contain", "shnoo")

        // TODO: Add logout
    })
})
