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

## Release Steps

### Publishing Release (x.y.z)
**No RC version is published.**

1. Create release branch from develop: `git checkout develop && git pull && git checkout -b release/x.y.z`
2. Edit version files:
   - Update `wrappers/js-web/package.json` version to `x.y.z`
   - Update `wrappers/js/package.json` version to `x.y.z`
3. Ensure Jenkinsfile has `"release/*"` in `Publish artifacts` stage branches
4. Commit: `git commit -m "Publish x.y.z"`
5. Push: `git push -u origin release/x.y.z`
6. Follow the steps to merge into master:
   - Create temporary branch `publish-x.y.z` from `release/x.y.z`
   - Merge master into it: `git merge origin/master`
   - Create PR to master and merge
   - **Go artifact is published via tag**: `git tag x.y.z origin/master && git push origin x.y.z`
7. Merge back to develop with next version if needed

### Publishing Hotfix (x.y.z-n)
1. Checkout/create `support/x.y.z` branch from tag
2. Cherry-pick required commits
3. Update version in both package.json files to `x.y.z-n`
4. Ensure Jenkinsfile has `support/*` in `Publish artifacts` stage (for first HF only)
5. Push and create PR to merge into `support/x.y.z`

## Dependencies Impact
After releasing amf-custom-validator, update the version in:
- amf-custom-validator-scalajs (before its release)
- apb (before its release)

## Notes
- Does not need to follow the release times of the rest of the projects
- Released only when needed
