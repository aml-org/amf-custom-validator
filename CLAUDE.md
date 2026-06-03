# amf-custom-validator

**Repository**: https://github.com/aml-org-emu/amf-custom-validator

## Release Process
**Does NOT follow the regular release process exactly** - no RC versions.

## Release Position
- Released **outside the regular release line** (when needed)
- Must be adopted in **amf-custom-validator-scalajs** before its release if there's a new version

## Artifacts Published
This project publishes **3 artifacts**:
1. **Go artifact** - simply a git tag on a commit
2. **JS web artifact** - published to NPM
3. **JS node artifact** - published to NPM

## Files & Locations

### Version Files
- `wrappers/js-web/package.json` - for the web version
- `wrappers/js/package.json` - for the node version

### Dependencies
None specified in the regular AMF release line

## Jenkins Configuration

### Stages to Modify
- **Publishing stages**: `Publish artifacts`

### Jenkinsfile Updates
- **For Release**: Ensure `"release/*"` is in branch list in `Publish artifacts` stage

## Special Instructions
- **No RC version is published** for this repository
- Go artifact is published via git tag after merging to master
- JS artifacts (web and node) are published to NPM via Jenkins

## Dependencies Impact
After releasing amf-custom-validator, update the version in:
- amf-custom-validator-scalajs (before its release)
- apb (before its release)

## Notes
- Does not need to follow the release times of the rest of the projects
- Released only when needed
