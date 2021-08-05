# 2. Remove focusNode from individual traces

Date: 2021-08-05

## Status

Accepted

## Context

In the validation report we are including both `shacl:focusNode` properties at the individual trace level and at the
validation result level.

From the SHACL spec:
> Each validation result has exactly one value for the property sh:focusNode that is equal to the focus node that has caused the result


## Decision

Remove `sh:focusNode` from individual traces

## Consequences

None
