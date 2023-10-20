//
// Opentememtry has tree env variable to control the sampling ratio:
//		- OTEL_TRACES_SAMPLER: It has 6 possible value which three of them are important.
//					* `always_on` traces all the requests. always_off traces nothing.
//					* `traceidratio` traces with a ratio.
//				On dev, I suggest using always_on to gather all the data. on production, we can use
//				traceidratio to decrease the load.
//		- OTEL_TRACES_SAMPLER_ARG: It’s float number between 0 and 1. 0 means tracing is off. 1 means tracing
//				100% of the request. It is sampling ratio. e.g., 0.1 means 10% of request will be sampled.
//

package trace
