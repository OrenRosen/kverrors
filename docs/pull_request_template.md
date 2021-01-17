# Type of change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Refactor
- [ ] Remove unused code
- [ ] Test

# Description, Motivation and Context

Please include a summary of the change and which issue is fixed. Please also include relevant motivation and context. List any dependencies that are required for this change.

# Link to resources

- [Jira Ticket]()
- [Documentation]()

# Screenshots (if appropriate):

# Testing

- [ ] This code wasn't tested
  - ( Please explain why) 
- [ ] This code was tested
  - (Please explain how)
- [ ] An additional manual test should run before merge+deploy.
- [ ] An additional manual test should run after merge+deploy. 

### Bug fix
- [ ] A unit test was added for proving my fix is effective.
- [ ] An automation test was added for proving my fix is effective.

### New feature
- [ ] A unit test was added for checking new logic.
- [ ] An automation test was added for checking new logic.

### Refactor
- [ ] The refactor is covered fully in unit-test / automation-tests.

### Breaking change
- [ ] New tests were added for checking new logic
- [ ] Should keep the test for deprecated logic until it will be removed.

# Removing unused code

Please explain how you've made sure the removed code isn't used.

# Monitoring

- [ ] The code introduce new warning/error reporting outside of a normal failed request.
  - (Please explain which)
  - [ ] I've made sure the code won't have duplicate error logs.
- [ ] The code will impact system performance
  - [ ] Negative: 
  - [ ] Positive:
  
# dependancies?
- [ ] The code introduce new external dependencies
  - (Please explain which) 

# Is the code introduce new warning/error reporting

Meaning reports that doesn't go by the normal flow of a failed request. 

# Checklist:

- [ ] I have performed a self-review of my own code
- [ ] I have made corresponding changes to the documentation
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published in downstream modules