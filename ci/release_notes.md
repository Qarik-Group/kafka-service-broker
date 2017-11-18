* BREAKING: run the web app/service broker using subcommand `run-broker`
* new subcommand `sanity-test-topic-plan` consumes the `topic` service plan credentials JSON and performs sanity test
* new subcommand `sanity-test-shared-plan` consumes the `shared` service plan credentials JSON and performs sanity test
* `bin/sanity-test` will run a temporary broker if `$SANITY_TEST_RUN_BROKER` is set
