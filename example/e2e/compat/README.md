# Backwards compatibility e2e testing

1. Find latest tagged revision(s). 
  - Should also be parameterizable.
  - If multiple occur, select largest. 
  - TODO: major version backwards compatibility matrix.
  - If the current commit is the the latest tagged revision, search for an older tagged revision.
2. Create build using tagged revision. Run using tagged revision's configuration.
3. Create build using latest revision. Run using latest revision's configura