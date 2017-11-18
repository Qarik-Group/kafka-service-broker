package cmd

// BrokerOpts describes the flags/options for the CLI
type BrokerOpts struct {
	Version bool `short:"v" long:"version" description:"Show version"`

	// Commands
	RunBroker           RunBrokerOpts           `command:"run-broker" alias:"b" alias:"bkr" alias:"broker" description:"Run the service broker web app"`
	SanityTestTopicPlan SanityTestTopicPlanOpts `command:"sanity-test-topic-plan" description:"Consume 'topic' service plan credentials JSON via STDIN and perform sanity tests"`
}

// Opts carries all the user provided options (from flags or env vars)
var Opts BrokerOpts
