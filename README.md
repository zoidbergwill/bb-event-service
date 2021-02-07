# NOTICE: This project has been archived!

As discussed in [bb-browser#21](https://github.com/buildbarn/bb-browser/issues/21),
support for Build Event Streams is going to be removed from Buildbarn.
We invite people in the community to build a more complete solution for
ingesting and displaying Build Event Streams!

# Buildbarn Event Service [![Build status](https://badge.buildkite.com/0c488b2f6dc0f10dc8736412d49997b17e2dcb92bc59846a40.svg)](https://buildkite.com/buildbarn/bb-event-service)

Buildbarn Event Service is an implementation of
[Bazel's Build Event Service](https://docs.bazel.build/versions/master/build-event-protocol.html#the-build-event-service).
In short, this service can be used to log and store all sorts of
information regarding invocations of Bazel. Information includes the
list of targets analyzed, whether they could be built, output files that
were generated and test results.

This service heavily builds on top of other Buildbarn components.
Instead of having its own storage layer, it can store its results into
[the Buildbarn storage daemon](https://github.com/buildbarn/bb-storage).
Build event streams are stored in the Content Addressable Storage (CAS),
while entries that permit lookups of these streams by Bazel invocation
ID (a per-invocation UUID) are stored in the Action Cache (AC).

[Buildbarn Browser](https://github.com/buildbarn/bb-browser) has
integrated support for displaying build event streams stored in the
CAS/AC by visiting `/build_events/<instance>/<uuid>`. For builds that
also used remote execution, Buildbarn Browser is capable of linking
directly to output files of build actions, making it a valuable tool for
sharing build results with colleagues. By invoking Bazel with
`--bes_results_url=...`, it may automatically print URLs pointing to
Buildbarn Browser at the start of every invocation.

Prebuilt container images of Buildbarn Event Service may be found on
[Docker Hub](https://hub.docker.com/r/buildbarn/bb-event-service).
Examples of how the Buildbarn Event Service may be deployed and used can
be found in [the Buildbarn deployments repository](https://github.com/buildbarn/bb-deployments).

## Updating the bazel git patch

Inspired by https://github.com/EdSchouten/bazel/commit/ae0bd93386ccf9c7e7481829dc9fbb8179b38396

```
git@github.com:zoidbergwill/bazel.git
cd bazel
git remote add upstream git@github.com:bazelbuild/bazel.git
git checkout zoid/bb-event-service
git diff-index -p upstream/master . > build_event_stream.diff
git pull upstream master
```

A better guide lives in rules_go [here](https://github.com/bazelbuild/rules_go/wiki/Updating-dependencies#how-to-update)

## Stuff I manually need to sprinkle in redis v8 and redisext at the moment for some reason

```
        "@io_opentelemetry_go_otel//api/global",
        "@io_opentelemetry_go_otel//api/metric",
        "@io_opentelemetry_go_otel//api/trace",
```

A better option might be teaching bb_storage's go_dependencies about ignoring already imported repos, and then changing redis/v8 and redisext to be http_archive's with patches that add correct BUILD files.
